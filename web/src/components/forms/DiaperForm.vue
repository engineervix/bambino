<template>
  <v-form ref="form" @submit.prevent="handleSubmit">
    <!-- Date and time inputs -->
    <v-row class="mb-4">
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
        
        <!-- Validation message for diaper type -->
        <v-alert
          v-if="!formData.wet && !formData.dirty && showDiaperTypeError"
          type="warning"
          variant="tonal"
          density="compact"
          class="mt-2"
        >
          Please select wet, dirty, or both
        </v-alert>
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

    <!-- Submit button -->
    <v-btn
      type="submit"
      color="primary"
      :loading="loading"
      :disabled="loading || (!formData.wet && !formData.dirty)"
      block
    >
      {{ editMode ? 'Update Diaper' : 'Save' }}
    </v-btn>
  </v-form>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useActivityStore } from '@/stores/activity'
import { useErrorHandling } from '@/composables/useErrorHandling'
import { combineDateAndTime, getCurrentDate, getCurrentTime, getDateString, getTimeString } from '@/utils/datetime'
import { validationRules, validateDateTime } from '@/utils/validation'

const props = defineProps({
  activity: {
    type: Object,
    default: null
  },
  editMode: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['success', 'cancel'])

// Stores
const activityStore = useActivityStore()

// Error handling
const { error: formError, loading, handleError, clearError: clearFormError, withErrorHandling } = useErrorHandling()

// Form state
const form = ref(null)
const showDiaperTypeError = ref(false)

// Initialize form data from props or defaults
const initializeFormData = () => {
  if (props.editMode && props.activity) {
    const activity = props.activity
    const startTime = new Date(activity.start_time)
    return {
      date: getDateString(startTime),
      time: getTimeString(startTime),
      wet: activity.diaper_data?.wet || false,
      dirty: activity.diaper_data?.dirty || false,
      color: activity.diaper_data?.color || null,
      consistency: activity.diaper_data?.consistency || null,
      notes: activity.notes || ''
    }
  }
  return {
    date: getCurrentDate(),
    time: getCurrentTime(),
    wet: false,
    dirty: false,
    color: null,
    consistency: null,
    notes: ''
  }
}

const formData = ref(initializeFormData())

// Validation rules
const rules = {
  required: validationRules.required,
  maxLength: validationRules.maxLength
}

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

// Hide diaper type error when user makes a selection
watch([() => formData.value.wet, () => formData.value.dirty], () => {
  if (formData.value.wet || formData.value.dirty) {
    showDiaperTypeError.value = false
  }
})

// Watch for prop changes to reinitialize form data
watch(() => props.activity, () => {
  if (props.editMode && props.activity) {
    formData.value = initializeFormData()
  }
}, { deep: true })

// Submit form
async function handleSubmit() {
  // Reset error indicators
  showDiaperTypeError.value = false
  
  // Validate form
  const { valid } = await form.value.validate()
  if (!valid) return
  
  // Custom validation for diaper type
  if (!formData.value.wet && !formData.value.dirty) {
    showDiaperTypeError.value = true
    handleError({
      title: 'Invalid Selection',
      message: 'Please select wet, dirty, or both'
    })
    return
  }
  
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
      type: 'diaper',
      start_time: activityDateTime,
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
    
    let response
    if (props.editMode && props.activity) {
      // Update existing activity
      response = await activityStore.updateActivity(props.activity.id, activityData)
    } else {
      // Create new activity
      response = await activityStore.createActivity(activityData)
    }
    
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