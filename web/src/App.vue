<template>
  <v-app>
    <!-- Main content -->
    <v-main>
      <!-- Desktop / tablet navigation drawer -->
      <v-navigation-drawer
        v-if="showDrawer"
        permanent
        rail
        class="elevation-1"
        width="72"
      >
        <v-list density="compact" nav>
          <v-list-item
            v-for="item in navItems"
            :key="item.value"
            :to="item.to"
            :active="activeTab === item.value"
            rounded="xl"
          >
            <template #prepend>
              <v-icon>{{ activeTab === item.value ? item.iconActive : item.iconInactive }}</v-icon>
            </template>
            <v-list-item-title>{{ item.label }}</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-navigation-drawer>

      <router-view />
    </v-main>

    <!-- Bottom navigation (only show when authenticated and not on login page) -->
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
const showDrawer = computed(() => isAuthenticated.value && !isLoginPage.value && display.mdAndUp.value)

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