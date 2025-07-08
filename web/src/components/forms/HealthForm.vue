<template>
  <v-form ref="form" @submit.prevent="handleSubmit">
    <!-- Record type -->
    <v-select
      v-model="formData.record_type"
      :items="recordTypes"
      label="Type"
      variant="outlined"
      density="compact"
      class="mb-4"
    />

    <!-- Date and time -->
    <v-row class="mb-4">
      <v-col cols="7">
        <v-text-field
          v-model="formData.date"
          label="Date"
          type="date"
          variant="outlined"
          density="compact"
        />
      </v-col>
      <v-col cols="5">
        <v-text-field
          v-model="formData.time"
          label="Time"
          type="time"
          variant="outlined"
          density="compact"
        />
      </v-col>
    </v-row>

    <!-- Provider -->
    <v-text-field
      v-model="formData.provider"
      label="Provider/Location"
      variant="outlined"
      density="compact"
      class="mb-4"
    />

    <!-- Vaccine-specific fields -->
    <v-expand-transition>
      <div v-if="formData.record_type === 'vaccine'">
        <v-text-field
          v-model="formData.vaccine_name"
          label="Vaccine Name"
          variant="outlined"
          density="compact"
          class="mb-4"
          :rules="[v => !!v || 'Vaccine name is required']"
        />
      </div>
    </v-expand-transition>

    <!-- Illness-specific fields -->
    <v-expand-transition>
      <div v-if="formData.record_type === 'illness'">
        <v-textarea
          v-model="formData.symptoms"
          label="Symptoms"
          variant="outlined"
          rows="2"
          density="compact"
          class="mb-3"
        />
        
        <v-textarea
          v-model="formData.treatment"
          label="Treatment"
          variant="outlined"
          rows="2"
          density="compact"
          class="mb-4"
        />
      </div>
    </v-expand-transition>

    <!-- Notes -->
    <v-textarea
      v-model="formData.notes"
      label="Notes"
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
      block
    >
      Save Health Record
    </v-btn>
  </v-form>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useActivityStore } from '@/stores/activity'
import { combineDateAndTime, getCurrentDate, getCurrentTime } from '@/utils/datetime'

const emit = defineEmits(['success', 'cancel'])

const activityStore = useActivityStore()

// Form state
const form = ref(null)
const loading = ref(false)
const formError = ref(null)
const formData = ref({
  record_type: 'checkup',
  date: getCurrentDate(),
  time: getCurrentTime(),
  provider: '',
  vaccine_name: '',
  symptoms: '',
  treatment: '',
  notes: ''
})

// Record type options
const recordTypes = [
  { title: 'Checkup', value: 'checkup' },
  { title: 'Vaccine', value: 'vaccine' },
  { title: 'Illness', value: 'illness' }
]

// Clear type-specific fields when record type changes
watch(() => formData.value.record_type, (newType) => {
  if (newType !== 'vaccine') {
    formData.value.vaccine_name = ''
  }
  if (newType !== 'illness') {
    formData.value.symptoms = ''
    formData.value.treatment = ''
  }
})

// Submit form
async function handleSubmit() {
  const { valid } = await form.value.validate()
  if (!valid) return

  loading.value = true
  formError.value = null
  
  const activityDateTime = combineDateAndTime(formData.value.date, formData.value.time)
  
  const activityData = {
    type: 'health',
    start_time: activityDateTime,
    notes: formData.value.notes,
    health_data: {
      record_type: formData.value.record_type,
      provider: formData.value.provider
    }
  }
  
  if (formData.value.record_type === 'vaccine' && formData.value.vaccine_name) {
    activityData.health_data.vaccine_name = formData.value.vaccine_name
  }
  
  if (formData.value.record_type === 'illness') {
    if (formData.value.symptoms) {
      activityData.health_data.symptoms = formData.value.symptoms
    }
    if (formData.value.treatment) {
      activityData.health_data.treatment = formData.value.treatment
    }
  }
  
  const result = await activityStore.createActivity(activityData)
  
  if (result.success) {
    emit('success', result.data)
  } else {
    formError.value = result.error || 'Failed to save health record'
  }
  
  loading.value = false
}
</script>