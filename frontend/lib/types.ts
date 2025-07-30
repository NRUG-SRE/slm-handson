export interface Product {
  id: string
  name: string
  description: string
  price: number
  imageUrl: string
  stock: number
}

export interface CartItem {
  productId: string
  product: Product
  quantity: number
}

export interface Cart {
  items: CartItem[]
  totalAmount: number
}

export interface Order {
  id: string
  items: CartItem[]
  totalAmount: number
  status: 'pending' | 'completed' | 'failed'
  createdAt: string
}