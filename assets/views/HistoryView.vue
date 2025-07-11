<template>
  <v-container>
    <!-- Header -->
    <div class="d-flex align-center mb-4">
      <h1 class="text-h4">History</h1>
      <v-spacer></v-spacer>
      <v-btn icon variant="text" @click="showFilters = !showFilters" :color="hasActiveFilters ? 'primary' : undefined">
        <v-icon>mdi-filter-variant</v-icon>
      </v-btn>
    </div>

    <!-- Filters -->
    <v-expand-transition>
      <v-card v-if="showFilters" variant="outlined" class="mb-4">
        <v-card-text>
          <v-row>
            <!-- Date Range -->
            <v-col cols="12" sm="6">
              <v-text-field
                v-model="filters.startDate"
                label="Start Date"
                type="date"
                variant="outlined"
                density="compact"
                clearable
              />
            </v-col>
            <v-col cols="12" sm="6">
              <v-text-field
                v-model="filters.endDate"
                label="End Date"
                type="date"
                variant="outlined"
                density="compact"
                clearable
              />
            </v-col>

            <!-- Activity Type Filter -->
            <v-col cols="12" sm="6">
              <v-select
                v-model="filters.type"
                :items="activityTypeOptions"
                label="Activity Type"
                variant="outlined"
                density="compact"
                clearable
              />
            </v-col>

            <!-- Quick Date Presets -->
            <v-col cols="12" sm="6">
              <div class="d-flex flex-wrap gap-2">
                <v-btn
                  v-for="preset in datePresets"
                  :key="preset.key"
                  size="small"
                  variant="outlined"
                  @click="applyDatePreset(preset)"
                >
                  {{ preset.label }}
                </v-btn>
              </div>
            </v-col>
          </v-row>

          <!-- Filter Actions -->
          <div class="d-flex gap-2 mt-2">
            <v-btn color="primary" @click="applyFilters" :loading="loading"> Apply Filters </v-btn>
            <v-btn variant="outlined" @click="clearFilters"> Clear </v-btn>
          </div>
        </v-card-text>
      </v-card>
    </v-expand-transition>

    <!-- Activities List -->
    <div v-if="loading && activities.length === 0" class="text-center py-8">
      <v-progress-circular indeterminate />
      <p class="mt-2 text-grey">Loading activities...</p>
    </div>

    <div v-else-if="activities.length === 0" class="text-center py-8">
      <v-icon size="64" class="text-grey mb-4">mdi-calendar-blank</v-icon>
      <h3 class="text-h6 mb-2">No Activities Found</h3>
      <p class="text-grey">
        {{ hasActiveFilters ? "Try adjusting your filters" : "Start tracking activities to see them here" }}
      </p>
    </div>

    <div v-else>
      <!-- Activity Cards -->
      <div class="activity-list">
        <v-card v-for="activity in activities" :key="activity.id" class="mb-3" @click="viewActivity(activity)">
          <v-card-text class="pa-4">
            <div class="d-flex align-start">
              <!-- Activity Icon -->
              <v-avatar :color="getActivityColor(activity.type)" size="40" class="mr-3">
                <v-icon color="white">{{ getActivityIcon(activity.type) }}</v-icon>
              </v-avatar>

              <!-- Activity Content -->
              <div class="flex-grow-1">
                <div class="d-flex align-center mb-1">
                  <h4 class="text-h6">{{ getActivityTitle(activity.type) }}</h4>
                  <v-spacer></v-spacer>
                  <span class="text-caption text-grey">
                    {{ formatActivityDate(activity.start_time) }}
                  </span>
                </div>

                <!-- Activity Details -->
                <div class="activity-details">
                  <ActivityDetails :activity="activity" />
                </div>

                <!-- Notes -->
                <p v-if="activity.notes" class="text-body-2 text-grey mt-2">
                  {{ activity.notes }}
                </p>
              </div>

              <!-- Action Menu -->
              <v-menu>
                <template v-slot:activator="{ props }">
                  <v-btn icon variant="text" size="small" v-bind="props" @click.stop>
                    <v-icon>mdi-dots-vertical</v-icon>
                  </v-btn>
                </template>
                <v-list>
                  <v-list-item @click="editActivity(activity)">
                    <template v-slot:prepend>
                      <v-icon>mdi-pencil</v-icon>
                    </template>
                    <v-list-item-title>Edit</v-list-item-title>
                  </v-list-item>
                  <v-list-item @click="confirmDelete(activity)">
                    <template v-slot:prepend>
                      <v-icon>mdi-delete</v-icon>
                    </template>
                    <v-list-item-title>Delete</v-list-item-title>
                  </v-list-item>
                </v-list>
              </v-menu>
            </div>
          </v-card-text>
        </v-card>
      </div>

      <!-- Load More / Pagination -->
      <div v-if="pagination.totalPages > 1" class="text-center mt-4">
        <v-btn v-if="pagination.page < pagination.totalPages" variant="outlined" @click="loadMore" :loading="loading">
          Load More Activities
        </v-btn>

        <p class="text-caption text-grey mt-2">Showing {{ activities.length }} of {{ pagination.total }} activities</p>
      </div>
    </div>

    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="showDeleteDialog" max-width="400">
      <v-card>
        <v-card-title>Delete Activity</v-card-title>
        <v-card-text>
          Are you sure you want to delete this {{ activityToDelete?.type }} activity? This action cannot be undone.
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn @click="showDeleteDialog = false">Cancel</v-btn>
          <v-btn color="error" @click="deleteActivity" :loading="loading"> Delete </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Activity Detail Dialog -->
    <v-dialog v-model="showDetailDialog" max-width="500" scrollable>
      <v-card v-if="selectedActivity">
        <v-card-title class="d-flex align-center">
          <v-icon :color="getActivityColor(selectedActivity.type)" class="mr-2">
            {{ getActivityIcon(selectedActivity.type) }}
          </v-icon>
          {{ getActivityTitle(selectedActivity.type) }}
          <v-spacer></v-spacer>
          <v-btn icon variant="text" @click="showDetailDialog = false">
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-card-title>

        <v-divider></v-divider>

        <v-card-text class="pa-4">
          <ActivityDetails :activity="selectedActivity" :detailed="true" />

          <v-divider v-if="selectedActivity.notes" class="my-3"></v-divider>

          <div v-if="selectedActivity.notes">
            <h4 class="text-subtitle-2 mb-2">Notes</h4>
            <p class="text-body-2">{{ selectedActivity.notes }}</p>
          </div>
        </v-card-text>

        <v-card-actions>
          <v-btn color="primary" variant="text" @click="editActivity(selectedActivity)"> Edit </v-btn>
          <v-spacer></v-spacer>
          <v-btn color="error" variant="text" @click="confirmDelete(selectedActivity)"> Delete </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Edit Activity Dialog -->
    <v-dialog v-model="showEditDialog" max-width="500" persistent scrollable>
      <v-card v-if="activityToEdit">
        <v-card-title class="d-flex align-center">
          <v-icon :color="getActivityColor(activityToEdit.type)" class="mr-2">
            {{ getActivityIcon(activityToEdit.type) }}
          </v-icon>
          <span>Edit {{ getActivityTitle(activityToEdit.type) }}</span>
          <v-spacer></v-spacer>
          <v-btn icon variant="text" @click="closeEditDialog">
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-card-title>

        <v-divider></v-divider>

        <v-card-text class="pa-4">
          <!-- Show warning for timer activities -->
          <v-alert
            v-if="isTimerActivity(activityToEdit) && !activityToEdit.end_time"
            type="warning"
            variant="tonal"
            class="mb-4"
          >
            This activity is currently running. You can only edit notes and some details.
          </v-alert>

          <!-- Dynamic edit form component -->
          <component
            v-if="currentEditFormComponent"
            :is="currentEditFormComponent"
            :activity="activityToEdit"
            :edit-mode="true"
            :has-timer="getActivityTypeConfig(activityToEdit.type)?.hasTimer"
            @success="handleEditSuccess"
            @cancel="closeEditDialog"
          />
        </v-card-text>
      </v-card>
    </v-dialog>

    <!-- Success Snackbar -->
    <v-snackbar v-model="showSuccess" color="success" :timeout="3000" location="top">
      <v-icon start>mdi-check-circle</v-icon>
      {{ successMessage }}
    </v-snackbar>
  </v-container>
</template>

<script setup>
import { ref, computed, onMounted, markRaw } from "vue";
import { useActivityStore } from "@/stores/activity";
import { storeToRefs } from "pinia";
import { formatActivityDate } from "@/utils/datetime";
import { format, subDays, startOfWeek, startOfMonth } from "date-fns";

// Import the activity details component
import ActivityDetails from "@/components/activity/ActivityDetails.vue";

// Import form components for editing
import FeedForm from "@/components/forms/FeedForm.vue";
import PumpForm from "@/components/forms/PumpForm.vue";
import DiaperForm from "@/components/forms/DiaperForm.vue";
import SleepForm from "@/components/forms/SleepForm.vue";
import GrowthForm from "@/components/forms/GrowthForm.vue";
import HealthForm from "@/components/forms/HealthForm.vue";
import MilestoneForm from "@/components/forms/MilestoneForm.vue";

// Mark components as raw to avoid reactivity overhead
const editFormComponents = {
  feed: markRaw(FeedForm),
  pump: markRaw(PumpForm),
  diaper: markRaw(DiaperForm),
  sleep: markRaw(SleepForm),
  growth: markRaw(GrowthForm),
  health: markRaw(HealthForm),
  milestone: markRaw(MilestoneForm),
};

const activityStore = useActivityStore();
const { activities, loading, pagination, activityTypes } = storeToRefs(activityStore);

// State
const showFilters = ref(false);
const showDeleteDialog = ref(false);
const showDetailDialog = ref(false);
const showEditDialog = ref(false);
const showSuccess = ref(false);
const successMessage = ref("");
const activityToDelete = ref(null);
const selectedActivity = ref(null);
const activityToEdit = ref(null);

// Filters
const filters = ref({
  startDate: null,
  endDate: null,
  type: null,
});

// Activity type options for filter
const activityTypeOptions = computed(() => {
  return activityTypes.value.map((type) => ({
    title: type.title,
    value: type.id,
  }));
});

// Check if any filters are active
const hasActiveFilters = computed(() => {
  return filters.value.startDate || filters.value.endDate || filters.value.type;
});

// Current edit form component
const currentEditFormComponent = computed(() => {
  if (!activityToEdit.value) return null;
  return editFormComponents[activityToEdit.value.type] || null;
});

// Date presets for quick filtering
const datePresets = [
  {
    key: "today",
    label: "Today",
    getDates: () => {
      const today = format(new Date(), "yyyy-MM-dd");
      return { startDate: today, endDate: today };
    },
  },
  {
    key: "yesterday",
    label: "Yesterday",
    getDates: () => {
      const yesterday = format(subDays(new Date(), 1), "yyyy-MM-dd");
      return { startDate: yesterday, endDate: yesterday };
    },
  },
  {
    key: "last7days",
    label: "Last 7 Days",
    getDates: () => ({
      startDate: format(subDays(new Date(), 6), "yyyy-MM-dd"),
      endDate: format(new Date(), "yyyy-MM-dd"),
    }),
  },
  {
    key: "thisweek",
    label: "This Week",
    getDates: () => ({
      startDate: format(startOfWeek(new Date(), { weekStartsOn: 1 }), "yyyy-MM-dd"),
      endDate: format(new Date(), "yyyy-MM-dd"),
    }),
  },
  {
    key: "thismonth",
    label: "This Month",
    getDates: () => ({
      startDate: format(startOfMonth(new Date()), "yyyy-MM-dd"),
      endDate: format(new Date(), "yyyy-MM-dd"),
    }),
  },
];

// Methods
function getActivityColor(type) {
  const config = activityTypes.value.find((at) => at.id === type);
  return config?.color || "grey";
}

function getActivityIcon(type) {
  const config = activityTypes.value.find((at) => at.id === type);
  return config?.icon || "mdi-circle";
}

function getActivityTitle(type) {
  const config = activityTypes.value.find((at) => at.id === type);
  return config?.title || type;
}

function getActivityTypeConfig(type) {
  return activityTypes.value.find((at) => at.id === type);
}

function isTimerActivity(activity) {
  const config = getActivityTypeConfig(activity.type);
  if (!config?.hasTimer) {
    return false;
  }

  // For feeds, only breast feeding is a timer activity. Bottle feeds are not.
  if (activity.type === "feed") {
    return activity.source === "left" || activity.source === "right";
  }

  return true;
}

function applyDatePreset(preset) {
  const dates = preset.getDates();
  filters.value.startDate = dates.startDate;
  filters.value.endDate = dates.endDate;
  applyFilters();
}

async function applyFilters() {
  const params = {};

  if (filters.value.startDate) {
    params.start_date = filters.value.startDate;
  }
  if (filters.value.endDate) {
    params.end_date = filters.value.endDate;
  }
  if (filters.value.type) {
    params.type = filters.value.type;
  }

  await activityStore.fetchActivities(params);
}

function clearFilters() {
  filters.value = {
    startDate: null,
    endDate: null,
    type: null,
  };
  activityStore.fetchActivities();
}

async function loadMore() {
  const params = {
    page: pagination.value.page + 1,
  };

  // Add current filters
  if (filters.value.startDate) params.start_date = filters.value.startDate;
  if (filters.value.endDate) params.end_date = filters.value.endDate;
  if (filters.value.type) params.type = filters.value.type;

  await activityStore.fetchActivities(params);
}

function viewActivity(activity) {
  selectedActivity.value = activity;
  showDetailDialog.value = true;
}

function editActivity(activity) {
  // Close any open dialogs
  showDetailDialog.value = false;
  showDeleteDialog.value = false;

  // Set the activity to edit
  activityToEdit.value = { ...activity }; // Clone to avoid mutations
  showEditDialog.value = true;
}

function closeEditDialog() {
  showEditDialog.value = false;
  // Clear after animation
  setTimeout(() => {
    activityToEdit.value = null;
  }, 300);
}

async function handleEditSuccess() {
  // Close edit dialog
  closeEditDialog();

  // Show success message
  successMessage.value = "Activity updated successfully";
  showSuccess.value = true;

  // Refresh the activities list to show updated data
  await applyFilters();
}

function confirmDelete(activity) {
  activityToDelete.value = activity;
  showDetailDialog.value = false;
  showDeleteDialog.value = true;
}

async function deleteActivity() {
  if (!activityToDelete.value) return;

  const result = await activityStore.deleteActivity(activityToDelete.value.id);

  if (result.success) {
    successMessage.value = "Activity deleted successfully";
    showSuccess.value = true;
  }

  showDeleteDialog.value = false;
  activityToDelete.value = null;
}

// Load activities on mount
onMounted(() => {
  activityStore.fetchActivities();
});
</script>

<style scoped>
.activity-list .v-card {
  cursor: pointer;
  transition: all 0.2s;
}

.activity-list .v-card:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.activity-details {
  font-size: 0.875rem;
}
</style>
