# 超軽量版：JavaScript直接実行（TypeScript不使用）
FROM node:18-slim

# 必要最小限のパッケージをインストール
RUN apt-get update && apt-get install -y \
    chromium \
    fonts-liberation \
    dumb-init \
    curl \
    --no-install-recommends \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean

# Playwrightの環境変数
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright
ENV PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
ENV PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH=/usr/bin/chromium

# 非rootユーザーを作成
RUN groupadd -r appuser && useradd -r -g appuser -G audio,video appuser \
    && mkdir -p /home/appuser /app \
    && chown -R appuser:appuser /home/appuser /app

WORKDIR /app

# package.jsonをコピーして依存関係をインストール
COPY package.json ./package.json
RUN npm install --only=production && npm cache clean --force

# JavaScriptファイルを直接コピー
COPY src/index.js ./index.js

# 環境変数
ENV TARGET_URL=http://frontend:3000
ENV ACCESS_INTERVAL=10
ENV DURATION=300
ENV HEADLESS=true
ENV CONCURRENT_USERS=1

# 非rootユーザーで実行
USER appuser

# ヘルスチェック
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD curl -f ${TARGET_URL} || exit 1

# dumb-initで実行
ENTRYPOINT ["dumb-init", "--"]
CMD ["node", "index.js"]