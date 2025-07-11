<template>
  <v-container class="pa-3 pa-sm-4">
    <v-row>
      <v-col cols="12">
        <h1 class="text-h5 text-sm-h4 mb-3 mb-sm-4">Trends &amp; Statistics</h1>
      </v-col>
    </v-row>

    <!-- Daily Totals -->
    <v-row>
      <v-col cols="12">
        <!-- Mobile-friendly header layout -->
        <div class="mb-4">
          <h2 class="text-h5 mb-3">
            {{ dailySummaryTitle }}
          </h2>

          <!-- Date Selection Controls - Mobile optimized -->
          <div class="d-flex flex-column flex-sm-row align-start align-sm-center ga-3">
            <!-- Error message if any -->
            <v-alert
              v-if="statsStore.error"
              type="warning"
              density="compact"
              class="mb-2 text-caption w-100"
              closable
              @click:close="statsStore.error = null"
            >
              {{ statsStore.error }}
            </v-alert>

            <!-- Date navigation row -->
            <div class="d-flex align-center">
              <v-btn
                icon="mdi-chevron-left"
                variant="text"
                size="large"
                @click="statsStore.changeDay('prev')"
                :disabled="statsStore.dailyLoading"
              ></v-btn>

              <!-- Date display with date picker -->
              <v-menu v-model="datePickerMenu" :close-on-content-click="false">
                <template v-slot:activator="{ props }">
                  <v-btn
                    v-bind="props"
                    variant="text"
                    class="text-subtitle-1 text-sm-h6 font-weight-regular px-2 px-sm-4"
                    :loading="statsStore.dailyLoading"
                    size="large"
                  >
                    {{ formattedDailyDate }}
                    <v-icon icon="mdi-calendar" class="ml-2" size="20"></v-icon>
                  </v-btn>
                </template>

                <v-date-picker
                  v-model="selectedDate"
                  @update:model-value="onDateSelected"
                  :min="minDate"
                  :max="today"
                  show-adjacent-months
                  color="primary"
                >
                  <template v-slot:footer>
                    <div class="pa-3 text-caption text-medium-emphasis">
                      <v-icon icon="mdi-information-outline" size="16" class="mr-1"></v-icon>
                      <div>Date range:</div>
                      <div class="mt-1" v-if="statsStore.babyBirthDate">
                        • From: {{ formatDate(new Date(statsStore.babyBirthDate), "MMM d, yyyy") }}
                      </div>
                      <div>
                        • To: {{ formatDate(new Date(today), "MMM d, yyyy") }}
                      </div>
                    </div>
                  </template>
                </v-date-picker>
              </v-menu>

              <v-btn
                icon="mdi-chevron-right"
                variant="text"
                size="large"
                @click="statsStore.changeDay('next')"
                :disabled="statsStore.isViewingToday || statsStore.dailyLoading"
              ></v-btn>
            </div>

            <!-- Controls row -->
            <div class="d-flex align-center ga-2">
              <v-progress-circular
                v-if="statsStore.dailyLoading"
                indeterminate
                size="24"
              ></v-progress-circular>

              <!-- Quick jump to today -->
              <v-btn
                v-if="!statsStore.isViewingToday"
                variant="outlined"
                size="small"
                class="text-caption"
                @click="jumpToToday"
                :disabled="statsStore.dailyLoading"
              >
                Today
              </v-btn>
            </div>
          </div>
        </div>
      </v-col>
    </v-row>

    <!-- Daily Totals Section -->
    <v-row class="mb-4">
      <!-- Feeds today -->
      <v-col cols="6" sm="6" md="3">
        <v-card elevation="2" class="pa-4 pa-sm-5 text-center h-100 position-relative overflow-hidden">
          <!-- Background gradient -->
          <div class="position-absolute w-100 h-100" style="top: 0; left: 0; background: linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(147, 51, 234, 0.1) 100%); border-radius: inherit;"></div>

          <div class="position-relative">
            <v-icon icon="mdi-baby-bottle" :size="$vuetify.display.xs ? 32 : 40" class="mb-2 mb-sm-3" color="accent1" />
            <div class="text-body-2 text-sm-subtitle-1 mb-1 mb-sm-2">{{ feedsLabel }}</div>
            <div class="text-h6 text-sm-h5 font-weight-bold mb-1 mb-sm-2">
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
      <v-col cols="6" sm="6" md="3">
        <v-card elevation="2" class="pa-4 pa-sm-5 text-center h-100 position-relative overflow-hidden">
          <!-- Background gradient -->
          <div class="position-absolute w-100 h-100" style="top: 0; left: 0; background: linear-gradient(135deg, rgba(34, 197, 94, 0.1) 0%, rgba(59, 130, 246, 0.1) 100%); border-radius: inherit;"></div>

          <div class="position-relative">
            <v-icon icon="mdi-pump" :size="$vuetify.display.xs ? 32 : 40" class="mb-2 mb-sm-3" color="accent1" />
            <div class="text-body-2 text-sm-subtitle-1 mb-1 mb-sm-2">{{ pumpingLabel }}</div>
            <div class="text-h6 text-sm-h5 font-weight-bold mb-1 mb-sm-2">
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
      <v-col cols="6" sm="6" md="3">
        <v-card elevation="2" class="pa-4 pa-sm-5 text-center h-100 position-relative overflow-hidden">
          <!-- Background gradient -->
          <div class="position-absolute w-100 h-100" style="top: 0; left: 0; background: linear-gradient(135deg, rgba(99, 102, 241, 0.1) 0%, rgba(236, 72, 153, 0.1) 100%); border-radius: inherit;"></div>

          <div class="position-relative">
            <v-icon icon="mdi-toilet" :size="$vuetify.display.xs ? 32 : 40" class="mb-2 mb-sm-3" color="accent2" />
            <div class="text-body-2 text-sm-subtitle-1 mb-1 mb-sm-2">{{ diapersLabel }}</div>
            <div class="text-h6 text-sm-h5 font-weight-bold mb-1 mb-sm-2">
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
      <v-col cols="6" sm="6" md="3">
        <v-card elevation="2" class="pa-4 pa-sm-5 text-center h-100 position-relative overflow-hidden">
          <!-- Background gradient -->
          <div class="position-absolute w-100 h-100" style="top: 0; left: 0; background: linear-gradient(135deg, rgba(168, 85, 247, 0.1) 0%, rgba(236, 72, 153, 0.1) 100%); border-radius: inherit;"></div>

          <div class="position-relative">
            <v-icon icon="mdi-sleep" :size="$vuetify.display.xs ? 32 : 40" class="mb-2 mb-sm-3" color="accent2" />
            <div class="text-body-2 text-sm-subtitle-1 mb-1 mb-sm-2">{{ sleepLabel }}</div>
            <div class="text-h6 text-sm-h5 font-weight-bold mb-1 mb-sm-2">
              <span v-if="statsStore.loading">…</span>
              <span v-else>{{ sleepToday }}</span>
            </div>
            <div v-if="sleepSessionsToday > 0" class="text-caption text-medium-emphasis">
              {{ sleepSessionsToday }} session{{ sleepSessionsToday > 1 ? 's' : '' }}
            </div>
          </div>
        </v-card>
      </v-col>
    </v-row>

    <!-- Recent Activity Section -->
    <v-row class="mb-4">
      <v-col cols="12">
        <h2 class="text-h5 mb-3 mb-sm-4 d-flex align-center">
          <v-icon icon="mdi-clock-outline" size="20" class="mr-2" color="primary"></v-icon>
          Recent Activity
        </h2>
      </v-col>
    </v-row>
    <v-row class="mb-4">
      <!-- Currently sleeping / Last sleep -->
      <v-col cols="12" sm="6">
        <v-card
          elevation="1"
          class="pa-4 pa-sm-5 h-100 position-relative overflow-hidden"
          :class="{ 'border-thin': true }"
          style="border-color: rgba(var(--v-theme-primary), 0.2);"
        >
          <!-- Subtle background pattern -->
          <div class="position-absolute w-100 h-100" style="top: 0; left: 0; background: linear-gradient(135deg, rgba(245, 158, 11, 0.06) 0%, rgba(168, 85, 247, 0.06) 100%); border-radius: inherit;"></div>

          <div class="position-relative d-flex align-center">
            <!-- Status indicator -->
            <div class="mr-4">
              <div class="position-relative">
                <v-icon
                  icon="mdi-sleep"
                  :size="$vuetify.display.xs ? 40 : 48"
                  :color="statsStore.recent?.currently_sleeping ? 'success' : 'amber'"
                />
                <v-badge
                  v-if="statsStore.recent?.currently_sleeping"
                  dot
                  color="success"
                  class="position-absolute"
                  style="top: -2px; right: -2px;"
                ></v-badge>
              </div>
            </div>

            <!-- Content -->
            <div class="flex-grow-1">
              <div class="text-subtitle-1 text-sm-h6 mb-1 font-weight-medium">{{ sleepingTitle }}</div>
              <div class="text-h6 text-sm-h5 font-weight-bold mb-1" :class="statsStore.recent?.currently_sleeping ? 'text-success' : ''">
                <span v-if="statsStore.loading">…</span>
                <span v-else>{{ sleepingDisplay }}</span>
              </div>
              <div v-if="lastSleepDurationDisplay" class="text-caption text-medium-emphasis">
                Duration: {{ lastSleepDurationDisplay }}
              </div>
              <div v-else-if="statsStore.recent?.currently_sleeping" class="text-caption text-success">
                <v-icon icon="mdi-circle" size="8" class="mr-1"></v-icon>
                Sleep in progress
              </div>
            </div>
          </div>
        </v-card>
      </v-col>

      <!-- Last Feed -->
      <v-col cols="12" sm="6">
        <v-card
          elevation="1"
          class="pa-4 pa-sm-5 h-100 position-relative overflow-hidden"
          :class="{ 'border-thin': true }"
          style="border-color: rgba(var(--v-theme-primary), 0.2);"
        >
          <!-- Subtle background pattern -->
          <div class="position-absolute w-100 h-100" style="top: 0; left: 0; background: linear-gradient(135deg, rgba(34, 197, 94, 0.06) 0%, rgba(245, 158, 11, 0.06) 100%); border-radius: inherit;"></div>

          <div class="position-relative d-flex align-center">
            <!-- Icon -->
            <div class="mr-4">
              <v-icon icon="mdi-baby-bottle" :size="$vuetify.display.xs ? 40 : 48" color="accent1" />
            </div>

            <!-- Content -->
            <div class="flex-grow-1">
              <div class="text-subtitle-1 text-sm-h6 mb-1 font-weight-medium">Last fed</div>
              <div class="text-h6 text-sm-h5 font-weight-bold mb-1">
                <span v-if="statsStore.loading">…</span>
                <span v-else>{{ lastFedDisplay }}</span>
              </div>
              <div v-if="lastFeedAmount" class="text-caption text-medium-emphasis">
                Amount: {{ lastFeedAmount }}
              </div>
            </div>
          </div>
        </v-card>
      </v-col>
    </v-row>

    <!-- Past 7 Days Averages -->
    <v-row class="mt-4 mt-sm-6">
      <v-col cols="12">
        <h2 class="text-h5 mb-3 mb-sm-4">Past 7 Days</h2>
      </v-col>
    </v-row>
    <v-row class="mb-4">
      <!-- Weekly activities -->
      <v-col cols="12" sm="6" md="4">
        <v-card elevation="2" class="pa-4 pa-sm-6 text-center h-100 position-relative overflow-hidden">
          <!-- Background gradient -->
          <div class="position-absolute w-100 h-100" style="top: 0; left: 0; background: linear-gradient(135deg, rgba(236, 72, 153, 0.1) 0%, rgba(99, 102, 241, 0.1) 100%); border-radius: inherit;"></div>

          <!-- Background decorative icon - hidden on mobile -->
          <div class="position-absolute d-none d-sm-block" style="top: -10px; right: -10px; opacity: 0.08; transform: rotate(15deg);">
            <v-icon icon="mdi-chart-line" size="120" color="accent2" />
          </div>

          <div class="position-relative">
            <v-icon icon="mdi-calendar-week" :size="$vuetify.display.xs ? 40 : 48" class="mb-3 mb-sm-4" color="accent2" />
            <div class="text-subtitle-1 text-sm-h6 mb-2 mb-sm-3 font-weight-medium">Activities (past 7 days)</div>
            <div class="text-h4 text-sm-h3 font-weight-bold mb-2 mb-sm-3" style="color: #ec4899;">
              <span v-if="statsStore.loading">…</span>
              <span v-else>{{ activitiesWeek }}</span>
            </div>
            <div class="text-body-2 text-medium-emphasis mb-1 mb-sm-2">
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
      <v-col cols="12" sm="6" md="4">
        <v-card elevation="2" class="pa-4 pa-sm-6 text-center h-100 position-relative overflow-hidden">
          <!-- Background gradient -->
          <div class="position-absolute w-100 h-100" style="top: 0; left: 0; background: linear-gradient(135deg, rgba(245, 158, 11, 0.1) 0%, rgba(168, 85, 247, 0.1) 100%); border-radius: inherit;"></div>

          <!-- Background decorative icon - hidden on mobile -->
          <div class="position-absolute d-none d-sm-block" style="top: -15px; right: -15px; opacity: 0.08; transform: rotate(-15deg);">
            <v-icon icon="mdi-sleep" size="140" color="amber" />
          </div>

          <div class="position-relative">
            <v-icon icon="mdi-sleep" :size="$vuetify.display.xs ? 40 : 48" class="mb-3 mb-sm-4" color="amber" />
            <div class="text-subtitle-1 text-sm-h6 mb-2 mb-sm-3 font-weight-medium">Average Daily Sleep</div>
            <div v-if="statsStore.loading" class="text-center py-3 py-sm-4">
              <v-progress-circular indeterminate color="amber" :size="$vuetify.display.xs ? 32 : 40" />
            </div>
            <div v-else-if="statsStore.error" class="text-center text-error py-3 py-sm-4">
              <v-icon icon="mdi-alert-circle" class="mb-2" />
              <div>{{ statsStore.error }}</div>
            </div>
            <div v-else>
              <div class="text-h4 text-sm-h2 font-weight-bold mb-2 mb-sm-3" style="color: #f59e0b">
                {{ avgSleep }}
              </div>
              <div class="text-body-2 text-medium-emphasis mb-1 mb-sm-2">
                per night
              </div>
              <div class="text-caption text-medium-emphasis">
                <v-icon icon="mdi-moon-waning-crescent" size="16" class="mr-1" />
                Based on past 7 days of data
              </div>
            </div>
          </div>
        </v-card>
      </v-col>

      <!-- Avg Counts -->
      <v-col cols="12" md="4">
        <v-card elevation="2" class="pa-4 pa-sm-6 h-100 position-relative overflow-hidden">
          <!-- Background gradient -->
          <div class="position-absolute w-100 h-100" style="top: 0; left: 0; background: linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(34, 197, 94, 0.1) 100%); border-radius: inherit;"></div>

          <!-- Background decorative icon - hidden on mobile -->
          <div class="position-absolute d-none d-sm-block" style="top: -10px; left: -10px; opacity: 0.08; transform: rotate(-10deg);">
            <v-icon icon="mdi-chart-bar" size="130" color="primary" />
          </div>

          <div class="position-relative">
            <div class="text-center mb-3 mb-sm-4">
              <v-icon icon="mdi-chart-bar" :size="$vuetify.display.xs ? 40 : 48" color="primary" />
            </div>
            <div class="text-subtitle-1 text-sm-h6 mb-3 mb-sm-4 font-weight-medium text-center">Average Daily Counts</div>

            <div v-if="statsStore.loading" class="text-center py-6 py-sm-10">
              <v-progress-circular indeterminate color="primary" :size="$vuetify.display.xs ? 32 : 40" />
            </div>
            <div v-else-if="statsStore.error" class="text-center text-error py-6 py-sm-10">
              <v-icon icon="mdi-alert-circle" class="mb-2" />
              <div>{{ statsStore.error }}</div>
            </div>
            <template v-else>
              <div style="height: 150px" class="d-block d-sm-none">
                <!-- Mobile: Show simplified data instead of chart -->
                <div class="d-flex justify-space-around align-center h-100">
                  <div class="text-center">
                    <div class="text-h6 font-weight-bold mb-1" style="color: #6366f1;">
                      {{ avgCountsChartData.datasets[0]?.data[0] || '0' }}
                    </div>
                    <div class="text-caption">Diapers</div>
                  </div>
                  <div class="text-center">
                    <div class="text-h6 font-weight-bold mb-1" style="color: #ec4899;">
                      {{ avgCountsChartData.datasets[0]?.data[1] || '0' }}
                    </div>
                    <div class="text-caption">Feeds</div>
                  </div>
                </div>
              </div>
              <div style="height: 200px" class="d-none d-sm-block">
                <BarChart :chart-data="avgCountsChartData" :chart-options="chartOptions" />
              </div>
            </template>
          </div>
        </v-card>
      </v-col>
    </v-row>

    <!-- Past 7 Days Sleep Trend Chart -->
    <v-row class="mb-4">
      <v-col cols="12">
        <v-card elevation="2" class="position-relative overflow-hidden">
          <!-- Background gradient -->
          <div class="position-absolute w-100 h-100" style="top: 0; left: 0; background: linear-gradient(135deg, rgba(245, 158, 11, 0.05) 0%, rgba(236, 72, 153, 0.05) 100%); border-radius: inherit;"></div>

          <div class="position-relative">
            <v-card-title class="text-subtitle-1 text-sm-h6 pb-2 pb-sm-3">Sleep Trend (Past 7 Days)</v-card-title>
            <v-card-text class="pa-3 pa-sm-4">
              <div v-if="statsStore.loading" class="text-center py-6 py-sm-10">
                <v-progress-circular indeterminate color="primary" :size="$vuetify.display.xs ? 32 : 40" />
              </div>
              <div v-else-if="statsStore.error" class="text-center text-error py-6 py-sm-10">
                {{ statsStore.error }}
              </div>
              <div v-else :style="{ height: $vuetify.display.xs ? '250px' : '300px' }">
                <LineChart :chart-data="weeklySleepTrendData" :chart-options="mobileChartOptions" />
              </div>
            </v-card-text>
          </div>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useStatsStore } from "@/stores/stats";
import { formatTimeAgo } from "@/utils/datetime";
import { format as formatDate, isToday } from "date-fns";
import { useDisplay } from "vuetify";
import BarChart from "@/components/charts/BarChart.vue";
import LineChart from "@/components/charts/LineChart.vue";

const statsStore = useStatsStore();
const { xs } = useDisplay();

// Date picker state
const datePickerMenu = ref(false);
const selectedDate = ref(null);
const today = new Date().toISOString().split('T')[0]; // Format for v-date-picker

// Fetch stats on mount
onMounted(() => {
  statsStore.fetchInitialStats();
});

// Date picker methods
const minDate = computed(() => {
  // Use baby's birth date if available
  if (statsStore.babyBirthDate) {
    return statsStore.babyBirthDate;
  }

  // Fallback: limit to 1 year ago if nothing else is available
  const oneYearAgo = new Date();
  oneYearAgo.setFullYear(oneYearAgo.getFullYear() - 1);
  return oneYearAgo.toISOString().split('T')[0];
});

const onDateSelected = (date) => {
  if (date) {
    selectedDate.value = date;
    statsStore.setDailyDate(new Date(date));
    datePickerMenu.value = false;
  }
};

const jumpToToday = () => {
  statsStore.setDailyDate(new Date());
};

const dailySummaryTitle = computed(() => {
  return isToday(statsStore.dailyDate) ? "Today's Summary" : "Day Summary";
});

const formattedDailyDate = computed(() => {
  return formatDate(statsStore.dailyDate, "MMM d, yyyy");
});

// Dynamic text helpers
const feedsLabel = computed(() => {
  return isToday(statsStore.dailyDate) ? "Feeds today" : "Feeds this day";
});

const pumpingLabel = computed(() => {
  return isToday(statsStore.dailyDate) ? "Pumping today" : "Pumping this day";
});

const diapersLabel = computed(() => {
  return isToday(statsStore.dailyDate) ? "Diapers today" : "Diapers this day";
});

const sleepLabel = computed(() => {
  return isToday(statsStore.dailyDate) ? "Sleep today" : "Sleep this day";
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

// Mobile-optimized chart options with smaller text and better touch targets
const mobileChartOptions = computed(() => ({
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
      ticks: {
        font: {
          size: xs.value ? 10 : 12,
        },
      },
    },
    x: {
      grid: {
        display: false,
      },
      ticks: {
        font: {
          size: xs.value ? 10 : 12,
        },
        maxRotation: xs.value ? 45 : 0,
      },
    },
  },
  elements: {
    point: {
      radius: xs.value ? 3 : 4,
      hoverRadius: xs.value ? 5 : 6,
    },
    line: {
      borderWidth: xs.value ? 2 : 3,
    },
  },
}));

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
