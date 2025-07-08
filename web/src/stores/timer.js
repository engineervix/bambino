import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useTimerStore = defineStore('timer', () => {
  const activeTimers = ref({
    feed: null,
    pump: null,
    sleep: null
  })

  function startTimer(type, activityId) {
    activeTimers.value[type] = {
      startTime: new Date(),
      activityId: activityId,
      duration: 0
    }
  }

  function stopTimer(type) {
    activeTimers.value[type] = null
  }

  function getActiveTimer(type) {
    return activeTimers.value[type]
  }

  function hasActiveTimer(type) {
    return !!activeTimers.value[type]
  }

  return { 
    activeTimers, 
    startTimer, 
    stopTimer,
    getActiveTimer,
    hasActiveTimer
  }
})