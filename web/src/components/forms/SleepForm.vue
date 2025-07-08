<template>
  <v-form ref="form" @submit.prevent="handleSubmit">
    <!-- Timer or manual entry -->
    <v-card v-if="!useTimer" variant="outlined" class="mb-4">
      <v-card-text>
        <!-- Start time -->
        <v-text-field
          v-model="formData.startTime"
          label="Start Time"
          type="time"
          variant="outlined"
          density="compact"
          class="mb-3"
        />

        <!-- End time (optional) -->
        <v-text-field
          v-model="formData.endTime"
          label="End Time (optional)"
          type="time"
          variant="outlined"
          density="compact"
          clearable
        />
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

    <!-- Quality (only if ending sleep) -->
    <v-card v-if="useTimer || formData.endTime" variant="outlined" class="mb-4">
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
  startTime: new Date().toTimeString().slice(0, 5),
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

// Timer display
const timerDisplay = computed(() => {
  const timer = timerStore.activeTimers.sleep
  if (!timer) return '00:00'
  
  const elapsed = Math.floor((Date.now() - new Date(timer.startTime).getTime()) / 1000)
  const hours = Math.floor(elapsed / 3600)
  const minutes = Math.floor((elapsed % 3600) / 60)
  
  if (hours > 0) {
    return `${hours}h ${minutes}m`
  }
  return `${minutes} min`
})

// Check for active timer on mount
onMounted(() => {
  if (timerStore.activeTimers.sleep) {
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
  
  const result = await activityStore.startTimer('sleep', {
    sleep_data: {
      location: formData.value.location
    },
    notes: formData.value.notes
  })
  
  if (result.success) {
    timerStore.startTimer('sleep', result.data.id)
    useTimer.value = true
    startTimerDisplay()
  }
  
  loading.value = false
}

// Stop timer
async function stopTimer() {
  loading.value = true
  
  const timer = timerStore.activeTimers.sleep
  if (!timer || !timer.activityId) {
    loading.value = false
    return
  }
  
  const result = await activityStore.stopTimer(timer.activityId, {
    quality: formData.value.quality,
    notes: formData.value.notes
  })
  
  if (result.success) {
    timerStore.stopTimer('sleep')
    emit('success', result.data)
  }
  
  loading.value = false
}

// Submit form
async function handleSubmit() {
  loading.value = true
  
  const activityData = {
    type: 'sleep',
    start_time: new Date(`${new Date().toDateString()} ${formData.value.startTime}`),
    notes: formData.value.notes,
    sleep_data: {
      location: formData.value.location
    }
  }
  
  if (formData.value.endTime) {
    activityData.end_time = new Date(`${new Date().toDateString()} ${formData.value.endTime}`)
    activityData.sleep_data.quality = formData.value.quality
  }
  
  const result = await activityStore.createActivity(activityData)
  
  if (result.success) {
    emit('success', result.data)
  }
  
  loading.value = false
}
</script>