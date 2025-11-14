# Requirements Document

## Project Description (Input)
ハンズオンガイドをHTMLで作成したい。

◆概要
このハンズオンでは、リアルなデモ用ECサイトを題材に、
Service Levelを設計・管理・運用するシミュレーションの体験を提供します。
New Relicだけでなく、Service Levelを学びたい方を意識した構成にしており、
体験を通じて、以下を学ぶことができます。

座学
    Service Levels を管理する必要性
    品質に対する考え方

実技
    Service Level Indicator（SLI）の選定
    Service Level Objective（SLO）の設計
    New Relic Service Level Managementを活用したSLO管理
    エラーバジェットの運用
    APM/RUMを活用したオブザーバビリティの活用
◆聴講対象者
継続的にオンラインサービス／デジタルプロダクトを開発・改善・運用している方
    プロダクトオーナー
    プロジェクトマネージャー
    SRE
    開発者
◆カリキュラム
    Service Levelの座学（20min）
        AI時代にこそやりたいService Level Managementってどんなもの？
        可用性 vs 機能のエクササイズ
        SLOが高すぎたり、低すぎたりするとどうなるか
        SLA / SLI / SLOとは何か
        SLAとSLOは別物
        SLIとSLOの作り方
デモ用ECサイトを動かそう（30min）
    GitHub CodeSpacesを利用したアプリケーションの立ち上げ
    New Relic API Keyの払出しと適用
    デモ用ECサイトをユーザー目線でいっぱい触ろう（動作理解）
    New RelicのAPMを見てみよう
    New RelicのRUMを見てみよう
休憩（10min）
SLOを作ろう(30min)
    このアプリケーションのユーザージャーニーを考えよう
    SLIを設計しよう
    New Relic Service Level ManagementでSLOを管理しよう
    エラー率・レスポンスタイムを調整してSLO違反を体験して意思決定をしてみよう

## Introduction
本ドキュメントは、New Relic Service Level ManagementのハンズオンガイドをHTML形式で提供するための要件を定義します。このガイドは、プロダクトオーナー、プロジェクトマネージャー、SRE、開発者などがService Level（SLI/SLO/SLA）の概念を学び、実際のデモECサイトを使用してSLO設計・管理を体験できるようにすることを目的としています。

## Requirements

### Requirement 1: HTMLガイドの基本構造
**Objective:** ハンズオン参加者として、わかりやすく構造化されたHTMLガイドを閲覧したい。そうすることで、順序立ててカリキュラムを進めることができる。

#### Acceptance Criteria
1. HTMLガイドシステムは、単一HTMLファイルまたは複数ページで構成されること
2. HTMLガイドシステムは、PCブラウザで適切に表示されること
3. HTMLガイドシステムは、目次（Table of Contents）を提供し、各セクションへのリンクを含むこと
4. HTMLガイドシステムは、視覚的に読みやすいタイポグラフィとレイアウトを使用すること
5. HTMLガイドシステムは、コードブロック、画像、リンクを適切に表示できること

### Requirement 2: 座学セクション - Service Level Management基礎
**Objective:** ハンズオン参加者として、Service Level Managementの理論的背景を理解したい。そうすることで、実技演習の意義を理解できる。

#### Acceptance Criteria
1. HTMLガイドは、座学セクションをプレゼンテーション形式（スライドライクなレイアウト）で提供すること
2. When 参加者が座学セクションにアクセスするとき、HTMLガイドは「AI時代におけるService Level Management」の説明を表示すること
3. HTMLガイドは、従来の組織内成果物（月次・週次報告、改善提案書、障害報告書、再発防止策、品質向上レポート等）が最も重要なステークホルダー（ユーザー）の視点を欠いている問題を指摘すること
4. HTMLガイドは、ユーザー視点とプロダクト視点の違いを明確に説明すること（ユーザーが求める「高い可用性による良質なユーザー体験」vs プロダクトが求める「ユーザー数、利用量、単価」）
5. HTMLガイドは、Service Level Managementにおいてユーザー視点を最優先すべき理由を解説すること
6. HTMLガイドは、「可用性 vs 機能」のトレードオフに関する解説を含むこと
7. HTMLガイドは、可用性vs機能のトレードオフ解説において、以下の対比例を提示すること：
   - Aのアプリ：タスクの登録と削除の機能しかないが、絶対に期待した通りの操作ができる
   - Bのアプリ：NotionやBacklogのように魅力的な機能を兼ね揃えているが、3回に1回くらいの頻度で操作が失敗する
8. HTMLガイドは、上記の例を通じて信頼性を基準としたバランスをとることの重要性を訴求すること
9. HTMLガイドは、SLOレベル（99.9%, 99.99%, 99.999%等）の具体的な意味とダウンタイム許容時間の比較表を提供すること
10. HTMLガイドは、SLOに9を一つ追加するたびに必要となる労力とコストの増加について解説すること
11. HTMLガイドは、SLA、SLI、SLOの定義と違いを明確に説明すること
12. HTMLガイドは、SLIとSLOの設計方法に関するガイドラインを提供すること
13. HTMLガイドは、スライド間を移動するためのナビゲーション機能（次へ/前へ）を提供すること
14. HTMLガイドは、座学の推奨時間（20分）を明示すること

### Requirement 3: 環境セットアップセクション
**Objective:** ハンズオン参加者として、デモECサイトを迅速にセットアップしたい。そうすることで、実技演習をスムーズに開始できる。

#### Acceptance Criteria
1. When 参加者が環境セットアップセクションにアクセスするとき、HTMLガイドはGitHub Codespacesの起動手順を段階的に説明すること
2. HTMLガイドは、New Relic License Keyの取得方法と`.env`ファイルへの設定手順を含むこと
3. HTMLガイドは、`docker compose up -d --build`コマンドを含む起動コマンドを明示すること
4. HTMLガイドは、フロントエンド（ポート3000）とAPI（ポート8080）へのアクセス方法を説明すること
5. HTMLガイドは、デモECサイトの動作確認方法（商品閲覧、カート追加、決済フロー）を案内すること
6. HTMLガイドは、New Relic APMおよびRUMでのデータ確認手順を含むこと
7. HTMLガイドは、環境セットアップの推奨時間（30分）を明示すること

### Requirement 4: SLO設定・管理セクション
**Objective:** ハンズオン参加者として、実際のアプリケーションでSLOを設計・設定・管理する実技を体験したい。そうすることで、実務でのSLM運用スキルを習得できる。

#### Acceptance Criteria
1. When 参加者がSLO設定セクションにアクセスするとき、HTMLガイドはユーザージャーニーの特定方法を説明すること
2. HTMLガイドは、SLI設計のワークフロー（成功率ベース、レスポンスタイムベース）を段階的に解説すること
3. HTMLガイドは、New Relic Service Level ManagementでのSLO作成手順を含むこと
4. HTMLガイドは、環境変数（ERROR_RATE、RESPONSE_TIME_MAX）を変更してパフォーマンス劣化をシミュレートする手順を提供すること
5. HTMLガイドは、Playwrightベースの自動ユーザージャーニー負荷生成の起動コマンド（`docker compose --profile playwright up playwright-generator`）を含むこと
6. HTMLガイドは、SLO違反時のエラーバジェット運用とアラート設定の解説を含むこと
7. HTMLガイドは、SLO設定・管理の推奨時間（30分）を明示すること

### Requirement 5: コードスニペットとコマンド表示
**Objective:** ハンズオン参加者として、実行すべきコマンドを明確に把握したい。そうすることで、操作ミスを減らし効率的に演習を進められる。

#### Acceptance Criteria
1. When HTMLガイドがターミナルコマンドを表示するとき、HTMLガイドはシンタックスハイライト付きのコードブロックを使用すること
2. HTMLガイドは、各コマンドの目的を簡潔に説明するコメントまたは説明文を含むこと
3. Where コピー機能が実装される場合、HTMLガイドは「コピー」ボタンを各コードブロックに提供すること
4. HTMLガイドは、`.env`ファイルの設定例を具体的に示すこと
5. HTMLガイドは、エラー発生時のトラブルシューティング用コマンド（`docker compose logs`等）を含むこと

### Requirement 6: 視覚的な補助資料
**Objective:** ハンズオン参加者として、図表やスクリーンショットを参照したい。そうすることで、操作手順を視覚的に理解できる。

#### Acceptance Criteria
1. Where 図表が必要な場合、HTMLガイドはアーキテクチャ図（フロントエンド、バックエンド、New Relic連携）を含むこと
2. Where スクリーンショットが提供される場合、HTMLガイドはNew Relic UIでのSLO設定画面のスクリーンショットを含むこと
3. HTMLガイドは、ユーザージャーニーフロー図（TOPページ→商品詳細→カート→決済）を含むこと
4. HTMLガイドは、SLO違反時のNew Relicダッシュボード例を視覚的に示すこと
5. HTMLガイドは、全ての画像に代替テキスト（alt属性）を設定すること

### Requirement 7: ナビゲーションと進捗管理
**Objective:** ハンズオン参加者として、自分の進捗状況を把握したい。そうすることで、時間配分を調整できる。

#### Acceptance Criteria
1. HTMLガイドは、各セクションの推奨所要時間を明示すること
2. HTMLガイドは、「前へ」「次へ」ナビゲーションボタンを各セクションの末尾に提供すること
3. Where 複数ページ構成の場合、HTMLガイドは現在位置を示すブレッドクラムまたは進捗インジケーターを表示すること
4. HTMLガイドは、休憩時間（10分）のセクションを明示すること

### Requirement 8: 技術スタックとブラウザ互換性
**Objective:** ハンズオン参加者として、様々な環境でガイドを閲覧したい。そうすることで、使用するデバイスに制約を受けずに参加できる。

#### Acceptance Criteria
1. HTMLガイドは、Chrome最新版で正常に表示されること
2. HTMLガイドは、CSSフレームワークとしてTailwind CSSを採用すること
3. If JavaScriptが使用される場合、HTMLガイドは基本機能（コンテンツ表示）がJavaScript無効時にも動作すること
