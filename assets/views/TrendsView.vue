<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-4">Trends &amp; Statistics</h1>
      </v-col>
    </v-row>

    <!-- Daily Totals -->
    <v-row>
      <v-col cols="12">
        <h2 class="text-h5 mb-4">Today's Summary</h2>
      </v-col>
    </v-row>
    <v-row class="mb-4">
      <!-- Feeds today -->
      <v-col cols="12" sm="6" md="4" lg="3">
        <v-card elevation="2" class="pa-5 text-center h-100 position-relative overflow-hidden">
          <!-- Background gradient -->
          <div class="position-absolute w-100 h-100" style="top: 0; left: 0; background: linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(147, 51, 234, 0.1) 100%); border-radius: inherit;"></div>

          <div class="position-relative">
            <v-icon icon="mdi-baby-bottle" size="40" class="mb-3" color="accent1" />
            <div class="text-subtitle-1 mb-2">Feeds today</div>
            <div class="text-h5 font-weight-bold mb-2">
              <span v-if="statsStore.loading">…</span>
              <span v-else>{{ feedsToday }}</span>
            </div>
            <div class="text-caption text-medium-emphasis">
              <span v-if="totalFeedAmountToday">{{ totalFeedAmountToday }} total</span>
              <span v-else>{{ feedsToday === "1" ? "feed" : "feeds" }}</span>
            </div>
          </div>
        </v-card>
      </v-col>

      <!-- Pumping today -->
      <v-col cols="12" sm="6" md="4" lg="3">
        <v-card elevation="2" class="pa-5 text-center h-100 position-relative overflow-hidden">
          <!-- Background gradient -->
          <div class="position-absolute w-100 h-100" style="top: 0; left: 0; background: linear-gradient(135deg, rgba(34, 197, 94, 0.1) 0%, rgba(59, 130, 246, 0.1) 100%); border-radius: inherit;"></div>

          <div class="position-relative">
            <v-icon icon="mdi-pump" size="40" class="mb-3" color="accent1" />
            <div class="text-subtitle-1 mb-2">Pumping today</div>
            <div class="text-h5 font-weight-bold mb-2">
              <span v-if="statsStore.loading">…</span>
              <span v-else>{{ pumpingToday }}</span>
            </div>
            <div class="text-caption text-medium-emphasis">
              <span v-if="totalPumpAmountToday">{{ totalPumpAmountToday }} total</span>
              <span v-else>{{ pumpingToday === "1" ? "session" : "sessions" }}</span>
            </div>
          </div>
        </v-card>
      </v-col>

      <!-- Diapers today -->
      <v-col cols="12" sm="6" md="4" lg="3">
        <v-card elevation="2" class="pa-5 text-center h-100 position-relative overflow-hidden">
          <!-- Background gradient -->
          <div class="position-absolute w-100 h-100" style="top: 0; left: 0; background: linear-gradient(135deg, rgba(99, 102, 241, 0.1) 0%, rgba(236, 72, 153, 0.1) 100%); border-radius: inherit;"></div>

          <div class="position-relative">
            <v-icon icon="mdi-toilet" size="40" class="mb-3" color="accent2" />
            <div class="text-subtitle-1 mb-2">Diapers today</div>
            <div class="text-h5 font-weight-bold mb-2">
              <span v-if="statsStore.loading">…</span>
              <span v-else>{{ diapersToday }}</span>
            </div>
            <div class="text-caption text-medium-emphasis">
              <span v-if="statsStore.loading">Loading...</span>
              <span v-else-if="diaperBreakdown">
                {{ diaperBreakdown.wet }} wet, {{ diaperBreakdown.dirty }} dirty
              </span>
              <span v-else>
                {{ diapersToday === "1" ? "diaper" : "diapers" }} changed
              </span>
            </div>
          </div>
        </v-card>
      </v-col>

      <!-- Sleep today -->
      <v-col cols="12" sm="6" md="4" lg="3">
        <v-card elevation="2" class="pa-5 text-center h-100 position-relative overflow-hidden">
          <!-- Background gradient -->
          <div class="position-absolute w-100 h-100" style="top: 0; left: 0; background: linear-gradient(135deg, rgba(168, 85, 247, 0.1) 0%, rgba(236, 72, 153, 0.1) 100%); border-radius: inherit;"></div>

          <div class="position-relative">
            <v-icon icon="mdi-sleep" size="40" class="mb-3" color="accent2" />
            <div class="text-subtitle-1 mb-2">Sleep today</div>
            <div class="text-h5 font-weight-bold mb-2">
              <span v-if="statsStore.loading">…</span>
              <span v-else>{{ sleepToday }}</span>
            </div>
            <div v-if="sleepSessionsToday > 0" class="text-caption text-medium-emphasis">
              {{ sleepSessionsToday }} session{{ sleepSessionsToday > 1 ? 's' : '' }}
            </div>
          </div>
        </v-card>
      </v-col>

      <!-- Currently sleeping -->
      <v-col cols="12" sm="6" md="4" lg="3">
        <v-card elevation="2" class="pa-5 text-center h-100 position-relative overflow-hidden">
          <!-- Background gradient -->
          <div class="position-absolute w-100 h-100" style="top: 0; left: 0; background: linear-gradient(135deg, rgba(245, 158, 11, 0.1) 0%, rgba(168, 85, 247, 0.1) 100%); border-radius: inherit;"></div>

          <div class="position-relative">
            <v-icon icon="mdi-sleep" size="40" class="mb-3" color="accent1" />
            <div class="text-subtitle-1 mb-2">{{ sleepingTitle }}</div>
            <div class="text-h5 font-weight-bold mb-2">
              <span v-if="statsStore.loading">…</span>
              <span v-else>{{ sleepingDisplay }}</span>
            </div>
            <div v-if="lastSleepDurationDisplay" class="text-caption text-medium-emphasis">
              {{ lastSleepDurationDisplay }}
            </div>
          </div>
        </v-card>
      </v-col>

      <!-- Last Feed -->
      <v-col cols="12" sm="6" md="4" lg="3">
        <v-card elevation="2" class="pa-5 text-center h-100 position-relative overflow-hidden">
          <!-- Background gradient -->
          <div class="position-absolute w-100 h-100" style="top: 0; left: 0; background: linear-gradient(135deg, rgba(34, 197, 94, 0.1) 0%, rgba(245, 158, 11, 0.1) 100%); border-radius: inherit;"></div>

          <div class="position-relative">
            <v-icon icon="mdi-baby-bottle" size="40" class="mb-3" color="accent1" />
            <div class="text-subtitle-1 mb-2">Last fed</div>
            <div class="text-h5 font-weight-bold mb-2">
              <span v-if="statsStore.loading">…</span>
              <span v-else>{{ lastFedDisplay }}</span>
            </div>
            <div v-if="lastFeedAmount" class="text-caption text-medium-emphasis">
              {{ lastFeedAmount }}
            </div>
          </div>
        </v-card>
      </v-col>
    </v-row>

    <!-- Weekly Averages -->
    <v-row class="mt-6">
      <v-col cols="12">
        <h2 class="text-h5 mb-4">Weekly Trends</h2>
      </v-col>
    </v-row>
    <v-row class="mb-4">
      <!-- Weekly activities -->
      <v-col cols="12" md="4">
        <v-card elevation="2" class="pa-6 text-center h-100 position-relative overflow-hidden">
          <!-- Background gradient -->
          <div class="position-absolute w-100 h-100" style="top: 0; left: 0; background: linear-gradient(135deg, rgba(236, 72, 153, 0.1) 0%, rgba(99, 102, 241, 0.1) 100%); border-radius: inherit;"></div>

          <!-- Background decorative icon -->
          <div class="position-absolute" style="top: -10px; right: -10px; opacity: 0.08; transform: rotate(15deg);">
            <v-icon icon="mdi-chart-line" size="120" color="accent2" />
          </div>

          <div class="position-relative">
            <v-icon icon="mdi-calendar-week" size="48" class="mb-4" color="accent2" />
            <div class="text-h6 mb-3 font-weight-medium">Activities this week</div>
            <div class="text-h3 font-weight-bold mb-3" style="color: #ec4899;">
              <span v-if="statsStore.loading">…</span>
              <span v-else>{{ activitiesWeek }}</span>
            </div>
            <div class="text-body-2 text-medium-emphasis mb-2">
              total activities
            </div>
            <div class="text-caption text-medium-emphasis">
              <v-icon icon="mdi-trending-up" size="16" class="mr-1" />
              {{ Math.round((activitiesWeek / 7) * 10) / 10 }} per day average
            </div>
          </div>
        </v-card>
      </v-col>

      <!-- Avg Sleep -->
      <v-col cols="12" md="4">
        <v-card elevation="2" class="pa-6 text-center h-100 position-relative overflow-hidden">
          <!-- Background gradient -->
          <div class="position-absolute w-100 h-100" style="top: 0; left: 0; background: linear-gradient(135deg, rgba(245, 158, 11, 0.1) 0%, rgba(168, 85, 247, 0.1) 100%); border-radius: inherit;"></div>

          <!-- Background decorative icon -->
          <div class="position-absolute" style="top: -15px; right: -15px; opacity: 0.08; transform: rotate(-15deg);">
            <v-icon icon="mdi-sleep" size="140" color="amber" />
          </div>

          <div class="position-relative">
            <v-icon icon="mdi-sleep" size="48" class="mb-4" color="amber" />
            <div class="text-h6 mb-3 font-weight-medium">Average Daily Sleep</div>
            <div v-if="statsStore.loading" class="text-center py-4">
              <v-progress-circular indeterminate color="amber" size="40" />
            </div>
            <div v-else-if="statsStore.error" class="text-center text-error py-4">
              <v-icon icon="mdi-alert-circle" class="mb-2" />
              <div>{{ statsStore.error }}</div>
            </div>
            <div v-else>
              <div class="text-h2 font-weight-bold mb-3" style="color: #f59e0b">
                {{ avgSleep }}
              </div>
              <div class="text-body-2 text-medium-emphasis mb-2">
                per night
              </div>
              <div class="text-caption text-medium-emphasis">
                <v-icon icon="mdi-moon-waning-crescent" size="16" class="mr-1" />
                Based on this week's data
              </div>
            </div>
          </div>
        </v-card>
      </v-col>

      <!-- Avg Counts -->
      <v-col cols="12" md="4">
        <v-card elevation="2" class="pa-6 h-100 position-relative overflow-hidden">
          <!-- Background gradient -->
          <div class="position-absolute w-100 h-100" style="top: 0; left: 0; background: linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(34, 197, 94, 0.1) 100%); border-radius: inherit;"></div>

          <!-- Background decorative icon -->
          <div class="position-absolute" style="top: -10px; left: -10px; opacity: 0.08; transform: rotate(-10deg);">
            <v-icon icon="mdi-chart-bar" size="130" color="primary" />
          </div>

          <div class="position-relative">
            <div class="text-center mb-4">
              <v-icon icon="mdi-chart-bar" size="48" color="primary" />
            </div>
            <div class="text-h6 mb-4 font-weight-medium text-center">Average Daily Counts</div>

            <div v-if="statsStore.loading" class="text-center py-10">
              <v-progress-circular indeterminate color="primary" size="40" />
            </div>
            <div v-else-if="statsStore.error" class="text-center text-error py-10">
              <v-icon icon="mdi-alert-circle" class="mb-2" />
              <div>{{ statsStore.error }}</div>
            </div>
            <div v-else style="height: 200px">
              <BarChart :chart-data="avgCountsChartData" :chart-options="chartOptions" />
            </div>
          </div>
        </v-card>
      </v-col>
    </v-row>

    <!-- Weekly Sleep Trend Chart -->
    <v-row class="mb-4">
      <v-col cols="12">
        <v-card elevation="2" class="position-relative overflow-hidden">
          <!-- Background gradient -->
          <div class="position-absolute w-100 h-100" style="top: 0; left: 0; background: linear-gradient(135deg, rgba(245, 158, 11, 0.05) 0%, rgba(236, 72, 153, 0.05) 100%); border-radius: inherit;"></div>

          <div class="position-relative">
            <v-card-title class="text-h6 pb-3">Weekly Sleep Trend</v-card-title>
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
          </div>
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

// Fetch stats on mount
onMounted(() => {
  statsStore.fetchStats();
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

const diaperBreakdown = computed(() => {
  return statsStore.daily?.diaper_breakdown || null;
});

const feedsToday = computed(() => {
  return statsStore.daily?.counts?.feed ?? "0";
});

const totalFeedAmountToday = computed(() => {
  const amount = statsStore.daily?.totals?.feed_amount_ml;
  if (!amount || amount === 0) return "";
  return `${Math.round(amount)} ml`;
});

const sleepToday = computed(() => {
  const hours = statsStore.daily?.totals?.sleep_hours;
  if (!hours || hours === 0) return "0h";
  return formatHoursMinutes(hours);
});

const sleepSessionsToday = computed(() => {
  return statsStore.daily?.counts?.sleep ?? 0;
});

const pumpingToday = computed(() => {
  return statsStore.daily?.counts?.pump ?? "0";
});

const totalPumpAmountToday = computed(() => {
  const amount = statsStore.daily?.totals?.pump_amount_ml;
  if (!amount || amount === 0) return "";
  return `${Math.round(amount)} ml`;
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
      const hours = diffMs / 3600000; // ms in an hour
      return formatHoursMinutes(hours);
    }
    return "Ongoing";
  }

  // Not sleeping → time ago it ended
  const lastSleep = statsStore.recent.last_sleep;
  if (lastSleep?.ended) {
    return formatTimeAgo(lastSleep.ended);
  }
  return "—";
});

const lastSleepDurationDisplay = computed(() => {
  if (statsStore.recent?.currently_sleeping) return "";
  const dur = statsStore.recent?.last_sleep?.duration_hours;
  if (dur !== undefined && dur !== null) {
    return formatHoursMinutes(dur);
  }
  return "";
});

const sleepingTitle = computed(() => {
  if (statsStore.recent?.currently_sleeping) {
    return "Currently sleeping";
  }
  return "Last sleep";
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
