export interface Product {
  id: string
  name: string
  description: string
  price: number
  imageUrl: string
  stock: number
}

export interface CartItem {
  id: string
  productId: string
  product: Product
  quantity: number
  createdAt: string
  updatedAt: string
}

export interface Cart {
  id: string
  items: CartItem[]
  totalAmount: number
  createdAt: string
  updatedAt: string
}

export interface Order {
  id: string
  items: CartItem[]
  totalAmount: number
  status: 'pending' | 'completed' | 'failed'
  createdAt: string
}