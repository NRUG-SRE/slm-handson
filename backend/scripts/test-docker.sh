#!/bin/bash

# Docker環境でのテスト実行スクリプト
# CI/CDパイプラインや開発環境での統一テスト実行

set -e

# カラー出力の設定
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🐳 Docker環境でのテスト実行${NC}"
echo "================================"

# プロジェクトのルートディレクトリに移動
cd "$(dirname "$0")/.."

# Docker引数の設定
DOCKER_GO_VERSION=${GO_VERSION:-"1.21"}
DOCKER_ARGS="--rm -v $(pwd):/app -w /app"

echo -e "${YELLOW}📋 設定情報${NC}"
echo "   Goバージョン: $DOCKER_GO_VERSION"
echo "   作業ディレクトリ: $(pwd)"
echo ""

# テスト実行関数（Docker版）
run_docker_test() {
    local test_name="$1"
    local test_path="$2"
    local description="$3"
    
    echo -e "${YELLOW}📋 $test_name${NC}: $description"
    echo "   Docker実行: golang:$DOCKER_GO_VERSION"
    echo "   パス: $test_path"
    
    if docker run $DOCKER_ARGS golang:$DOCKER_GO_VERSION go test "$test_path" -v -timeout 30s; then
        echo -e "${GREEN}✅ $test_name: PASS${NC}"
        echo ""
        return 0
    else
        echo -e "${RED}❌ $test_name: FAIL${NC}"
        echo ""
        return 1
    fi
}

# テスト実行結果を記録
FAILED_TESTS=()
PASSED_TESTS=()

# 依存関係のダウンロード
echo -e "${BLUE}📦 依存関係の解決${NC}"
docker run $DOCKER_ARGS golang:$DOCKER_GO_VERSION go mod download
echo ""

# 1. Domain層のテスト
echo -e "${BLUE}🏛️  Domain層テスト${NC}"
if run_docker_test "Domain Entities" "./internal/domain/entity/..." "ビジネスエンティティとビジネスルールのテスト"; then
    PASSED_TESTS+=("Domain Entities")
else
    FAILED_TESTS+=("Domain Entities")
fi

# 2. UseCase層のテスト
echo -e "${BLUE}⚙️  UseCase層テスト${NC}"
if run_docker_test "UseCase Business Logic" "./internal/usecase/..." "ビジネスロジックとユースケースのテスト"; then
    PASSED_TESTS+=("UseCase Business Logic")
else
    FAILED_TESTS+=("UseCase Business Logic")
fi

# 3. Infrastructure層のテスト
echo -e "${BLUE}🗄️  Infrastructure層テスト${NC}"
if run_docker_test "Infrastructure Persistence" "./internal/infrastructure/..." "データ永続化と外部サービス統合のテスト"; then
    PASSED_TESTS+=("Infrastructure Persistence")
else
    FAILED_TESTS+=("Infrastructure Persistence")
fi

# 4. Interface層のテスト
echo -e "${BLUE}🌐 Interface層テスト${NC}"
if run_docker_test "Interface Handlers" "./internal/interface/..." "HTTPハンドラーとAPIエンドポイントのテスト"; then
    PASSED_TESTS+=("Interface Handlers")
else
    FAILED_TESTS+=("Interface Handlers")
fi

# 5. 統合テスト (制限付き)
echo -e "${BLUE}🔄 統合テスト${NC}"
echo -e "${YELLOW}⚠️  注意: 統合テストは現在制限付きで実行されます${NC}"
if run_docker_test "E2E Integration (Limited)" "./test/integration/..." "統合テスト（制限付き実行）"; then
    PASSED_TESTS+=("E2E Integration")
else
    FAILED_TESTS+=("E2E Integration")
fi

# 6. 全体テスト実行とカバレッジ
echo -e "${BLUE}📊 全体テストとカバレッジ${NC}"
echo -e "${YELLOW}📋 Full Test Suite${NC}: 全てのテストとカバレッジレポート生成"

if docker run $DOCKER_ARGS golang:$DOCKER_GO_VERSION bash -c "
    go test ./... -coverprofile=coverage.out -timeout 60s && 
    go tool cover -func=coverage.out | tail -n 1 &&
    echo 'カバレッジレポートが生成されました: coverage.out'
"; then
    echo -e "${GREEN}✅ Full Test Suite: PASS${NC}"
    PASSED_TESTS+=("Full Test Suite")
else
    echo -e "${RED}❌ Full Test Suite: FAIL${NC}"
    FAILED_TESTS+=("Full Test Suite")
fi

echo ""

# 最終結果のサマリー
echo -e "${BLUE}🏁 Docker テスト実行結果${NC}"
echo "================================"

if [ ${#PASSED_TESTS[@]} -gt 0 ]; then
    echo -e "${GREEN}✅ 成功したテスト (${#PASSED_TESTS[@]}):${NC}"
    for test in "${PASSED_TESTS[@]}"; do
        echo "   • $test"
    done
    echo ""
fi

if [ ${#FAILED_TESTS[@]} -gt 0 ]; then
    echo -e "${RED}❌ 失敗したテスト (${#FAILED_TESTS[@]}):${NC}"
    for test in "${FAILED_TESTS[@]}"; do
        echo "   • $test"
    done
    echo ""
fi

echo "合計テスト: $((${#PASSED_TESTS[@]} + ${#FAILED_TESTS[@]}))"
echo -e "成功: ${GREEN}${#PASSED_TESTS[@]}${NC}"
echo -e "失敗: ${RED}${#FAILED_TESTS[@]}${NC}"

# CI環境の検出
if [ -n "$CI" ]; then
    echo ""
    echo -e "${BLUE}🤖 CI環境で実行中${NC}"
    
    # CI環境向けの追加レポート
    if [ -f "coverage.out" ]; then
        echo -e "${GREEN}📊 カバレッジファイルをアーティファクトとして保存してください: coverage.out${NC}"
    fi
fi

# 終了コード
if [ ${#FAILED_TESTS[@]} -eq 0 ]; then
    echo ""
    echo -e "${GREEN}🎉 Docker環境での全テストが成功しました！${NC}"
    exit 0
else
    echo ""
    echo -e "${RED}💥 ${#FAILED_TESTS[@]} 個のテストが失敗しました。${NC}"
    exit 1
fi