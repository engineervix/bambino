import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import apiClient from '@/api/client'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(null)
  const loading = ref(false)
  const error = ref(null)

  // Getters
  const isAuthenticated = computed(() => !!user.value)
  const username = computed(() => user.value?.username || '')

  // Actions
  async function login(credentials) {
    loading.value = true
    error.value = null
    
    try {
      const response = await apiClient.post('/auth/login', credentials)
      await checkAuth() // Get full user data
      router.push('/')
      return { success: true }
    } catch (err) {
      error.value = err.message
      return { success: false, error: err.message }
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    loading.value = true
    
    try {
      await apiClient.post('/auth/logout')
      user.value = null
      router.push('/login')
    } catch (err) {
      // Even if logout fails, clear local state
      user.value = null
      router.push('/login')
    } finally {
      loading.value = false
    }
  }

  async function checkAuth() {
    try {
      const response = await apiClient.get('/auth/me')
      user.value = response.data
      return true
    } catch (err) {
      user.value = null
      return false
    }
  }

  // Clear error
  function clearError() {
    error.value = null
  }

  return {
    // State
    user,
    loading,
    error,
    // Getters
    isAuthenticated,
    username,
    // Actions
    login,
    logout,
    checkAuth,
    clearError
  }
})