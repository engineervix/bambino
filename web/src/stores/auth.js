import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import apiClient from '@/api/client'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(null)
  const babies = ref([])
  const selectedBaby = ref(null)
  const loading = ref(false)
  const error = ref(null)

  // Getters
  const isAuthenticated = computed(() => !!user.value)
  const username = computed(() => user.value?.username || '')
  const currentBaby = computed(() => selectedBaby.value)
  const hasBaby = computed(() => babies.value.length > 0)

  // Actions
  async function login(credentials) {
    loading.value = true
    error.value = null
    
    try {
      const response = await apiClient.post('/auth/login', credentials)
      await checkAuth() // Get full user data and babies
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
      babies.value = []
      selectedBaby.value = null
      router.push('/login')
    } catch (err) {
      // Even if logout fails, clear local state
      user.value = null
      babies.value = []
      selectedBaby.value = null
      router.push('/login')
    } finally {
      loading.value = false
    }
  }

  async function checkAuth() {
    try {
      const response = await apiClient.get('/auth/me')
      user.value = response.data
      
      // Load babies
      await loadBabies()
      
      return true
    } catch (err) {
      user.value = null
      babies.value = []
      selectedBaby.value = null
      return false
    }
  }

  async function loadBabies() {
    try {
      const response = await apiClient.get('/babies')
      babies.value = response.data
      
      // Auto-select first baby if available
      if (babies.value.length > 0 && !selectedBaby.value) {
        selectedBaby.value = babies.value[0]
      }
      
      return { success: true, data: response.data }
    } catch (err) {
      error.value = err.message
      return { success: false, error: err.message }
    }
  }

  function selectBaby(baby) {
    selectedBaby.value = baby
    // Could persist selection to localStorage
    if (baby) {
      localStorage.setItem('selectedBabyId', baby.id)
    }
  }

  // Load selected baby from localStorage on init
  function loadSelectedBaby() {
    const savedBabyId = localStorage.getItem('selectedBabyId')
    if (savedBabyId && babies.value.length > 0) {
      const baby = babies.value.find(b => b.id === savedBabyId)
      if (baby) {
        selectedBaby.value = baby
      }
    }
  }

  // Clear error
  function clearError() {
    error.value = null
  }

  return {
    // State
    user,
    babies,
    selectedBaby,
    loading,
    error,
    // Getters
    isAuthenticated,
    username,
    currentBaby,
    hasBaby,
    // Actions
    login,
    logout,
    checkAuth,
    loadBabies,
    selectBaby,
    loadSelectedBaby,
    clearError
  }
})