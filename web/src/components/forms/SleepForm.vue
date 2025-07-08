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
              :rules="[rules.required]"
            />
          </v-col>
          <v-col cols="5">
            <v-text-field
              v-model="formData.startTime"
              label="Start Time"
              type="time"
              variant="outlined"
              density="compact"
              :rules="[rules.required]"
            />
          </v-col>
        </v-row>

        <!-- End date and time (optional) -->
        <v-expand-transition>
          <v-row v-if="showEndTime" class="mb-3">
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
        </v-expand-transition>

        <!-- Toggle for end time -->
        <div class="d-flex justify-space-between align-center mb-3">
          <v-btn
            v-if="!showEndTime"
            variant="outlined"
            size="small"
            @click="enableEndTime"
          >
            <v-icon start>mdi-clock-plus</v-icon>
            Add End Time
          </v-btn>
          
          <v-btn
            v-else
            variant="text"
            size="small"
            color="error"
            @click="clearEndTime"
          >
            <v-icon start>mdi-close</v-icon>
            Remove End Time
          </v-btn>
        </div>
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
    <v-expand-transition>
      <v-card v-if="useTimer || showEndTime" variant="outlined" class="mb-4">
        <v-card-text>
          <p class="text-body-2 mb-2">Sleep Quality</p>
          <v-rating
            v-model="formData.quality"
            color="yellow-darken-2"
            hover
            size="large"
          />
          <div class="text-caption text-center mt-1">
            {{ qualityLabels[formData.quality - 1] || 'Not rated' }}
          </div>
        </v-card-text>
      </v-card>
    </v-expand-transition>

    <!-- Notes -->
    <v-textarea
      v-model="formData.notes"
      label="Notes (optional)"
      variant="outlined"
      rows="2"
      density="compact"
      :rules="[rules.maxLength(1000)]"
      class="mb-4"
    />

    <!-- Standardized error display -->
    <v-alert
      v-if="formError"
      type="error"
      variant="tonal"
      class="mb-4"
      closable
      @click:close="clearFormError"
    >
      <div v-if="formError.title" class="font-weight-medium mb-1">{{ formError.title }}</div>
      <div>{{ formError.message || formError }}</div>
    </v-alert>

    <!-- Actions -->
    <div class="d-flex gap-2">
      <v-btn
        v-if="!useTimer && hasTimer"
        variant="outlined"
        color="sleep"
        @click="startTimer"
        :loading="loading"
        :disabled="loading"
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
        :disabled="loading"
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
        :disabled="loading"
        block
      >
        Save
      </v-btn>
    </div>
  </v-form>
</template>

<script setup>
import { ref, computed, onUnmounted, watch } from 'vue'
import { useTimerStore } from '@/stores/timer'
import { useActivityStore } from '@/stores/activity'
import { useErrorHandling } from '@/composables/useErrorHandling'
import { combineDateAndTime, getCurrentDate, getCurrentTime } from '@/utils/datetime'
import { validationRules, validateDateTime } from '@/utils/validation'

const props = defineProps({
  hasTimer: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['success', 'cancel'])

// Stores
const timerStore = useTimerStore()
const activityStore = useActivityStore()

// Error handling
const { error: formError, loading, handleError, clearError: clearFormError, withErrorHandling } = useErrorHandling()

// Form state
const form = ref(null)
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

// Validation rules
const rules = {
  required: validationRules.required,
  maxLength: validationRules.maxLength
}

// Location options
const locationOptions = [
  { title: 'Crib', value: 'crib' },
  { title: 'Bassinet', value: 'bassinet' },
  { title: 'Car Seat', value: 'car_seat' },
  { title: 'Stroller', value: 'stroller' },
  { title: 'Parent Bed', value: 'parent_bed' },
  { title: 'Other', value: 'other' }
]

// Quality labels
const qualityLabels = [
  'Poor',
  'Fair', 
  'Good',
  'Very Good',
  'Excellent'
]

// Timer display with persistent calculation - shows hours for longer sleeps
const timerDisplay = computed(() => {
  const timer = timerStore.getActiveTimer('sleep')
  if (!timer) return '00:00'
  return timerStore.getFormattedDuration('sleep')
})

// Watch for active timer changes
watch(() => timerStore.hasActiveTimer('sleep'), (hasTimer) => {
  useTimer.value = hasTimer
  if (hasTimer) {
    startTimerDisplay()
  } else {
    stopTimerDisplay()
  }
}, { immediate: true })

onUnmounted(() => {
  stopTimerDisplay()
})

// Timer display management
function startTimerDisplay() {
  if (timerInterval.value) {
    clearInterval(timerInterval.value)
  }
  timerInterval.value = setInterval(() => {
    // Force reactivity update
  }, 1000)
}

function stopTimerDisplay() {
  if (timerInterval.value) {
    clearInterval(timerInterval.value)
    timerInterval.value = null
  }
}

// Enable end time input
function enableEndTime() {
  showEndTime.value = true
  // Set end time to current time as default
  formData.value.endTime = getCurrentTime()
  // If it's a different day, update end date
  const now = new Date()
  if (now.toDateString() !== new Date(formData.value.startDate).toDateString()) {
    formData.value.endDate = getCurrentDate()
  }
}

// Clear end time
function clearEndTime() {
  showEndTime.value = false
  formData.value.endTime = null
}

// Start timer
async function startTimer() {
  const result = await withErrorHandling(async () => {
    const response = await activityStore.startTimer('sleep', {
      sleep_data: {
        location: formData.value.location
      },
      notes: formData.value.notes
    })
    
    if (!response.success) {
      throw new Error(response.error)
    }
    
    const success = timerStore.startTimer('sleep', response.data.id)
    if (!success) {
      throw new Error('Failed to start local timer')
    }
    
    return response.data
  })
  
  if (!result.success) {
    handleError({
      title: 'Timer Start Failed',
      message: result.error
    })
  }
}

// Stop timer
async function stopTimer() {
  const timer = timerStore.getActiveTimer('sleep')
  if (!timer?.activityId) {
    handleError({
      title: 'Timer Error',
      message: 'No active timer found'
    })
    return
  }
  
  const result = await withErrorHandling(async () => {
    const stopData = {}
    if (formData.value.quality) {
      stopData.quality = formData.value.quality
    }
    if (formData.value.notes) {
      stopData.notes = formData.value.notes
    }
    
    const response = await activityStore.stopTimer(timer.activityId, stopData)
    if (!response.success) {
      throw new Error(response.error)
    }
    
    timerStore.stopTimer('sleep')
    return response.data
  })
  
  if (result.success) {
    emit('success', result.data)
  }
}

// Submit form
async function handleSubmit() {
  // Validate form
  const { valid } = await form.value.validate()
  if (!valid) return
  
  // Validate start date/time
  const startDateTimeError = validateDateTime(formData.value.startDate, formData.value.startTime)
  if (startDateTimeError) {
    handleError({
      title: 'Invalid Start Time',
      message: startDateTimeError
    })
    return
  }
  
  // Validate end date/time if provided
  if (showEndTime.value && formData.value.endTime) {
    const endDateTimeError = validateDateTime(formData.value.endDate, formData.value.endTime)
    if (endDateTimeError) {
      handleError({
        title: 'Invalid End Time',
        message: endDateTimeError
      })
      return
    }
    
    // Check that end is after start
    const startDateTime = combineDateAndTime(formData.value.startDate, formData.value.startTime)
    const endDateTime = combineDateAndTime(formData.value.endDate, formData.value.endTime)
    
    if (endDateTime <= startDateTime) {
      handleError({
        title: 'Invalid Time Range',
        message: 'End time must be after start time'
      })
      return
    }
  }
  
  const result = await withErrorHandling(async () => {
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
      activityData.end_time = endDateTime
      
      if (formData.value.quality) {
        activityData.sleep_data.quality = formData.value.quality
      }
    }
    
    const response = await activityStore.createActivity(activityData)
    if (!response.success) {
      throw new Error(response.error)
    }
    
    return response.data
  })
  
  if (result.success) {
    emit('success', result.data)
  }
}
</script>