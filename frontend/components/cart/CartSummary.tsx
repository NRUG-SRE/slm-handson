'use client'

import Link from 'next/link'
import { Cart } from '@/lib/types'

interface CartSummaryProps {
  cart: Cart
}

export default function CartSummary({ cart }: CartSummaryProps) {
  const itemCount = cart.items.reduce((sum, item) => sum + item.quantity, 0)
  const subtotal = cart.totalAmount
  const shipping = subtotal > 10000 ? 0 : 500 // 10,000円以上で送料無料
  const tax = Math.floor(subtotal * 0.1) // 消費税10%
  const total = subtotal + shipping + tax

  return (
    <div className="bg-white rounded-lg shadow-md">
      <div className="p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">注文内容</h2>
        
        <div className="space-y-2 mb-4">
          <div className="flex justify-between">
            <span className="text-gray-600">商品数</span>
            <span className="font-medium">{itemCount}点</span>
          </div>
          
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

        {/* 送料無料まであといくらかの表示 */}
        {shipping > 0 && (
          <div className="bg-blue-50 p-3 rounded-md mb-4">
            <p className="text-sm text-blue-700">
              あと¥{(10000 - subtotal).toLocaleString()}で送料無料！
            </p>
          </div>
        )}

        {/* 注文進行ボタン */}
        <div className="space-y-3">
          <Link
            href="/checkout"
            className="w-full bg-blue-600 text-white text-center py-3 px-4 rounded-md hover:bg-blue-700 transition-colors font-medium block"
          >
            レジに進む
          </Link>
          
          <Link
            href="/"
            className="w-full bg-gray-100 text-gray-700 text-center py-2 px-4 rounded-md hover:bg-gray-200 transition-colors block"
          >
            買い物を続ける
          </Link>
        </div>

        {/* 配送・返品情報 */}
        <div className="mt-6 pt-4 border-t text-sm text-gray-600">
          <div className="space-y-2">
            <div className="flex items-center">
              <svg className="w-4 h-4 mr-2 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
              <span>1-3営業日でお届け</span>
            </div>
            
            <div className="flex items-center">
              <svg className="w-4 h-4 mr-2 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
              <span>30日間返品保証</span>
            </div>
            
            <div className="flex items-center">
              <svg className="w-4 h-4 mr-2 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
              <span>安全な決済システム</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}