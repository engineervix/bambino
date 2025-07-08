import { defineStore } from 'pinia'
import { ref } from 'vue'
import apiClient from '@/api/client'

export const useActivityStore = defineStore('activity', () => {
  // State
  const activities = ref([])
  const loading = ref(false)
  const error = ref(null)
  const currentActivity = ref(null) // For edit operations
  const pagination = ref({
    page: 1,
    pageSize: 20,
    total: 0,
    totalPages: 0
  })

  // Activity types configuration (keep your existing config)
  const activityTypes = ref([
    {
      id: 'feed',
      title: 'Feed',
      description: 'Track a feeding session',
      icon: 'mdi-bottle-baby',
      color: 'feed',
      hasTimer: true
    },
    // ... rest of your activity types
  ])

  // Enhanced Actions
  async function createActivity(activityData) {
    loading.value = true
    error.value = null
    
    try {
      const response = await apiClient.post('/activities', activityData)
      
      // Add to beginning of activities list if we're on first page
      if (pagination.value.page === 1) {
        activities.value.unshift(response.data)
        // Remove last item to maintain page size
        if (activities.value.length > pagination.value.pageSize) {
          activities.value.pop()
        }
      }
      
      return { success: true, data: response.data }
    } catch (err) {
      error.value = err.message
      return { success: false, error: err.message }
    } finally {
      loading.value = false
    }
  }

  // NEW: Get single activity
  async function getActivity(id) {
    loading.value = true
    error.value = null
    
    try {
      const response = await apiClient.get(`/activities/${id}`)
      currentActivity.value = response.data
      return { success: true, data: response.data }
    } catch (err) {
      error.value = err.message
      return { success: false, error: err.message }
    } finally {
      loading.value = false
    }
  }

  // NEW: Update activity
  async function updateActivity(id, activityData) {
    loading.value = true
    error.value = null
    
    try {
      const response = await apiClient.put(`/activities/${id}`, activityData)
      
      // Update in activities list
      const index = activities.value.findIndex(a => a.id === id)
      if (index !== -1) {
        activities.value[index] = response.data
      }
      
      currentActivity.value = response.data
      return { success: true, data: response.data }
    } catch (err) {
      error.value = err.message
      return { success: false, error: err.message }
    } finally {
      loading.value = false
    }
  }

  // NEW: Delete activity
  async function deleteActivity(id) {
    loading.value = true
    error.value = null
    
    try {
      await apiClient.delete(`/activities/${id}`)
      
      // Remove from activities list
      const index = activities.value.findIndex(a => a.id === id)
      if (index !== -1) {
        activities.value.splice(index, 1)
        pagination.value.total--
      }
      
      return { success: true }
    } catch (err) {
      error.value = err.message
      return { success: false, error: err.message }
    } finally {
      loading.value = false
    }
  }

  // Enhanced fetch with better pagination support
  async function fetchActivities(params = {}) {
    loading.value = true
    error.value = null
    
    try {
      const response = await apiClient.get('/activities', { params })
      
      activities.value = response.data.activities
      pagination.value = {
        page: response.data.page,
        pageSize: response.data.page_size,
        total: response.data.total,
        totalPages: response.data.total_pages
      }
      
      return { success: true, data: response.data }
    } catch (err) {
      error.value = err.message
      return { success: false, error: err.message }
    } finally {
      loading.value = false
    }
  }

  // Timer operations (keep your existing implementation)
  async function startTimer(type, initialData = {}) {
    loading.value = true
    error.value = null
    
    try {
      const response = await apiClient.post('/activities/timer/start', {
        type,
        ...initialData
      })
      return { success: true, data: response.data }
    } catch (err) {
      error.value = err.message
      return { success: false, error: err.message }
    } finally {
      loading.value = false
    }
  }

  async function stopTimer(activityId, additionalData = {}) {
    loading.value = true
    error.value = null
    
    try {
      const response = await apiClient.put(`/activities/timer/${activityId}/stop`, additionalData)
      return { success: true, data: response.data }
    } catch (err) {
      error.value = err.message
      return { success: false, error: err.message }
    } finally {
      loading.value = false
    }
  }

  // Statistics
  async function getRecentStats() {
    try {
      const response = await apiClient.get('/stats/recent')
      return { success: true, data: response.data }
    } catch (err) {
      return { success: false, error: err.message }
    }
  }

  async function getDailyStats(date = null) {
    try {
      const params = date ? { date } : {}
      const response = await apiClient.get('/stats/daily', { params })
      return { success: true, data: response.data }
    } catch (err) {
      return { success: false, error: err.message }
    }
  }

  async function getWeeklyStats(week = null) {
    try {
      const params = week ? { week } : {}
      const response = await apiClient.get('/stats/weekly', { params })
      return { success: true, data: response.data }
    } catch (err) {
      return { success: false, error: err.message }
    }
  }

  // Utility functions
  function clearError() {
    error.value = null
  }

  function clearCurrentActivity() {
    currentActivity.value = null
  }

  function getActivityTypeConfig(type) {
    return activityTypes.value.find(at => at.id === type)
  }

  return {
    // State
    activities,
    loading,
    error,
    currentActivity,
    pagination,
    activityTypes,
    
    // Actions
    createActivity,
    getActivity,
    updateActivity,
    deleteActivity,
    fetchActivities,
    startTimer,
    stopTimer,
    getRecentStats,
    getDailyStats,
    getWeeklyStats,
    
    // Utilities
    clearError,
    clearCurrentActivity,
    getActivityTypeConfig
  }
})