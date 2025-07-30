import ProductList from '@/components/product/ProductList'

export default function Home() {
  return (
    <div className="container mx-auto px-4 py-8">
      <section className="mb-8">
        <h1 className="text-4xl font-bold text-gray-900 mb-4">
          SLM ハンズオン ECサイト
        </h1>
        <p className="text-lg text-gray-600 mb-2">
          New Relic Service Level Managementを学ぶためのデモアプリケーション
        </p>
        <p className="text-sm text-gray-500">
          商品を選んでカートに追加し、購入フローを体験してください。
        </p>
      </section>

      <section>
        <h2 className="text-2xl font-semibold text-gray-800 mb-6">
          商品一覧
        </h2>
        <ProductList />
      </section>
    </div>
  )
}