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
        <!-- Time input -->
        <v-text-field
          v-model="formData.time"
          label="Time"
          type="time"
          variant="outlined"
          density="compact"
          class="mb-3"
        />

        <!-- Amount -->
        <v-text-field
          v-model.number="formData.amount_ml"
          label="Amount (ml)"
          type="number"
          variant="outlined"
          density="compact"
          append-inner-icon="mdi-baby-bottle"
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
      class="mb-4"
    />

    <!-- Actions -->
    <div class="d-flex gap-2">
      <v-btn
        v-if="!useTimer && hasTimer"
        variant="outlined"
        color="pump"
        @click="startTimer"
        :loading="loading"
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
        block
      >
        Save
      </v-btn>
    </div>
  </v-form>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useTimerStore } from '@/stores/timer'
import { useActivityStore } from '@/stores/activity'

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
const useTimer = ref(false)
const timerInterval = ref(null)
const formData = ref({
  breast: 'both',
  time: new Date().toTimeString().slice(0, 5),
  amount_ml: null,
  duration_minutes: null,
  notes: ''
})

// Timer display
const timerDisplay = computed(() => {
  const timer = timerStore.activeTimers.pump
  if (!timer) return '00:00'
  
  const elapsed = Math.floor((Date.now() - new Date(timer.startTime).getTime()) / 1000)
  const minutes = Math.floor(elapsed / 60)
  const seconds = elapsed % 60
  return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
})

// Check for active timer on mount
onMounted(() => {
  if (timerStore.activeTimers.pump) {
    useTimer.value = true
    startTimerDisplay()
  }
})

onUnmounted(() => {
  if (timerInterval.value) {
    clearInterval(timerInterval.value)
  }
})

// Start timer display update
function startTimerDisplay() {
  timerInterval.value = setInterval(() => {
    // Force reactivity update
  }, 1000)
}

// Start timer
async function startTimer() {
  loading.value = true
  
  const result = await activityStore.startTimer('pump', {
    pump_data: {
      breast: formData.value.breast
    },
    notes: formData.value.notes
  })
  
  if (result.success) {
    timerStore.startTimer('pump', result.data.id)
    useTimer.value = true
    startTimerDisplay()
  }
  
  loading.value = false
}

// Stop timer
async function stopTimer() {
  loading.value = true
  
  const timer = timerStore.activeTimers.pump
  if (!timer || !timer.activityId) {
    loading.value = false
    return
  }
  
  const result = await activityStore.stopTimer(timer.activityId, {
    amount_ml: formData.value.amount_ml,
    notes: formData.value.notes
  })
  
  if (result.success) {
    timerStore.stopTimer('pump')
    emit('success', result.data)
  }
  
  loading.value = false
}

// Submit form
async function handleSubmit() {
  loading.value = true
  
  const activityData = {
    type: 'pump',
    start_time: new Date(`${new Date().toDateString()} ${formData.value.time}`),
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
    const startTime = new Date(activityData.start_time)
    activityData.end_time = new Date(startTime.getTime() + formData.value.duration_minutes * 60000)
  }
  
  const result = await activityStore.createActivity(activityData)
  
  if (result.success) {
    emit('success', result.data)
  }
  
  loading.value = false
}
</script>