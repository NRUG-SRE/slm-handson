'use client'

import { useState } from 'react'
import Link from 'next/link'
import { CartItem as CartItemType } from '@/lib/types'

interface CartItemProps {
  item: CartItemType
  onQuantityChange: (itemId: string, newQuantity: number) => void
  onRemove: (itemId: string) => void
}

export default function CartItem({ item, onQuantityChange, onRemove }: CartItemProps) {
  const [isUpdating, setIsUpdating] = useState(false)
  const [showRemoveConfirm, setShowRemoveConfirm] = useState(false)

  const handleQuantityChange = async (newQuantity: number) => {
    if (newQuantity < 1 || newQuantity > item.product.stock || isUpdating) {
      return
    }

    setIsUpdating(true)
    try {
      await onQuantityChange(item.id, newQuantity)
    } finally {
      setIsUpdating(false)
    }
  }

  const handleRemove = async () => {
    setIsUpdating(true)
    try {
      await onRemove(item.id)
    } finally {
      setIsUpdating(false)
      setShowRemoveConfirm(false)
    }
  }

  const subtotal = item.product.price * item.quantity

  return (
    <div className="flex items-center space-x-4 py-4 border-b border-gray-200 last:border-b-0">
      {/* 商品画像 */}
      <div className="flex-shrink-0">
        <Link href={`/products/${item.product.id}`}>
          <div className="w-20 h-20 bg-gray-200 rounded-md overflow-hidden">
            <img
              src={item.product.imageUrl || '/images/placeholder.png'}
              alt={item.product.name}
              className="w-full h-full object-cover object-center hover:opacity-75 transition-opacity"
            />
          </div>
        </Link>
      </div>

      {/* 商品情報 */}
      <div className="flex-grow">
        <div className="flex justify-between">
          <div>
            <Link 
              href={`/products/${item.product.id}`}
              className="text-lg font-semibold text-gray-900 hover:text-blue-600 transition-colors"
            >
              {item.product.name}
            </Link>
            <p className="text-sm text-gray-500 mt-1 truncate">
              {item.product.description}
            </p>
            <p className="text-lg font-bold text-gray-900 mt-2">
              ¥{item.product.price.toLocaleString()}
            </p>
          </div>

          {/* 小計 */}
          <div className="text-right">
            <p className="text-lg font-bold text-gray-900">
              ¥{subtotal.toLocaleString()}
            </p>
          </div>
        </div>

        {/* 数量コントロールと削除ボタン */}
        <div className="flex items-center justify-between mt-4">
          <div className="flex items-center space-x-3">
            <span className="text-sm text-gray-700">数量:</span>
            
            {/* 数量変更コントロール */}
            <div className="flex items-center space-x-2">
              <button
                onClick={() => handleQuantityChange(item.quantity - 1)}
                disabled={item.quantity <= 1 || isUpdating}
                className="w-8 h-8 rounded-md border border-gray-300 flex items-center justify-center hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                -
              </button>
              
              <span className="w-12 text-center font-medium">
                {isUpdating ? (
                  <div className="w-4 h-4 mx-auto animate-spin rounded-full border-2 border-blue-600 border-t-transparent"></div>
                ) : (
                  item.quantity
                )}
              </span>
              
              <button
                onClick={() => handleQuantityChange(item.quantity + 1)}
                disabled={item.quantity >= item.product.stock || isUpdating}
                className="w-8 h-8 rounded-md border border-gray-300 flex items-center justify-center hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                +
              </button>
            </div>

            <span className="text-xs text-gray-500">
              (在庫: {item.product.stock}個)
            </span>
          </div>

          {/* 削除ボタン */}
          <div className="flex items-center space-x-2">
            {!showRemoveConfirm ? (
              <button
                onClick={() => setShowRemoveConfirm(true)}
                disabled={isUpdating}
                className="text-red-600 hover:text-red-800 text-sm font-medium disabled:opacity-50"
              >
                削除
              </button>
            ) : (
              <div className="flex items-center space-x-2">
                <span className="text-sm text-red-600">削除しますか?</span>
                <button
                  onClick={handleRemove}
                  disabled={isUpdating}
                  className="bg-red-600 text-white text-xs px-2 py-1 rounded hover:bg-red-700 disabled:opacity-50"
                >
                  はい
                </button>
                <button
                  onClick={() => setShowRemoveConfirm(false)}
                  disabled={isUpdating}
                  className="bg-gray-300 text-gray-700 text-xs px-2 py-1 rounded hover:bg-gray-400 disabled:opacity-50"
                >
                  いいえ
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}