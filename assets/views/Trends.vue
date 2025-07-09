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
  </v-container>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useStatsStore } from '@/stores/stats'
import { formatTimeAgo } from '@/utils/datetime'

const statsStore = useStatsStore()

// Fetch stats on mount (if not already loaded)
onMounted(() => {
  if (!statsStore.lastUpdated) {
    statsStore.fetchStats()
  }
})

// --- Computed helpers -------------------------------------------------

const lastFedDisplay = computed(() => {
  const feed = statsStore.recent?.last_feed
  if (!feed) return '—'
  return formatTimeAgo(feed.time)
})

const lastFeedAmount = computed(() => {
  const feed = statsStore.recent?.last_feed
  if (!feed?.amount_ml) return ''
  return `${feed.amount_ml} ml`
})

const diapersToday = computed(() => {
  return statsStore.daily?.counts?.diaper ?? '0'
})

function formatHoursMinutes(hoursFloat) {
  if (hoursFloat === undefined || hoursFloat === null) return ''
  const totalMinutes = Math.round(hoursFloat * 60)
  const h = Math.floor(totalMinutes / 60)
  const m = totalMinutes % 60
  return `${h > 0 ? h + 'h ' : ''}${m}m`
}

const sleepingDisplay = computed(() => {
  if (!statsStore.recent) return '—'

  // Ongoing sleep
  if (statsStore.recent.currently_sleeping) {
    const startISO = statsStore.daily?.last_activities?.sleep
    if (startISO) {
      const diffMs = Date.now() - new Date(startISO).getTime()
      const minutes = Math.floor(diffMs / 60000)
      const h = Math.floor(minutes / 60)
      const m = minutes % 60
      return `${h > 0 ? h + 'h ' : ''}${m}m`
    }
    return 'Sleeping'
  }

  // Not sleeping → last sleep duration
  const dur = statsStore.recent.last_sleep?.duration_hours
  if (dur !== undefined && dur !== null) {
    return formatHoursMinutes(dur)
  }
  return '—'
})

const activitiesWeek = computed(() => {
  const avgs = statsStore.weekly?.daily_averages
  if (!avgs) return '0'
  let total = 0
  for (const [key, value] of Object.entries(avgs)) {
    if (key.endsWith('_per_day') && !key.includes('hours') && !key.includes('amount')) {
      total += value * 7
    }
  }
  return Math.round(total)
})
</script>