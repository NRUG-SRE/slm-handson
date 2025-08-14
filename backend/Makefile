# SLM ハンズオン バックエンド Makefile
# テスト実行とCI/CD タスクの自動化

.PHONY: help test test-docker test-unit test-integration test-coverage \
        build clean lint fmt vet security deps check-deps \
        docker-build docker-test docker-run coverage-html

# デフォルトのGoバージョン
GO_VERSION ?= 1.21

# Docker設定
DOCKER_IMAGE_NAME = slm-handson-backend
DOCKER_TAG ?= latest

# カバレッジ設定
COVERAGE_FILE = coverage.out
COVERAGE_HTML = coverage.html

help: ## ヘルプメッセージを表示
	@echo "SLM ハンズオン バックエンド - 利用可能なコマンド:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "テスト関連:"
	@echo "  make test              - 全てのテストを実行"
	@echo "  make test-unit         - ユニットテストのみ実行"
	@echo "  make test-integration  - 統合テストのみ実行"
	@echo "  make test-coverage     - カバレッジ付きテスト実行"
	@echo ""
	@echo "Docker関連:"
	@echo "  make docker-test       - Docker環境でテスト実行"
	@echo "  make docker-build      - Dockerイメージをビルド"
	@echo "  make docker-run        - Dockerコンテナを実行"

# =============================================================================
# 依存関係管理
# =============================================================================

deps: ## 依存関係をダウンロード
	@echo "📦 依存関係をダウンロード中..."
	go mod download
	go mod verify

check-deps: ## 依存関係の確認
	@echo "🔍 依存関係の確認..."
	go mod verify
	go mod tidy

# =============================================================================
# ビルド関連
# =============================================================================

build: ## アプリケーションをビルド
	@echo "🔨 アプリケーションビルド中..."
	@if command -v go >/dev/null 2>&1; then \
		go build -v -o bin/server ./cmd/server; \
	else \
		echo "⚠️  Goがインストールされていません。Dockerを使用します..."; \
		docker run --rm -v $(PWD):/app -w /app golang:$(GO_VERSION) go build -v -o bin/server ./cmd/server; \
	fi

clean: ## ビルド成果物をクリーンアップ
	@echo "🧹 クリーンアップ中..."
	rm -rf bin/
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	docker rmi $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) 2>/dev/null || true

# =============================================================================
# テスト関連
# =============================================================================

test: ## 全てのテストを実行（推奨）
	@echo "🧪 全体テスト実行中..."
	./scripts/test.sh

test-docker: ## Docker環境でテスト実行
	@echo "🐳 Docker環境でテスト実行中..."
	./scripts/test-docker.sh

test-unit: ## ユニットテストのみ実行
	@echo "🧪 ユニットテスト実行中..."
	@if command -v go >/dev/null 2>&1; then \
		echo "Goを使用してテスト実行..."; \
		echo "Domain層テスト..."; \
		go test -v -timeout 30s ./internal/domain/entity/...; \
		echo ""; \
		echo "UseCase層テスト..."; \
		go test -v -timeout 30s ./internal/usecase/...; \
		echo ""; \
		echo "Infrastructure層テスト..."; \
		go test -v -timeout 30s ./internal/infrastructure/...; \
		echo ""; \
		echo "Interface層テスト..."; \
		go test -v -timeout 30s ./internal/interface/...; \
	else \
		echo "⚠️  Goがインストールされていません。Dockerを使用します..."; \
		echo "Domain層テスト..."; \
		docker run --rm -v $(PWD):/app -w /app golang:$(GO_VERSION) go test -v -timeout 30s ./internal/domain/entity/...; \
		echo ""; \
		echo "UseCase層テスト..."; \
		docker run --rm -v $(PWD):/app -w /app golang:$(GO_VERSION) go test -v -timeout 30s ./internal/usecase/...; \
		echo ""; \
		echo "Infrastructure層テスト..."; \
		docker run --rm -v $(PWD):/app -w /app golang:$(GO_VERSION) go test -v -timeout 30s ./internal/infrastructure/...; \
		echo ""; \
		echo "Interface層テスト..."; \
		docker run --rm -v $(PWD):/app -w /app golang:$(GO_VERSION) go test -v -timeout 30s ./internal/interface/...; \
	fi

test-integration: ## 統合テストのみ実行
	@echo "🔄 統合テスト実行中..."
	@if command -v go >/dev/null 2>&1; then \
		go test -v -timeout 60s ./test/integration/...; \
	else \
		echo "⚠️  Goがインストールされていません。Dockerを使用します..."; \
		docker run --rm -v $(PWD):/app -w /app golang:$(GO_VERSION) go test -v -timeout 60s ./test/integration/...; \
	fi

test-coverage: $(COVERAGE_FILE) ## カバレッジ付きテスト実行
	@echo "📊 カバレッジレポート生成中..."
	go tool cover -func=$(COVERAGE_FILE) | tail -n 1

$(COVERAGE_FILE): ## カバレッジファイル生成
	@echo "📊 カバレッジ測定中..."
	go test -coverprofile=$(COVERAGE_FILE) -timeout 60s ./...

coverage-html: $(COVERAGE_FILE) ## HTMLカバレッジレポート生成
	@echo "🌐 HTMLカバレッジレポート生成中..."
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "📄 HTMLレポートが生成されました: $(COVERAGE_HTML)"

# =============================================================================
# コード品質
# =============================================================================

lint: ## コードリンティング実行
	@echo "🔍 コードリンティング中..."
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint がインストールされていません。インストール方法:"; \
		echo "curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin v1.54.2"; \
		exit 1; \
	}
	golangci-lint run --timeout=5m

fmt: ## コードフォーマット実行
	@echo "📝 コードフォーマット中..."
	gofmt -w .
	go mod tidy

vet: ## go vet 実行
	@echo "🔍 go vet 実行中..."
	go vet ./...

check-fmt: ## フォーマットチェック
	@echo "📝 フォーマットチェック中..."
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "以下のファイルがフォーマットされていません:"; \
		gofmt -l .; \
		echo "実行してください: make fmt"; \
		exit 1; \
	else \
		echo "✅ 全ファイルが適切にフォーマットされています"; \
	fi

security: ## セキュリティスキャン実行
	@echo "🔒 セキュリティスキャン中..."
	@command -v gosec >/dev/null 2>&1 || { \
		echo "gosec がインストールされていません。インストール方法:"; \
		echo "go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
		exit 1; \
	}
	gosec ./...

# =============================================================================
# Docker関連
# =============================================================================

docker-build: ## Dockerイメージをビルド
	@echo "🐳 Dockerイメージビルド中..."
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) .

docker-run: docker-build ## Dockerコンテナを実行
	@echo "🚀 Dockerコンテナ起動中..."
	docker run --rm -p 8080:8080 \
		-e NEW_RELIC_API_KEY="" \
		-e NEW_RELIC_APP_NAME="slm-handson-api-local" \
		$(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

docker-test-run: docker-build ## テスト用Dockerコンテナを実行
	@echo "🧪 テスト用Dockerコンテナ起動中..."
	docker run --rm -p 8080:8080 --name slm-test-app \
		-e NEW_RELIC_API_KEY="" \
		-e NEW_RELIC_APP_NAME="slm-handson-api-test" \
		$(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

# =============================================================================
# CI/CD関連
# =============================================================================

ci-test: deps check-fmt vet test-coverage ## CI環境用テスト実行
	@echo "🤖 CI環境でのテスト実行完了"

ci-test-full: deps build check-fmt vet lint security test-coverage ## CI環境用フルテスト実行
	@echo "🤖 CI環境でのフルテスト実行完了"

# =============================================================================
# 開発支援
# =============================================================================

dev-setup: ## 開発環境セットアップ
	@echo "⚙️ 開発環境セットアップ中..."
	go mod download
	go mod verify
	@echo "✅ 開発環境の準備が完了しました"
	@echo ""
	@echo "次のステップ:"
	@echo "1. テスト実行: make test"
	@echo "2. アプリケーション起動: make docker-run"
	@echo "3. 開発用コマンド: make help"

run: build ## ローカルでアプリケーション実行
	@echo "🚀 ローカルでアプリケーション起動中..."
	./bin/server

# =============================================================================
# 情報表示
# =============================================================================

info: ## プロジェクト情報を表示
	@echo "📋 SLM ハンズオン バックエンド プロジェクト情報"
	@echo "================================"
	@echo "Goバージョン: $$(go version)"
	@echo "プロジェクトディレクトリ: $$(pwd)"
	@echo "Dockerイメージ: $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)"
	@echo ""
	@echo "アーキテクチャ:"
	@echo "- Domain層: ビジネスエンティティとルール"
	@echo "- UseCase層: ビジネスロジック"
	@echo "- Interface層: HTTPハンドラー"
	@echo "- Infrastructure層: データ永続化"
	@echo ""
	@echo "テスト構成:"
	@echo "- ユニットテスト: 各層のテスト"
	@echo "- 統合テスト: E2Eテスト"
	@echo "- カバレッジ: 全体的なテストカバレッジ"

# スクリプトの実行権限を確保
scripts/test.sh:
	chmod +x scripts/test.sh

scripts/test-docker.sh:
	chmod +x scripts/test-docker.sh