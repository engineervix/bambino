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

    <!-- Milestone type -->
    <v-combobox
      v-model="formData.milestone_type"
      :items="milestoneTypes"
      label="Milestone"
      variant="outlined"
      density="compact"
      :rules="[rules.required, rules.maxLength(50)]"
      placeholder="Select or type a custom milestone"
      class="mb-4"
    />

    <!-- Description -->
    <v-textarea
      v-model="formData.description"
      label="Description"
      variant="outlined"
      rows="3"
      density="compact"
      :rules="[rules.maxLength(500)]"
      placeholder="Tell us about this special moment..."
      class="mb-4"
    />

    <!-- Photo note (future feature) -->
    <v-alert type="info" variant="tonal" density="compact" class="mb-4">
      <template v-slot:prepend>
        <v-icon>mdi-camera</v-icon>
      </template>
      <div class="text-body-2">
        <strong>💡 Tip:</strong> Don't forget to take a photo of this special moment! Photo uploads will be available in
        a future update.
      </div>
    </v-alert>

    <!-- Standardized error display -->
    <v-alert v-if="formError" type="error" variant="tonal" class="mb-4" closable @click:close="clearFormError">
      <div v-if="formError.title" class="font-weight-medium mb-1">{{ formError.title }}</div>
      <div>{{ formError.message || formError }}</div>
    </v-alert>

    <!-- Submit button -->
    <v-btn type="submit" color="milestone" :loading="loading" :disabled="loading" block class="text-none">
      <v-icon start>mdi-party-popper</v-icon>
      {{ editMode ? "Update Milestone" : "Save Milestone" }}
    </v-btn>
  </v-form>
</template>

<script setup>
import { ref, watch } from "vue";
import { useActivityStore } from "@/stores/activity";
import { useErrorHandling } from "@/composables/useErrorHandling";
import { getCurrentDate, getDateString } from "@/utils/datetime";
import { validationRules } from "@/utils/validation";

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
      date: getDateString(startTime),
      milestone_type: activity.milestone_data?.milestone_type || "",
      description: activity.milestone_data?.description || "",
    };
  }
  return {
    date: getCurrentDate(),
    milestone_type: "",
    description: "",
  };
};

const formData = ref(initializeFormData());

// Validation rules
const rules = {
  required: validationRules.required,
  maxLength: validationRules.maxLength,
};

// Common milestone types (users can add custom ones via combobox)
const milestoneTypes = [
  "First Smile",
  "First Laugh",
  "First Word",
  "First Steps",
  "First Tooth",
  "Rolled Over",
  "Sat Up",
  "Crawled",
  "Stood Up",
  "Pulled to Stand",
  'Said "Mama"',
  'Said "Dada"',
  "Clapped Hands",
  "Waved Bye-Bye",
  "Slept Through Night",
  "First Solid Food",
  "Used Sippy Cup",
  "First Hair Cut",
  "First Day at Daycare",
  "First Trip",
  "First Holiday",
  "First Birthday",
  "Walked Independently",
  "Climbed Stairs",
  "Used Potty",
  "First Swimming",
  "First Beach Visit",
  "Met Grandparents",
  "First Book Read",
  "Played Peek-a-Boo",
  "First Toy Preference",
];

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

  // Additional validation
  if (!formData.value.milestone_type.trim()) {
    handleError({
      title: "Missing Milestone",
      message: "Please specify what milestone this is",
    });
    return;
  }

  const result = await withErrorHandling(async () => {
    // Milestones typically happen at a memorable time, so we'll set to noon
    const activityDateTime = new Date(`${formData.value.date}T12:00:00`);

    const activityData = {
      type: "milestone",
      start_time: activityDateTime,
      notes: "", // Using description field instead of notes for milestones
      milestone_data: {
        milestone_type: formData.value.milestone_type.trim(),
        description: formData.value.description.trim(),
      },
    };

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

<style scoped>
/* Prevent text transformation on the button */
.text-none {
  text-transform: none !important;
}

/* Add some celebration feeling */
.v-btn.text-milestone {
  background: linear-gradient(45deg, #ff5722, #ff9800);
}
</style>
