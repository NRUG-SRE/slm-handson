// New Relic RUM ユーティリティ関数

declare global {
  interface Window {
    newrelic?: {
      setPageViewName?: (name: string) => void;
      addPageAction: (name: string, attributes?: Record<string, any>) => void;
      setCustomAttribute: (name: string, value: string | number | boolean) => void;
      noticeError: (error: Error, customAttributes?: Record<string, any>) => void;
      setUserId?: (userId: string) => void;
      finished: boolean;
    };
  }
}

/**
 * New Relicがロードされているかチェック
 */
export const isNewRelicAvailable = (): boolean => {
  return typeof window !== 'undefined' && !!window.newrelic;
};

/**
 * ページビューを手動で記録
 */
export const trackPageView = (pageName: string): void => {
  if (isNewRelicAvailable() && window.newrelic!.setPageViewName) {
    window.newrelic!.setPageViewName(pageName);
  }
};

/**
 * カスタムページアクションを記録
 */
export const trackPageAction = (actionName: string, attributes?: Record<string, any>): void => {
  if (isNewRelicAvailable()) {
    window.newrelic!.addPageAction(actionName, {
      timestamp: Date.now(),
      url: window.location.href,
      ...attributes
    });
  }
};

/**
 * ECサイト固有のアクション追跡
 */
export const trackECommerceAction = {
  // 商品閲覧
  viewProduct: (productId: string, productName: string, price: number) => {
    trackPageAction('product_view', {
      productId,
      productName,
      price,
      category: 'ecommerce'
    });
  },

  // カートに追加
  addToCart: (productId: string, productName: string, quantity: number, price: number) => {
    trackPageAction('add_to_cart', {
      productId,
      productName,
      quantity,
      price,
      totalValue: quantity * price,
      category: 'ecommerce'
    });
  },

  // カートから削除
  removeFromCart: (productId: string, quantity: number) => {
    trackPageAction('remove_from_cart', {
      productId,
      quantity,
      category: 'ecommerce'
    });
  },

  // 決済開始
  beginCheckout: (cartTotal: number, itemCount: number) => {
    trackPageAction('begin_checkout', {
      cartTotal,
      itemCount,
      category: 'ecommerce'
    });
  },

  // 購入完了
  completePurchase: (orderId: string, total: number, itemCount: number) => {
    trackPageAction('purchase_complete', {
      orderId,
      total,
      itemCount,
      category: 'ecommerce',
      conversionEvent: true
    });
  }
};

/**
 * エラーを手動で記録
 */
export const trackError = (error: Error, context?: Record<string, any>): void => {
  if (isNewRelicAvailable()) {
    window.newrelic!.noticeError(error, {
      timestamp: Date.now(),
      url: window.location.href,
      ...context
    });
  }
};

/**
 * カスタム属性を設定
 */
export const setCustomAttribute = (name: string, value: string | number | boolean): void => {
  if (isNewRelicAvailable()) {
    window.newrelic!.setCustomAttribute(name, value);
  }
};

/**
 * ユーザーIDを設定
 */
export const setUserId = (userId: string): void => {
  if (isNewRelicAvailable() && window.newrelic!.setUserId) {
    window.newrelic!.setUserId(userId);
  }
};