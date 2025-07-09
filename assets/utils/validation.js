/**
 * Common validation rules for forms
 */
export const validationRules = {
  required: (v) => !!v || "This field is required",

  positiveNumber: (v) => {
    if (!v) return true; // Allow empty for optional fields
    return v > 0 || "Must be a positive number";
  },

  positiveInteger: (v) => {
    if (!v) return true;
    return (Number.isInteger(Number(v)) && v > 0) || "Must be a positive whole number";
  },

  maxLength: (max) => (v) => {
    if (!v) return true;
    return v.length <= max || `Cannot exceed ${max} characters`;
  },

  email: (v) => {
    if (!v) return true;
    const pattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return pattern.test(v) || "Invalid email address";
  },

  minLength: (min) => (v) => {
    if (!v) return true;
    return v.length >= min || `Must be at least ${min} characters`;
  },

  range: (min, max) => (v) => {
    if (!v) return true;
    const num = Number(v);
    return (num >= min && num <= max) || `Must be between ${min} and ${max}`;
  },
};

/**
 * Validates date/time combinations
 */
export function validateDateTime(startDate, startTime, endDate, endTime) {
  if (!startDate || !startTime) {
    return "Start date and time are required";
  }

  if (endDate && endTime) {
    const start = new Date(`${startDate}T${startTime}`);
    const end = new Date(`${endDate}T${endTime}`);

    if (end <= start) {
      return "End time must be after start time";
    }
  }

  return null;
}

/**
 * Validates that at least one field in an object has a value
 */
export function validateAtLeastOne(obj, fieldNames, message = "At least one field is required") {
  const hasValue = fieldNames.some((field) => {
    const value = obj[field];
    return value !== null && value !== undefined && value !== "";
  });

  return hasValue ? null : message;
}

/**
 * Standard error handler for form submissions
 */
export function handleFormError(error) {
  if (typeof error === "string") {
    return error;
  }

  if (error?.response?.data?.message) {
    return error.response.data.message;
  }

  if (error?.message) {
    return error.message;
  }

  return "An unexpected error occurred";
}
