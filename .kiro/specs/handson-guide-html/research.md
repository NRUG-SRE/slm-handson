# Research & Design Decisions

---
**Purpose**: Capture discovery findings, architectural investigations, and rationale that inform the technical design.
---

## Summary
- **Feature**: `handson-guide-html`
- **Discovery Scope**: 新規機能 (New Feature)
- **Key Findings**:
  - Tailwind CSS v4 CDN版を使用することで、ビルドプロセス不要のシンプルな実装が可能
  - 座学セクションのプレゼンテーション形式には、外部フレームワーク不要のカスタムCSS実装を採用
  - 単一HTMLファイル構成により、デプロイと配布が容易

## Research Log

### Tailwind CSS CDN選定
- **Context**: 要件でTailwind CSSの使用が指定されているが、ビルドプロセスを避けたい
- **Sources Consulted**:
  - https://tailwindcss.com/docs/installation/play-cdn
  - https://tailkits.com/blog/tailwind-css-v4-cdn-setup/
  - https://www.jsdelivr.com/package/npm/tailwindcss
- **Findings**:
  - Tailwind CSS v4では`@tailwindcss/browser@4` CDN版が利用可能
  - CDN URL: `https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4`
  - カスタム設定は`<style type="text/tailwindcss">`タグで可能
  - 開発用途に適しており、本ハンズオンの要件に合致
- **Implications**:
  - npm/Node.js不要でHTMLファイルのみで完結
  - @themeディレクティブでカラーパレットやフォントのカスタマイズ可能
  - 本番環境向けではないが、ハンズオンガイド用途では問題なし

### プレゼンテーションフレームワーク調査
- **Context**: 座学セクションをスライド形式で提供する必要がある
- **Sources Consulted**:
  - https://www.jqueryscript.net/blog/best-html-presentation-framework.html
  - https://sli.dev/guide/why (Slidev)
  - reveal.js, impress.js等の調査
- **Findings**:
  - 主要な選択肢:
    - reveal.js (37.6k stars): 最も成熟、機能豊富
    - Slidev (31.5k stars): Viteベース、開発者向け
    - impress.js (37.5k stars): 3Dトランジション
  - シンプルな要件の場合、カスタムCSS実装も有効
- **Implications**:
  - reveal.js等の外部フレームワークは機能過多でオーバーヘッド大
  - カスタムCSS + JavaScriptによる軽量実装を選択
  - キーボードナビゲーション（矢印キー）のみで十分

### HTMLガイド配置場所
- **Context**: プロジェクト内でのHTMLガイドの配置先を決定
- **Findings**:
  - 現在`docs/`ディレクトリは存在しない
  - README.mdがプロジェクトルートに配置されている
  - フロントエンド/バックエンドは別ディレクトリで分離
- **Implications**:
  - `docs/`ディレクトリを新規作成し、HTMLガイドを配置
  - GitHub Pagesでの公開も容易
  - 静的ファイル（画像、スタイル）も`docs/`配下に配置

## Architecture Pattern Evaluation

| Option | Description | Strengths | Risks / Limitations | Notes |
|--------|-------------|-----------|---------------------|-------|
| 単一HTMLファイル | 全コンテンツを1つのHTMLファイルに集約 | デプロイ簡単、依存なし、オフライン動作 | ファイルサイズ増加、メンテナンス性 | 要件1.1で許可されている |
| 複数HTMLページ | セクションごとにHTMLファイル分割 | メンテナンス性向上、読み込み高速化 | ナビゲーション複雑化、相対パス管理 | スケール時の選択肢 |
| reveal.js使用 | 既製プレゼンテーションフレームワーク | 機能豊富、プラグインエコシステム | 学習コスト、カスタマイズ制約 | オーバースペック |
| カスタムCSS実装 | Tailwind CSSベースのカスタムスライド | 軽量、完全制御、要件に最適化 | 初期実装コスト | 選択した方式 |

## Design Decisions

### Decision: 単一HTMLファイル + セクションベースナビゲーション
- **Context**: 要件1.1で「単一HTMLファイルまたは複数ページ」が許可されている
- **Alternatives Considered**:
  1. 単一HTMLファイル — 全セクションを1ファイルに集約、アンカーリンクでナビゲーション
  2. 複数HTMLページ — セクションごとにファイル分割、HTMLリンクで遷移
- **Selected Approach**: 単一HTMLファイル + セクションID + スムーススクロール
- **Rationale**:
  - デプロイが最もシンプル（1ファイルをホスティングするだけ）
  - オフライン動作が保証される
  - CDN依存を最小化（Tailwind CSSのみ）
  - 目次からのアンカーリンクで十分なナビゲーション体験
- **Trade-offs**:
  - Benefits: デプロイ容易、依存最小、オフライン対応
  - Compromises: ファイルサイズ増加（ただし画像は外部参照で軽減）
- **Follow-up**: 将来的にコンテンツが増加した場合は複数ページ化を検討

### Decision: Vite + Tailwind CSS v4ビルド方式の詳細
- **Vite Setup**:
  ```bash
  npm create vite@latest docs -- --template vanilla
  cd docs
  npm install -D tailwindcss @tailwindcss/vite
  ```
- **Vite設定** (`vite.config.js`):
  ```javascript
  import { defineConfig } from 'vite'
  import tailwindcss from '@tailwindcss/vite'

  export default defineConfig({
    plugins: [tailwindcss()],
    build: {
      outDir: 'dist',
      rollupOptions: {
        input: {
          main: './index.html'
        }
      }
    }
  })
  ```
- **CSSインポート** (`style.css`):
  ```css
  @import "tailwindcss";
  ```
- **ビルドコマンド**:
  - 開発: `npm run dev` (HMR対応)
  - 本番: `npm run build` → `dist/`ディレクトリに最適化されたHTML/CSS/JS出力
- **デプロイ**: GitHub Actionsで自動ビルド、GitHub Pagesで配信
  ```yaml
  # .github/workflows/deploy.yml
  name: Deploy to GitHub Pages
  on:
    push:
      branches: [ main ]
  jobs:
    build:
      runs-on: ubuntu-latest
      steps:
        - uses: actions/checkout@v4
        - uses: actions/setup-node@v4
          with:
            node-version: '20'
        - working-directory: ./docs
          run: npm ci && npm run build
        - uses: actions/upload-pages-artifact@v3
          with:
            path: ./docs/dist
    deploy:
      needs: build
      permissions:
        pages: write
        id-token: write
      uses: actions/deploy-pages@v4
  ```
- **重要**: Viteはビルド時のみ動作し、サーバーサイドレンダリング（SSR）は行わない。最終成果物は完全な静的ファイル。GitHub Pagesで無料ホスティング。

### Decision: カスタムCSS実装による座学スライド
- **Context**: 要件2.1で座学セクションをプレゼンテーション形式で提供
- **Alternatives Considered**:
  1. reveal.js — 成熟したフレームワーク、機能豊富
  2. カスタムCSS実装 — Tailwind CSSベース、軽量
- **Selected Approach**: Tailwind CSS + カスタムJavaScript（キーボードナビゲーション）
- **Rationale**:
  - reveal.jsは機能過多（3Dトランジション、プラグインシステム等不要）
  - 要件はシンプル（スライド表示 + 次へ/前へナビゲーション）
  - Tailwind CSSで十分なスタイリングが可能
  - JavaScriptは最小限（キーボードイベントハンドリングのみ）
- **Trade-offs**:
  - Benefits: 軽量、完全制御、依存最小
  - Compromises: 初期実装コスト（ただし要件がシンプルなため工数小）
- **Follow-up**: スライドアニメーション（フェード等）を検討可能

### Decision: シンタックスハイライトにPrism.jsを採用
- **Context**: 要件5.1でシンタックスハイライト付きコードブロックが必要
- **Alternatives Considered**:
  1. Prism.js — 軽量、多言語対応、テーマ豊富
  2. Highlight.js — 自動言語検出、やや重い
  3. 手書き`<pre><code>` — シンタックスハイライトなし
- **Selected Approach**: Prism.js CDN版（Bash, YAML対応）
- **Rationale**:
  - Bashコマンドとenv/YAML設定のハイライトが必要
  - Prism.jsは軽量で必要な言語のみロード可能
  - CDN版でビルド不要
- **Trade-offs**:
  - Benefits: 軽量、見やすい、コピーボタンプラグイン対応
  - Compromises: 外部CDN依存（ただしフォールバック不要）
- **Follow-up**: コピーボタンプラグイン（prism-copy-to-clipboard）を追加検討

## Risks & Mitigations

- **Risk 1: Tailwind CSS CDNのロード失敗** — Mitigation: CDNフォールバック不要（オフライン時は基本スタイルのみ表示）
- **Risk 2: 単一HTMLファイルのサイズ増加** — Mitigation: 画像は外部参照、インライン画像を避ける
- **Risk 3: ブラウザ互換性** — Mitigation: モダンブラウザのみをターゲット（要件8.1）、polyfill不要

## References
- [Tailwind CSS v4 CDN Setup](https://tailwindcss.com/docs/installation/play-cdn) — 公式ドキュメント、CDN使用方法
- [Prism.js Documentation](https://prismjs.com/) — シンタックスハイライトライブラリ
- [MDN: Smooth Scrolling](https://developer.mozilla.org/en-US/docs/Web/CSS/scroll-behavior) — CSSスムーススクロール実装
