<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-4">Trends &amp; Statistics</h1>
      </v-col>

      <!-- Last Feed -->
      <v-col cols="12" sm="6" md="3">
        <v-card elevation="1" class="pa-4 text-center">
          <v-icon icon="mdi-baby-bottle" size="36" class="mb-2" color="accent1" />
          <div class="text-subtitle-1">Last fed</div>
          <div class="text-h5 font-weight-bold">
            <span v-if="statsStore.loading">…</span>
            <span v-else>{{ lastFedDisplay }}</span>
          </div>
          <div v-if="lastFeedAmount" class="text-caption text-medium-emphasis">
            {{ lastFeedAmount }}
          </div>
        </v-card>
      </v-col>

      <!-- Diapers today -->
      <v-col cols="12" sm="6" md="3">
        <v-card elevation="1" class="pa-4 text-center">
          <v-icon icon="mdi-toilet" size="36" class="mb-2" color="accent2" />
          <div class="text-subtitle-1">Diapers today</div>
          <div class="text-h5 font-weight-bold">
            <span v-if="statsStore.loading">…</span>
            <span v-else>{{ diapersToday }}</span>
          </div>
        </v-card>
      </v-col>

      <!-- Sleeping -->
      <v-col cols="12" sm="6" md="3">
        <v-card elevation="1" class="pa-4 text-center">
          <v-icon icon="mdi-sleep" size="36" class="mb-2" color="accent1" />
          <div class="text-subtitle-1">Currently sleeping</div>
          <div class="text-h5 font-weight-bold">
            <span v-if="statsStore.loading">…</span>
            <span v-else>{{ sleepingDisplay }}</span>
          </div>
        </v-card>
      </v-col>

      <!-- Weekly activities -->
      <v-col cols="12" sm="6" md="3">
        <v-card elevation="1" class="pa-4 text-center">
          <v-icon icon="mdi-calendar-week" size="36" class="mb-2" color="accent2" />
          <div class="text-subtitle-1">Activities this week</div>
          <div class="text-h5 font-weight-bold">
            <span v-if="statsStore.loading">…</span>
            <span v-else>{{ activitiesWeek }}</span>
          </div>
        </v-card>
      </v-col>
    </v-row>

    <!-- Weekly Averages -->
    <v-row>
      <!-- Avg Counts -->
      <v-col cols="12" md="6">
        <v-card elevation="1">
          <v-card-title>Average Daily Counts</v-card-title>
          <v-card-text>
            <div v-if="statsStore.loading" class="text-center py-10">
              <v-progress-circular indeterminate color="primary" />
            </div>
            <div v-else-if="statsStore.error" class="text-center text-error py-10">
              {{ statsStore.error }}
            </div>
            <div v-else style="height: 250px">
              <BarChart :chart-data="avgCountsChartData" :chart-options="chartOptions" />
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Avg Sleep -->
      <v-col cols="12" md="6">
        <v-card elevation="1" class="fill-height">
          <v-card-title>Average Daily Sleep</v-card-title>
          <v-card-text class="d-flex align-center justify-center">
            <div v-if="statsStore.loading" class="text-center">
              <v-progress-circular indeterminate color="primary" />
            </div>
            <div v-else-if="statsStore.error" class="text-center text-error">
              {{ statsStore.error }}
            </div>
            <div v-else class="text-h3 font-weight-bold text-center" style="color: #f59e0b">
              {{ avgSleep }}
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Weekly Sleep Trend Chart -->
    <v-row>
      <v-col cols="12">
        <v-card elevation="1">
          <v-card-title>Weekly Sleep Trend</v-card-title>
          <v-card-text>
            <div v-if="statsStore.loading" class="text-center py-10">
              <v-progress-circular indeterminate color="primary" />
            </div>
            <div v-else-if="statsStore.error" class="text-center text-error py-10">
              {{ statsStore.error }}
            </div>
            <div v-else style="height: 300px">
              <LineChart :chart-data="weeklySleepTrendData" :chart-options="chartOptions" />
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { computed, onMounted } from "vue";
import { useStatsStore } from "@/stores/stats";
import { formatTimeAgo } from "@/utils/datetime";
import BarChart from "@/components/charts/BarChart.vue";
import LineChart from "@/components/charts/LineChart.vue";

const statsStore = useStatsStore();

// Fetch stats on mount (if not already loaded)
onMounted(() => {
  if (!statsStore.lastUpdated) {
    statsStore.fetchStats();
  }
});

// --- Chart Data -------------------------------------------------------

const avgCountsChartData = computed(() => {
  const averages = statsStore.weekly?.daily_averages;
  if (!averages) return { labels: [], datasets: [] };

  const labels = ["Diapers", "Feeds"];
  const data = [averages.diaper_per_day?.toFixed(1) || 0, averages.feed_per_day?.toFixed(1) || 0];

  return {
    labels,
    datasets: [
      {
        label: "Daily Average",
        backgroundColor: ["#6366f1", "#ec4899"],
        data,
      },
    ],
  };
});

const avgSleep = computed(() => {
  const hours = statsStore.weekly?.daily_averages?.sleep_hours_per_day;
  if (hours === undefined || hours === null) return "—";
  return formatHoursMinutes(hours);
});

const weeklySleepTrendData = computed(() => {
  const breakdown = statsStore.weekly?.daily_breakdown;
  if (!breakdown || breakdown.length === 0) {
    return { labels: [], datasets: [] };
  }

  const labels = breakdown.map((d) => new Date(d.date).toLocaleDateString(undefined, { weekday: "short" }));
  const data = breakdown.map((d) => d.sleep_duration_hours.toFixed(1));

  return {
    labels,
    datasets: [
      {
        label: "Sleep Hours",
        borderColor: "#f59e0b",
        backgroundColor: "rgba(245, 158, 11, 0.2)",
        tension: 0.1,
        fill: true,
        data,
      },
    ],
  };
});

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: false,
    },
  },
  scales: {
    y: {
      beginAtZero: true,
      grid: {
        color: "rgba(255, 255, 255, 0.1)",
      },
    },
    x: {
      grid: {
        display: false,
      },
    },
  },
};

// --- Computed helpers -------------------------------------------------

const lastFedDisplay = computed(() => {
  const feed = statsStore.recent?.last_feed;
  if (!feed) return "—";
  return formatTimeAgo(feed.time);
});

const lastFeedAmount = computed(() => {
  const feed = statsStore.recent?.last_feed;
  if (!feed?.amount_ml) return "";
  return `${feed.amount_ml} ml`;
});

const diapersToday = computed(() => {
  return statsStore.daily?.counts?.diaper ?? "0";
});

function formatHoursMinutes(hoursFloat) {
  if (hoursFloat === undefined || hoursFloat === null) return "";
  const totalMinutes = Math.round(hoursFloat * 60);
  const h = Math.floor(totalMinutes / 60);
  const m = totalMinutes % 60;
  return `${h > 0 ? h + "h " : ""}${m}m`;
}

const sleepingDisplay = computed(() => {
  if (!statsStore.recent) return "—";

  // Ongoing sleep
  if (statsStore.recent.currently_sleeping) {
    const startISO = statsStore.daily?.last_activities?.sleep;
    if (startISO) {
      const diffMs = Date.now() - new Date(startISO).getTime();
      const minutes = Math.floor(diffMs / 60000);
      const h = Math.floor(minutes / 60);
      const m = minutes % 60;
      return `${h > 0 ? h + "h " : ""}${m}m`;
    }
    return "Sleeping";
  }

  // Not sleeping → last sleep duration
  const dur = statsStore.recent.last_sleep?.duration_hours;
  if (dur !== undefined && dur !== null) {
    return formatHoursMinutes(dur);
  }
  return "—";
});

const activitiesWeek = computed(() => {
  const avgs = statsStore.weekly?.daily_averages;
  if (!avgs) return "0";
  let total = 0;
  for (const [key, value] of Object.entries(avgs)) {
    if (key.endsWith("_per_day") && !key.includes("hours") && !key.includes("amount")) {
      total += value * 7;
    }
  }
  return Math.round(total);
});
</script>
