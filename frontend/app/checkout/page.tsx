'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { cartApi, orderApi } from '@/lib/api'
import { Cart, Order } from '@/lib/types'

export default function CheckoutPage() {
  const router = useRouter()
  const [cart, setCart] = useState<Cart | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [isProcessing, setIsProcessing] = useState(false)
  const [order, setOrder] = useState<Order | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchCart()
  }, [])

  const fetchCart = async () => {
    try {
      setIsLoading(true)
      const cartData = await cartApi.getCart()
      setCart(cartData)
      
      // カートが空の場合はカートページにリダイレクト
      if (!cartData.items || cartData.items.length === 0) {
        router.push('/cart')
        return
      }
    } catch (error) {
      console.error('Failed to fetch cart:', error)
      setError('カート情報の取得に失敗しました')
    } finally {
      setIsLoading(false)
    }
  }

  const handlePlaceOrder = async () => {
    if (!cart || cart.items.length === 0) {
      setError('カートが空です')
      return
    }

    try {
      setIsProcessing(true)
      setError(null)
      
      // 注文を作成（バックエンドでカートも空になる）
      const createdOrder = await orderApi.createOrder()
      setOrder(createdOrder)
      
      // カート更新イベントを発火してヘッダーのカート数を更新
      window.dispatchEvent(new CustomEvent('cartUpdated'))
      
    } catch (error) {
      console.error('Failed to create order:', error)
      setError('注文の作成に失敗しました。もう一度お試しください。')
      setIsProcessing(false)
    }
  }

  if (isLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-center items-center min-h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      </div>
    )
  }

  // 注文完了画面
  if (order) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-2xl mx-auto">
          <div className="bg-green-50 border border-green-200 rounded-lg p-6 mb-6">
            <div className="flex items-center mb-4">
              <svg className="w-6 h-6 text-green-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
              <h1 className="text-2xl font-bold text-green-800">注文が完了しました！</h1>
            </div>
            <p className="text-green-700 mb-4">
              ご注文ありがとうございます。注文確認メールを送信いたします。
            </p>
            <div className="bg-white p-4 rounded border">
              <p className="text-sm text-gray-600 mb-2">注文番号</p>
              <p className="font-mono text-lg font-bold text-gray-900">{order.id}</p>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow-md p-6 mb-6">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">注文内容</h2>
            <div className="space-y-4">
              {order.items.map((item) => (
                <div key={item.id} className="flex justify-between items-center py-2 border-b border-gray-100 last:border-b-0">
                  <div>
                    <p className="font-medium text-gray-900">{item.product.name}</p>
                    <p className="text-sm text-gray-500">数量: {item.quantity}</p>
                  </div>
                  <p className="font-medium text-gray-900">¥{(item.price * item.quantity).toLocaleString()}</p>
                </div>
              ))}
              
              <div className="border-t pt-4">
                <div className="flex justify-between items-center">
                  <span className="text-lg font-bold text-gray-900">合計</span>
                  <span className="text-lg font-bold text-gray-900">¥{order.totalAmount.toLocaleString()}</span>
                </div>
              </div>
            </div>
          </div>

          <div className="flex space-x-4">
            <Link
              href="/"
              className="flex-1 bg-blue-600 text-white text-center py-3 px-4 rounded-md hover:bg-blue-700 transition-colors font-medium"
            >
              ショッピングを続ける
            </Link>
            <button
              onClick={() => window.print()}
              className="flex-1 bg-gray-100 text-gray-700 text-center py-3 px-4 rounded-md hover:bg-gray-200 transition-colors font-medium"
            >
              注文内容を印刷
            </button>
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-2xl mx-auto">
          <div className="bg-red-50 border border-red-200 rounded-lg p-6">
            <div className="flex items-center mb-2">
              <svg className="w-5 h-5 text-red-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <h2 className="text-lg font-semibold text-red-800">エラーが発生しました</h2>
            </div>
            <p className="text-red-700 mb-4">{error}</p>
            <div className="flex space-x-4">
              <Link
                href="/cart"
                className="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700 transition-colors"
              >
                カートに戻る
              </Link>
              <button
                onClick={fetchCart}
                className="bg-gray-100 text-gray-700 px-4 py-2 rounded-md hover:bg-gray-200 transition-colors"
              >
                再試行
              </button>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (!cart || cart.items.length === 0) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-2xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">カートが空です</h1>
          <p className="text-gray-600 mb-6">決済を行うには、まず商品をカートに追加してください。</p>
          <Link
            href="/"
            className="bg-blue-600 text-white px-6 py-3 rounded-md hover:bg-blue-700 transition-colors font-medium"
          >
            商品を見る
          </Link>
        </div>
      </div>
    )
  }

  const subtotal = cart.totalAmount
  const shipping = subtotal > 10000 ? 0 : 500
  const tax = Math.floor(subtotal * 0.1)
  const total = subtotal + shipping + tax

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-3xl font-bold text-gray-900 mb-8">決済</h1>
        
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* 注文内容 */}
          <div className="bg-white rounded-lg shadow-md p-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">注文内容</h2>
            
            <div className="space-y-4 mb-6">
              {cart.items.map((item) => (
                <div key={item.id} className="flex items-center space-x-4 py-3 border-b border-gray-100 last:border-b-0">
                  <div className="w-16 h-16 bg-gray-200 rounded-md overflow-hidden flex-shrink-0">
                    <img
                      src={item.product.imageUrl || '/images/placeholder.png'}
                      alt={item.product.name}
                      className="w-full h-full object-cover object-center"
                    />
                  </div>
                  <div className="flex-grow">
                    <h3 className="font-medium text-gray-900">{item.product.name}</h3>
                    <p className="text-sm text-gray-500">¥{item.product.price.toLocaleString()} × {item.quantity}</p>
                  </div>
                  <div className="text-right">
                    <p className="font-medium text-gray-900">¥{(item.product.price * item.quantity).toLocaleString()}</p>
                  </div>
                </div>
              ))}
            </div>

            {/* 金額内訳 */}
            <div className="space-y-2 pt-4 border-t border-gray-200">
              <div className="flex justify-between">
                <span className="text-gray-600">小計</span>
                <span className="font-medium">¥{subtotal.toLocaleString()}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">送料</span>
                <span className="font-medium">
                  {shipping === 0 ? (
                    <span className="text-green-600">無料</span>
                  ) : (
                    `¥${shipping.toLocaleString()}`
                  )}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">消費税(10%)</span>
                <span className="font-medium">¥{tax.toLocaleString()}</span>
              </div>
              <div className="border-t pt-2">
                <div className="flex justify-between">
                  <span className="text-lg font-bold text-gray-900">合計</span>
                  <span className="text-lg font-bold text-gray-900">¥{total.toLocaleString()}</span>
                </div>
              </div>
            </div>
          </div>

          {/* 決済情報（ハンズオン用簡略版） */}
          <div className="bg-white rounded-lg shadow-md p-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">決済情報</h2>
            
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
              <div className="flex items-start">
                <svg className="w-5 h-5 text-blue-600 mr-2 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <div>
                  <p className="text-sm font-medium text-blue-800 mb-1">ハンズオン用簡略版</p>
                  <p className="text-sm text-blue-700">
                    実際のクレジットカード情報は不要です。「注文を確定する」ボタンを押すだけで注文が完了します。
                  </p>
                </div>
              </div>
            </div>

            {/* 配送先情報（デモ用） */}
            <div className="mb-6">
              <h3 className="text-lg font-medium text-gray-900 mb-3">配送先</h3>
              <div className="bg-gray-50 p-4 rounded-lg">
                <p className="text-sm text-gray-600">デモユーザー様</p>
                <p className="text-sm text-gray-600">〒123-4567</p>
                <p className="text-sm text-gray-600">東京都サンプル区テスト町1-2-3</p>
                <p className="text-sm text-gray-600">サンプルアパート101</p>
              </div>
            </div>

            {/* 決済方法（デモ用） */}
            <div className="mb-6">
              <h3 className="text-lg font-medium text-gray-900 mb-3">決済方法</h3>
              <div className="bg-gray-50 p-4 rounded-lg">
                <p className="text-sm text-gray-600">デモ決済（ハンズオン用）</p>
                <p className="text-xs text-gray-500 mt-1">実際の決済は行われません</p>
              </div>
            </div>

            {/* 注文確定ボタン */}
            <div className="space-y-4">
              <button
                onClick={handlePlaceOrder}
                disabled={isProcessing}
                className={`w-full py-4 px-6 rounded-md font-medium text-lg transition-colors ${
                  isProcessing
                    ? 'bg-gray-400 text-gray-700 cursor-not-allowed'
                    : 'bg-green-600 text-white hover:bg-green-700'
                }`}
              >
                {isProcessing ? (
                  <div className="flex items-center justify-center">
                    <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white mr-2"></div>
                    処理中...
                  </div>
                ) : (
                  '注文を確定する'
                )}
              </button>
              
              <Link
                href="/cart"
                className="block w-full bg-gray-100 text-gray-700 text-center py-3 px-4 rounded-md hover:bg-gray-200 transition-colors font-medium"
              >
                カートに戻る
              </Link>
            </div>

            {/* 注意事項 */}
            <div className="mt-6 pt-6 border-t border-gray-200">
              <h4 className="text-sm font-medium text-gray-900 mb-2">ご注意</h4>
              <ul className="text-xs text-gray-600 space-y-1">
                <li>• これはNew Relic SLMハンズオン用のデモサイトです</li>
                <li>• 実際の商品配送や決済は行われません</li>
                <li>• 入力された情報は保存されません</li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}