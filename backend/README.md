# SLM Handson Backend - Go APIサーバー

## 概要

このバックエンドアプリケーションは、New Relic Service Level Management (SLM) のハンズオン用に設計されたECサイトのAPIサーバーです。クリーンアーキテクチャの原則に従って実装され、New Relic APMと統合してサーバーサイドのパフォーマンス監視を実現しています。

## アーキテクチャ概要

### クリーンアーキテクチャの採用理由

本プロジェクトでは以下の理由からクリーンアーキテクチャを採用しています：

1. **ビジネスロジックの独立性**: ドメイン層がフレームワークやデータベースに依存しない
2. **テスタビリティ**: 各層が独立しているため、単体テストが容易
3. **保守性**: 責務が明確に分離され、変更の影響範囲を限定できる
4. **柔軟性**: インフラ層の実装を容易に切り替え可能（例：インメモリDB → RDB）

### 依存関係の方向

```
外部インターフェース（HTTP/CLI）
        ↓
interface層（ハンドラー/プレゼンター）
        ↓
usecase層（アプリケーションビジネスロジック）
        ↓
domain層（エンティティ/ビジネスルール）
        ↑
infrastructure層（外部サービス実装）
```

重要な原則：
- **依存性逆転の原則（DIP）**: 上位層は下位層に依存せず、両者とも抽象に依存
- **内側の層は外側を知らない**: domain層は他の層の存在を知らない
- **インターフェースは内側で定義**: リポジトリインターフェースはdomain層で定義

## ディレクトリ構成

```
backend/
├── cmd/
│   └── server/
│       └── main.go                # アプリケーションエントリーポイント
│                                  # - サーバー起動、DI設定、グレースフルシャットダウン
│
├── internal/                      # 外部パッケージから参照されない内部実装
│   │
│   ├── domain/                   # 【コア層】ビジネスの中核
│   │   ├── entity/               # ビジネスエンティティ
│   │   │   ├── product.go       # 商品エンティティ（ID、名前、価格、在庫等）
│   │   │   ├── cart.go          # カートエンティティ（商品と数量のマップ）
│   │   │   ├── order.go         # 注文エンティティ（注文詳細、合計金額等）
│   │   │   └── errors.go        # ドメイン固有のエラー定義
│   │   │
│   │   └── repository/           # リポジトリインターフェース（抽象）
│   │       ├── product_repository.go  # 商品データアクセスの抽象定義
│   │       ├── cart_repository.go     # カートデータアクセスの抽象定義
│   │       └── order_repository.go    # 注文データアクセスの抽象定義
│   │
│   ├── usecase/                  # 【アプリケーション層】ビジネスユースケース
│   │   ├── product_usecase.go   # 商品関連のビジネスロジック
│   │   │                        # - 商品一覧取得、詳細取得
│   │   ├── cart_usecase.go      # カート操作のビジネスロジック
│   │   │                        # - 商品追加、数量変更、削除、合計計算
│   │   └── order_usecase.go     # 注文処理のビジネスロジック
│   │                            # - 注文作成、在庫確認、カートクリア
│   │
│   ├── interface/               # 【インターフェースアダプター層】外部との境界
│   │   └── api/                # HTTP API実装
│   │       ├── router.go       # ルーティング設定、ミドルウェア適用
│   │       │
│   │       ├── handler/        # HTTPハンドラー（コントローラー）
│   │       │   ├── health_handler.go   # ヘルスチェック（/health）
│   │       │   ├── product_handler.go  # 商品API（GET /api/products/*）
│   │       │   ├── cart_handler.go     # カートAPI（GET/POST/PUT /api/cart/*）
│   │       │   ├── order_handler.go    # 注文API（GET/POST /api/orders）
│   │       │   ├── swagger_handler.go  # API仕様書配信（/api/docs）
│   │       │   └── constants.go        # ハンドラー共通の定数定義
│   │       │
│   │       ├── middleware/     # HTTPミドルウェア
│   │       │   ├── cors.go           # CORS設定（クロスオリジン対応）
│   │       │   └── monitoring.go     # New Relic APMトランザクション追跡
│   │       │
│   │       └── presenter/      # レスポンスフォーマッター
│   │           └── response.go       # 統一的なJSONレスポンス生成
│   │
│   └── infrastructure/         # 【インフラストラクチャ層】技術的詳細
│       ├── persistence/        # データ永続化の実装
│       │   └── memory/        # インメモリDB実装（デモ用）
│       │       ├── product_repository.go  # 商品リポジトリ実装
│       │       ├── cart_repository.go     # カートリポジトリ実装
│       │       └── order_repository.go    # 注文リポジトリ実装
│       │
│       └── monitoring/        # 監視・計測
│           └── newrelic.go    # New Relic APMエージェント初期化
│
└── pkg/                       # 外部パッケージから参照可能な共有コード
    ├── config/
    │   └── config.go         # 環境変数管理、設定値の構造体
    │
    └── utils/
        └── random.go         # ランダム遅延生成（SLO違反シミュレーション用）
```

## 各層の責務と実装詳細

### 1. Domain層（エンティティとビジネスルール）

**責務**: ビジネスの核となる概念とルールを表現

- **依存**: なし（最も内側の層）
- **実装内容**:
  - `entity/`: ビジネスオブジェクト（Product, Cart, Order）
  - `repository/`: データアクセスの抽象インターフェース

```go
// domain/entity/product.go の例
type Product struct {
    ID          string
    Name        string
    Description string
    Price       float64
    Stock       int
    ImageURL    string
}

// domain/repository/product_repository.go の例
type ProductRepository interface {
    FindAll() ([]*entity.Product, error)
    FindByID(id string) (*entity.Product, error)
    UpdateStock(id string, quantity int) error
}
```

### 2. UseCase層（アプリケーションビジネスロジック）

**責務**: アプリケーション固有のビジネスロジックを実装

- **依存**: Domain層のみ
- **実装内容**:
  - トランザクション制御
  - 複数エンティティの協調
  - ビジネスルールの実行

```go
// usecase/order_usecase.go の例
type OrderUseCase struct {
    orderRepo   repository.OrderRepository
    productRepo repository.ProductRepository
    cartRepo    repository.CartRepository
}

func (u *OrderUseCase) CreateOrder(cartID string) (*entity.Order, error) {
    // 1. カート取得
    // 2. 在庫確認
    // 3. 注文作成
    // 4. 在庫更新
    // 5. カートクリア
}
```

### 3. Interface層（コントローラー/プレゼンター）

**責務**: 外部インターフェースとUseCaseの橋渡し

- **依存**: UseCase層、Domain層
- **実装内容**:
  - HTTPリクエスト/レスポンス処理
  - 入力値バリデーション
  - エラーハンドリング
  - レスポンスフォーマット変換

```go
// interface/api/handler/product_handler.go の例
type ProductHandler struct {
    useCase *usecase.ProductUseCase
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
    products, err := h.useCase.GetAllProducts()
    // HTTPレスポンスへの変換
}
```

### 4. Infrastructure層（技術的実装）

**責務**: 技術的な詳細の実装

- **依存**: Domain層（インターフェースを実装）
- **実装内容**:
  - データベースアクセス（現在はインメモリ）
  - 外部API連携
  - モニタリング（New Relic APM）

```go
// infrastructure/persistence/memory/product_repository.go の例
type ProductRepository struct {
    mu       sync.RWMutex
    products map[string]*entity.Product
}

func (r *ProductRepository) FindAll() ([]*entity.Product, error) {
    // インメモリストアから商品一覧を返す
}
```

## 主要機能

### APIエンドポイント

| エンドポイント | メソッド | 説明 |
|-------------|---------|------|
| `/health` | GET | ヘルスチェック |
| `/api/products` | GET | 商品一覧取得 |
| `/api/products/{id}` | GET | 商品詳細取得 |
| `/api/cart` | GET | カート内容取得 |
| `/api/cart/items` | POST | カートに商品追加 |
| `/api/cart/items/{id}` | PUT | カート内商品の数量変更 |
| `/api/cart/items/{id}` | DELETE | カート内商品の削除 |
| `/api/orders` | GET | 注文一覧取得 |
| `/api/orders` | POST | 注文作成 |
| `/api/docs` | GET | Swagger UI |

### New Relic APM統合

- **自動計測**: HTTPトランザクション、データベースクエリ（将来）
- **カスタムセグメント**: ビジネスロジックの詳細追跡
- **エラー追跡**: アプリケーションエラーの自動収集
- **パフォーマンスメトリクス**: レスポンスタイム、スループット

### パフォーマンス調整機能（SLOデモ用）

環境変数による動的な振る舞い変更：

| 環境変数 | 説明 | デフォルト値 |
|---------|------|------------|
| `ERROR_RATE` | エラー発生率（0.0-1.0） | 0.1 |
| `RESPONSE_TIME_MIN` | 最小レスポンス時間（ms） | 50 |
| `RESPONSE_TIME_MAX` | 最大レスポンス時間（ms） | 500 |
| `SLOW_ENDPOINT_RATE` | 遅延エンドポイント発生率 | 0.2 |

## 起動方法

### ローカル開発

```bash
# 依存関係のインストール
go mod download

# アプリケーション起動
go run cmd/server/main.go
```

### Docker使用

```bash
# イメージビルド
docker build -t slm-handson-api .

# コンテナ起動
docker run -p 8080:8080 \
  -e NEW_RELIC_API_KEY=your-key \
  -e NEW_RELIC_APP_NAME=slm-handson-api \
  slm-handson-api
```

### Docker Compose（推奨）

```bash
# プロジェクトルートから
docker-compose up -d api-server
```

## 開発ガイドライン

### 新機能追加時の手順

1. **Domain層**: 必要に応じてエンティティやリポジトリインターフェースを定義
2. **UseCase層**: ビジネスロジックを実装
3. **Infrastructure層**: リポジトリの具体実装を作成
4. **Interface層**: HTTPハンドラーを実装
5. **Router**: エンドポイントを登録
6. **main.go**: 必要に応じてDI設定を追加

### コーディング規約

- **エラーハンドリング**: 早期リターンパターンを使用
- **並行処理**: sync.Mutex/RWMutexで適切に保護
- **コンテキスト**: context.Contextを適切に伝播
- **ログ**: 構造化ログを使用（将来的に）

### テスト方針

- **単体テスト**: 各層ごとにモックを使用してテスト
- **統合テスト**: HTTPハンドラーレベルでのテスト
- **E2Eテスト**: Docker Composeで全体を起動してテスト

## トラブルシューティング

### New Relic APMにデータが表示されない

1. `NEW_RELIC_API_KEY`が正しく設定されているか確認
2. `NEW_RELIC_APP_NAME`が設定されているか確認
3. ネットワーク接続を確認（New Relicエンドポイントへの接続）

### パフォーマンスが不安定

1. `ERROR_RATE`などの環境変数を確認
2. `docker logs api-server`でエラーログを確認
3. メモリ使用量を確認（インメモリDBのため）

## ライセンス

このプロジェクトは教育目的のデモアプリケーションです。