const API_BASE = import.meta.env.VITE_API_URL || ''

export const api = async (endpoint: string, options: RequestInit = {}) => {
  const token = localStorage.getItem('token')
  
  const config: RequestInit = {
    headers: {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` }),
    },
    ...options,
  }

  // Правильная обработка путей
  let fullUrl = endpoint
  if (!endpoint.startsWith('/')) {
    fullUrl = '/' + endpoint
  }

  try {
    const response = await fetch(`${API_BASE}${fullUrl}`, config)
    
    if (!response.ok) {
      let errorMessage = 'Something went wrong'
      try {
        const errorData = await response.json()
        errorMessage = errorData.message || errorMessage
      } catch (e) {}
      throw new Error(errorMessage)
    }

    return await response.json()
  } catch (err: any) {
    console.error('API Error:', err)
    if (err.message === 'Failed to fetch') {
      throw new Error('Cannot connect to server. Check if Docker is running.')
    }
    throw err
  }
}