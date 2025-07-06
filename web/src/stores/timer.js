import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useTimerStore = defineStore('timer', () => {
  const activeTimers = ref({
    feed: null,
    pump: null,
    sleep: null
  })

  function startTimer(type) {
    activeTimers.value[type] = {
      startTime: new Date(),
      duration: 0
    }
  }

  function stopTimer(type) {
    activeTimers.value[type] = null
  }

  return { activeTimers, startTimer, stopTimer }
})
