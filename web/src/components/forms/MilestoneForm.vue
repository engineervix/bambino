<template>
  <v-form ref="form" @submit.prevent="handleSubmit">
    <!-- Date -->
    <v-text-field
      v-model="formData.date"
      label="Date"
      type="date"
      variant="outlined"
      density="compact"
      class="mb-4"
    />

    <!-- Milestone type -->
    <v-combobox
      v-model="formData.milestone_type"
      :items="milestoneTypes"
      label="Milestone"
      variant="outlined"
      density="compact"
      class="mb-4"
      :rules="[v => !!v || 'Milestone type is required']"
    />

    <!-- Description -->
    <v-textarea
      v-model="formData.description"
      label="Description"
      variant="outlined"
      rows="3"
      density="compact"
      class="mb-4"
      placeholder="Tell us about this special moment..."
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
      color="milestone"
      :loading="loading"
      block
    >
      <v-icon start>mdi-party-popper</v-icon>
      Save Milestone
    </v-btn>
  </v-form>
</template>

<script setup>
import { ref } from 'vue'
import { useActivityStore } from '@/stores/activity'

const emit = defineEmits(['success', 'cancel'])

const activityStore = useActivityStore()

// Form state
const form = ref(null)
const loading = ref(false)
const formError = ref(null)
const formData = ref({
  date: new Date().toISOString().split('T')[0],
  milestone_type: '',
  description: ''
})

// Common milestone types (users can add custom ones via combobox)
const milestoneTypes = [
  'First Smile',
  'First Laugh',
  'First Word',
  'First Steps',
  'First Tooth',
  'Rolled Over',
  'Sat Up',
  'Crawled',
  'Stood Up',
  'Slept Through Night',
  'First Solid Food',
  'First Hair Cut',
  'First Day at Daycare',
  'First Trip',
  'First Holiday',
  'First Birthday'
]

// Submit form
async function handleSubmit() {
  const { valid } = await form.value.validate()
  if (!valid) return

  loading.value = true
  formError.value = null
  
  const activityData = {
    type: 'milestone',
    start_time: new Date(`${formData.value.date}T12:00:00`), // Default to noon
    notes: '', // Using description field instead
    milestone_data: {
      milestone_type: formData.value.milestone_type,
      description: formData.value.description
    }
  }
  
  const result = await activityStore.createActivity(activityData)
  
  if (result.success) {
    emit('success', result.data)
  } else {
    formError.value = result.error || 'Failed to save milestone'
  }
  
  loading.value = false
}
</script>

<style scoped>
/* Add some festive styling for milestones */
.v-btn {
  text-transform: none;
}
</style>