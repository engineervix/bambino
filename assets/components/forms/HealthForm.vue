<template>
  <v-form ref="form" @submit.prevent="handleSubmit">
    <!-- Record type -->
    <v-select
      v-model="formData.record_type"
      :items="recordTypes"
      label="Type"
      variant="outlined"
      density="compact"
      :rules="[rules.required]"
      class="mb-4"
    />

    <!-- Date and time -->
    <v-row dense class="mb-4">
      <v-col cols="6">
        <v-text-field
          v-model="formData.date"
          label="Date"
          type="date"
          variant="outlined"
          density="compact"
          :rules="[rules.required]"
        />
      </v-col>
      <v-col cols="6">
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

    <!-- Provider -->
    <v-text-field
      v-model="formData.provider"
      label="Provider/Location"
      variant="outlined"
      density="compact"
      :rules="[rules.maxLength(100)]"
      placeholder="Dr. Smith, City Hospital, etc."
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
          :rules="[rules.required, rules.maxLength(100)]"
          placeholder="DTaP, Rotavirus, etc."
          class="mb-4"
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
          :rules="[rules.maxLength(500)]"
          placeholder="Fever, cough, runny nose, etc."
          class="mb-3"
        />

        <v-textarea
          v-model="formData.treatment"
          label="Treatment"
          variant="outlined"
          rows="2"
          density="compact"
          :rules="[rules.maxLength(500)]"
          placeholder="Medication, rest, monitoring, etc."
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
      :rules="[rules.maxLength(1000)]"
      placeholder="Additional notes, observations, follow-up instructions..."
      class="mb-4"
    />

    <!-- Standardized error display -->
    <v-alert v-if="formError" type="error" variant="tonal" class="mb-4" closable @click:close="clearFormError">
      <div v-if="formError.title" class="font-weight-medium mb-1">{{ formError.title }}</div>
      <div>{{ formError.message || formError }}</div>
    </v-alert>

    <!-- Submit button -->
    <v-btn type="submit" color="primary" :loading="loading" :disabled="loading" block>
      <v-icon start>{{ getRecordIcon() }}</v-icon>
      {{ editMode ? "Update Health Record" : "Save Health Record" }}
    </v-btn>
  </v-form>
</template>

<script setup>
import { ref, watch, computed } from "vue";
import { useActivityStore } from "@/stores/activity";
import { useErrorHandling } from "@/composables/useErrorHandling";
import { combineDateAndTime, getCurrentDate, getCurrentTime, getDateString, getTimeString } from "@/utils/datetime";
import { validationRules, validateDateTime } from "@/utils/validation";

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

// Initialize form data from props or defaults
const initializeFormData = () => {
  if (props.editMode && props.activity) {
    const activity = props.activity;
    const startTime = new Date(activity.start_time);
    return {
      record_type: activity.health_data?.record_type || "checkup",
      date: getDateString(startTime),
      time: getTimeString(startTime),
      provider: activity.health_data?.provider || "",
      vaccine_name: activity.health_data?.vaccine_name || "",
      symptoms: activity.health_data?.symptoms || "",
      treatment: activity.health_data?.treatment || "",
      notes: activity.notes || "",
    };
  }
  return {
    record_type: "checkup",
    date: getCurrentDate(),
    time: getCurrentTime(),
    provider: "",
    vaccine_name: "",
    symptoms: "",
    treatment: "",
    notes: "",
  };
};

const formData = ref(initializeFormData());

// Validation rules
const rules = {
  required: validationRules.required,
  maxLength: validationRules.maxLength,
};

// Record type options
const recordTypes = [
  { title: "Checkup", value: "checkup" },
  { title: "Vaccine", value: "vaccine" },
  { title: "Illness", value: "illness" },
];

// Get appropriate icon for the record type
const getRecordIcon = computed(() => {
  return () => {
    switch (formData.value.record_type) {
      case "vaccine":
        return "mdi-needle";
      case "illness":
        return "mdi-thermometer";
      case "checkup":
      default:
        return "mdi-stethoscope";
    }
  };
});

// Clear type-specific fields when record type changes
watch(
  () => formData.value.record_type,
  (newType) => {
    if (newType !== "vaccine") {
      formData.value.vaccine_name = "";
    }
    if (newType !== "illness") {
      formData.value.symptoms = "";
      formData.value.treatment = "";
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
  // Validate form
  const { valid } = await form.value.validate();
  if (!valid) return;

  // Custom validation for vaccine records
  if (formData.value.record_type === "vaccine" && !formData.value.vaccine_name.trim()) {
    handleError({
      title: "Missing Vaccine Information",
      message: "Vaccine name is required for vaccine records",
    });
    return;
  }

  // Validate date/time
  const dateTimeError = validateDateTime(formData.value.date, formData.value.time);
  if (dateTimeError) {
    handleError({
      title: "Invalid Date/Time",
      message: dateTimeError,
    });
    return;
  }

  const result = await withErrorHandling(async () => {
    const activityDateTime = combineDateAndTime(formData.value.date, formData.value.time);

    const activityData = {
      type: "health",
      start_time: activityDateTime,
      notes: formData.value.notes,
      health_data: {
        record_type: formData.value.record_type,
        provider: formData.value.provider,
      },
    };

    if (formData.value.record_type === "vaccine" && formData.value.vaccine_name) {
      activityData.health_data.vaccine_name = formData.value.vaccine_name;
    }

    if (formData.value.record_type === "illness") {
      if (formData.value.symptoms) {
        activityData.health_data.symptoms = formData.value.symptoms;
      }
      if (formData.value.treatment) {
        activityData.health_data.treatment = formData.value.treatment;
      }
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
</script>
