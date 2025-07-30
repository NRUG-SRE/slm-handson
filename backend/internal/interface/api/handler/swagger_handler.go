package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SwaggerHandler struct{}

func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

// ServeSwaggerYAML serves the swagger.yaml file
func (h *SwaggerHandler) ServeSwaggerYAML(c *gin.Context) {
	// swagger.yamlファイルの内容を直接返す
	swaggerContent := `openapi: 3.0.3
info:
  title: New Relic SLM Handson API
  description: |
    New Relic Service Level Management (SLM) ハンズオン用のECサイトAPIです。
    
    このAPIは以下の機能を提供します：
    - 商品一覧・詳細の取得
    - ショッピングカート操作
    - 注文処理
    - SLMデモ用のエラー生成エンドポイント
    
    ## レスポンス形式
    すべてのAPIレスポンスは以下の形式で返されます：
    '''json
    {
      "success": true,
      "data": { /* 実際のデータ */ }
    }
    '''
    
    エラー時は以下の形式：
    '''json
    {
      "success": false,
      "error": {
        "code": "ERROR_CODE",
        "message": "エラーメッセージ"
      }
    }
    '''
  version: 1.0.0
  contact:
    name: NRUG-SRE
    url: https://github.com/NRUG-SRE/slm-handson

servers:
  - url: http://localhost:8080
    description: Local development server

paths:
  /health:
    get:
      summary: ヘルスチェック
      description: APIサーバーの稼働状況を確認するヘルスチェックエンドポイント
      tags:
        - Health
      responses:
        '200':
          description: サーバーが正常に稼働中
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    example: "ok"
                  timestamp:
                    type: string
                    format: date-time
              examples:
                success:
                  summary: 正常レスポンス
                  value:
                    status: "ok"
                    timestamp: "2025-07-30T04:20:19Z"

  /api/products:
    get:
      summary: 商品一覧取得
      description: ECサイトの商品一覧を取得します。SLI測定の対象エンドポイントです。
      tags:
        - Products
      responses:
        '200':
          description: 商品一覧の取得に成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: array
                    items:
                      type: object
                      properties:
                        id:
                          type: string
                          format: uuid
                        name:
                          type: string
                        description:
                          type: string
                        price:
                          type: integer
                        imageUrl:
                          type: string
                        stock:
                          type: integer
                        createdAt:
                          type: string
                          format: date-time
                        updatedAt:
                          type: string
                          format: date-time
              examples:
                success:
                  summary: 商品一覧の取得成功
                  value:
                    success: true
                    data:
                      - id: "fa03a45e-8138-41b7-b2af-ed63881fb5f9"
                        name: "ワイヤレスヘッドホン"
                        description: "高音質なノイズキャンセリング機能付きワイヤレスヘッドホン"
                        price: 25000
                        imageUrl: "/images/headphones.jpg"
                        stock: 10
                        createdAt: "2025-07-30T04:20:19Z"
                        updatedAt: "2025-07-30T04:20:19Z"
                      - id: "5cfac614-a3d7-4230-a2a4-5b633834d1d2"
                        name: "スマートウォッチ"
                        description: "フィットネストラッキング機能付きの最新スマートウォッチ"
                        price: 35000
                        imageUrl: "/images/smartwatch.jpg"
                        stock: 5
                        createdAt: "2025-07-30T04:20:19Z"
                        updatedAt: "2025-07-30T04:20:19Z"
        '500':
          description: サーバー内部エラー（SLMデモ用のランダムエラーを含む）
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: false
                  error:
                    type: object
                    properties:
                      code:
                        type: string
                      message:
                        type: string
              examples:
                server_error:
                  summary: SLMデモ用のランダムエラー
                  value:
                    success: false
                    error:
                      code: "INTERNAL_SERVER_ERROR"
                      message: "Failed to get products"

  /api/products/{id}:
    get:
      summary: 商品詳細取得
      description: 指定されたIDの商品詳細情報を取得します
      tags:
        - Products
      parameters:
        - name: id
          in: path
          required: true
          description: 商品ID
          schema:
            type: string
            format: uuid
          example: "fa03a45e-8138-41b7-b2af-ed63881fb5f9"
      responses:
        '200':
          description: 商品詳細の取得に成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: object
                    properties:
                      id:
                        type: string
                        format: uuid
                      name:
                        type: string
                      description:
                        type: string
                      price:
                        type: integer
                      imageUrl:
                        type: string
                      stock:
                        type: integer
                      createdAt:
                        type: string
                        format: date-time
                      updatedAt:
                        type: string
                        format: date-time
              examples:
                success:
                  summary: 商品詳細の取得成功
                  value:
                    success: true
                    data:
                      id: "fa03a45e-8138-41b7-b2af-ed63881fb5f9"
                      name: "ワイヤレスヘッドホン"
                      description: "高音質なノイズキャンセリング機能付きワイヤレスヘッドホン"
                      price: 25000
                      imageUrl: "/images/headphones.jpg"
                      stock: 10
                      createdAt: "2025-07-30T04:20:19Z"
                      updatedAt: "2025-07-30T04:20:19Z"
        '404':
          description: 指定された商品が見つからない
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: false
                  error:
                    type: object
                    properties:
                      code:
                        type: string
                      message:
                        type: string
              examples:
                not_found:
                  summary: 商品が見つからない
                  value:
                    success: false
                    error:
                      code: "NOT_FOUND"
                      message: "Product not found"

  /api/cart:
    get:
      summary: カート内容取得
      description: 現在のカート内容を取得します。ユーザージャーニーの重要な部分です。
      tags:
        - Cart
      responses:
        '200':
          description: カート内容の取得に成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: object
                    properties:
                      id:
                        type: string
                      items:
                        type: array
                        items:
                          type: object
                          properties:
                            id:
                              type: string
                              format: uuid
                            productId:
                              type: string
                              format: uuid
                            product:
                              type: object
                            quantity:
                              type: integer
                            subtotal:
                              type: integer
                            addedAt:
                              type: string
                              format: date-time
                      totalAmount:
                        type: integer
                      itemCount:
                        type: integer
                      updatedAt:
                        type: string
                        format: date-time
              examples:
                with_items:
                  summary: アイテムが入ったカート
                  value:
                    success: true
                    data:
                      id: "default"
                      items:
                        - id: "item-1"
                          productId: "fa03a45e-8138-41b7-b2af-ed63881fb5f9"
                          product:
                            id: "fa03a45e-8138-41b7-b2af-ed63881fb5f9"
                            name: "ワイヤレスヘッドホン"
                            price: 25000
                          quantity: 2
                          subtotal: 50000
                          addedAt: "2025-07-30T04:25:00Z"
                      totalAmount: 50000
                      itemCount: 2
                      updatedAt: "2025-07-30T04:25:00Z"
                empty_cart:
                  summary: 空のカート
                  value:
                    success: true
                    data:
                      id: "default"
                      items: []
                      totalAmount: 0
                      itemCount: 0
                      updatedAt: "2025-07-30T04:20:19Z"

  /api/cart/items:
    post:
      summary: 商品をカートに追加
      description: 指定された商品をカートに追加します
      tags:
        - Cart
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - productId
                - quantity
              properties:
                productId:
                  type: string
                  format: uuid
                  description: 商品ID
                quantity:
                  type: integer
                  minimum: 1
                  description: 追加する数量
            examples:
              add_headphones:
                summary: ヘッドホンを2個追加
                value:
                  productId: "fa03a45e-8138-41b7-b2af-ed63881fb5f9"
                  quantity: 2
      responses:
        '200':
          description: カートへの追加に成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: object
              examples:
                success:
                  summary: カート追加成功
                  value:
                    success: true
                    data:
                      id: "default"
                      items:
                        - id: "item-1"
                          productId: "fa03a45e-8138-41b7-b2af-ed63881fb5f9"
                          quantity: 2
                          subtotal: 50000
                      totalAmount: 50000
                      itemCount: 2
        '404':
          description: 指定された商品が見つからない
          content:
            application/json:
              examples:
                not_found:
                  summary: 商品が見つからない
                  value:
                    success: false
                    error:
                      code: "NOT_FOUND"
                      message: "Product not found"
        '422':
          description: 在庫不足
          content:
            application/json:
              examples:
                insufficient_stock:
                  summary: 在庫不足エラー
                  value:
                    success: false
                    error:
                      code: "UNPROCESSABLE_ENTITY"
                      message: "Insufficient stock"

  /api/orders:
    post:
      summary: 注文作成
      description: カート内容を基に注文を作成します。最も重要なビジネスKPIを測定するエンドポイントです。
      tags:
        - Orders
      responses:
        '201':
          description: 注文作成に成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: object
                    properties:
                      id:
                        type: string
                        format: uuid
                      items:
                        type: array
                      totalAmount:
                        type: integer
                      status:
                        type: string
                      createdAt:
                        type: string
                        format: date-time
              examples:
                success:
                  summary: 注文作成成功
                  value:
                    success: true
                    data:
                      id: "order-12345"
                      items:
                        - id: "orderitem-1"
                          productId: "fa03a45e-8138-41b7-b2af-ed63881fb5f9"
                          productName: "ワイヤレスヘッドホン"
                          price: 25000
                          quantity: 2
                          subtotal: 50000
                      totalAmount: 50000
                      status: "completed"
                      createdAt: "2025-07-30T04:30:00Z"
        '422':
          description: カートが空、または在庫不足
          content:
            application/json:
              examples:
                empty_cart:
                  summary: 空のカート
                  value:
                    success: false
                    error:
                      code: "UNPROCESSABLE_ENTITY"
                      message: "Cart is empty"

  /api/v1/error:
    get:
      summary: SLMデモ用エラー生成
      description: |
        SLM（Service Level Management）デモンストレーション用のエンドポイントです。
        環境変数 ERROR_RATE に基づいてランダムにエラーを発生させます。
        
        - ERROR_RATE=0.1 → 10%の確率でエラー
        - ERROR_RATE=0.5 → 50%の確率でエラー
      tags:
        - Demo
      responses:
        '200':
          description: 正常レスポンス
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: object
                    properties:
                      message:
                        type: string
                      timestamp:
                        type: string
                        format: date-time
              examples:
                success:
                  summary: SLMデモ正常レスポンス
                  value:
                    success: true
                    data:
                      message: "This is a demo endpoint for SLM"
                      timestamp: "2025-07-30T04:35:00Z"
        '500':
          description: SLMデモ用の意図的なエラー
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: false
                  error:
                    type: object
                    properties:
                      code:
                        type: string
                      message:
                        type: string
              examples:
                demo_error:
                  summary: SLMデモエラー
                  value:
                    success: false
                    error:
                      code: "INTERNAL_SERVER_ERROR"
                      message: "Demo error for SLM testing"

tags:
  - name: Health
    description: ヘルスチェック関連
  - name: Products
    description: 商品管理
  - name: Cart
    description: ショッピングカート
  - name: Orders
    description: 注文管理
  - name: Demo
    description: SLMデモ用エンドポイント`
    
	c.Data(http.StatusOK, "application/x-yaml; charset=utf-8", []byte(swaggerContent))
}

// SwaggerUI serves a simple Swagger UI HTML page
func (h *SwaggerHandler) SwaggerUI(c *gin.Context) {
	html := `<!DOCTYPE html>
<html>
<head>
  <title>New Relic SLM Handson API Documentation</title>
  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui.css" />
  <style>
    html {
      box-sizing: border-box;
      overflow: -moz-scrollbars-vertical;
      overflow-y: scroll;
    }
    *, *:before, *:after {
      box-sizing: inherit;
    }
    body {
      margin:0;
      background: #fafafa;
    }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function() {
      const ui = SwaggerUIBundle({
        url: '/api/docs/swagger.yaml',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      });
    };
  </script>
</body>
</html>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}