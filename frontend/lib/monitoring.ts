'use client'

declare global {
  interface Window {
    newrelic?: {
      setPageViewName?: (name: string) => void
      addPageAction: (name: string, attributes?: Record<string, any>) => void
      setCustomAttribute: (name: string, value: string | number | boolean) => void
      noticeError: (error: Error, customAttributes?: Record<string, any>) => void
      setUserId?: (userId: string) => void
      finished: boolean
    }
  }
}

export class NewRelicMonitoring {
  private static instance: NewRelicMonitoring
  private isInitialized = false
  private sessionId: string

  private constructor() {
    // セッションIDを生成（ユーザー識別用）
    this.sessionId = this.generateSessionId()
  }

  public static getInstance(): NewRelicMonitoring {
    if (!NewRelicMonitoring.instance) {
      NewRelicMonitoring.instance = new NewRelicMonitoring()
    }
    return NewRelicMonitoring.instance
  }

  private generateSessionId(): string {
    return `session_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
  }

  public init(): void {
    if (this.isInitialized || typeof window === 'undefined') {
      return
    }

    // New Relic Browser Agentの初期化
    // 実際の初期化は layout.tsx で New Relic スクリプトタグによって行われる
    
    // セッションIDをNew Relicに設定
    if (window.newrelic && window.newrelic.setUserId) {
      window.newrelic.setUserId(this.sessionId)
    }
    
    // セッション情報をカスタム属性として設定
    this.setUserAttribute('sessionId', this.sessionId)
    this.setUserAttribute('clientTimestamp', Date.now())
    
    this.isInitialized = true
  }

  public trackPageView(pageName: string, attributes?: Record<string, any>): void {
    if (typeof window !== 'undefined' && window.newrelic) {
      window.newrelic.addPageAction('PageView', {
        page: pageName,
        ...attributes,
      })
    }
  }

  public trackUserAction(actionName: string, attributes?: Record<string, any>): void {
    if (typeof window !== 'undefined' && window.newrelic) {
      window.newrelic.addPageAction(actionName, attributes)
    }
  }

  public trackPurchase(orderId: string, amount: number, items: any[]): void {
    this.trackUserAction('Purchase', {
      orderId,
      amount,
      itemCount: items.length,
      items: items.map(item => ({
        productId: item.productId,
        quantity: item.quantity,
        price: item.product?.price
      }))
    })
  }

  public trackAddToCart(productId: string, quantity: number, price: number): void {
    this.trackUserAction('AddToCart', {
      productId,
      quantity,
      price,
      totalValue: quantity * price
    })
  }

  public trackRemoveFromCart(productId: string, quantity: number): void {
    this.trackUserAction('RemoveFromCart', {
      productId,
      quantity
    })
  }

  public trackError(error: Error, context?: string): void {
    if (typeof window !== 'undefined' && window.newrelic) {
      if (context) {
        window.newrelic.setCustomAttribute('errorContext', context)
      }
      window.newrelic.noticeError(error)
    }
  }

  public setUserAttribute(key: string, value: string | number | boolean): void {
    if (typeof window !== 'undefined' && window.newrelic) {
      window.newrelic.setCustomAttribute(key, value)
    }
  }

  public getSessionId(): string {
    return this.sessionId
  }
}

// シングルトンインスタンスをエクスポート
export const monitoring = NewRelicMonitoring.getInstance()

// React Hook for easy monitoring
export function useMonitoring() {
  return {
    trackPageView: monitoring.trackPageView.bind(monitoring),
    trackUserAction: monitoring.trackUserAction.bind(monitoring),
    trackPurchase: monitoring.trackPurchase.bind(monitoring),
    trackAddToCart: monitoring.trackAddToCart.bind(monitoring),
    trackRemoveFromCart: monitoring.trackRemoveFromCart.bind(monitoring),
    trackError: monitoring.trackError.bind(monitoring),
    setUserAttribute: monitoring.setUserAttribute.bind(monitoring),
  }
}

