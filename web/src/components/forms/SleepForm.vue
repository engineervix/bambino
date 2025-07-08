<template>
  <v-form ref="form" @submit.prevent="handleSubmit">
    <!-- Timer or manual entry -->
    <v-card v-if="!useTimer" variant="outlined" class="mb-4">
      <v-card-text>
        <!-- Start date and time -->
        <v-row class="mb-3">
          <v-col cols="7">
            <v-text-field
              v-model="formData.startDate"
              label="Start Date"
              type="date"
              variant="outlined"
              density="compact"
            />
          </v-col>
          <v-col cols="5">
            <v-text-field
              v-model="formData.startTime"
              label="Start Time"
              type="time"
              variant="outlined"
              density="compact"
            />
          </v-col>
        </v-row>

        <!-- End date and time (optional) -->
        <v-row v-if="showEndTime">
          <v-col cols="7">
            <v-text-field
              v-model="formData.endDate"
              label="End Date"
              type="date"
              variant="outlined"
              density="compact"
            />
          </v-col>
          <v-col cols="5">
            <v-text-field
              v-model="formData.endTime"
              label="End Time"
              type="time"
              variant="outlined"
              density="compact"
              clearable
              @click:clear="clearEndTime"
            />
          </v-col>
        </v-row>

        <!-- Toggle for end time -->
        <v-btn
          v-if="!showEndTime"
          variant="outlined"
          size="small"
          @click="enableEndTime"
          class="mb-3"
        >
          <v-icon start>mdi-clock-plus</v-icon>
          Add End Time
        </v-btn>
      </v-card-text>
    </v-card>

    <!-- Timer display -->
    <v-card v-else variant="tonal" color="sleep" class="mb-4">
      <v-card-text class="text-center">
        <v-icon size="48" class="mb-2">mdi-sleep</v-icon>
        <p class="text-h4">{{ timerDisplay }}</p>
        <p class="text-body-2">Sleep Timer Running</p>
      </v-card-text>
    </v-card>

    <!-- Location -->
    <v-select
      v-model="formData.location"
      :items="locationOptions"
      label="Location"
      variant="outlined"
      density="compact"
      class="mb-4"
    />

    <!-- Quality (only if ending sleep or manual entry with end time) -->
    <v-card v-if="useTimer || showEndTime" variant="outlined" class="mb-4">
      <v-card-text>
        <p class="text-body-2 mb-2">Sleep Quality</p>
        <v-rating
          v-model="formData.quality"
          color="yellow-darken-2"
          hover
          size="large"
        />
      </v-card-text>
    </v-card>

    <!-- Notes -->
    <v-textarea
      v-model="formData.notes"
      label="Notes (optional)"
      variant="outlined"
      rows="2"
      density="compact"
      class="mb-4"
    />

    <!-- Error display -->
    <v-alert
      v-if="formError"
      type="error"
      variant="tonal"
      class="mb-4"
      closable
      @click:close="formError = null"
    >
      {{ formError }}
    </v-alert>

    <!-- Actions -->
    <div class="d-flex gap-2">
      <v-btn
        v-if="!useTimer && hasTimer"
        variant="outlined"
        color="sleep"
        @click="startTimer"
        :loading="loading"
        block
      >
        <v-icon start>mdi-timer</v-icon>
        Start Sleep Timer
      </v-btn>
      
      <v-btn
        v-else-if="useTimer"
        color="success"
        @click="stopTimer"
        :loading="loading"
        block
      >
        <v-icon start>mdi-stop</v-icon>
        End Sleep
      </v-btn>
      
      <v-btn
        v-else
        type="submit"
        color="primary"
        :loading="loading"
        block
      >
        Save
      </v-btn>
    </div>
  </v-form>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useTimerStore } from '@/stores/timer'
import { useActivityStore } from '@/stores/activity'
import { combineDateAndTime, getCurrentDate, getCurrentTime } from '@/utils/datetime'

const props = defineProps({
  hasTimer: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['success', 'cancel'])

const timerStore = useTimerStore()
const activityStore = useActivityStore()

// Form state
const form = ref(null)
const loading = ref(false)
const formError = ref(null)
const useTimer = ref(false)
const timerInterval = ref(null)
const showEndTime = ref(false)
const formData = ref({
  startDate: getCurrentDate(),
  startTime: getCurrentTime(),
  endDate: getCurrentDate(),
  endTime: null,
  location: 'crib',
  quality: 3,
  notes: ''
})

// Location options
const locationOptions = [
  { title: 'Crib', value: 'crib' },
  { title: 'Bassinet', value: 'bassinet' },
  { title: 'Car Seat', value: 'car_seat' },
  { title: 'Stroller', value: 'stroller' },
  { title: 'Parent Bed', value: 'parent_bed' },
  { title: 'Other', value: 'other' }
]

// Timer display with persistent calculation - shows hours for longer sleeps
const timerDisplay = computed(() => {
  const timer = timerStore.activeTimers.sleep
  if (!timer) return '00:00'
  
  const elapsed = timerStore.getTimerDuration('sleep')
  const hours = Math.floor(elapsed / 3600)
  const minutes = Math.floor((elapsed % 3600) / 60)
  
  if (hours > 0) {
    return `${hours}h ${minutes}m`
  }
  return `${minutes} min`
})

// Watch for active timer changes (handles persistence restore)
watch(() => timerStore.activeTimers.sleep, (newTimer) => {
  if (newTimer) {
    useTimer.value = true
    startTimerDisplay()
  } else {
    useTimer.value = false
    if (timerInterval.value) {
      clearInterval(timerInterval.value)
    }
  }
}, { immediate: true })

onMounted(() => {
  // Timer state is handled by the watcher
})

onUnmounted(() => {
  if (timerInterval.value) {
    clearInterval(timerInterval.value)
  }
})

// Start timer display update
function startTimerDisplay() {
  if (timerInterval.value) {
    clearInterval(timerInterval.value)
  }
  timerInterval.value = setInterval(() => {
    // Force reactivity update - the computed will recalculate
  }, 1000)
}

// Enable end time input
function enableEndTime() {
  showEndTime.value = true
  // Set end time to current time as default
  formData.value.endTime = getCurrentTime()
}

// Clear end time
function clearEndTime() {
  showEndTime.value = false
  formData.value.endTime = null
}

// Start timer
async function startTimer() {
  loading.value = true
  formError.value = null
  
  const result = await activityStore.startTimer('sleep', {
    sleep_data: {
      location: formData.value.location
    },
    notes: formData.value.notes
  })
  
  if (result.success) {
    timerStore.startTimer('sleep', result.data.id)
  } else {
    formError.value = result.error || 'Failed to start timer'
  }
  
  loading.value = false
}

// Stop timer
async function stopTimer() {
  loading.value = true
  formError.value = null
  
  const timer = timerStore.activeTimers.sleep
  if (!timer || !timer.activityId) {
    formError.value = 'No active timer found'
    loading.value = false
    return
  }
  
  const data = {}
  if (formData.value.quality) {
    data.quality = formData.value.quality
  }
  if (formData.value.notes) {
    data.notes = formData.value.notes
  }
  
  const result = await activityStore.stopTimer(timer.activityId, data)
  
  if (result.success) {
    timerStore.stopTimer('sleep')
    emit('success', result.data)
  } else {
    formError.value = result.error || 'Failed to stop timer'
  }
  
  loading.value = false
}

// Submit form
async function handleSubmit() {
  loading.value = true
  formError.value = null
  
  const startDateTime = combineDateAndTime(formData.value.startDate, formData.value.startTime)
  
  const activityData = {
    type: 'sleep',
    start_time: startDateTime,
    notes: formData.value.notes,
    sleep_data: {
      location: formData.value.location
    }
  }
  
  if (showEndTime.value && formData.value.endTime) {
    const endDateTime = combineDateAndTime(formData.value.endDate, formData.value.endTime)
    
    // Validate end time is after start time
    if (endDateTime <= startDateTime) {
      formError.value = 'End time must be after start time'
      loading.value = false
      return
    }
    
    activityData.end_time = endDateTime
    if (formData.value.quality) {
      activityData.sleep_data.quality = formData.value.quality
    }
  }
  
  const result = await activityStore.createActivity(activityData)
  
  if (result.success) {
    emit('success', result.data)
  } else {
    formError.value = result.error || 'Failed to save activity'
  }
  
  loading.value = false
}
</script>