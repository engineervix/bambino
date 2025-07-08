<template>
  <v-form ref="formRef" @submit.prevent="handleSubmit">
    <slot 
      :formData="formData"
      :loading="loading"
      :error="error"
      :clearError="clearError"
      :updateFormData="updateFormData"
    />
    
    <!-- Standardized error display -->
    <v-alert
      v-if="error"
      type="error"
      variant="tonal"
      class="mb-4"
      closable
      @click:close="clearError"
    >
      <template v-if="typeof error === 'object' && error.title">
        <div class="font-weight-medium">{{ error.title }}</div>
        <div>{{ error.message }}</div>
      </template>
      <template v-else>
        {{ error }}
      </template>
    </v-alert>
    
    <!-- Submit button slot -->
    <slot 
      name="submit"
      :loading="loading"
      :disabled="submitDisabled"
      :handleSubmit="handleSubmit"
    >
      <v-btn
        type="submit"
        :color="submitColor"
        :loading="loading"
        :disabled="submitDisabled"
        block
      >
        <v-icon v-if="submitIcon" start>{{ submitIcon }}</v-icon>
        {{ submitText }}
      </v-btn>
    </slot>
  </v-form>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useFormHandling } from '@/composables/useErrorHandling'

const props = defineProps({
  // Form configuration
  initialData: {
    type: Object,
    default: () => ({})
  },
  
  // Submit configuration
  submitText: {
    type: String,
    default: 'Save'
  },
  
  submitColor: {
    type: String,
    default: 'primary'
  },
  
  submitIcon: {
    type: String,
    default: null
  },
  
  submitDisabled: {
    type: Boolean,
    default: false
  },
  
  // Validation configuration
  validateForm: {
    type: Boolean,
    default: true
  },
  
  // Success message
  successMessage: {
    type: String,
    default: 'Saved successfully!'
  }
})

const emit = defineEmits(['submit', 'success', 'error'])

// Form handling composable
const { error, loading, clearError, handleFormSubmit } = useFormHandling()

// Form state
const formRef = ref(null)
const formData = ref({ ...props.initialData })

// Update form data helper
function updateFormData(key, value) {
  if (typeof key === 'object') {
    // Update multiple fields
    Object.assign(formData.value, key)
  } else {
    // Update single field
    formData.value[key] = value
  }
}

// Handle form submission
async function handleSubmit() {
  const result = await handleFormSubmit(
    formRef,
    () => {
      // Emit submit event and return the promise
      return new Promise((resolve, reject) => {
        emit('submit', {
          data: formData.value,
          resolve,
          reject
        })
      })
    },
    {
      validateForm: props.validateForm,
      successMessage: props.successMessage,
      onSuccess: (result, message) => {
        emit('success', { data: result, message })
      },
      onError: (errorMessage, originalError) => {
        emit('error', { message: errorMessage, error: originalError })
      }
    }
  )
  
  return result
}

// Reset form
function resetForm() {
  formData.value = { ...props.initialData }
  clearError()
  if (formRef.value) {
    formRef.value.resetValidation()
  }
}

// Validate form manually
async function validateForm() {
  if (!formRef.value) return { valid: true }
  return await formRef.value.validate()
}

// Expose methods for parent components
defineExpose({
  formData,
  loading,
  error,
  resetForm,
  validateForm,
  clearError,
  updateFormData
})
</script>