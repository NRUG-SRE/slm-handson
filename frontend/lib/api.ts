import { Product, Cart, CartItem, Order } from './types'
import { monitoring } from './monitoring'

// ブラウザからはNext.jsのプロキシAPI経由でアクセス
const API_BASE_URL = typeof window !== 'undefined' ? '/api' : (process.env.INTERNAL_API_URL || 'http://api-server:8080/api')

class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message)
    this.name = 'ApiError'
  }
}

// バックエンドのレスポンス形式
interface BackendResponse<T> {
  success: boolean
  data?: T
  error?: {
    code: string
    message: string
  }
}

async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`
  
  // New Relic distributed tracing用のヘッダーを追加
  const distributedTracingHeaders: Record<string, string> = {}
  
  // New Relic Browser Agent v3では、W3C Trace Contextが自動的に処理される
  // カスタムヘッダーとして情報を送信
  if (typeof window !== 'undefined' && (window as any).newrelic) {
    try {
      // New Relicの現在のトレース情報を取得（利用可能な場合）
      if ((window as any).newrelic.getTrace) {
        const traceInfo = (window as any).newrelic.getTrace()
        if (traceInfo && traceInfo.traceId) {
          distributedTracingHeaders['X-NewRelic-Trace'] = traceInfo.traceId
        }
      }
      
      // セッション情報も含める
      distributedTracingHeaders['X-NewRelic-Browser'] = 'true'
    } catch (error) {
      console.debug('New Relic tracing header creation failed:', error)
    }
  }
  
  // セッションIDをヘッダーに追加（RUMとAPMの連携用）
  if (typeof window !== 'undefined') {
    const sessionId = monitoring.getSessionId()
    if (sessionId) {
      distributedTracingHeaders['X-Session-ID'] = sessionId
    }
  }
  
  const config: RequestInit = {
    headers: {
      'Content-Type': 'application/json',
      ...distributedTracingHeaders,
      ...options.headers,
    },
    ...options,
  }

  try {
    const response = await fetch(url, config)

    if (!response.ok) {
      throw new ApiError(response.status, `HTTP Error: ${response.status}`)
    }

    const backendResponse: BackendResponse<T> = await response.json()
    
    // バックエンドのレスポンス形式に対応
    if (!backendResponse.success) {
      const errorMessage = backendResponse.error?.message || 'Unknown server error'
      throw new ApiError(response.status, errorMessage)
    }

    if (backendResponse.data === undefined) {
      throw new Error('No data received from server')
    }

    return backendResponse.data
  } catch (error) {
    if (error instanceof ApiError) {
      throw error
    }
    throw new Error('Network error occurred')
  }
}

export const productApi = {
  // 商品一覧取得
  getProducts: async (): Promise<Product[]> => {
    return apiRequest<Product[]>('/products')
  },

  // 商品詳細取得
  getProduct: async (id: string): Promise<Product> => {
    return apiRequest<Product>(`/products/${id}`)
  },
}

export const cartApi = {
  // カート内容取得
  getCart: async (): Promise<Cart> => {
    return apiRequest<Cart>('/cart')
  },

  // 商品をカートに追加
  addToCart: async (productId: string, quantity: number): Promise<Cart> => {
    return apiRequest<Cart>('/cart/items', {
      method: 'POST',
      body: JSON.stringify({ productId, quantity }),
    })
  },

  // カート内商品の数量変更
  updateCartItem: async (itemId: string, quantity: number): Promise<Cart> => {
    return apiRequest<Cart>(`/cart/items/${itemId}`, {
      method: 'PUT',
      body: JSON.stringify({ quantity }),
    })
  },
}

export const orderApi = {
  // 注文作成（カート内容から自動で注文作成）
  createOrder: async (): Promise<Order> => {
    return apiRequest<Order>('/orders', {
      method: 'POST',
    })
  },

  // 注文詳細取得
  getOrder: async (id: string): Promise<Order> => {
    return apiRequest<Order>(`/orders/${id}`)
  },

  // 全注文一覧取得（管理者用・ハンズオン確認用）
  getOrders: async (): Promise<Order[]> => {
    return apiRequest<Order[]>('/orders')
  },
}

// ヘルスチェックAPI
export const healthApi = {
  check: async (): Promise<{ status: string }> => {
    return apiRequest<{ status: string }>('/health')
  },
}

// エラーAPIのテストエンドポイント
export const demoApi = {
  triggerError: async (): Promise<any> => {
    return apiRequest<any>('/v1/error')
  },
}

export { ApiError }