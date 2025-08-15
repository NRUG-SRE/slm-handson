# マルチステージビルド: ビルドステージ
FROM golang:1.21-alpine AS builder

# 作業ディレクトリを設定
WORKDIR /app

# Go modulesをコピーしてダウンロード
COPY go.mod ./
RUN go mod download

# ソースコードをコピー
COPY . .

# バイナリをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o access-generator .

# 実行ステージ: 軽量な本番イメージ
FROM alpine:latest

# SSL証明書とタイムゾーン情報を追加
RUN apk --no-cache add ca-certificates tzdata curl

# 作業ディレクトリを設定
WORKDIR /root/

# ビルドされたバイナリをコピー
COPY --from=builder /app/access-generator .

# 環境変数のデフォルト値
ENV TARGET_URL=http://frontend:3000
ENV ACCESS_INTERVAL=10
ENV DURATION=300

# ヘルスチェック用のコマンド
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f ${TARGET_URL} || exit 1

# メインコマンド
CMD ["./access-generator"]