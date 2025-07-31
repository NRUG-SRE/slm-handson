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

export interface OrderItem {
  id: string
  productId: string
  product: Product
  price: number
  quantity: number
  createdAt: string
}

export interface Order {
  id: string
  items: OrderItem[]
  totalAmount: number
  status: 'pending' | 'processing' | 'completed' | 'cancelled'
  createdAt: string
  updatedAt: string
}