'use client'

import { useEffect, useState } from 'react'
import ProductCard from './ProductCard'
import { Product } from '@/lib/types'
import { productApi, ApiError } from '@/lib/api'
import { useMonitoring } from '@/lib/monitoring'

// モックデータ（API接続失敗時のフォールバック）
const mockProducts: Product[] = [
  {
    id: '1',
    name: 'ワイヤレスヘッドホン',
    description: '高音質なノイズキャンセリング機能付きワイヤレスヘッドホン',
    price: 25000,
    imageUrl: '/images/headphones.jpg',
    stock: 10
  },
  {
    id: '2',
    name: 'スマートウォッチ',
    description: 'フィットネストラッキング機能付きの最新スマートウォッチ',
    price: 35000,
    imageUrl: '/images/smartwatch.jpg',
    stock: 5
  },
  {
    id: '3',
    name: 'ポータブルスピーカー',
    description: '防水機能付きの高音質Bluetoothスピーカー',
    price: 12000,
    imageUrl: '/images/speaker.jpg',
    stock: 15
  },
  {
    id: '4',
    name: 'ワイヤレスキーボード',
    description: '人間工学に基づいたデザインのワイヤレスキーボード',
    price: 8500,
    imageUrl: '/images/keyboard.jpg',
    stock: 20
  },
  {
    id: '5',
    name: '4K Webカメラ',
    description: 'リモートワークに最適な高画質Webカメラ',
    price: 15000,
    imageUrl: '/images/webcam.jpg',
    stock: 8
  },
  {
    id: '6',
    name: 'USB-C ハブ',
    description: '7つのポートを備えた多機能USB-Cハブ',
    price: 6500,
    imageUrl: '/images/usb-hub.jpg',
    stock: 0
  }
]

export default function ProductList() {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [useMockData, setUseMockData] = useState(false)
  const { trackPageView, trackError } = useMonitoring()

  useEffect(() => {
    trackPageView('ProductList')
    fetchProducts()
  }, [])

  const fetchProducts = async () => {
    try {
      setLoading(true)
      setError(null)
      
      // 実際のAPI呼び出しを試行
      const data = await productApi.getProducts()
      setProducts(data)
      setUseMockData(false)
    } catch (err) {
      console.warn('API接続に失敗しました。モックデータを使用します:', err)
      
      // API接続失敗時はモックデータを使用
      setProducts(mockProducts)
      setUseMockData(true)
      
      if (err instanceof ApiError) {
        trackError(err, 'ProductList API Error')
      } else if (err instanceof Error) {
        trackError(err, 'ProductList Network Error')
      }
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="text-center py-8">
        <p className="text-red-600">{error}</p>
        <button 
          onClick={fetchProducts}
          className="mt-4 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        >
          再試行
        </button>
      </div>
    )
  }

  return (
    <div>
      {useMockData && (
        <div className="mb-4 p-3 bg-yellow-100 border border-yellow-400 rounded-md">
          <p className="text-yellow-800 text-sm">
            ⚠️ API接続に失敗したため、デモ用のモックデータを表示しています
          </p>
        </div>
      )}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {products.map((product) => (
          <ProductCard key={product.id} product={product} />
        ))}
      </div>
    </div>
  )
}