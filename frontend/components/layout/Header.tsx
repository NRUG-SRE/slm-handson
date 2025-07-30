'use client'

import Link from 'next/link'
import { useState, useEffect } from 'react'
import { cartApi } from '@/lib/api'

export default function Header() {
  const [cartItemCount, setCartItemCount] = useState(0)

  const fetchCartItemCount = async () => {
    try {
      const cart = await cartApi.getCart()
      const totalItems = cart.items.reduce((sum, item) => sum + item.quantity, 0)
      setCartItemCount(totalItems)
    } catch (error) {
      console.error('Failed to fetch cart item count:', error)
      setCartItemCount(0)
    }
  }

  useEffect(() => {
    fetchCartItemCount()
    
    // カート更新イベントリスナーを追加
    const handleCartUpdate = () => {
      fetchCartItemCount()
    }
    
    window.addEventListener('cartUpdated', handleCartUpdate)
    
    return () => {
      window.removeEventListener('cartUpdated', handleCartUpdate)
    }
  }, [])

  return (
    <header className="bg-blue-600 text-white shadow-lg">
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <Link href="/" className="text-2xl font-bold hover:text-blue-50 transition-colors">
              SLM ECサイト
            </Link>
            <div className="flex items-center space-x-2 text-blue-100">
              <img 
                src="/images/nrug-sre.png" 
                alt="NRUG-SRE" 
                className="w-8 h-8 object-contain"
              />
              <span className="text-sm font-medium">NRUG-SRE支部 Presents</span>
            </div>
          </div>
          
          <nav className="flex items-center space-x-6">
            <Link href="/" className="hover:text-blue-50 transition-colors">
              商品一覧
            </Link>
            <Link 
              href="/cart" 
              className="flex items-center space-x-2 hover:text-blue-50 transition-colors"
            >
              <svg 
                className="w-6 h-6" 
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path 
                  strokeLinecap="round" 
                  strokeLinejoin="round" 
                  strokeWidth={2} 
                  d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" 
                />
              </svg>
              <span>カート</span>
              {cartItemCount > 0 && (
                <span className="bg-blue-50 text-blue-700 rounded-full px-2 py-1 text-xs font-bold">
                  {cartItemCount}
                </span>
              )}
            </Link>
          </nav>
        </div>
      </div>
    </header>
  )
}