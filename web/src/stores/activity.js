import { defineStore } from 'pinia'
import { ref } from 'vue'
import apiClient from '@/api/client'

export const useActivityStore = defineStore('activity', () => {
  // State
  const activities = ref([])
  const loading = ref(false)
  const error = ref(null)

  // Activity types configuration
  const activityTypes = ref([
    {
      id: 'feed',
      title: 'Feed',
      description: 'Track a feeding session',
      icon: 'mdi-bottle-baby',
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
      hasTimer: false,
      subTypes: ['weight', 'height', 'head']
    },
    {
      id: 'health',
      title: 'Health',
      description: 'Medical records',
      icon: 'mdi-medical-bag',
      color: 'health',
      hasTimer: false,
      subTypes: ['checkup', 'vaccine', 'illness']
    },
    {
      id: 'milestone',
      title: 'Baby Firsts',
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
    
    try {
      const response = await apiClient.post('/activities', activityData)
      activities.value.unshift(response.data)
      return { success: true, data: response.data }
    } catch (err) {
      error.value = err.message
      return { success: false, error: err.message }
    } finally {
      loading.value = false
    }
  }

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

  async function fetchActivities(params = {}) {
    loading.value = true
    error.value = null
    
    try {
      const response = await apiClient.get('/activities', { params })
      activities.value = response.data.activities
      return { success: true, data: response.data }
    } catch (err) {
      error.value = err.message
      return { success: false, error: err.message }
    } finally {
      loading.value = false
    }
  }

  async function getRecentStats() {
    try {
      const response = await apiClient.get('/stats/recent')
      return { success: true, data: response.data }
    } catch (err) {
      return { success: false, error: err.message }
    }
  }

  return {
    // State
    activities,
    loading,
    error,
    activityTypes,
    // Actions
    createActivity,
    startTimer,
    stopTimer,
    fetchActivities,
    getRecentStats
  }
})