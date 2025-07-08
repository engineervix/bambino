<template>
  <v-form ref="form" @submit.prevent="handleSubmit">
    <!-- Breast selection -->
    <v-btn-toggle
      v-model="formData.breast"
      mandatory
      color="pump"
      class="mb-4"
      block
    >
      <v-btn value="left" block>
        <v-icon start>mdi-chevron-left</v-icon>
        Left
      </v-btn>
      <v-btn value="right" block>
        <v-icon start>mdi-chevron-right</v-icon>
        Right
      </v-btn>
      <v-btn value="both" block>
        <v-icon start>mdi-chevron-double-left</v-icon>
        Both
      </v-btn>
    </v-btn-toggle>

    <!-- Timer or manual entry -->
    <v-card v-if="!useTimer" variant="outlined" class="mb-4">
      <v-card-text>
        <!-- Date and time inputs -->
        <v-row class="mb-3">
          <v-col cols="7">
            <v-text-field
              v-model="formData.date"
              label="Date"
              type="date"
              variant="outlined"
              density="compact"
              :rules="[rules.required]"
            />
          </v-col>
          <v-col cols="5">
            <v-text-field
              v-model="formData.time"
              label="Time"
              type="time"
              variant="outlined"
              density="compact"
              :rules="[rules.required]"
            />
          </v-col>
        </v-row>

        <!-- Amount -->
        <v-text-field
          v-model.number="formData.amount_ml"
          label="Amount (ml)"
          type="number"
          variant="outlined"
          density="compact"
          append-inner-icon="mdi-baby-bottle"
          :rules="[rules.positiveNumber]"
          class="mb-3"
        />

        <!-- Duration -->
        <v-text-field
          v-model.number="formData.duration_minutes"
          label="Duration (minutes)"
          type="number"
          variant="outlined"
          density="compact"
          append-inner-icon="mdi-timer"
          :rules="[rules.positiveInteger]"
        />
      </v-card-text>
    </v-card>

    <!-- Timer display -->
    <v-card v-else variant="tonal" color="pump" class="mb-4">
      <v-card-text class="text-center">
        <v-icon size="48" class="mb-2">mdi-timer</v-icon>
        <p class="text-h4">{{ timerDisplay }}</p>
        <p class="text-body-2">Pump Timer Running</p>
        
        <!-- Amount input while timer running -->
        <v-text-field
          v-model.number="formData.amount_ml"
          label="Amount (ml)"
          type="number"
          variant="outlined"
          density="compact"
          append-inner-icon="mdi-baby-bottle"
          :rules="[rules.positiveNumber]"
          class="mt-3"
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
        color="pump"
        @click="startTimer"
        :loading="loading"
        :disabled="loading"
        block
      >
        <v-icon start>mdi-timer</v-icon>
        Start Timer
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
        Stop & Save
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
const formData = ref({
  breast: 'both',
  date: getCurrentDate(),
  time: getCurrentTime(),
  amount_ml: null,
  duration_minutes: null,
  notes: ''
})

// Validation rules
const rules = {
  required: validationRules.required,
  positiveNumber: validationRules.positiveNumber,
  positiveInteger: validationRules.positiveInteger,
  maxLength: validationRules.maxLength
}

// Timer display with persistent calculation
const timerDisplay = computed(() => {
  const timer = timerStore.getActiveTimer('pump')
  if (!timer) return '00:00'
  return timerStore.getFormattedDuration('pump')
})

// Watch for active timer changes
watch(() => timerStore.hasActiveTimer('pump'), (hasTimer) => {
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

// Start timer
async function startTimer() {
  const result = await withErrorHandling(async () => {
    const response = await activityStore.startTimer('pump', {
      pump_data: {
        breast: formData.value.breast
      },
      notes: formData.value.notes
    })
    
    if (!response.success) {
      throw new Error(response.error)
    }
    
    const success = timerStore.startTimer('pump', response.data.id)
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
  const timer = timerStore.getActiveTimer('pump')
  if (!timer?.activityId) {
    handleError({
      title: 'Timer Error',
      message: 'No active timer found'
    })
    return
  }
  
  const result = await withErrorHandling(async () => {
    const stopData = {}
    if (formData.value.amount_ml) {
      stopData.amount_ml = formData.value.amount_ml
    }
    if (formData.value.notes) {
      stopData.notes = formData.value.notes
    }
    
    const response = await activityStore.stopTimer(timer.activityId, stopData)
    if (!response.success) {
      throw new Error(response.error)
    }
    
    timerStore.stopTimer('pump')
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
  
  // Validate date/time
  const dateTimeError = validateDateTime(formData.value.date, formData.value.time)
  if (dateTimeError) {
    handleError({
      title: 'Invalid Date/Time',
      message: dateTimeError
    })
    return
  }
  
  const result = await withErrorHandling(async () => {
    const activityDateTime = combineDateAndTime(formData.value.date, formData.value.time)
    
    const activityData = {
      type: 'pump',
      start_time: activityDateTime,
      notes: formData.value.notes,
      pump_data: {
        breast: formData.value.breast
      }
    }
    
    if (formData.value.amount_ml) {
      activityData.pump_data.amount_ml = formData.value.amount_ml
    }
    
    if (formData.value.duration_minutes) {
      activityData.pump_data.duration_minutes = formData.value.duration_minutes
      // Calculate end time
      activityData.end_time = new Date(activityDateTime.getTime() + formData.value.duration_minutes * 60000)
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