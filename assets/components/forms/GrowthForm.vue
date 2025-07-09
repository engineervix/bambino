<template>
  <v-form ref="form" @submit.prevent="handleSubmit">
    <!-- Date input -->
    <v-text-field
      v-model="formData.date"
      label="Date"
      type="date"
      variant="outlined"
      density="compact"
      :rules="[rules.required]"
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
          :rules="[rules.positiveNumber]"
          class="mb-3"
          placeholder="e.g., 7.5"
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
          :rules="[rules.positiveNumber]"
          class="mb-3"
          placeholder="e.g., 65.5"
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
          :rules="[rules.positiveNumber]"
          placeholder="e.g., 42.0"
        />

        <!-- At least one measurement required message -->
        <v-alert
          v-if="!hasAnyMeasurement && showMeasurementError"
          type="warning"
          variant="tonal"
          density="compact"
          class="mt-2"
        >
          Please provide at least one measurement
        </v-alert>
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
      placeholder="Doctor visit notes, growth observations, etc."
      class="mb-4"
    />

    <!-- Standardized error display -->
    <v-alert v-if="formError" type="error" variant="tonal" class="mb-4" closable @click:close="clearFormError">
      <div v-if="formError.title" class="font-weight-medium mb-1">{{ formError.title }}</div>
      <div>{{ formError.message || formError }}</div>
    </v-alert>

    <!-- Submit button -->
    <v-btn type="submit" color="primary" :loading="loading" :disabled="loading || !hasAnyMeasurement" block>
      <v-icon start>mdi-chart-line</v-icon>
      {{ editMode ? "Update Measurements" : "Save Measurements" }}
    </v-btn>
  </v-form>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useActivityStore } from "@/stores/activity";
import { useErrorHandling } from "@/composables/useErrorHandling";
import { getCurrentDate, getDateString } from "@/utils/datetime";
import { validationRules, validateAtLeastOne } from "@/utils/validation";

const props = defineProps({
  activity: {
    type: Object,
    default: null,
  },
  editMode: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["success", "cancel"]);

// Stores
const activityStore = useActivityStore();

// Error handling
const { error: formError, loading, handleError, clearError: clearFormError, withErrorHandling } = useErrorHandling();

// Form state
const form = ref(null);
const showMeasurementError = ref(false);

// Initialize form data from props or defaults
const initializeFormData = () => {
  if (props.editMode && props.activity) {
    const activity = props.activity;
    const startTime = new Date(activity.start_time);
    return {
      date: getDateString(startTime),
      weight_kg: activity.growth_data?.weight_kg || null,
      height_cm: activity.growth_data?.height_cm || null,
      head_circumference_cm: activity.growth_data?.head_circumference_cm || null,
      notes: activity.notes || "",
    };
  }
  return {
    date: getCurrentDate(),
    weight_kg: null,
    height_cm: null,
    head_circumference_cm: null,
    notes: "",
  };
};

const formData = ref(initializeFormData());

// Validation rules
const rules = {
  required: validationRules.required,
  positiveNumber: validationRules.positiveNumber,
  maxLength: validationRules.maxLength,
};

// Check if at least one measurement is provided
const hasAnyMeasurement = computed(() => {
  return formData.value.weight_kg || formData.value.height_cm || formData.value.head_circumference_cm;
});

// Hide measurement error when user provides a measurement
watch(
  () => hasAnyMeasurement.value,
  (hasValue) => {
    if (hasValue) {
      showMeasurementError.value = false;
    }
  },
);

// Watch for prop changes to reinitialize form data
watch(
  () => props.activity,
  () => {
    if (props.editMode && props.activity) {
      formData.value = initializeFormData();
    }
  },
  { deep: true },
);

// Submit form
async function handleSubmit() {
  // Reset error indicators
  showMeasurementError.value = false;

  // Validate form
  const { valid } = await form.value.validate();
  if (!valid) return;

  // Check if at least one measurement is provided
  if (!hasAnyMeasurement.value) {
    showMeasurementError.value = true;
    handleError({
      title: "Missing Measurements",
      message: "Please provide at least one measurement (weight, height, or head circumference)",
    });
    return;
  }

  // Validate that provided measurements are reasonable
  const validationError = validateMeasurements();
  if (validationError) {
    handleError({
      title: "Invalid Measurement",
      message: validationError,
    });
    return;
  }

  const result = await withErrorHandling(async () => {
    // Growth measurements typically happen at a consistent time (like doctor visits)
    // so we'll set it to noon on the selected date
    const activityDateTime = new Date(`${formData.value.date}T12:00:00`);

    const activityData = {
      type: "growth",
      start_time: activityDateTime,
      notes: formData.value.notes,
      growth_data: {},
    };

    if (formData.value.weight_kg) {
      activityData.growth_data.weight_kg = formData.value.weight_kg;
    }

    if (formData.value.height_cm) {
      activityData.growth_data.height_cm = formData.value.height_cm;
    }

    if (formData.value.head_circumference_cm) {
      activityData.growth_data.head_circumference_cm = formData.value.head_circumference_cm;
    }

    let response;
    if (props.editMode && props.activity) {
      // Update existing activity
      response = await activityStore.updateActivity(props.activity.id, activityData);
    } else {
      // Create new activity
      response = await activityStore.createActivity(activityData);
    }

    if (!response.success) {
      throw new Error(response.error);
    }

    return response.data;
  });

  if (result.success) {
    emit("success", result.data);
  }
}

// Validate measurement values are reasonable
function validateMeasurements() {
  const { weight_kg, height_cm, head_circumference_cm } = formData.value;

  // Weight validation (reasonable range for babies/toddlers: 0.5kg - 50kg)
  if (weight_kg && (weight_kg < 0.5 || weight_kg > 50)) {
    return "Weight should be between 0.5kg and 50kg";
  }

  // Height validation (reasonable range: 20cm - 150cm)
  if (height_cm && (height_cm < 20 || height_cm > 150)) {
    return "Height should be between 20cm and 150cm";
  }

  // Head circumference validation (reasonable range: 20cm - 60cm)
  if (head_circumference_cm && (head_circumference_cm < 20 || head_circumference_cm > 60)) {
    return "Head circumference should be between 20cm and 60cm";
  }

  return null;
}
</script>
