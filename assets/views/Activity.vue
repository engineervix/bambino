<template>
  <div>
    <!-- Header with baby info -->
    <v-sheet
      class="hero-banner mb-8"
      :height="bannerHeight"
      elevation="0"
      :style="{ opacity: headerOpacity }"
    >
      <v-container class="fill-height d-flex align-center px-6" fluid>
        <div class="d-flex flex-column align-center justify-center text-center w-100">
          <v-avatar :size="avatarSize" class="mb-4 elevation-2">
            <v-img src="/baby.svg" cover />
          </v-avatar>
          <h1 class="text-h4 font-weight-bold mb-2">
            {{ currentBaby ? currentBaby.name : 'Bambino' }}
          </h1>
          <p class="text-subtitle-2 mb-0">
            {{ currentBaby ? currentBaby.age_display : 'No profile' }} • {{ currentDate }}
          </p>
        </div>
      </v-container>

      <!-- Decorative wave -->
      <svg class="wave" viewBox="0 0 1440 320" preserveAspectRatio="none">
        <path
          d="M0,160L48,154.7C96,149,192,139,288,138.7C384,139,480,149,576,170.7C672,192,768,224,864,202.7C960,181,1056,107,1152,85.3C1248,64,1344,96,1392,112L1440,128L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z"
          fill="currentColor"
          opacity="0.12"
        />
      </svg>
    </v-sheet>

    <!-- Activity cards -->
    <v-container>
      <v-row>
        <!-- Main activity types - 2 columns on larger screens -->
        <v-col
          v-for="activity in mainActivities"
          :key="activity.id"
          cols="12"
          sm="6"
          lg="4"
          class="py-3"
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
              <v-expansion-panel-title color="growth-bg">
                <v-icon :color="'growth'" class="mr-3">mdi-human-male-height</v-icon>
                <span class="text-h6">Growth</span>
              </v-expansion-panel-title>
              <v-expansion-panel-text>
                <v-list>
                  <v-list-item @click="handleQuickAdd({ id: 'growth' })">
                    <template v-slot:prepend>
                      <v-icon :color="'growth'">mdi-human-male-height</v-icon>
                    </template>
                    <v-list-item-title>Record Measurements</v-list-item-title>
                    <v-list-item-subtitle>Weight, height, head circumference</v-list-item-subtitle>
                  </v-list-item>
                </v-list>
              </v-expansion-panel-text>
            </v-expansion-panel>

            <!-- Health section -->
            <v-expansion-panel>
              <v-expansion-panel-title color="health-bg">
                <v-icon :color="'health'" class="mr-3">mdi-medical-bag</v-icon>
                <span class="text-h6">Health</span>
              </v-expansion-panel-title>
              <v-expansion-panel-text>
                <v-list>
                  <v-list-item @click="handleQuickAdd({ id: 'health', subType: 'checkup' })">
                    <template v-slot:prepend>
                      <v-icon :color="'health'">mdi-stethoscope</v-icon>
                    </template>
                    <v-list-item-title>Medical Checkup</v-list-item-title>
                    <v-list-item-subtitle>Doctor visit, routine checkup</v-list-item-subtitle>
                  </v-list-item>

                  <v-list-item @click="handleQuickAdd({ id: 'health', subType: 'vaccine' })">
                    <template v-slot:prepend>
                      <v-icon :color="'health'">mdi-needle</v-icon>
                    </template>
                    <v-list-item-title>Vaccination</v-list-item-title>
                    <v-list-item-subtitle>Record vaccines received</v-list-item-subtitle>
                  </v-list-item>

                  <v-list-item @click="handleQuickAdd({ id: 'health', subType: 'illness' })">
                    <template v-slot:prepend>
                      <v-icon :color="'health'">mdi-thermometer</v-icon>
                    </template>
                    <v-list-item-title>Illness</v-list-item-title>
                    <v-list-item-subtitle>Symptoms, treatment</v-list-item-subtitle>
                  </v-list-item>
                </v-list>
              </v-expansion-panel-text>
            </v-expansion-panel>
          </v-expansion-panels>
        </v-col>
      </v-row>
    </v-container>

    <!-- Quick add dialog -->
    <v-fade-transition>
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
    </v-fade-transition>

    <!-- Success snackbar -->
    <v-snackbar
      v-model="showSuccess"
      color="success"
      :timeout="3000"
      location="bottom"
    >
      <v-icon start>mdi-check-circle</v-icon>
      Activity saved successfully!
    </v-snackbar>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, markRaw } from 'vue'
import { onUnmounted } from 'vue'
import { format } from 'date-fns'
import { useActivityStore } from '@/stores/activity'
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import ActivityCard from '@/components/activity/ActivityCard.vue'
import { useDisplay } from 'vuetify'

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

// Hero banner opacity
const headerOpacity = ref(1)

// Responsive sizes
const display = useDisplay()
const bannerHeight = computed(() => (display.mdAndUp.value ? 240 : 180))
const avatarSize = computed(() => (display.mdAndUp.value ? 96 : 72))

function handleScroll() {
  headerOpacity.value = Math.max(0, 1 - window.scrollY / 180)
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll, { passive: true })
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})

// Current date display
const currentDate = computed(() => {
  return format(new Date(), 'EEEE, MMM d')
})

// Current form component
const currentFormComponent = computed(() => {
  if (!currentActivity.value) return null
  return formComponents[currentActivity.value.id] || null
})

// Main activities (cards) - This was the issue!
const mainActivities = computed(() => {
  return [
    {
      id: 'feed',
      title: 'Feed',
      description: 'Track a feeding session',
      icon: 'mdi-baby-bottle',
      color: 'feed',
      hasTimer: true
    },
    {
      id: 'pump',
      title: 'Pump',
      description: 'Track a pumping session',
      icon: 'mdi-mother-nurse',
      color: 'pump',
      hasTimer: true
    },
    {
      id: 'diaper',
      title: 'Diaper',
      description: 'Track a diaper change',
      icon: 'mdi-baby',
      color: 'diaper',
      hasTimer: false
    },
    {
      id: 'sleep',
      title: 'Sleep',
      description: 'Track a sleep session',
      icon: 'mdi-sleep',
      color: 'sleep',
      hasTimer: true
    },
    {
      id: 'milestone',
      title: 'Baby Firsts',
      description: 'Track memorable moments',
      icon: 'mdi-party-popper',
      color: 'milestone',
      hasTimer: false
    }
  ]
})

// Handlers
function handleActivityClick(activity) {
  // For now, same as quick add
  handleQuickAdd(activity)
}

function handleQuickAdd(activity) {
  // Find the full activity config from store or use the passed activity
  const fullActivity = activityStore.activityTypes.find(a => a.id === activity.id) || activity
  currentActivity.value = fullActivity
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

<style scoped>
.hero-banner {
  position: relative;
  background: linear-gradient(135deg, rgba(var(--v-theme-primary),0.35) 0%, rgba(var(--v-theme-accent1),0.35) 100%);
  color: white;
  overflow: hidden;
  /* wave overscroll handles bottom edge */
}

.wave {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 60px;
  pointer-events: none;
  color: rgba(255,255,255,0.5);
}
</style>
