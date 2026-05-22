export interface User {
  id: number
  username: string
  isAdmin: boolean
  createdAt: string
}

export interface Card {
  id: number
  number: string
  balance: string
  isBlocked: boolean
  ownerName: string
  keyValue: string
}

export interface Key {
  id: number
  value: string
  type: 'A' | 'B'
  description: string
}

export interface Terminal {
  id: number
  serialNumber: string
  address: string
  name: string
}

export interface Transaction {
  id: number
  cardNumber: string
  price: string
  terminalSerialNumber: string
}

export interface AuthResponse {
  token: string
}

export interface ApiError {
  message: string
}