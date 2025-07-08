<template>
  <v-container class="pa-0">
    <!-- Header with baby info -->
    <v-sheet class="pa-4 mb-4" color="surface">
      <div v-if="currentBaby" class="d-flex align-center">
        <v-avatar size="56" class="mr-3">
          <v-icon size="large">mdi-baby-face</v-icon>
        </v-avatar>
        <div>
          <h1 class="text-h5">{{ currentBaby.name }}</h1>
          <p class="text-body-2 text-grey">{{ currentBaby.age_display }} • {{ currentDate }}</p>
        </div>
      </div>
      <div v-else class="text-center py-4">
        <v-icon size="48" class="mb-2 text-grey">mdi-baby-face-outline</v-icon>
        <p class="text-body-1 text-grey">No baby profile found</p>
        <v-btn color="primary" variant="tonal" size="small" to="/account">
          Create Baby Profile
        </v-btn>
      </div>
    </v-sheet>

    <!-- Activity cards -->
    <v-container>
      <v-row>
        <!-- Main activity types -->
        <v-col 
          v-for="activity in mainActivities" 
          :key="activity.id"
          cols="12"
          class="py-2"
        >
          <activity-card
            :title="activity.title"
            :description="activity.description"
            :icon="activity.icon"
            :color="activity.color"
            @click="handleActivityClick(activity)"
            @add="handleQuickAdd(activity)"
          />
        </v-col>

        <!-- Expandable sections -->
        <v-col cols="12" class="py-2">
          <v-expansion-panels variant="accordion">
            <!-- Growth section -->
            <v-expansion-panel>
              <v-expansion-panel-title color="growth">
                <v-icon class="mr-3">mdi-human-male-height</v-icon>
                <span class="text-h6">Growth</span>
              </v-expansion-panel-title>
              <v-expansion-panel-text>
                <v-list>
                  <v-list-item
                    v-for="type in growthTypes"
                    :key="type.id"
                    @click="handleQuickAdd({ id: 'growth', subType: type.id })"
                  >
                    <template v-slot:prepend>
                      <v-icon>{{ type.icon }}</v-icon>
                    </template>
                    <v-list-item-title>{{ type.title }}</v-list-item-title>
                  </v-list-item>
                </v-list>
              </v-expansion-panel-text>
            </v-expansion-panel>

            <!-- Health section -->
            <v-expansion-panel>
              <v-expansion-panel-title color="health">
                <v-icon class="mr-3">mdi-medical-bag</v-icon>
                <span class="text-h6">Health</span>
              </v-expansion-panel-title>
              <v-expansion-panel-text>
                <v-list>
                  <v-list-item
                    v-for="type in healthTypes"
                    :key="type.id"
                    @click="handleQuickAdd({ id: 'health', subType: type.id })"
                  >
                    <template v-slot:prepend>
                      <v-icon>{{ type.icon }}</v-icon>
                    </template>
                    <v-list-item-title>{{ type.title }}</v-list-item-title>
                  </v-list-item>
                </v-list>
              </v-expansion-panel-text>
            </v-expansion-panel>
          </v-expansion-panels>
        </v-col>
      </v-row>
    </v-container>

    <!-- Quick add dialog -->
    <v-dialog v-model="showQuickAdd" max-width="500" persistent scrollable>
      <v-card>
        <v-card-title class="d-flex align-center">
          <v-icon class="mr-2" :color="currentActivity?.color">{{ currentActivity?.icon }}</v-icon>
          <span>{{ currentActivity?.title }}</span>
          <v-spacer></v-spacer>
          <v-btn icon variant="text" @click="closeDialog">
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-card-title>
        
        <v-divider></v-divider>
        
        <v-card-text class="pa-4">
          <!-- Dynamic form component -->
          <component
            v-if="currentFormComponent"
            :is="currentFormComponent"
            :has-timer="currentActivity?.hasTimer"
            @success="handleFormSuccess"
            @cancel="closeDialog"
          />
        </v-card-text>
      </v-card>
    </v-dialog>

    <!-- Success snackbar -->
    <v-snackbar
      v-model="showSuccess"
      color="success"
      :timeout="3000"
      location="top"
    >
      <v-icon start>mdi-check-circle</v-icon>
      Activity saved successfully!
    </v-snackbar>
  </v-container>
</template>

<script setup>
import { ref, computed, onMounted, markRaw } from 'vue'
import { format } from 'date-fns'
import { useActivityStore } from '@/stores/activity'
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import ActivityCard from '@/components/activity/ActivityCard.vue'

// Form components
import FeedForm from '@/components/forms/FeedForm.vue'
import PumpForm from '@/components/forms/PumpForm.vue'
import DiaperForm from '@/components/forms/DiaperForm.vue'
import SleepForm from '@/components/forms/SleepForm.vue'
import GrowthForm from '@/components/forms/GrowthForm.vue'
import HealthForm from '@/components/forms/HealthForm.vue'
import MilestoneForm from '@/components/forms/MilestoneForm.vue'

// Mark components as raw to avoid reactivity overhead
const formComponents = {
  feed: markRaw(FeedForm),
  pump: markRaw(PumpForm),
  diaper: markRaw(DiaperForm),
  sleep: markRaw(SleepForm),
  growth: markRaw(GrowthForm),
  health: markRaw(HealthForm),
  milestone: markRaw(MilestoneForm)
}

const activityStore = useActivityStore()
const authStore = useAuthStore()
const { currentBaby } = storeToRefs(authStore)

// State
const showQuickAdd = ref(false)
const showSuccess = ref(false)
const currentActivity = ref(null)

// Current date display
const currentDate = computed(() => {
  return format(new Date(), 'EEEE, MMM d')
})

// Current form component
const currentFormComponent = computed(() => {
  if (!currentActivity.value) return null
  return formComponents[currentActivity.value.id] || null
})

// Main activities (cards)
const mainActivities = computed(() => {
  return activityStore.activityTypes.filter(a => 
    ['feed', 'pump', 'diaper', 'sleep', 'milestone'].includes(a.id)
  )
})

// Growth types
const growthTypes = [
  { id: 'weight', title: 'Weight', icon: 'mdi-scale' },
  { id: 'height', title: 'Height', icon: 'mdi-human-male-height-variant' },
  { id: 'head', title: 'Head Size', icon: 'mdi-head' }
]

// Health types
const healthTypes = [
  { id: 'medical', title: 'Medical', icon: 'mdi-doctor' },
  { id: 'vaccine', title: 'Vaccine', icon: 'mdi-needle' }
]

// Handlers
function handleActivityClick(activity) {
  // For now, same as quick add
  handleQuickAdd(activity)
}

function handleQuickAdd(activity) {
  // Find the full activity config from store
  const fullActivity = activityStore.activityTypes.find(a => a.id === activity.id)
  currentActivity.value = fullActivity || activity
  showQuickAdd.value = true
}

function handleFormSuccess(data) {
  // Close dialog
  closeDialog()
  
  // Show success message
  showSuccess.value = true
  
  // Optionally refresh stats
  activityStore.getRecentStats()
}

function closeDialog() {
  showQuickAdd.value = false
  // Clear current activity after animation
  setTimeout(() => {
    currentActivity.value = null
  }, 300)
}

// Load recent stats on mount
onMounted(async () => {
  await activityStore.getRecentStats()
})
</script>