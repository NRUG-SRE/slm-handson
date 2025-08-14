#!/bin/bash

# テスト実行スクリプト
# Clean Architecture に従ったレイヤー別テスト実行

set -e

# カラー出力の設定
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== SLM ハンズオン バックエンド テストスイート ===${NC}"
echo ""

# プロジェクトのルートディレクトリに移動
cd "$(dirname "$0")/.."

# テスト実行関数
run_test() {
    local test_name="$1"
    local test_path="$2"
    local description="$3"
    
    echo -e "${YELLOW}📋 $test_name${NC}: $description"
    echo "   パス: $test_path"
    
    if go test "$test_path" -v -timeout 30s; then
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

# 1. Domain層のテスト (Entity)
echo -e "${BLUE}🏛️  Domain層テスト${NC}"
if run_test "Domain Entities" "./internal/domain/entity/..." "ビジネスエンティティと ビジネスルールのテスト"; then
    PASSED_TESTS+=("Domain Entities")
else
    FAILED_TESTS+=("Domain Entities")
fi

# 2. UseCase層のテスト (ビジネスロジック)
echo -e "${BLUE}⚙️  UseCase層テスト${NC}"
if run_test "UseCase Business Logic" "./internal/usecase/..." "ビジネスロジックと ユースケースのテスト"; then
    PASSED_TESTS+=("UseCase Business Logic")
else
    FAILED_TESTS+=("UseCase Business Logic")
fi

# 3. Infrastructure層のテスト (データ永続化)
echo -e "${BLUE}🗄️  Infrastructure層テスト${NC}"
if run_test "Infrastructure Persistence" "./internal/infrastructure/..." "データ永続化と 外部サービス統合のテスト"; then
    PASSED_TESTS+=("Infrastructure Persistence")
else
    FAILED_TESTS+=("Infrastructure Persistence")
fi

# 4. Interface層のテスト (HTTPハンドラー)
echo -e "${BLUE}🌐 Interface層テスト${NC}"
if run_test "Interface Handlers" "./internal/interface/..." "HTTPハンドラーと APIエンドポイントのテスト"; then
    PASSED_TESTS+=("Interface Handlers")
else
    FAILED_TESTS+=("Interface Handlers")
fi

# 5. 統合テスト (E2E)
echo -e "${BLUE}🔄 統合テスト (E2E)${NC}"
if run_test "E2E Integration" "./test/integration/..." "エンドツーエンド統合テスト"; then
    PASSED_TESTS+=("E2E Integration")
else
    FAILED_TESTS+=("E2E Integration")
fi

# 6. 全体テスト実行（カバレッジ付き）
echo -e "${BLUE}📊 カバレッジテスト${NC}"
echo -e "${YELLOW}📋 Coverage Report${NC}: 全体のテストカバレッジレポート生成"

if go test ./... -coverprofile=coverage.out -timeout 60s; then
    echo -e "${GREEN}✅ Coverage Report: PASS${NC}"
    
    # カバレッジレポートの表示
    echo ""
    echo -e "${BLUE}📈 テストカバレッジサマリー${NC}"
    go tool cover -func=coverage.out | tail -n 1
    
    # HTMLレポート生成
    go tool cover -html=coverage.out -o coverage.html
    echo -e "${GREEN}📄 HTMLカバレッジレポート: coverage.html${NC}"
    
    PASSED_TESTS+=("Coverage Report")
else
    echo -e "${RED}❌ Coverage Report: FAIL${NC}"
    FAILED_TESTS+=("Coverage Report")
fi

echo ""

# 最終結果のサマリー
echo -e "${BLUE}🏁 テスト実行結果サマリー${NC}"
echo "=============================="

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

# 終了コード
if [ ${#FAILED_TESTS[@]} -eq 0 ]; then
    echo ""
    echo -e "${GREEN}🎉 全てのテストが成功しました！${NC}"
    exit 0
else
    echo ""
    echo -e "${RED}💥 ${#FAILED_TESTS[@]} 個のテストが失敗しました。${NC}"
    exit 1
fi