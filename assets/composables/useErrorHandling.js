import { ref } from 'vue'

/**
 * Composable for standardized error handling across forms
 */
export function useErrorHandling() {
  const error = ref(null)
  const loading = ref(false)

  /**
   * Handles API errors and extracts user-friendly messages
   */
  function handleError(err) {
    console.error('Error occurred:', err)
    
    // Clear any previous errors
    error.value = null
    
    if (typeof err === 'string') {
      error.value = err
      return err
    }
    
    // Handle axios/HTTP errors
    if (err?.response) {
      const status = err.response.status
      const data = err.response.data
      
      // Extract message from various response formats
      let message = data?.message || data?.error || data?.detail
      
      // Handle specific HTTP status codes
      if (!message) {
        switch (status) {
          case 400:
            message = 'Invalid request. Please check your input.'
            break
          case 401:
            message = 'Authentication required. Please log in.'
            break
          case 403:
            message = 'Permission denied.'
            break
          case 404:
            message = 'Resource not found.'
            break
          case 422:
            message = 'Validation failed. Please check your input.'
            break
          case 429:
            message = 'Too many requests. Please try again later.'
            break
          case 500:
            message = 'Server error. Please try again.'
            break
          default:
            message = `Request failed (${status})`
        }
      }
      
      error.value = message
      return message
    }
    
    // Handle network errors
    if (err?.request) {
      const message = 'Network error. Please check your connection.'
      error.value = message
      return message
    }
    
    // Handle timeout errors
    if (err?.code === 'ECONNABORTED') {
      const message = 'Request timed out. Please try again.'
      error.value = message
      return message
    }
    
    // Handle validation errors (from our validation utils)
    if (err?.type === 'validation') {
      error.value = err.message
      return err.message
    }
    
    // Fallback for unknown errors
    const message = err?.message || 'An unexpected error occurred'
    error.value = message
    return message
  }

  /**
   * Clears the current error
   */
  function clearError() {
    error.value = null
  }

  /**
   * Sets loading state
   */
  function setLoading(isLoading) {
    loading.value = isLoading
  }

  /**
   * Wrapper for async operations with standardized error handling
   */
  async function withErrorHandling(asyncOperation, options = {}) {
    const { 
      loadingState = true,
      clearPreviousError = true,
      onSuccess,
      onError 
    } = options

    if (clearPreviousError) {
      clearError()
    }

    if (loadingState) {
      setLoading(true)
    }

    try {
      const result = await asyncOperation()
      
      if (onSuccess) {
        onSuccess(result)
      }
      
      return { success: true, data: result, error: null }
    } catch (err) {
      const errorMessage = handleError(err)
      
      if (onError) {
        onError(errorMessage, err)
      }
      
      return { success: false, data: null, error: errorMessage }
    } finally {
      if (loadingState) {
        setLoading(false)
      }
    }
  }

  /**
   * Creates a validation error
   */
  function createValidationError(message) {
    const validationError = new Error(message)
    validationError.type = 'validation'
    return validationError
  }

  return {
    // State
    error,
    loading,
    
    // Methods
    handleError,
    clearError,
    setLoading,
    withErrorHandling,
    createValidationError
  }
}

/**
 * Specific composable for form handling
 */
export function useFormHandling() {
  const errorHandling = useErrorHandling()

  /**
   * Standardized form submission wrapper
   */
  async function handleFormSubmit(formRef, submitFunction, options = {}) {
    const { 
      validateForm = true,
      successMessage = 'Saved successfully!',
      onSuccess,
      onError 
    } = options

    // Validate form if reference provided
    if (validateForm && formRef?.value) {
      const { valid } = await formRef.value.validate()
      if (!valid) {
        return { success: false, error: 'Please fix the errors above' }
      }
    }

    // Execute submission with error handling
    return await errorHandling.withErrorHandling(
      submitFunction,
      {
        onSuccess: (result) => {
          if (onSuccess) {
            onSuccess(result, successMessage)
          }
        },
        onError
      }
    )
  }

  return {
    ...errorHandling,
    handleFormSubmit
  }
}