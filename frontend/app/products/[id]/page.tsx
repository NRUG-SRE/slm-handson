'use client'

import { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import Link from 'next/link'
import { Product } from '@/lib/types'
import { productApi, cartApi } from '@/lib/api'
import { trackECommerceAction, trackError } from '@/lib/newrelic'

export default function ProductDetailPage() {
  const params = useParams()
  const id = params.id as string
  
  const [product, setProduct] = useState<Product | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [quantity, setQuantity] = useState(1)
  const [isAddingToCart, setIsAddingToCart] = useState(false)
  const [cartMessage, setCartMessage] = useState<string | null>(null)

  useEffect(() => {
    const fetchProduct = async () => {
      try {
        setLoading(true)
        const data = await productApi.getProduct(id)
        setProduct(data)
        
        // New Relic: 商品閲覧を記録
        trackECommerceAction.viewProduct(data.id, data.name, data.price)
      } catch (err) {
        console.error('商品情報の取得に失敗しました:', err)
        setError('商品情報の取得に失敗しました')
      } finally {
        setLoading(false)
      }
    }

    fetchProduct()
  }, [id])

  const handleQuantityChange = (value: number) => {
    if (value >= 1 && product && value <= product.stock) {
      setQuantity(value)
    }
  }

  const handleAddToCart = async () => {
    if (!product) return

    try {
      setIsAddingToCart(true)
      setCartMessage(null)
      
      await cartApi.addToCart(product.id, quantity)
      
      // New Relic: カート追加イベントを記録
      trackECommerceAction.addToCart(product.id, product.name, quantity, product.price)
      
      // カート更新イベントを発火
      window.dispatchEvent(new CustomEvent('cartUpdated'))
      
      setCartMessage(`${product.name}を${quantity}個カートに追加しました`)
      
      // 3秒後にメッセージをクリア
      setTimeout(() => {
        setCartMessage(null)
      }, 3000)
    } catch (err) {
      console.error('カートへの追加に失敗しました:', err)
      setCartMessage('カートへの追加に失敗しました')
    } finally {
      setIsAddingToCart(false)
    }
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  if (error || !product) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
          <p>{error || '商品が見つかりませんでした'}</p>
        </div>
        <Link href="/" className="mt-4 inline-block text-blue-600 hover:text-blue-800">
          ← TOPページへ戻る
        </Link>
      </div>
    )
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* パンくずリスト */}
      <nav className="text-sm mb-6">
        <ol className="flex items-center space-x-2">
          <li>
            <Link href="/" className="text-gray-500 hover:text-gray-700">
              TOP
            </Link>
          </li>
          <li>
            <span className="mx-2 text-gray-400">/</span>
          </li>
          <li className="text-gray-900">{product.name}</li>
        </ol>
      </nav>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        {/* 商品画像 */}
        <div className="relative w-full pt-[100%] bg-gray-200 rounded-lg overflow-hidden">
          <img
            src={product.imageUrl}
            alt={product.name}
            className="absolute inset-0 w-full h-full object-center object-cover"
          />
        </div>

        {/* 商品情報 */}
        <div>
          <h1 className="text-3xl font-extrabold text-gray-900 mb-4">
            {product.name}
          </h1>
          
          <p className="text-gray-600 mb-6">
            {product.description}
          </p>

          <div className="mb-6">
            <span className="text-3xl font-bold text-gray-900">
              ¥{product.price.toLocaleString()}
            </span>
            <span className="text-sm text-gray-500 ml-2">(税込)</span>
          </div>

          {/* 在庫情報 */}
          <div className="mb-6">
            {product.stock > 0 ? (
              <div>
                <span className="text-sm text-green-600">在庫あり</span>
                <span className="text-sm text-gray-500 ml-2">
                  (残り{product.stock}個)
                </span>
              </div>
            ) : (
              <span className="text-sm text-red-600">在庫切れ</span>
            )}
          </div>

          {/* 数量選択 */}
          {product.stock > 0 && (
            <div className="mb-6">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                数量
              </label>
              <div className="flex items-center space-x-3">
                <button
                  onClick={() => handleQuantityChange(quantity - 1)}
                  disabled={quantity <= 1}
                  className="w-10 h-10 rounded-md border border-gray-300 flex items-center justify-center hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  -
                </button>
                <input
                  type="number"
                  value={quantity}
                  onChange={(e) => handleQuantityChange(parseInt(e.target.value) || 1)}
                  min="1"
                  max={product.stock}
                  className="w-20 text-center border border-gray-300 rounded-md py-2 px-3"
                />
                <button
                  onClick={() => handleQuantityChange(quantity + 1)}
                  disabled={quantity >= product.stock}
                  className="w-10 h-10 rounded-md border border-gray-300 flex items-center justify-center hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  +
                </button>
              </div>
            </div>
          )}

          {/* カートに追加ボタン */}
          <button
            onClick={handleAddToCart}
            disabled={product.stock === 0 || isAddingToCart}
            className="w-full bg-blue-600 text-white py-3 px-8 rounded-md font-medium hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors duration-200"
          >
            {isAddingToCart ? (
              <span className="flex items-center justify-center">
                <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                処理中...
              </span>
            ) : product.stock === 0 ? (
              '在庫切れ'
            ) : (
              'カートに追加'
            )}
          </button>

          {/* カート追加メッセージ */}
          {cartMessage && (
            <div className={`mt-4 p-3 rounded-md ${
              cartMessage.includes('失敗') 
                ? 'bg-red-100 text-red-700' 
                : 'bg-green-100 text-green-700'
            }`}>
              {cartMessage}
            </div>
          )}

          {/* 商品詳細情報 */}
          <div className="mt-8 border-t pt-8">
            <h2 className="text-lg font-medium text-gray-900 mb-4">商品詳細</h2>
            <dl className="space-y-3">
              <div className="flex justify-between">
                <dt className="text-sm text-gray-500">商品ID</dt>
                <dd className="text-sm text-gray-900">{product.id}</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-sm text-gray-500">カテゴリ</dt>
                <dd className="text-sm text-gray-900">電子機器</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-sm text-gray-500">配送目安</dt>
                <dd className="text-sm text-gray-900">1-3営業日</dd>
              </div>
            </dl>
          </div>
        </div>
      </div>

      {/* 関連アクション */}
      <div className="mt-12 flex justify-between items-center">
        <Link 
          href="/" 
          className="text-blue-600 hover:text-blue-800 flex items-center"
        >
          <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          商品一覧へ戻る
        </Link>
        
        <Link 
          href="/cart" 
          className="text-blue-600 hover:text-blue-800 flex items-center"
        >
          カートを見る
          <svg className="w-5 h-5 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14 5l7 7m0 0l-7 7m7-7H3" />
          </svg>
        </Link>
      </div>
    </div>
  )
}