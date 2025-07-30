import Link from 'next/link'
import { Product } from '@/lib/types'
import { useMonitoring } from '@/lib/monitoring'

interface ProductCardProps {
  product: Product
}

export default function ProductCard({ product }: ProductCardProps) {
  const { trackUserAction } = useMonitoring()

  const handleProductClick = () => {
    trackUserAction('ProductView', {
      productId: product.id,
      productName: product.name,
      price: product.price,
      inStock: product.stock > 0
    })
  }
  return (
    <div className="bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow duration-300">
      <Link href={`/products/${product.id}`} onClick={handleProductClick}>
        <div className="aspect-w-1 aspect-h-1 w-full overflow-hidden rounded-t-lg bg-gray-200">
          <img
            src={product.imageUrl || '/images/placeholder.png'}
            alt={product.name}
            className="h-64 w-full object-cover object-center hover:opacity-75 transition-opacity"
          />
        </div>
      </Link>
      
      <div className="p-4">
        <Link href={`/products/${product.id}`} onClick={handleProductClick}>
          <h3 className="text-lg font-semibold text-gray-900 hover:text-blue-600 transition-colors">
            {product.name}
          </h3>
        </Link>
        
        <p className="mt-1 text-sm text-gray-500 line-clamp-2">
          {product.description}
        </p>
        
        <div className="mt-4 flex items-center justify-between">
          <p className="text-2xl font-bold text-gray-900">
            ¥{product.price.toLocaleString()}
          </p>
          
          {product.stock > 0 ? (
            <span className="text-sm text-green-600 font-medium">
              在庫あり
            </span>
          ) : (
            <span className="text-sm text-red-600 font-medium">
              在庫なし
            </span>
          )}
        </div>
        
        <Link 
          href={`/products/${product.id}`}
          className="mt-4 block w-full bg-blue-600 text-white text-center py-2 px-4 rounded-md hover:bg-blue-700 transition-colors"
          onClick={handleProductClick}
        >
          詳細を見る
        </Link>
      </div>
    </div>
  )
}