import axios from 'axios'
import router from '@/router'

// Create axios instance with default config
const apiClient = axios.create({
  baseURL: '/api',
  timeout: 30000,
  withCredentials: true, // Important for session cookies
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor
apiClient.interceptors.request.use(
  config => {
    // Could add loading indicators here
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// Response interceptor
apiClient.interceptors.response.use(
  response => {
    return response
  },
  error => {
    if (error.response) {
      // Handle 401 Unauthorized
      if (error.response.status === 401) {
        // Import auth store dynamically to avoid circular imports
        import('@/stores/auth').then(({ useAuthStore }) => {
          const authStore = useAuthStore()
          // Clear auth state and redirect to login
          authStore.clearAuthState()
          router.push('/login')
        })
      }
      
      // Return error with message from backend
      const message = error.response.data?.message || error.response.data?.error || 'An error occurred'
      error.message = message
    } else if (error.request) {
      error.message = 'No response from server'
    } else {
      error.message = 'Request failed'
    }
    
    return Promise.reject(error)
  }
)

export default apiClient