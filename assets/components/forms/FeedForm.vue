<template>
  <v-form ref="form" @submit.prevent="handleSubmit">
    <!-- Feed type selection -->
    <v-btn-toggle
      v-model="formData.feed_type"
      mandatory
      color="primary"
      density="comfortable"
      class="w-100 mb-4 d-flex"
      :disabled="isTimerRunning"
    >
      <v-btn value="bottle" class="flex-grow-1">
        <v-icon start>mdi-bottle-baby</v-icon>
        Bottle
      </v-btn>
      <v-btn value="breast_left" class="flex-grow-1">
        <v-icon start>mdi-chevron-left</v-icon>
        Left
      </v-btn>
      <v-btn value="breast_right" class="flex-grow-1">
        <v-icon start>mdi-chevron-right</v-icon>
        Right
      </v-btn>
    </v-btn-toggle>

    <!-- Timer or manual entry -->
    <v-card v-if="!useTimer" variant="outlined" class="mb-4">
      <v-card-text>
        <!-- Date and time inputs -->
        <v-row dense class="mb-3">
          <v-col cols="6">
            <v-text-field
              v-model="formData.date"
              label="Date"
              type="date"
              variant="outlined"
              density="compact"
              :rules="[rules.required]"
              :disabled="isTimerRunning"
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
              :disabled="isTimerRunning"
            />
          </v-col>
        </v-row>

        <!-- Amount input (for bottle) -->
        <v-text-field
          v-if="formData.feed_type === 'bottle'"
          v-model.number="formData.amount_ml"
          label="Amount (ml)"
          type="number"
          variant="outlined"
          density="compact"
          append-inner-icon="mdi-baby-bottle"
          :rules="[rules.positiveNumber]"
          class="mb-3"
        />

        <!-- Duration input (for breast) -->
        <v-text-field
          v-if="formData.feed_type !== 'bottle'"
          v-model.number="formData.duration_minutes"
          label="Duration (minutes)"
          type="number"
          variant="outlined"
          density="compact"
          append-inner-icon="mdi-timer"
          :rules="[rules.positiveInteger]"
        />
      </v-card-text>
    </v-card>

    <!-- Timer display -->
    <v-card v-else variant="tonal" color="primary" class="mb-4">
      <v-card-text class="text-center">
        <v-icon size="48" class="mb-2">mdi-timer</v-icon>
        <p class="text-h4">{{ timerDisplay }}</p>
        <p class="text-body-2">Timer Running</p>

        <!-- Amount input while timer running (for bottle) -->
        <v-text-field
          v-if="formData.feed_type === 'bottle'"
          v-model.number="formData.amount_ml"
          label="Amount (ml)"
          type="number"
          variant="outlined"
          density="compact"
          append-inner-icon="mdi-baby-bottle"
          :rules="[rules.positiveNumber]"
          class="mt-3"
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
      :rules="[rules.maxLength(1000)]"
      class="mb-4"
    />

    <!-- Standardized error display -->
    <v-alert v-if="formError" type="error" variant="tonal" class="mb-4" closable @click:close="clearFormError">
      <div v-if="formError.title" class="font-weight-medium mb-1">{{ formError.title }}</div>
      <div>{{ formError.message || formError }}</div>
    </v-alert>

    <!-- Actions -->
    <div class="d-flex gap-2">
      <!-- This button handles saving manual entries and updates -->
      <v-btn
        v-if="!useTimer"
        type="submit"
        color="primary"
        :loading="loading"
        :disabled="loading || isTimerRunning"
        :block="editMode || !feedSupportsTimer"
        class="flex-grow-1"
      >
        {{ editMode ? "Update Feed" : "Save" }}
      </v-btn>

      <!-- This button starts the timer, shown only for new entries -->
      <v-btn
        v-if="!useTimer && !editMode && feedSupportsTimer"
        variant="outlined"
        color="primary"
        @click="startTimer"
        :loading="loading"
        :disabled="loading"
        class="flex-grow-1"
      >
        <v-icon start>mdi-timer</v-icon>
        Start Timer
      </v-btn>

      <!-- This button stops the timer -->
      <v-btn v-if="useTimer" color="success" @click="stopTimer" :loading="loading" :disabled="loading" block>
        <v-icon start>mdi-stop</v-icon>
        Stop & Save
      </v-btn>
    </div>
  </v-form>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from "vue";
import { useTimerStore } from "@/stores/timer";
import { useActivityStore } from "@/stores/activity";
import { useErrorHandling } from "@/composables/useErrorHandling";
import { combineDateAndTime, getCurrentDate, getCurrentTime, getDateString, getTimeString } from "@/utils/datetime";
import { validationRules, validateDateTime } from "@/utils/validation";

const props = defineProps({
  hasTimer: {
    type: Boolean,
    default: true,
  },
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
const timerStore = useTimerStore();
const activityStore = useActivityStore();

// Error handling
const { error: formError, loading, handleError, clearError: clearFormError, withErrorHandling } = useErrorHandling();

// Form state
const form = ref(null);
const useTimer = ref(false);
const timerInterval = ref(null);

// Determine if the current feed type supports a timer
const feedSupportsTimer = computed(() => {
  if (!props.hasTimer) return false;
  return formData.value.feed_type === "breast_left" || formData.value.feed_type === "breast_right";
});

const isTimerRunning = computed(() => {
  if (!props.editMode || !props.activity) return false;
  // A running activity has no end_time.
  return feedSupportsTimer.value && !props.activity.end_time;
});

// Initialize form data from props or defaults
const initializeFormData = () => {
  if (props.editMode && props.activity) {
    const activity = props.activity;
    const startTime = new Date(activity.start_time);
    return {
      feed_type: activity.feed_data?.feed_type || "bottle",
      date: getDateString(startTime),
      time: getTimeString(startTime),
      amount_ml: activity.feed_data?.amount_ml || null,
      duration_minutes: activity.feed_data?.duration_minutes || null,
      notes: activity.notes || "",
    };
  }
  return {
    feed_type: "bottle",
    date: getCurrentDate(),
    time: getCurrentTime(),
    amount_ml: null,
    duration_minutes: null,
    notes: "",
  };
};

const formData = ref(initializeFormData());

// Validation rules
const rules = {
  required: validationRules.required,
  positiveNumber: validationRules.positiveNumber,
  positiveInteger: validationRules.positiveInteger,
  maxLength: validationRules.maxLength,
};

// Timer display with persistent calculation
const timerDisplay = computed(() => {
  if (!timerStore.hasActiveTimer("feed")) return "00:00";
  return timerStore.formattedDurations.feed || "00:00";
});

// Watch for active timer changes (only if not in edit mode)
watch(
  () => timerStore.hasActiveTimer("feed"),
  (hasTimer) => {
    if (!props.editMode) {
      useTimer.value = hasTimer;
      if (hasTimer) {
        startTimerDisplay();
      } else {
        stopTimerDisplay();
      }
    }
  },
  { immediate: true },
);

onUnmounted(() => {
  stopTimerDisplay();
});

// Initialize on mount
onMounted(() => {
  // Don't show timer in edit mode
  if (props.editMode) {
    useTimer.value = false;
  }
});

// Timer display management
function startTimerDisplay() {
  if (timerInterval.value) {
    clearInterval(timerInterval.value);
  }
  timerInterval.value = setInterval(() => {
    // Force reactivity update
  }, 1000);
}

function stopTimerDisplay() {
  if (timerInterval.value) {
    clearInterval(timerInterval.value);
    timerInterval.value = null;
  }
}

// Start timer
async function startTimer() {
  const result = await withErrorHandling(async () => {
    const response = await activityStore.startTimer("feed", {
      feed_data: {
        feed_type: formData.value.feed_type,
      },
      notes: formData.value.notes,
    });

    if (!response.success) {
      throw new Error(response.error);
    }

    const success = timerStore.startTimer("feed", response.data.id);
    if (!success) {
      throw new Error("Failed to start local timer");
    }

    return response.data;
  });

  if (!result.success) {
    handleError({
      title: "Timer Start Failed",
      message: result.error,
    });
  }
}

// Stop timer
async function stopTimer() {
  const timer = timerStore.getActiveTimer("feed");
  if (!timer?.activityId) {
    handleError({
      title: "Timer Error",
      message: "No active timer found",
    });
    return;
  }

  const result = await withErrorHandling(async () => {
    const stopData = {};
    if (formData.value.feed_type === "bottle" && formData.value.amount_ml) {
      stopData.amount_ml = formData.value.amount_ml;
    }
    if (formData.value.notes) {
      stopData.notes = formData.value.notes;
    }

    const response = await activityStore.stopTimer(timer.activityId, stopData);
    if (!response.success) {
      throw new Error(response.error);
    }

    timerStore.stopTimer("feed");
    return response.data;
  });

  if (result.success) {
    emit("success", result.data);
  }
}

// Submit form
async function handleSubmit() {
  // Validate form
  const { valid } = await form.value.validate();
  if (!valid) return;

  // Custom validation for date/time
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
      type: "feed",
      start_time: activityDateTime,
      notes: formData.value.notes,
      feed_data: {
        feed_type: formData.value.feed_type,
      },
    };

    if (formData.value.feed_type === "bottle" && formData.value.amount_ml) {
      activityData.feed_data.amount_ml = formData.value.amount_ml;
    }

    if (formData.value.feed_type !== "bottle" && formData.value.duration_minutes) {
      activityData.feed_data.duration_minutes = formData.value.duration_minutes;
      activityData.end_time = new Date(activityDateTime.getTime() + formData.value.duration_minutes * 60000);
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
</script>
