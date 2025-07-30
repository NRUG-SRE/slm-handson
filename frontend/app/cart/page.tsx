'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { cartApi } from '@/lib/api'
import { useNewRelicMonitoring } from '@/lib/monitoring'
import { Cart } from '@/lib/types'
import CartItem from '@/components/cart/CartItem'
import CartSummary from '@/components/cart/CartSummary'

export default function CartPage() {
  const monitoring = useNewRelicMonitoring()
  const [cart, setCart] = useState<Cart | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchCart = async () => {
    try {
      setLoading(true)
      const cartData = await cartApi.getCart()
      setCart(cartData)
      
      // New Relic: カートページビューを記録
      monitoring.trackPageView('Cart', {
        itemCount: cartData.items.length,
        totalAmount: cartData.totalAmount
      })
    } catch (err) {
      console.error('カート情報の取得に失敗しました:', err)
      setError('カート情報の取得に失敗しました')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchCart()
  }, [])

  const handleQuantityChange = async (itemId: string, newQuantity: number) => {
    if (!cart) return

    try {
      // APIを使用して数量変更
      const updatedCart = await cartApi.updateCartItem(itemId, newQuantity)
      setCart(updatedCart)
      
      // カート更新イベントを発火
      window.dispatchEvent(new CustomEvent('cartUpdated'))
      
      // New Relic: 数量変更イベントを記録
      monitoring.trackUserAction('UpdateCartQuantity', {
        itemId,
        newQuantity,
        productId: cart.items.find(item => item.id === itemId)?.productId
      })
    } catch (err) {
      console.error('数量変更に失敗しました:', err)
      // エラー時は元のデータを再取得
      fetchCart()
    }
  }

  const handleRemoveItem = async (itemId: string) => {
    if (!cart) return

    try {
      // New Relic: アイテム削除イベントを記録（削除前に記録）
      const removedItem = cart.items.find(item => item.id === itemId)
      if (removedItem) {
        monitoring.trackRemoveFromCart(removedItem.productId, removedItem.quantity)
      }
      
      // 数量を0にすることで削除（バックエンドの仕様に合わせる）
      const updatedCart = await cartApi.updateCartItem(itemId, 0)
      setCart(updatedCart)
      
      // カート更新イベントを発火
      window.dispatchEvent(new CustomEvent('cartUpdated'))
    } catch (err) {
      console.error('アイテム削除に失敗しました:', err)
      // エラー時は元のデータを再取得
      fetchCart()
    }
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
          <p>{error}</p>
        </div>
        <Link href="/" className="mt-4 inline-block text-blue-600 hover:text-blue-800">
          ← TOPページへ戻る
        </Link>
      </div>
    )
  }

  if (!cart || cart.items.length === 0) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="text-center">
          <div className="mb-8">
            <svg className="mx-auto h-24 w-24 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          </div>
          <h2 className="text-2xl font-bold text-gray-900 mb-4">カートは空です</h2>
          <p className="text-gray-600 mb-8">
            商品を追加して、お買い物をお楽しみください。
          </p>
          <Link 
            href="/" 
            className="bg-blue-600 text-white px-6 py-3 rounded-md hover:bg-blue-700 transition-colors"
          >
            商品を見る
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">ショッピングカート</h1>
      
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* カートアイテム一覧 */}
        <div className="lg:col-span-2">
          <div className="bg-white rounded-lg shadow-md">
            <div className="p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">
                カート内の商品 ({cart.items.length}件)
              </h2>
              
              <div className="space-y-4">
                {cart.items.map((item) => (
                  <CartItem
                    key={item.id}
                    item={item}
                    onQuantityChange={handleQuantityChange}
                    onRemove={handleRemoveItem}
                  />
                ))}
              </div>
            </div>
          </div>
        </div>

        {/* カートサマリー */}
        <div className="lg:col-span-1">
          <CartSummary cart={cart} />
        </div>
      </div>

      {/* 関連アクション */}
      <div className="mt-8 flex justify-between items-center">
        <Link 
          href="/" 
          className="text-blue-600 hover:text-blue-800 flex items-center"
        >
          <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          買い物を続ける
        </Link>
      </div>
    </div>
  )
}