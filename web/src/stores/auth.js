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
  const authChecked = ref(false) // Track if we've done initial auth check

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
      
      // Initialize timers after successful auth
      await initializeUserSession()
      
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
      // Clear timers before logout
      const { useTimerStore } = await import('@/stores/timer')
      const timerStore = useTimerStore()
      timerStore.clearAllTimers()
      
      await apiClient.post('/auth/logout')
      clearAuthState()
      router.push('/login')
    } catch (err) {
      // Even if logout fails, clear local state
      clearAuthState()
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
      
      authChecked.value = true
      return true
    } catch (err) {
      user.value = null
      babies.value = []
      selectedBaby.value = null
      authChecked.value = true
      return false
    }
  }

  // Initial auth check - only called once on app load
  async function initializeAuth() {
    if (authChecked.value) {
      return isAuthenticated.value
    }
    
    const isAuth = await checkAuth()
    if (isAuth) {
      await initializeUserSession()
    }
    return isAuth
  }

  // Initialize user session data (timers, baby selection)
  async function initializeUserSession() {
    try {
      // Load selected baby from localStorage
      loadSelectedBaby()
      
      // Initialize timers
      const { useTimerStore } = await import('@/stores/timer')
      const timerStore = useTimerStore()
      await timerStore.initializeTimers()
    } catch (error) {
      console.warn('Failed to initialize user session:', error)
    }
  }

  // Fast auth check - uses local state if already checked
  function isAuthenticatedFast() {
    return authChecked.value ? isAuthenticated.value : null
  }

  // Force re-check auth (for specific scenarios like session timeout)
  async function recheckAuth() {
    authChecked.value = false
    return await checkAuth()
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
    // Persist selection to localStorage
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

  // Clear auth state (for session timeout/errors)
  function clearAuthState() {
    user.value = null
    babies.value = []
    selectedBaby.value = null
    authChecked.value = false
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
    authChecked,
    // Getters
    isAuthenticated,
    username,
    currentBaby,
    hasBaby,
    // Actions
    login,
    logout,
    checkAuth,
    initializeAuth,
    isAuthenticatedFast,
    recheckAuth,
    loadBabies,
    selectBaby,
    loadSelectedBaby,
    clearAuthState,
    clearError
  }
})
