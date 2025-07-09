<template>
  <v-app>
    <!-- Top navigation bar (desktop) -->
    <v-app-bar
      v-if="showTopBar"
      app
      color="surface"
      elevation="1"
      density="comfortable"
    >
      <v-toolbar-title class="text-h6 font-weight-medium d-flex align-center">
        <v-img
          src="/baby.svg"
          alt="Logo"
          width="32"
          height="32"
          cover
          class="d-none d-md-flex"
        />
      </v-toolbar-title>

      <v-spacer />

      <div class="d-flex align-center">
        <v-btn
          v-for="item in navItems"
          :key="'top-' + item.value"
          :to="item.to"
          variant="text"
          :class="{ 'text-primary': activeTab === item.value }"
          size="small"
        >
          <v-icon class="me-1">{{ activeTab === item.value ? item.iconActive : item.iconInactive }}</v-icon>
          <span class="text-capitalize d-none d-lg-inline">{{ item.label }}</span>
        </v-btn>
      </div>
    </v-app-bar>

    <!-- Main content -->
    <v-main>
      <router-view />
    </v-main>

    <!-- Bottom navigation (only show when authenticated and on small screens) -->
    <v-bottom-navigation
      v-if="showBottomNav"
      v-model="activeTab"
      grow
      bg-color="surface"
    >
      <v-btn
        v-for="item in navItems"
        :key="'bottom-' + item.value"
        :value="item.value"
        :to="item.to"
      >
        <v-icon>{{ activeTab === item.value ? item.iconActive : item.iconInactive }}</v-icon>
        <span>{{ item.label }}</span>
      </v-btn>
    </v-bottom-navigation>
  </v-app>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useDisplay } from 'vuetify'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'

const route = useRoute()
const authStore = useAuthStore()
const { isAuthenticated } = storeToRefs(authStore)

const activeTab = ref('activity')

const isLoginPage = computed(() => route.name === 'login')

// Vuetify display helper
const display = useDisplay()

const showBottomNav = computed(() => isAuthenticated.value && !isLoginPage.value && display.smAndDown.value)
const showTopBar = computed(() => isAuthenticated.value && !isLoginPage.value && display.mdAndUp.value)

const navItems = [
  {
    value: 'activity',
    to: '/',
    iconActive: 'mdi-star-four-points',
    iconInactive: 'mdi-star-four-points-outline',
    label: 'Activity'
  },
  {
    value: 'history',
    to: '/history',
    iconActive: 'mdi-calendar',
    iconInactive: 'mdi-calendar-outline',
    label: 'History'
  },
  {
    value: 'trends',
    to: '/trends',
    iconActive: 'mdi-chart-line',
    iconInactive: 'mdi-chart-line', // no outline variant
    label: 'Trends'
  },
  {
    value: 'account',
    to: '/account',
    iconActive: 'mdi-account',
    iconInactive: 'mdi-account-outline',
    label: 'Account'
  }
]

// Initialize auth state on app load
onMounted(async () => {
  await authStore.initializeAuth()

  // Initialize timers if authenticated
  if (authStore.isAuthenticated) {
    const { useTimerStore } = await import('@/stores/timer')
    const timerStore = useTimerStore()
    await timerStore.initializeTimers()
  }
})

// Update active tab based on route
watch(() => route.name, (newRouteName) => {
  if (newRouteName) {
    activeTab.value = newRouteName
  }
}, { immediate: true })
</script>

<style>
/* Global styles */
.v-bottom-navigation {
  position: fixed !important;
}

/* Adjust main content to account for bottom nav */
.v-main {
  padding-bottom: 56px !important;
}

/* Remove padding when no bottom nav */
.v-main:has(~ .v-bottom-navigation:not([style*="display: none"])) {
  padding-bottom: 56px !important;
}
</style>