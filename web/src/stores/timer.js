import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import apiClient from '@/api/client'

export const useTimerStore = defineStore('timer', () => {
  const activeTimers = ref({
    feed: null,
    pump: null,
    sleep: null
  })

  // Debounced persistence to avoid excessive localStorage writes
  let persistenceTimeout = null
  const STORAGE_KEY = 'baby-tracker-timers'
  const PERSISTENCE_DELAY = 300 // ms

  // Computed properties for easy access
  const hasActiveTimers = computed(() => {
    return Object.values(activeTimers.value).some(timer => timer !== null)
  })

  const activeTimerCount = computed(() => {
    return Object.values(activeTimers.value).filter(timer => timer !== null).length
  })

  const getActiveTimerTypes = computed(() => {
    return Object.keys(activeTimers.value).filter(type => activeTimers.value[type] !== null)
  })

  function startTimer(type, activityId) {
    // Validate inputs
    if (!type || !activityId) {
      console.error('Invalid timer start parameters:', { type, activityId })
      return false
    }

    // Check if timer already exists for this type
    if (activeTimers.value[type]) {
      console.warn(`Timer already active for ${type}`)
      return false
    }

    const timer = {
      startTime: new Date().toISOString(),
      activityId: activityId,
      type: type
    }
    
    activeTimers.value[type] = timer
    debouncedPersist()
    
    console.log(`Timer started for ${type}:`, timer)
    return true
  }

  function stopTimer(type) {
    if (!activeTimers.value[type]) {
      console.warn(`No active timer found for ${type}`)
      return false
    }

    console.log(`Timer stopped for ${type}`)
    activeTimers.value[type] = null
    debouncedPersist()
    return true
  }

  function getActiveTimer(type) {
    return activeTimers.value[type]
  }

  function hasActiveTimer(type) {
    return !!activeTimers.value[type]
  }

  // Debounced persistence to avoid excessive localStorage writes
  function debouncedPersist() {
    if (persistenceTimeout) {
      clearTimeout(persistenceTimeout)
    }
    
    persistenceTimeout = setTimeout(() => {
      persistTimers()
      persistenceTimeout = null
    }, PERSISTENCE_DELAY)
  }

  // Immediate persistence (for critical operations)
  function persistTimers() {
    try {
      const timersToSave = {}
      Object.keys(activeTimers.value).forEach(type => {
        if (activeTimers.value[type]) {
          timersToSave[type] = activeTimers.value[type]
        }
      })
      
      if (Object.keys(timersToSave).length > 0) {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(timersToSave))
      } else {
        localStorage.removeItem(STORAGE_KEY)
      }
      
      console.log('Timers persisted:', timersToSave)
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
        
        // Validate loaded timers
        Object.keys(timers).forEach(type => {
          if (timers[type] && timers[type].activityId && timers[type].startTime) {
            // Validate that startTime is a valid date
            const startTime = new Date(timers[type].startTime)
            if (!isNaN(startTime.getTime())) {
              activeTimers.value[type] = timers[type]
            } else {
              console.warn(`Invalid timer startTime for ${type}:`, timers[type])
            }
          }
        })
        
        console.log('Timers loaded from storage:', activeTimers.value)
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
    const validationPromises = []

    for (const [type, timer] of Object.entries(currentTimers)) {
      if (!timer || !timer.activityId) continue

      const validationPromise = (async () => {
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
          console.log(`Timer ${type} is invalid - activity not found:`, error.message)
          activeTimers.value[type] = null
          hasChanges = true
        }
      })()

      validationPromises.push(validationPromise)
    }

    // Wait for all validations to complete
    await Promise.all(validationPromises)

    if (hasChanges) {
      persistTimers() // Use immediate persistence for validation changes
    }

    return hasChanges
  }

  // Initialize timers (call on app load)
  async function initializeTimers() {
    console.log('Initializing timers...')
    
    // First load from localStorage
    loadPersistedTimers()
    
    // Then validate with backend (only if we have timers)
    if (hasActiveTimers.value) {
      try {
        await validateTimers()
      } catch (error) {
        console.warn('Timer validation failed:', error)
      }
    }
    
    console.log('Timer initialization complete')
  }

  // Clear all timers (useful for logout)
  function clearAllTimers() {
    console.log('Clearing all timers')
    activeTimers.value = {
      feed: null,
      pump: null,
      sleep: null
    }
    
    // Clear any pending persistence
    if (persistenceTimeout) {
      clearTimeout(persistenceTimeout)
      persistenceTimeout = null
    }
    
    localStorage.removeItem(STORAGE_KEY)
  }

  // Get timer duration in seconds
  function getTimerDuration(type) {
    const timer = activeTimers.value[type]
    if (!timer) return 0
    
    const startTime = new Date(timer.startTime)
    if (isNaN(startTime.getTime())) {
      console.warn(`Invalid startTime for ${type} timer:`, timer.startTime)
      return 0
    }
    
    return Math.floor((Date.now() - startTime.getTime()) / 1000)
  }

  // Get formatted duration string
  function getFormattedDuration(type) {
    const duration = getTimerDuration(type)
    const hours = Math.floor(duration / 3600)
    const minutes = Math.floor((duration % 3600) / 60)
    const seconds = duration % 60
    
    if (hours > 0) {
      return `${hours}h ${minutes}m`
    } else if (minutes > 0) {
      return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
    } else {
      return `0:${seconds.toString().padStart(2, '0')}`
    }
  }

  return { 
    // State
    activeTimers,
    
    // Computed
    hasActiveTimers,
    activeTimerCount,
    getActiveTimerTypes,
    
    // Actions
    startTimer, 
    stopTimer,
    getActiveTimer,
    hasActiveTimer,
    getTimerDuration,
    getFormattedDuration,
    
    // Lifecycle
    initializeTimers,
    validateTimers,
    clearAllTimers,
    
    // Manual persistence control
    persistTimers
  }
})