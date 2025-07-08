import { defineStore } from 'pinia'
import { ref } from 'vue'
import apiClient from '@/api/client'

export const useTimerStore = defineStore('timer', () => {
  const activeTimers = ref({
    feed: null,
    pump: null,
    sleep: null
  })

  const STORAGE_KEY = 'baby-tracker-timers'

  function startTimer(type, activityId) {
    const timer = {
      startTime: new Date().toISOString(), // Store as ISO string
      activityId: activityId,
      type: type
    }
    
    activeTimers.value[type] = timer
    persistTimers()
  }

  function stopTimer(type) {
    activeTimers.value[type] = null
    persistTimers()
  }

  function getActiveTimer(type) {
    return activeTimers.value[type]
  }

  function hasActiveTimer(type) {
    return !!activeTimers.value[type]
  }

  // Persist timers to localStorage
  function persistTimers() {
    try {
      const timersToSave = {}
      Object.keys(activeTimers.value).forEach(type => {
        if (activeTimers.value[type]) {
          timersToSave[type] = activeTimers.value[type]
        }
      })
      localStorage.setItem(STORAGE_KEY, JSON.stringify(timersToSave))
    } catch (error) {
      console.warn('Failed to persist timers:', error)
    }
  }

  // Load timers from localStorage
  function loadPersistedTimers() {
    try {
      const saved = localStorage.getItem(STORAGE_KEY)
      if (saved) {
        const timers = JSON.parse(saved)
        Object.keys(timers).forEach(type => {
          if (timers[type] && timers[type].activityId) {
            activeTimers.value[type] = timers[type]
          }
        })
      }
    } catch (error) {
      console.warn('Failed to load persisted timers:', error)
      // Clear corrupted data
      localStorage.removeItem(STORAGE_KEY)
    }
  }

  // Validate active timers with backend
  async function validateTimers() {
    const currentTimers = { ...activeTimers.value }
    let hasChanges = false

    for (const [type, timer] of Object.entries(currentTimers)) {
      if (!timer || !timer.activityId) continue

      try {
        // Check if activity still exists and has no end_time
        const response = await apiClient.get(`/activities/${timer.activityId}`)
        const activity = response.data

        // Timer is invalid if activity has an end_time
        if (activity.end_time) {
          console.log(`Timer ${type} is stale - activity already ended`)
          activeTimers.value[type] = null
          hasChanges = true
        }
      } catch (error) {
        // Activity doesn't exist or we can't access it
        console.log(`Timer ${type} is invalid - activity not found`)
        activeTimers.value[type] = null
        hasChanges = true
      }
    }

    if (hasChanges) {
      persistTimers()
    }

    return hasChanges
  }

  // Initialize timers (call on app load)
  async function initializeTimers() {
    // First load from localStorage
    loadPersistedTimers()
    
    // Then validate with backend
    await validateTimers()
  }

  // Clear all timers (useful for logout)
  function clearAllTimers() {
    activeTimers.value = {
      feed: null,
      pump: null,
      sleep: null
    }
    localStorage.removeItem(STORAGE_KEY)
  }

  // Get timer duration in seconds
  function getTimerDuration(type) {
    const timer = activeTimers.value[type]
    if (!timer) return 0
    
    const startTime = new Date(timer.startTime)
    return Math.floor((Date.now() - startTime.getTime()) / 1000)
  }

  return { 
    activeTimers, 
    startTimer, 
    stopTimer,
    getActiveTimer,
    hasActiveTimer,
    initializeTimers,
    validateTimers,
    clearAllTimers,
    getTimerDuration
  }
})