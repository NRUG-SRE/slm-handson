#!/bin/bash

# SLM ハンズオン用アクセス生成スクリプトのローカルテスト

set -e

echo "🚀 SLM ハンズオン ユーザーアクセス生成スクリプト ローカルテスト"
echo "=================================================================="

# 設定
export TARGET_URL=${TARGET_URL:-"http://localhost:3000"}
export ACCESS_INTERVAL=${ACCESS_INTERVAL:-"5"}
export DURATION=${DURATION:-"30"}

echo "設定:"
echo "  ターゲットURL: $TARGET_URL"
echo "  アクセス間隔: ${ACCESS_INTERVAL}秒"
echo "  実行時間: ${DURATION}秒"
echo ""

# Dockerを使用してGoプログラムを実行
echo "📡 ユーザーアクセス生成開始..."
echo "   (Ctrl+Cで中断可能)"
echo ""

# Dockerを使ってビルドして実行
docker run --rm -v $(pwd):/app -w /app \
    -e TARGET_URL="$TARGET_URL" \
    -e ACCESS_INTERVAL="$ACCESS_INTERVAL" \
    -e DURATION="$DURATION" \
    golang:1.21 go run main.go