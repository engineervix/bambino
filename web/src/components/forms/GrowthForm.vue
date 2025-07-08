<template>
  <v-form ref="form" @submit.prevent="handleSubmit">
    <!-- Date input -->
    <v-text-field
      v-model="formData.date"
      label="Date"
      type="date"
      variant="outlined"
      density="compact"
      class="mb-4"
    />

    <!-- Measurements -->
    <v-card variant="outlined" class="mb-4">
      <v-card-text>
        <p class="text-body-2 mb-3">Measurements</p>
        
        <!-- Weight -->
        <v-text-field
          v-model.number="formData.weight_kg"
          label="Weight (kg)"
          type="number"
          step="0.01"
          variant="outlined"
          density="compact"
          append-inner-icon="mdi-scale"
          class="mb-3"
          :rules="[v => v > 0 || !v || 'Weight must be positive']"
        />
        
        <!-- Height -->
        <v-text-field
          v-model.number="formData.height_cm"
          label="Height (cm)"
          type="number"
          step="0.1"
          variant="outlined"
          density="compact"
          append-inner-icon="mdi-human-male-height"
          class="mb-3"
          :rules="[v => v > 0 || !v || 'Height must be positive']"
        />
        
        <!-- Head circumference -->
        <v-text-field
          v-model.number="formData.head_circumference_cm"
          label="Head Circumference (cm)"
          type="number"
          step="0.1"
          variant="outlined"
          density="compact"
          append-inner-icon="mdi-head"
          :rules="[v => v > 0 || !v || 'Head circumference must be positive']"
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

    <!-- Submit button -->
    <v-btn
      type="submit"
      color="primary"
      :loading="loading"
      :disabled="!hasAnyMeasurement"
      block
    >
      Save Measurements
    </v-btn>
  </v-form>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useActivityStore } from '@/stores/activity'
import { getCurrentDate } from '@/utils/datetime'

const emit = defineEmits(['success', 'cancel'])

const activityStore = useActivityStore()

// Form state
const form = ref(null)
const loading = ref(false)
const formError = ref(null)
const formData = ref({
  date: getCurrentDate(),
  weight_kg: null,
  height_cm: null,
  head_circumference_cm: null,
  notes: ''
})

// Check if at least one measurement is provided
const hasAnyMeasurement = computed(() => {
  return formData.value.weight_kg || 
         formData.value.height_cm || 
         formData.value.head_circumference_cm
})

// Submit form
async function handleSubmit() {
  if (!hasAnyMeasurement.value) {
    return
  }

  const { valid } = await form.value.validate()
  if (!valid) return

  loading.value = true
  formError.value = null
  
  // Growth measurements typically happen at a consistent time (like doctor visits)
  // so we'll set it to noon on the selected date
  const activityDateTime = new Date(`${formData.value.date}T12:00:00`)
  
  const activityData = {
    type: 'growth',
    start_time: activityDateTime,
    notes: formData.value.notes,
    growth_data: {}
  }
  
  if (formData.value.weight_kg) {
    activityData.growth_data.weight_kg = formData.value.weight_kg
  }
  
  if (formData.value.height_cm) {
    activityData.growth_data.height_cm = formData.value.height_cm
  }
  
  if (formData.value.head_circumference_cm) {
    activityData.growth_data.head_circumference_cm = formData.value.head_circumference_cm
  }
  
  const result = await activityStore.createActivity(activityData)
  
  if (result.success) {
    emit('success', result.data)
  } else {
    formError.value = result.error || 'Failed to save measurements'
  }
  
  loading.value = false
}
</script>