# テストガイド - SLM ハンズオン バックエンド

このドキュメントは、SLM ハンズオン バックエンドプロジェクトのテスト戦略、実行方法、および実装内容について説明します。

## 📋 目次

- [テスト戦略概要](#テスト戦略概要)
- [テスト構成](#テスト構成)
- [実行方法](#実行方法)
- [各層のテスト詳細](#各層のテスト詳細)
- [CI/CD統合](#cicd統合)
- [トラブルシューティング](#トラブルシューティング)

## 🎯 テスト戦略概要

このプロジェクトは**クリーンアーキテクチャ**に基づいて設計されており、各層に適切なテストが実装されています。

### テストピラミッド

```
    🔺 E2E/統合テスト
   📊 Interface層テスト
  ⚙️ UseCase層テスト  
 🗄️ Infrastructure層テスト
🏛️ Domain層テスト（最も重要）
```

### 品質目標

- **テストカバレッジ**: 80%以上
- **テスト実行時間**: 2分以内
- **信頼性**: 非決定的テストの排除
- **保守性**: 読みやすく理解しやすいテスト

## 🏗️ テスト構成

### 1. Domain層テスト 🏛️

**場所**: `internal/domain/entity/*_test.go`  
**対象**: ビジネスエンティティとビジネスルール  
**重要度**: ★★★★★

```go
// 例: Product エンティティのテスト
func TestProduct_DecreaseStock(t *testing.T) {
    product := entity.NewProduct("テスト商品", "説明", 1000, "/image.jpg", 10)
    
    err := product.DecreaseStock(5)
    assert.NoError(t, err)
    assert.Equal(t, 5, product.Stock)
}
```

**テスト内容**:
- ✅ 商品の在庫管理ロジック
- ✅ カートのアイテム追加・削除・更新
- ✅ 注文のステータス遷移
- ✅ ビジネスルールの検証

### 2. UseCase層テスト ⚙️

**場所**: `internal/usecase/*_test.go`  
**対象**: ビジネスロジックと UseCase  
**重要度**: ★★★★☆

```go
// 例: ProductUseCase のテスト（モック使用）
func TestProductUseCase_GetAllProducts(t *testing.T) {
    mockRepo := &mocks.ProductRepositoryMock{}
    useCase := usecase.NewProductUseCase(mockRepo)
    
    mockRepo.SetGetAllResponse(expectedProducts, nil)
    
    products, err := useCase.GetAllProducts(context.Background())
    assert.NoError(t, err)
    assert.Len(t, products, 2)
}
```

**テスト内容**:
- ✅ 商品取得・検索ロジック
- ✅ カート操作ロジック
- ✅ 注文作成・管理ロジック
- ✅ エラーハンドリング
- ✅ モックを使用した依存関係の分離

### 3. Infrastructure層テスト 🗄️

**場所**: `internal/infrastructure/persistence/memory/*_test.go`  
**対象**: データ永続化とリポジトリ実装  
**重要度**: ★★★☆☆

```go
// 例: ProductRepository のテスト
func TestProductRepository_GetByID(t *testing.T) {
    repo := memory.NewProductRepository()
    
    product, err := repo.GetByID(context.Background(), productID)
    assert.NoError(t, err)
    assert.Equal(t, productID, product.ID)
}
```

**テスト内容**:
- ✅ CRUD操作の正確性
- ✅ 並行アクセス時の安全性
- ✅ データ整合性の確保
- ✅ 大量データでのパフォーマンス

### 4. Interface層テスト 🌐

**場所**: `internal/interface/api/handler/*_test.go`  
**対象**: HTTPハンドラーとAPIエンドポイント  
**重要度**: ★★★☆☆

```go
// 例: HealthHandler のテスト
func TestHealthHandler_HealthCheck(t *testing.T) {
    router := gin.New()
    handler := handler.NewHealthHandler()
    router.GET("/health", handler.HealthCheck)
    
    req, _ := http.NewRequest("GET", "/health", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
```

**テスト内容**:
- ✅ HTTPレスポンスの検証
- ✅ ステータスコードの確認
- ✅ JSON構造の検証
- ✅ エラーレスポンスの検証

### 5. 統合テスト (E2E) 🔄

**場所**: `test/integration/*_test.go`  
**対象**: エンドツーエンドのユーザージャーニー  
**重要度**: ★★★★☆

```go
// 例: 完全なECサイトフローのテスト
func TestE2E_CompleteUserJourney(t *testing.T) {
    app := setupTestApplication()
    
    // 1. 商品一覧取得
    // 2. 商品詳細取得  
    // 3. カートに追加
    // 4. 注文作成
    // 5. 注文確認
}
```

**テスト内容**:
- ✅ 完全なユーザージャーニー
- ✅ 複数APIの連携動作
- ✅ エラーシナリオ
- ✅ 並行アクセステスト

## 🚀 実行方法

### クイックスタート

```bash
# 1. 全テスト実行（推奨）
make test

# 2. Docker環境でテスト実行
make test-docker

# 3. カバレッジ付きテスト
make test-coverage
```

### 詳細な実行オプション

#### 1. レイヤー別テスト実行

```bash
# Domain層のみ
go test ./internal/domain/entity/... -v

# UseCase層のみ  
go test ./internal/usecase/... -v

# Infrastructure層のみ
go test ./internal/infrastructure/... -v

# Interface層のみ
go test ./internal/interface/... -v

# 統合テストのみ
go test ./test/integration/... -v
```

#### 2. 特定テストの実行

```bash
# 特定のテスト関数
go test ./internal/domain/entity -run TestProduct_DecreaseStock -v

# 特定のテストファイル
go test ./internal/domain/entity/product_test.go -v
```

#### 3. カバレッジレポート

```bash
# カバレッジ測定
go test ./... -coverprofile=coverage.out

# カバレッジ表示  
go tool cover -func=coverage.out

# HTMLレポート生成
go tool cover -html=coverage.out -o coverage.html
```

#### 4. Docker環境でのテスト

```bash
# Docker環境で全テスト実行
./scripts/test-docker.sh

# または
make test-docker
```

### 継続的インテグレーション

GitHub Actions ワークフローが `.github/workflows/test.yml` に設定されています：

- ✅ 複数Goバージョンでのテスト
- ✅ 並列実行による高速化
- ✅ カバレッジレポート生成
- ✅ セキュリティスキャン
- ✅ コード品質チェック

## 📊 テスト実行例

### 成功例

```bash
$ make test

🧪 SLM ハンズオン バックエンド テストスイート

🏛️  Domain層テスト
📋 Domain Entities: ビジネスエンティティとビジネスルールのテスト
✅ Domain Entities: PASS

⚙️  UseCase層テスト  
📋 UseCase Business Logic: ビジネスロジックとユースケースのテスト
✅ UseCase Business Logic: PASS

🗄️  Infrastructure層テスト
📋 Infrastructure Persistence: データ永続化と外部サービス統合のテスト
✅ Infrastructure Persistence: PASS

🌐 Interface層テスト
📋 Interface Handlers: HTTPハンドラーとAPIエンドポイントのテスト  
✅ Interface Handlers: PASS

🔄 統合テスト (E2E)
📋 E2E Integration: エンドツーエンド統合テスト
✅ E2E Integration: PASS

📊 カバレッジテスト
📋 Coverage Report: 全体のテストカバレッジレポート生成
✅ Coverage Report: PASS

🏁 テスト実行結果サマリー
==============================
✅ 成功したテスト (6):
   • Domain Entities
   • UseCase Business Logic  
   • Infrastructure Persistence
   • Interface Handlers
   • E2E Integration
   • Coverage Report

合計テスト: 6
成功: 6
失敗: 0

🎉 全てのテストが成功しました！
```

## 🔧 テスト設定とカスタマイズ

### 環境変数

テスト実行時に以下の環境変数を設定できます：

```bash
# New Relic監視を無効化（テスト用）
NEW_RELIC_API_KEY=""

# テストタイムアウト
GO_TEST_TIMEOUT=30s

# DockerのGoバージョン
GO_VERSION=1.21
```

### テスト用の設定

- **New Relic**: テスト時は無効化され、ログに "NEW_RELIC_API_KEY not set, New Relic monitoring disabled" が表示されます
- **Gin Framework**: テストモードで実行され、詳細なログは非表示になります
- **タイムアウト**: 各テストには適切なタイムアウトが設定されています

## 🐛 トラブルシューティング

### よくある問題と解決方法

#### 1. テストが失敗する

```bash
# 依存関係の問題
go mod tidy
go mod download

# キャッシュクリア
go clean -testcache
```

#### 2. Docker環境でのテスト失敗

```bash
# Dockerイメージの再ビルド
docker system prune -f
make docker-build
```

#### 3. カバレッジレポートが生成されない

```bash
# 手動でカバレッジ測定
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

#### 4. 統合テストでエラーが発生

統合テストは現在制限付きで実行されます。本格的な実行には追加の設定が必要です。

### ログの確認

テスト実行時のログは以下の場所で確認できます：

- **ローカル**: コンソール出力
- **CI/CD**: GitHub Actions のログタブ
- **Docker**: `docker logs <container_name>`

## 🎓 ベストプラクティス

### テスト作成時の指針

1. **AAA パターン**: Arrange, Act, Assert の順序でテストを構成
2. **1テスト1検証**: 1つのテスト関数では1つの観点のみを検証
3. **意味のある名前**: テスト名は「何を」「どのような条件で」「何を期待するか」を表現
4. **独立性**: テスト間で依存関係を持たない
5. **高速実行**: 外部依存を最小限に抑制

### テスト命名規則

```go
// ✅ 良い例
func TestProduct_DecreaseStock_WhenSufficientStock_ShouldUpdateStock(t *testing.T)
func TestCartUseCase_AddToCart_WhenProductNotFound_ShouldReturnError(t *testing.T)

// ❌ 悪い例  
func TestProduct1(t *testing.T)
func TestError(t *testing.T)
```

### モックの使用方針

- **Domain層**: モック不要（純粋な関数）
- **UseCase層**: リポジトリのモックを使用
- **Interface層**: 簡単な構造化レスポンステスト
- **Integration層**: 実際の依存関係を使用

## 📚 関連ドキュメント

- [README.md](./README.md) - プロジェクト概要
- [CLAUDE.md](./CLAUDE.md) - プロジェクト詳細仕様
- [Dockerfile](./Dockerfile) - Docker設定
- [Makefile](./Makefile) - ビルドとテスト自動化

## 🤝 貢献ガイド

新しいテストを追加する場合：

1. 適切な層にテストファイルを作成
2. 命名規則に従ったテスト名を使用
3. テストの目的と期待結果を明確に記述
4. 必要に応じてモックを実装
5. 全テストが成功することを確認

---

このテストガイドは、SLM ハンズオンプロジェクトの品質保証と開発効率向上を目的として作成されました。質問や改善提案がありましたら、メンテナーまでお声かけください。