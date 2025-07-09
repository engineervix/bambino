import { defineStore } from 'pinia'
import { ref } from 'vue'
import apiClient from '@/api/client'
import { useAuthStore } from './auth'

export const useActivityStore = defineStore('activity', () => {
  // State
  const activities = ref([])
  const loading = ref(false)
  const error = ref(null)
  const currentActivity = ref(null)
  const pagination = ref({
    page: 1,
    pageSize: 20,
    total: 0,
    totalPages: 0
  })

  // Activity types configuration
  const activityTypes = ref([
    {
      id: 'feed',
      title: 'Feed',
      description: 'Track a feeding session',
      icon: 'mdi-baby-bottle',
      color: 'feed',
      hasTimer: true
    },
    {
      id: 'pump',
      title: 'Pump',
      description: 'Track a pumping session',
      icon: 'mdi-mother-nurse',
      color: 'pump',
      hasTimer: true
    },
    {
      id: 'diaper',
      title: 'Diaper',
      description: 'Track a diaper change',
      icon: 'mdi-baby',
      color: 'diaper',
      hasTimer: false
    },
    {
      id: 'sleep',
      title: 'Sleep',
      description: 'Track a sleep session',
      icon: 'mdi-sleep',
      color: 'sleep',
      hasTimer: true
    },
    {
      id: 'growth',
      title: 'Growth',
      description: 'Record measurements',
      icon: 'mdi-human-male-height',
      color: 'growth',
      hasTimer: false
    },
    {
      id: 'health',
      title: 'Health',
      description: 'Medical records & vaccines',
      icon: 'mdi-medical-bag',
      color: 'health',
      hasTimer: false
    },
    {
      id: 'milestone',
      title: 'Milestone',
      description: 'Track memorable moments',
      icon: 'mdi-party-popper',
      color: 'milestone',
      hasTimer: false
    }
  ])

  // Actions
  async function createActivity(activityData) {
    loading.value = true
    error.value = null

    const authStore = useAuthStore()
    const babyId = authStore.currentBaby?.id

    if (!babyId) {
      error.value = 'No baby selected'
      loading.value = false
      return { success: false, error: 'No baby selected' }
    }

    try {
      const response = await apiClient.post('/activities', {
        ...activityData,
        baby_id: babyId
      })

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

  async function fetchActivities(params = {}) {
    loading.value = true
    error.value = null

    try {
      const response = await apiClient.get('/activities', { params })

      // Handle pagination properly
      if (params.page > 1) {
        // Append for pagination
        activities.value.push(...response.data.activities)
      } else {
        // Replace for new queries
        activities.value = response.data.activities
      }

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

  // Timer operations
  async function startTimer(type, initialData = {}) {
    loading.value = true
    error.value = null

    const authStore = useAuthStore()
    const babyId = authStore.currentBaby?.id

    if (!babyId) {
      error.value = 'No baby selected'
      loading.value = false
      return { success: false, error: 'No baby selected' }
    }

    try {
      const response = await apiClient.post('/activities/timer/start', {
        type,
        baby_id: babyId,
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

      // Update activity in list if it exists
      const index = activities.value.findIndex(a => a.id === activityId)
      if (index !== -1) {
        activities.value[index] = response.data
      }

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
