# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

これはNew Relic Service Level Management (SLM)のハンズオンプロジェクトで、Go APIサーバーとNew Relic APM統合を使用してSLO/SLI管理をデモンストレーションします。デプロイにはDocker Composeを使用します。

## 現在の状態

**🎉 プロジェクト完全実装済み** - 本番レベルのSLMハンズオン環境が稼働中：

### ✅ 実装完了（本番レベル）
- **Go APIサーバー**: New Relic APM統合、全APIエンドポイント、Swagger UI
- **Next.jsフロントエンド**: 全ページ実装、New Relic RUM統合、ECサイト機能完備
- **Docker構成**: マルチステージビルド、ヘルスチェック、環境変数対応
- **New Relic統合**: APM + RUM エンドツーエンド監視、セキュアな環境変数設定
- **SLMデモ機能**: エラー率・レスポンス時間調整、パフォーマンス劣化シミュレーション
- **ユーザージャーニー負荷生成**: Go実装、完全なECサイトフロー自動実行、リアルなユーザー行動シミュレーション
- **セキュリティ**: ハードコーディング排除、環境変数ベース設定

### 🎯 ハンズオン対応状況
- **環境セットアップ**: `.env`設定 + `docker-compose up -d`で完全起動
- **APM監視**: Go APIサーバーの全トランザクション・エラー・パフォーマンス計測
- **RUM監視**: Next.jsフロントエンドのページビュー・Ajax・エラー・Core Web Vitals計測
- **SLO設定**: New Relic UIでの実データベースSLI/SLO管理
- **パフォーマンス変化体験**: 環境変数による動的制御（ERROR_RATE等）
- **自動ユーザーアクセス**: プロファイル起動で継続的なユーザージャーニー実行

### ⚠️ 未実装（ハンズオン実施には不要）
- **詳細ドキュメント**: アーキテクチャ詳細説明（README.mdで十分）

## アプリケーション動作環境

### Docker Compose構成
```yaml
services:
  # APIサーバー (Go)
  api-server:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - NEW_RELIC_API_KEY
      - NEW_RELIC_APP_NAME=slm-handson-api
      - ERROR_RATE
      - RESPONSE_TIME_MIN
      - RESPONSE_TIME_MAX
      - SLOW_ENDPOINT_RATE
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # フロントエンドサーバー (Next.js)
  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_BASE_URL=http://api-server:8080/api
      - NEXT_PUBLIC_NEW_RELIC_BROWSER_KEY
      - NEXT_PUBLIC_NEW_RELIC_ACCOUNT_ID
      - NEXT_PUBLIC_NEW_RELIC_APPLICATION_ID
    depends_on:
      - api-server

  # SLOハンズオン用ユーザーアクセス生成器
  load-generator:
    build: ./scripts
    environment:
      - TARGET_URL=http://frontend:3000
      - ACCESS_INTERVAL=10
      - DURATION=3600
    depends_on:
      - frontend
      - api-server
    profiles:
      - load-test
```

### ユーザージャーニー負荷生成スクリプト仕様 ✅実装済み
- **実装言語**: Go（プロジェクト全体との一貫性）
- **完全なECサイトユーザージャーニー**:
  1. **TOPページ訪問** → GET /api/products（商品一覧取得）
  2. **商品詳細ページ表示** → GET /api/products/{id}（ランダム商品選択）
  3. **カート追加** → POST /api/cart/items（1-3個ランダム数量）
  4. **カートページ確認** → GET /api/cart（カート内容確認）
  5. **決済ページ表示** → GET /api/cart（注文内容再確認）
  6. **注文確定** → POST /api/orders（注文作成・決済完了）
- **リアルなユーザー行動シミュレーション**:
  - 各ステップ間で1-5秒のランダム思考時間
  - 商品選択、数量選択のランダム化
  - エラー処理と統計レポート（完了率、成功率）
- **設定可能パラメータ**:
  - `ACCESS_INTERVAL`: ジャーニー間隔（秒、デフォルト：10）
  - `DURATION`: 実行時間（秒、デフォルト：3600）
  - `TARGET_URL`: フロントエンドURL（Docker：http://frontend:3000）

## 開発コマンド

プロジェクトが実装された後、以下のコマンドを使用します：

```bash
# アプリケーション全体の起動
docker-compose up -d

# 個別サービスの起動
docker-compose up -d api-server
docker-compose up -d frontend

# ユーザージャーニー負荷生成の開始（プロファイル指定）
docker-compose --profile load-test up load-generator

# カスタム設定での実行
DURATION=300 ACCESS_INTERVAL=5 docker-compose --profile load-test up load-generator

# ログの確認
docker-compose logs -f
docker-compose logs -f api-server
docker-compose logs -f frontend

# アプリケーションの停止
docker-compose down

# ビルドして起動
docker-compose up --build -d

# 環境変数を.envファイルで設定して実行
# .envファイルを作成・編集後に実行
docker-compose up --build -d
```

## プロジェクト全体のディレクトリ構成

```
slm-handson/
├── docker-compose.yml              # Docker Compose設定ファイル
├── .env.example                    # 環境変数のサンプル
├── README.md                       # プロジェクト説明
├── CLAUDE.md                       # このファイル
│
├── backend/                        # Go APIサーバー
│   ├── Dockerfile                  # バックエンド用Dockerfile
│   ├── go.mod                      # Goモジュール定義
│   ├── go.sum                      # 依存関係のチェックサム
│   ├── .env.example                # バックエンド環境変数サンプル
│   ├── cmd/
│   │   └── server/
│   │       └── main.go            # アプリケーションエントリーポイント
│   ├── internal/
│   │   ├── domain/                # エンティティとビジネスルール
│   │   │   ├── entity/            # ドメインエンティティ
│   │   │   │   ├── product.go
│   │   │   │   ├── cart.go
│   │   │   │   └── order.go
│   │   │   └── repository/        # リポジトリインターフェース
│   │   │       ├── product_repository.go
│   │   │       ├── cart_repository.go
│   │   │       └── order_repository.go
│   │   ├── usecase/               # アプリケーションビジネスロジック
│   │   │   ├── product_usecase.go
│   │   │   ├── cart_usecase.go
│   │   │   └── order_usecase.go
│   │   ├── infrastructure/        # 外部サービスとの統合
│   │   │   ├── persistence/       # データ永続化の実装
│   │   │   │   └── memory/        # インメモリ実装
│   │   │   │       ├── product_repository.go
│   │   │   │       ├── cart_repository.go
│   │   │   │       └── order_repository.go
│   │   │   └── monitoring/        # New Relic APM統合
│   │   │       └── newrelic.go
│   │   └── interface/             # インターフェースアダプター
│   │       └── api/               # HTTPハンドラーとルーティング
│   │           ├── router.go      # ルーティング設定
│   │           ├── handler/       # APIハンドラー実装
│   │           │   ├── health_handler.go
│   │           │   ├── product_handler.go
│   │           │   ├── cart_handler.go
│   │           │   └── order_handler.go
│   │           ├── middleware/    # HTTPミドルウェア
│   │           │   ├── cors.go
│   │           │   ├── monitoring.go
│   │           │   └── performance.go
│   │           └── presenter/     # レスポンスフォーマッター
│   │               └── response.go
│   └── pkg/                       # 共有ユーティリティ
│       ├── config/                # 設定管理
│       │   └── config.go
│       └── utils/                 # ユーティリティ関数
│           └── random.go
│
├── frontend/                       # Next.js フロントエンド
│   ├── Dockerfile                  # フロントエンド用Dockerfile
│   ├── package.json               # npm依存関係
│   ├── package-lock.json
│   ├── next.config.js             # Next.js設定
│   ├── tailwind.config.js         # Tailwind CSS設定
│   ├── tsconfig.json              # TypeScript設定
│   ├── .env.example               # フロントエンド環境変数サンプル
│   ├── app/                       # App Router
│   │   ├── layout.tsx             # ルートレイアウト（New Relic RUM初期化）
│   │   ├── page.tsx               # TOPページ
│   │   ├── products/
│   │   │   └── [id]/
│   │   │       └── page.tsx       # 商品詳細ページ
│   │   ├── cart/
│   │   │   └── page.tsx           # カートページ
│   │   └── checkout/
│   │       └── page.tsx           # 決済ページ
│   ├── components/                # Reactコンポーネント
│   │   ├── layout/
│   │   │   ├── Header.tsx         # ヘッダーコンポーネント
│   │   │   └── Footer.tsx         # フッターコンポーネント
│   │   ├── product/
│   │   │   ├── ProductCard.tsx    # 商品カードコンポーネント
│   │   │   └── ProductList.tsx    # 商品一覧コンポーネント
│   │   └── cart/
│   │       ├── CartItem.tsx       # カートアイテムコンポーネント
│   │       └── CartSummary.tsx    # カート合計コンポーネント
│   ├── lib/                       # ライブラリとユーティリティ
│   │   ├── api.ts                 # APIクライアント
│   │   ├── types.ts               # TypeScript型定義
│   │   └── monitoring.ts          # New Relic RUM設定
│   └── public/                    # 静的ファイル
│       └── images/                # 商品画像など
│
├── scripts/ ✅実装済み              # ユーザージャーニー負荷生成スクリプト
│   ├── Dockerfile                  # Go用マルチステージDockerfile
│   ├── go.mod                      # Goモジュール定義
│   ├── main.go                     # ユーザージャーニーシミュレーションメイン
│   └── test_local.sh               # ローカルテスト用スクリプト
│
└── docs/ [⚠️未実装]               # ドキュメント
    ├── setup.md                   # セットアップガイド
    ├── architecture.md            # アーキテクチャ説明
    └── handson-guide.md           # ハンズオンガイド
```

## アーキテクチャガイドライン

このプロジェクトはクリーンアーキテクチャの原則に従って実装します：

1. **依存性の方向**：
   - 外側のレイヤーは内側のレイヤーに依存
   - domain層は他のレイヤーに依存しない
   - インターフェースは内側のレイヤーで定義し、外側で実装

2. **New Relic APM統合（バックエンド）**：
   - 環境変数: `NEW_RELIC_API_KEY` - New Relic License Key
   - infrastructure/monitoring でAPMエージェントを管理
   - ミドルウェアとしてHTTPハンドラーに適用
   - ユースケース層でのカスタムセグメント追跡
   - トランザクション追跡とカスタムイベント送信
   
   **パフォーマンス調整用環境変数**：
   - `ERROR_RATE` - エラー応答率の調整（0.0〜1.0、デフォルト: 0.1）
   - `RESPONSE_TIME_MIN` - 最小レスポンス時間（ミリ秒、デフォルト: 50）
   - `RESPONSE_TIME_MAX` - 最大レスポンス時間（ミリ秒、デフォルト: 500）
   - `SLOW_ENDPOINT_RATE` - 遅延エンドポイントの発生率（0.0〜1.0、デフォルト: 0.2）

3. **Docker設定**：
   - Goアプリケーション用のマルチステージDockerfile
   - 環境変数サポート付きのdocker-compose.yml
   - コンテナ監視用のヘルスチェックエンドポイント

4. **API仕様**：

   **ヘルスチェック**：
   - `GET /health` - ヘルスチェックエンドポイント

   **商品関連API（2つ）**：
   - `GET /api/products` - 商品一覧取得
   - `GET /api/products/{id}` - 商品詳細取得

   **カート機能API（3つ）**：
   - `GET /api/cart` - カート内容取得
   - `POST /api/cart/items` - 商品をカートに追加
   - `PUT /api/cart/items/{id}` - カート内商品の数量変更・削除（オプション）

   **注文・決済API（2つ）**：
   - `GET /api/orders` - 全注文一覧取得（管理者用・ハンズオン確認用）
   - `POST /api/orders` - 注文作成（決済処理も含む）

   **デモ用エンドポイント**：
   - `GET /api/v1/error` - ランダムにエラーを返すエンドポイント（エラー率調整用）

   **API仕様書エンドポイント**：
   - `GET /api/docs` - Swagger UI（ブラウザでAPIドキュメント閲覧）
   - `GET /api/docs/swagger.yaml` - OpenAPI 3.0.3仕様書（YAML形式）

5. **API仕様書（Swagger/OpenAPI）管理ポリシー**：

   **仕様書の場所と形式**：
   - OpenAPI 3.0.3準拠のYAML形式
   - バックエンドコード内に埋め込み（`internal/interface/api/handler/swagger_handler.go`）
   - プロジェクトルートの`swagger.yaml`は参考用（実際の配信はバックエンドから）

   **更新ルール**：
   - APIエンドポイントの追加・変更・削除時は必ずSwagger仕様も同時更新
   - レスポンススキーマ変更時は対応するexamplesも更新
   - エラーコード追加時は該当するエラーレスポンス例も追加

   **必須記載事項**：
   - 全エンドポイントの詳細説明（SLI/SLO測定対象の明示）
   - リクエスト・レスポンスの完全なスキーマ定義
   - 成功・エラー両方の具体的なレスポンス例（examples）
   - SLMデモ用エンドポイントの動作説明
   - 環境変数による動作変更の説明

   **品質基準**：
   - 全エンドポイントに実際のAPIレスポンスに基づくサンプルデータを提供
   - 日本語でのわかりやすい説明（ハンズオン参加者向け）
   - New Relic SLMの文脈に沿った説明の追加

## フロントエンド実装仕様

### 技術スタック
- **フレームワーク**: Next.js 14+ (App Router)
- **スタイリング**: Tailwind CSS
- **パッケージマネージャー**: npm
- **New Relic統合**: Real User Monitoring (RUM) - Browser Agent

### ページ仕様

1. **TOPページ (/) [✅実装済み]**: 
   - 商品一覧をグリッド表示
   - 各商品カードに画像、名前、価格、詳細リンクを表示
   - ヘッダーにカートアイコンとアイテム数を表示
   - NRUG-SREロゴとブランディングを表示

2. **商品詳細ページ (/products/[id]) [✅実装済み]**:
   - 商品画像、詳細情報、価格を表示
   - 数量選択とカートに追加ボタン
   - 在庫情報の表示
   - カート追加時のフィードバック

3. **カートページ (/cart) [✅実装済み]**:
   - カート内商品の一覧表示
   - 各アイテムの数量変更・削除機能
   - 合計金額、送料、税金の自動計算
   - 空カート時のメッセージ表示
   - レジに進むボタン（決済ページへのリンク）

4. **決済ページ (/checkout) [✅実装済み]**:
   - 注文内容のサマリー表示
   - 配送先・決済方法の表示（ハンズオン用簡略版）
   - 注文確定ボタン（カート内容から注文作成）
   - 注文完了画面とサンクスページ

### 環境変数（必須）
```
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api
NEXT_PUBLIC_NEW_RELIC_BROWSER_KEY=your-new-relic-browser-license-key（必須）
NEXT_PUBLIC_NEW_RELIC_ACCOUNT_ID=your-new-relic-account-id（必須）
NEXT_PUBLIC_NEW_RELIC_APPLICATION_ID=your-new-relic-application-id（必須）
```

### New Relic Real User Monitoring (RUM) 設定
- **セキュアな設定**: 環境変数から動的にスクリプト生成、ハードコーディング排除
- **Browser Agent**: layout.tsxで`beforeInteractive`戦略による初期化
- **自動追跡**: ページビュー、AJAXリクエスト、JavaScriptエラー、パフォーマンス
- **カスタムイベント**: ユーザーアクションの追跡（商品閲覧、カート操作、注文完了）
- **Core Web Vitals**: LCP、FID、CLS、TTFB等のパフォーマンス指標測定
- **エラーレポーティング**: フロントエンドエラーの自動収集とレポート
- **セッション追跡**: ユーザージャーニー分析とコンバージョン測定

### Tailwind設定
- レスポンシブデザイン対応
- ダークモード対応（オプション）
- カスタムカラーパレット使用

## ハンズオンシナリオ

### 1. 環境セットアップ（20分）
- **New Relic APIキー払出し**
  - New Relicアカウントへのアクセス確認
  - APIキーの生成と環境変数への設定
  
- **デモアプリケーションの稼働**
  - Docker Composeでアプリケーション起動
  - 動作確認（商品閲覧、カート追加、決済フロー）
  - 参加者によるアプリケーション操作時間
  
- **APM / Real User Monitoring計測確認**
  - New Relic UIでのAPMデータ確認
  - Browser (RUM)データの確認
  - 現状の計測内容の理解

### 2. SLM設定ハンズオン（40分）
- **ユーザージャーニーの設定**
  - ユーザージャーニーの概念説明
  - 推奨事項：SLOは1〜3つから開始
  - ディスカッション：「このECサイトで最も重要な機能は？」
  - 実装：購入完了までのジャーニー設定
  
- **SLIの設定**
  - SLI設計方法の解説
    - 計算式：成功したリクエスト数 / 全リクエスト数
    - レイテンシ：95パーセンタイル < 300ms のリクエスト数 / 全リクエスト数
  - 実装：
    - 可用性SLI（エラー率ベース）
    - パフォーマンスSLI（レスポンスタイムベース）
  
- **SLOの設定**
  - SLOレベルの意味を解説
    - 99.99% = 月間4.3分のダウンタイム許容
    - 99.9% = 月間43分のダウンタイム許容
  - 機能開発スピードへの影響の説明
  - 実装：各SLIに対するSLO目標値設定

### 3. SLO管理ハンズオン（30分）
- **自動ユーザージャーニー実行**
  ```bash
  # 継続的なユーザージャーニー負荷生成開始
  docker-compose --profile load-test up load-generator
  
  # カスタム設定での実行（5分間、10秒間隔）
  DURATION=300 ACCESS_INTERVAL=10 docker-compose --profile load-test up load-generator
  ```
  
- **Service Level変化の体験**
  - 環境変数変更によるパフォーマンス劣化シミュレーション
  ```bash
  export ERROR_RATE=0.3
  export RESPONSE_TIME_MAX=2000
  docker-compose up -d api-server
  ```
  - New Relic UIでのSLO違反確認
  
- **エラーバジェット運用シミュレーション**
  - エラーバジェット消費状況の確認
  - ディスカッション：
    - バジェット枯渇時の対応（機能凍結？改善優先？）
    - リスクとビジネス価値のバランス
  - アラート設定の実習

### 実施のポイント
- 各セクションで参加者の理解度を確認
- 実際の運用を想定した議論を促進
- SLMの理論と実践の両方をカバー

## 重要な実装ノート

- これはNew Relic SLMを学習するための教育用ハンズオンプロジェクトです
- SLIメトリクス（レスポンスタイム、エラー率、可用性）のデモンストレーションに焦点を当てる
- SLO違反をトリガーできる現実的なシナリオを実装する
- テレメトリーデータがNew Relicに適切に送信されることを確認する
- **バックエンド**: New Relic APMでサーバーサイドパフォーマンスを監視
- **フロントエンド**: New Relic Real User Monitoring (RUM)でクライアントサイドパフォーマンスを監視
- エンドツーエンドのトランザクション追跡による完全な可視性を実現

## 現在の実装状況

### ✅ 実装完了
- **バックエンド**: Go APIサーバー（全エンドポイント）
- **フロントエンド**: TOPページ、商品詳細ページ、カートページ、決済完了ページ
- **Docker環境**: docker-compose.yaml設定とコンテナ化
- **API仕様書**: Swagger/OpenAPI 3.0.3対応
- **監視統合**: New Relic APM/RUM統合
- **ユーザージャーニー負荷生成**: Go実装、完全なECサイトフロー自動実行
- **SVG画像**: サンプル商品画像（6種類）
- **NRUG-SREブランディング**: ヘッダーロゴとテキスト

### ⚠️ 未実装（ハンズオン実施には不要）
- **ドキュメント** (`docs/`): セットアップガイド等

### 🔧 ハンズオン実施可能な機能
- **完全なECサイトフロー**: TOPページ → 商品詳細 → カート追加 → カート確認 → 決済完了
- **自動ユーザージャーニー**: プロファイル起動による継続的な負荷生成とSLI/SLOデータ生成
- **SLO違反シミュレーション**: エラー率調整による障害体験
- **New Relic監視**: APM/RUMデータ確認とダッシュボード
- **API仕様確認**: Swagger UIでのAPIドキュメント閲覧

## 🚨 重要: CI/CDチェック必須事項

**バックエンドコードを改修した後は、以下のCI/CDチェックを必ず実行し、全てがパスすることを確認してからタスクを完了すること:**

### 必須チェック項目
1. **コードフォーマット** (`gofmt`)
   ```bash
   # Dockerを使用した実行
   cd backend
   docker run --rm -v $(pwd):/app -w /app golang:1.21 bash -c "gofmt -l ."
   # 出力が空であることを確認（フォーマットエラーがない）
   ```

2. **静的解析** (`go vet`)
   ```bash
   docker run --rm -v $(pwd):/app -w /app golang:1.21 go vet ./...
   ```

3. **静的チェック** (`staticcheck`)
   ```bash
   docker run --rm -v $(pwd):/app -w /app golang:1.21 bash -c \
     "go install honnef.co/go/tools/cmd/staticcheck@2023.1.6 && staticcheck ./..."
   ```

4. **ユニットテスト**
   ```bash
   make test-unit
   # または
   docker run --rm -v $(pwd):/app -w /app golang:1.21 go test -v ./internal/...
   ```

5. **ビルド確認**
   ```bash
   docker run --rm -v $(pwd):/app -w /app golang:1.21 go build -v ./cmd/server
   ```

### 一括チェックコマンド
```bash
# backend ディレクトリで実行
docker run --rm -v $(pwd):/app -w /app golang:1.21 bash -c "
  echo '=== CI/CDチェック開始 ===' &&
  echo '1. フォーマットチェック...' &&
  if [ -n \"\$(gofmt -l .)\" ]; then 
    echo '❌ フォーマットエラー:' && gofmt -l . && exit 1
  else 
    echo '✅ フォーマット: OK'
  fi &&
  echo '2. go vet...' &&
  go vet ./... && echo '✅ go vet: OK' &&
  echo '3. staticcheck...' &&
  go install honnef.co/go/tools/cmd/staticcheck@2023.1.6 2>/dev/null &&
  staticcheck ./... && echo '✅ staticcheck: OK' &&
  echo '4. ビルド...' &&
  go build -v ./cmd/server > /dev/null 2>&1 && echo '✅ ビルド: OK' &&
  echo '5. テスト...' &&
  go test ./internal/... > /dev/null 2>&1 && echo '✅ テスト: OK' &&
  echo '=== 全チェック完了 ✅ ==='
"
```

### GitHub Actions ワークフロー
これらのチェックは `.github/workflows/test.yml` で自動実行されます。
プッシュ前にローカルで上記チェックを実行し、CI/CDパイプラインでの失敗を防ぐこと。

### チェック失敗時の対処
- **フォーマットエラー**: `gofmt -w .` で自動修正
- **go vet エラー**: 指摘された問題を手動で修正
- **staticcheck エラー**: 警告内容に従って修正（nil pointer、deprecated関数など）
- **テスト失敗**: 失敗したテストをデバッグして修正
- **ビルドエラー**: コンパイルエラーを修正
