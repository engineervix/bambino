<template>
  <v-container class="pa-0">
    <!-- Header with baby info -->
    <v-sheet class="pa-4 mb-4" color="surface">
      <div class="d-flex align-center">
        <v-avatar size="56" class="mr-3">
          <v-icon size="large">mdi-baby-face</v-icon>
        </v-avatar>
        <div>
          <h1 class="text-h5">Baby</h1>
          <p class="text-body-2 text-grey">{{ currentDate }}</p>
        </div>
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
    <v-dialog v-model="showQuickAdd" max-width="500" persistent>
      <v-card>
        <v-card-title>
          <v-icon class="mr-2">{{ currentActivity?.icon }}</v-icon>
          {{ currentActivity?.title }}
        </v-card-title>
        <v-card-text>
          <p>Quick add form for {{ currentActivity?.id }} - To be implemented</p>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="showQuickAdd = false">Cancel</v-btn>
          <v-btn color="primary" text>Save</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { format } from 'date-fns'
import { useActivityStore } from '@/stores/activity'
import ActivityCard from '@/components/activity/ActivityCard.vue'

const activityStore = useActivityStore()

// State
const showQuickAdd = ref(false)
const currentActivity = ref(null)

// Current date display
const currentDate = computed(() => {
  return format(new Date(), 'EEEE, MMM d')
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
  currentActivity.value = activity
  showQuickAdd.value = true
}

// Load recent stats on mount
onMounted(async () => {
  await activityStore.getRecentStats()
})
</script>