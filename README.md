# New Relic Service Level Management ハンズオン

このプロジェクトは、New RelicのService Level Management（SLM）を活用してService Level Objective（SLO）やService Level Indicator（SLI）を管理するハンズオンを提供します。

## 概要

ECサイトをモデルとしたサンプルアプリケーション（Go APIサーバー + Next.jsフロントエンド）を使用して、New Relic APMとReal User Monitoring (RUM)によるエンドツーエンドのモニタリングとSLM機能の実践的な学習ができます。Docker Composeを使用して簡単に環境を構築し、実際のテレメトリーデータをNew Relicに送信しながらSLO/SLIの設定と管理を体験できます。

## 前提条件

- Docker および Docker Compose
- New Relicアカウント
- New Relic API Key（APM用）
- New Relic Browser API Key（RUM用）
- Go 1.21以上（ローカル開発時）
- Node.js 18以上（ローカル開発時）

## セットアップ

### 1. リポジトリのクローン

```bash
git clone https://github.com/NRUG-SRE/slm-handson.git
cd slm-handson
```

### 2. 環境変数の設定

`.env.example`をコピーして`.env`ファイルを作成し、必要な値を設定します：

```bash
cp .env.example .env
```

`.env`ファイルの内容：
```
# Backend (APM)
NEW_RELIC_API_KEY=your-api-key-here
NEW_RELIC_APP_NAME=slm-handson-api
ERROR_RATE=0.1
RESPONSE_TIME_MIN=50
RESPONSE_TIME_MAX=500
SLOW_ENDPOINT_RATE=0.2

# Frontend (RUM)
NEXT_PUBLIC_NEW_RELIC_BROWSER_KEY=your-browser-key-here
NEXT_PUBLIC_NEW_RELIC_ACCOUNT_ID=your-account-id
NEXT_PUBLIC_NEW_RELIC_APPLICATION_ID=your-app-id
```

### 3. アプリケーションの起動

Docker Composeを使用してサンプルアプリケーションを起動します：

```bash
docker-compose up -d
```

起動後、以下のURLでアクセスできます：
- フロントエンド: http://localhost:3000
- バックエンドAPI: http://localhost:8080/api
- ヘルスチェック: http://localhost:8080/health

## アーキテクチャ

### システム構成

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Browser   │────▶│  Frontend   │────▶│ API Server  │
│             │     │  (Next.js)  │     │    (Go)     │
└─────────────┘     └─────────────┘     └─────────────┘
      │                    │                    │
      │                    │                    │
      ▼                    ▼                    ▼
┌─────────────────────────────────────────────────────┐
│              New Relic Platform                     │
│  ┌─────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │   RUM   │  │     APM      │  │     SLM      │  │
│  └─────────┘  └──────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────┘
```

### コンポーネント

- **フロントエンド (Next.js)**
  - ECサイトのUI実装
  - Real User Monitoring (RUM) でクライアントサイドパフォーマンスを監視
  
- **APIサーバー (Go)**
  - クリーンアーキテクチャに基づく実装
  - 商品管理、カート、注文処理のREST API
  - New Relic APMでサーバーサイドパフォーマンスを監視

- **負荷生成スクリプト**
  - 実際のユーザー行動をシミュレート
  - SLO違反シナリオの再現

## ハンズオンシナリオ

### 1. 環境セットアップ（20分）
- New Relic APIキーの払い出しと設定
- デモアプリケーションの起動と動作確認
- APM / Real User Monitoringの計測確認

### 2. SLM設定ハンズオン（40分）
- **ユーザージャーニーの設定**
  - ECサイトで最も重要な機能の特定
  - 購入完了までのジャーニー設定
  
- **SLIの設定**
  - 可用性SLI（成功率ベース）
  - パフォーマンスSLI（レスポンスタイムベース）
  
- **SLOの設定**
  - 99.99% vs 99.9%の違いと影響
  - 適切な目標値の設定

### 3. SLO管理ハンズオン（30分）
- 擬似ユーザーアクセスによる負荷生成
- 環境変数によるService Level変化の体験
- エラーバジェット枯渇時の対応シミュレーション

## 主要なAPIエンドポイント

### 商品関連
- `GET /api/products` - 商品一覧
- `GET /api/products/{id}` - 商品詳細

### カート機能
- `GET /api/cart` - カート内容取得
- `POST /api/cart/items` - 商品をカートに追加
- `PUT /api/cart/items/{id}` - カート内商品の数量変更・削除

### 注文・決済
- `POST /api/orders` - 注文作成（決済処理含む）
- `GET /api/orders/{id}` - 注文確認

## 負荷生成とテスト

擬似ユーザーアクセスを生成してSLOの動作を確認：

```bash
# 通常の負荷テスト
docker-compose run load-generator

# パフォーマンス劣化シミュレーション
export ERROR_RATE=0.3
export RESPONSE_TIME_MAX=2000
docker-compose up -d api-server
```

## 開発コマンド

```bash
# ログの確認
docker-compose logs -f
docker-compose logs -f api-server
docker-compose logs -f frontend

# アプリケーションの停止
docker-compose down

# ビルドして起動
docker-compose up --build -d

# 個別サービスの再起動
docker-compose restart api-server
docker-compose restart frontend
```

## プロジェクト構成

```
slm-handson/
├── backend/          # Go APIサーバー
├── frontend/         # Next.jsフロントエンド  
├── scripts/          # 負荷生成スクリプト
├── docs/             # ドキュメント
├── docker-compose.yml
└── .env.example
```

## トラブルシューティング

- **ポート競合**: 3000番（フロントエンド）、8080番（API）が使用されていないか確認
- **New Relicにデータが表示されない**: 
  - API Keyが正しく設定されているか確認
  - アプリケーション名が正しいか確認
  - 数分待ってからリフレッシュ
- **Docker関連のエラー**: `docker-compose logs` でエラーメッセージを確認

## 参考資料

- [New Relic Service Level Management](https://docs.newrelic.com/docs/service-level-management/)
- [SLI/SLO設計のベストプラクティス](https://sre.google/sre-book/service-level-objectives/)
- [Error Budgetの運用方法](https://sre.google/workbook/error-budget-policy/)

## 貢献

イシューやプルリクエストは歓迎します。大きな変更を行う場合は、まずイシューを作成して変更内容について議論してください。

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。