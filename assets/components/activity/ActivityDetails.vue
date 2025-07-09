<template>
  <div class="activity-details">
    <!-- Feed Activity -->
    <div v-if="activity.type === 'feed' && activity.feed_data">
      <div class="d-flex align-center mb-1">
        <v-chip size="small" :color="getFeedTypeColor(activity.feed_data.feed_type)" class="mr-2">
          {{ getFeedTypeLabel(activity.feed_data.feed_type) }}
        </v-chip>
        <span v-if="activity.feed_data.amount_ml" class="text-body-2"> {{ activity.feed_data.amount_ml }}ml </span>
        <span v-if="activity.feed_data.duration_minutes" class="text-body-2 ml-2">
          {{ activity.feed_data.duration_minutes }} min
        </span>
      </div>

      <div v-if="detailed && activity.end_time" class="text-caption text-grey">
        Duration: {{ formatDurationFromTimes(activity.start_time, activity.end_time) }}
      </div>
    </div>

    <!-- Pump Activity -->
    <div v-else-if="activity.type === 'pump' && activity.pump_data">
      <div class="d-flex align-center mb-1">
        <v-chip size="small" color="pump" class="mr-2">
          {{ getPumpBreastLabel(activity.pump_data.breast) }}
        </v-chip>
        <span v-if="activity.pump_data.amount_ml" class="text-body-2"> {{ activity.pump_data.amount_ml }}ml </span>
        <span v-if="activity.pump_data.duration_minutes" class="text-body-2 ml-2">
          {{ activity.pump_data.duration_minutes }} min
        </span>
      </div>

      <div v-if="detailed && activity.end_time" class="text-caption text-grey">
        Duration: {{ formatDurationFromTimes(activity.start_time, activity.end_time) }}
      </div>
    </div>

    <!-- Diaper Activity -->
    <div v-else-if="activity.type === 'diaper' && activity.diaper_data">
      <div class="d-flex align-center mb-1">
        <v-chip v-if="activity.diaper_data.wet" size="small" color="blue" class="mr-1"> Wet </v-chip>
        <v-chip v-if="activity.diaper_data.dirty" size="small" color="brown" class="mr-1"> Dirty </v-chip>
      </div>

      <div
        v-if="detailed && (activity.diaper_data.color || activity.diaper_data.consistency)"
        class="text-caption text-grey"
      >
        <span v-if="activity.diaper_data.color">Color: {{ activity.diaper_data.color }}</span>
        <span v-if="activity.diaper_data.consistency" class="ml-2">
          Consistency: {{ activity.diaper_data.consistency }}
        </span>
      </div>
    </div>

    <!-- Sleep Activity -->
    <div v-else-if="activity.type === 'sleep' && activity.sleep_data">
      <div class="d-flex align-center mb-1">
        <span class="text-body-2">
          {{ activity.sleep_data.location || "Sleep" }}
        </span>
        <v-rating
          v-if="activity.sleep_data.quality"
          :model-value="activity.sleep_data.quality"
          readonly
          size="x-small"
          color="yellow-darken-2"
          class="ml-2"
        />
      </div>

      <div class="text-caption text-grey">
        <span v-if="activity.end_time">
          Duration: {{ formatDurationFromTimes(activity.start_time, activity.end_time) }}
        </span>
        <span v-else class="text-success"> Currently sleeping • {{ formatTimeAgo(activity.start_time) }} </span>
      </div>
    </div>

    <!-- Growth Activity -->
    <div v-else-if="activity.type === 'growth' && activity.growth_data">
      <div class="growth-measurements">
        <div v-if="activity.growth_data.weight_kg" class="d-flex align-center mb-1">
          <v-icon size="small" class="mr-1">mdi-scale</v-icon>
          <span class="text-body-2">{{ activity.growth_data.weight_kg }}kg</span>
        </div>

        <div v-if="activity.growth_data.height_cm" class="d-flex align-center mb-1">
          <v-icon size="small" class="mr-1">mdi-human-male-height</v-icon>
          <span class="text-body-2">{{ activity.growth_data.height_cm }}cm</span>
        </div>

        <div v-if="activity.growth_data.head_circumference_cm" class="d-flex align-center">
          <v-icon size="small" class="mr-1">mdi-head</v-icon>
          <span class="text-body-2">{{ activity.growth_data.head_circumference_cm }}cm head</span>
        </div>
      </div>
    </div>

    <!-- Health Activity -->
    <div v-else-if="activity.type === 'health' && activity.health_data">
      <div class="d-flex align-center mb-1">
        <v-chip size="small" :color="getHealthTypeColor(activity.health_data.record_type)" class="mr-2">
          {{ getHealthTypeLabel(activity.health_data.record_type) }}
        </v-chip>
        <span v-if="activity.health_data.vaccine_name" class="text-body-2">
          {{ activity.health_data.vaccine_name }}
        </span>
      </div>

      <div v-if="detailed && activity.health_data.provider" class="text-caption text-grey">
        Provider: {{ activity.health_data.provider }}
      </div>

      <div v-if="detailed && activity.health_data.symptoms" class="mt-2">
        <h5 class="text-subtitle-2">Symptoms</h5>
        <p class="text-body-2">{{ activity.health_data.symptoms }}</p>
      </div>

      <div v-if="detailed && activity.health_data.treatment" class="mt-2">
        <h5 class="text-subtitle-2">Treatment</h5>
        <p class="text-body-2">{{ activity.health_data.treatment }}</p>
      </div>
    </div>

    <!-- Milestone Activity -->
    <div v-else-if="activity.type === 'milestone' && activity.milestone_data">
      <div class="d-flex align-center mb-1">
        <v-chip size="small" color="milestone" class="mr-2">
          {{ activity.milestone_data.milestone_type }}
        </v-chip>
      </div>

      <div v-if="activity.milestone_data.description" class="text-body-2 mt-1">
        {{ activity.milestone_data.description }}
      </div>
    </div>

    <!-- Fallback for unknown activity types -->
    <div v-else class="text-caption text-grey">{{ activity.type }} activity</div>

    <!-- Time information (always shown in detailed view) -->
    <div v-if="detailed" class="mt-3 pt-2 border-t">
      <div class="text-caption text-grey">
        <div>Started: {{ formatDetailedTime(activity.start_time) }}</div>
        <div v-if="activity.end_time">Ended: {{ formatDetailedTime(activity.end_time) }}</div>
        <div v-if="activity.created_at !== activity.updated_at">
          Updated: {{ formatDetailedTime(activity.updated_at) }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { formatDuration, formatTimeAgo } from "@/utils/datetime";
import { format, parseISO } from "date-fns";

defineProps({
  activity: {
    type: Object,
    required: true,
  },
  detailed: {
    type: Boolean,
    default: false,
  },
});

// Feed type helpers
function getFeedTypeLabel(type) {
  const labels = {
    bottle: "Bottle",
    breast_left: "Left Breast",
    breast_right: "Right Breast",
    solid: "Solid Food",
  };
  return labels[type] || type;
}

function getFeedTypeColor(type) {
  const colors = {
    bottle: "blue",
    breast_left: "pink",
    breast_right: "pink",
    solid: "orange",
  };
  return colors[type] || "grey";
}

// Pump helpers
function getPumpBreastLabel(breast) {
  const labels = {
    left: "Left",
    right: "Right",
    both: "Both Breasts",
  };
  return labels[breast] || breast;
}

// Health type helpers
function getHealthTypeLabel(type) {
  const labels = {
    checkup: "Checkup",
    vaccine: "Vaccine",
    illness: "Illness",
  };
  return labels[type] || type;
}

function getHealthTypeColor(type) {
  const colors = {
    checkup: "green",
    vaccine: "blue",
    illness: "red",
  };
  return colors[type] || "grey";
}

// Time formatting helpers
function formatDurationFromTimes(startTime, endTime) {
  const start = new Date(startTime);
  const end = new Date(endTime);
  const durationMinutes = Math.floor((end - start) / (1000 * 60));
  return formatDuration(durationMinutes);
}

function formatDetailedTime(timeString) {
  const date = parseISO(timeString);
  return format(date, "MMM d, yyyy h:mm a");
}
</script>

<style scoped>
.activity-details {
  font-size: 0.875rem;
}

.growth-measurements .v-icon {
  opacity: 0.7;
}

.border-t {
  border-top: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}
</style>
