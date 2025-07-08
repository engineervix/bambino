<template>
  <v-form ref="form" @submit.prevent="handleSubmit">
    <!-- Time input -->
    <v-text-field
      v-model="formData.time"
      label="Time"
      type="time"
      variant="outlined"
      density="compact"
      class="mb-4"
    />

    <!-- Diaper type -->
    <v-card variant="outlined" class="mb-4">
      <v-card-text>
        <p class="text-body-2 mb-2">Diaper Type</p>
        <v-row>
          <v-col cols="6">
            <v-checkbox
              v-model="formData.wet"
              label="Wet"
              color="primary"
              hide-details
            />
          </v-col>
          <v-col cols="6">
            <v-checkbox
              v-model="formData.dirty"
              label="Dirty"
              color="primary"
              hide-details
            />
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- Additional details (shown if dirty) -->
    <v-expand-transition>
      <div v-if="formData.dirty">
        <v-select
          v-model="formData.color"
          :items="colorOptions"
          label="Color (optional)"
          variant="outlined"
          density="compact"
          clearable
          class="mb-3"
        />
        
        <v-select
          v-model="formData.consistency"
          :items="consistencyOptions"
          label="Consistency (optional)"
          variant="outlined"
          density="compact"
          clearable
          class="mb-4"
        />
      </div>
    </v-expand-transition>

    <!-- Notes -->
    <v-textarea
      v-model="formData.notes"
      label="Notes (optional)"
      variant="outlined"
      rows="2"
      density="compact"
      class="mb-4"
    />

    <!-- Submit button -->
    <v-btn
      type="submit"
      color="primary"
      :loading="loading"
      :disabled="!formData.wet && !formData.dirty"
      block
    >
      Save
    </v-btn>
  </v-form>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useActivityStore } from '@/stores/activity'

const emit = defineEmits(['success', 'cancel'])

const activityStore = useActivityStore()

// Form state
const form = ref(null)
const loading = ref(false)
const formData = ref({
  time: new Date().toTimeString().slice(0, 5),
  wet: false,
  dirty: false,
  color: null,
  consistency: null,
  notes: ''
})

// Options
const colorOptions = [
  { title: 'Yellow', value: 'yellow' },
  { title: 'Green', value: 'green' },
  { title: 'Brown', value: 'brown' },
  { title: 'Black', value: 'black' },
  { title: 'Red', value: 'red' },
  { title: 'White', value: 'white' }
]

const consistencyOptions = [
  { title: 'Liquid', value: 'liquid' },
  { title: 'Soft', value: 'soft' },
  { title: 'Normal', value: 'normal' },
  { title: 'Hard', value: 'hard' }
]

// Clear color/consistency if not dirty
watch(() => formData.value.dirty, (isDirty) => {
  if (!isDirty) {
    formData.value.color = null
    formData.value.consistency = null
  }
})

// Submit form
async function handleSubmit() {
  if (!formData.value.wet && !formData.value.dirty) {
    return
  }

  loading.value = true
  
  const activityData = {
    type: 'diaper',
    start_time: new Date(`${new Date().toDateString()} ${formData.value.time}`),
    notes: formData.value.notes,
    diaper_data: {
      wet: formData.value.wet,
      dirty: formData.value.dirty
    }
  }
  
  if (formData.value.color) {
    activityData.diaper_data.color = formData.value.color
  }
  
  if (formData.value.consistency) {
    activityData.diaper_data.consistency = formData.value.consistency
  }
  
  const result = await activityStore.createActivity(activityData)
  
  if (result.success) {
    emit('success', result.data)
  }
  
  loading.value = false
}
</script>