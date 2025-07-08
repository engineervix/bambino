import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/Login.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/',
      name: 'activity',
      component: () => import('../views/Activity.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/history',
      name: 'history',
      component: () => import('../views/History.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/trends',
      name: 'trends',
      component: () => import('../views/Trends.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/account',
      name: 'account',
      component: () => import('../views/Account.vue'),
      meta: { requiresAuth: true }
    }
  ]
})

// Navigation guard
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // Skip auth check for login page
  if (to.meta.requiresAuth === false) {
    next()
    return
  }

  // Fast check using local state if already authenticated
  const fastAuthCheck = authStore.isAuthenticatedFast()
  
  let isAuthenticated
  if (fastAuthCheck !== null) {
    // We have already checked auth, use local state
    isAuthenticated = fastAuthCheck
  } else {
    // First time or auth state unknown, do full check
    isAuthenticated = await authStore.initializeAuth()
  }
  
  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login')
  } else if (to.path === '/login' && isAuthenticated) {
    next('/')
  } else {
    // Load saved baby selection only if authenticated
    if (isAuthenticated) {
      authStore.loadSelectedBaby()
    }
    next()
  }
})

export default router