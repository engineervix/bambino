import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "@/stores/auth";

// Helper to capitalize first letter
function capitalize(word) {
  if (!word) return "";
  return word.charAt(0).toUpperCase() + word.slice(1);
}

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/login",
      name: "login",
      component: () => import("../views/UserLogin.vue"),
      meta: { requiresAuth: false },
    },
    {
      path: "/",
      name: "activity",
      component: () => import("../views/ActivityView.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/history",
      name: "history",
      component: () => import("../views/HistoryView.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/trends",
      name: "trends",
      component: () => import("../views/TrendsView.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/account",
      name: "account",
      component: () => import("../views/UserAccount.vue"),
      meta: { requiresAuth: true },
    },
    // 404 – keep as last route
    {
      path: "/:pathMatch(.*)*",
      name: "not-found",
      component: () => import("../views/NotFound.vue"),
      meta: { requiresAuth: false, title: "Not Found" },
    },
  ],
});

// Navigation guard
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore();

  // Skip auth check for login page
  if (to.meta.requiresAuth === false) {
    next();
    return;
  }

  // Fast check using local state if already authenticated
  const fastAuthCheck = authStore.isAuthenticatedFast();

  let isAuthenticated;
  if (fastAuthCheck !== null) {
    // We have already checked auth, use local state
    isAuthenticated = fastAuthCheck;
  } else {
    // First time or auth state unknown, do full check
    isAuthenticated = await authStore.initializeAuth();
  }

  if (to.meta.requiresAuth && !isAuthenticated) {
    next("/login");
  } else if (to.path === "/login" && isAuthenticated) {
    next("/");
  } else {
    // Load saved baby selection only if authenticated
    if (isAuthenticated) {
      authStore.loadSelectedBaby();
    }
    next();
  }
});

// Update document title after each navigation
router.afterEach((to) => {
  const appName = "Bambino";
  // Prefer explicit meta title, otherwise derive from route name
  const pageTitle = to.meta && to.meta.title ? to.meta.title : typeof to.name === "string" ? capitalize(to.name) : "";
  document.title = pageTitle ? `${appName} » ${pageTitle}` : appName;
});

export default router;
