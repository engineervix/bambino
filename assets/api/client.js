import axios from "axios";
import router from "@/router";

// Create axios instance with default config
const apiClient = axios.create({
  baseURL: "/api",
  timeout: 30000,
  withCredentials: true, // Important for session cookies
  headers: {
    "Content-Type": "application/json",
  },
});

// Track if we're already handling an auth error to prevent loops
let isHandlingAuthError = false;

// Request interceptor
apiClient.interceptors.request.use(
  (config) => {
    // Could add loading indicators here
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

// Response interceptor
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    if (error.response) {
      // Handle 401 Unauthorized
      if (error.response.status === 401 && !isHandlingAuthError) {
        isHandlingAuthError = true;

        try {
          // Import auth store dynamically to avoid circular imports
          const { useAuthStore } = await import("@/stores/auth");
          const authStore = useAuthStore();

          // Clear auth state and redirect to login
          authStore.clearAuthState();

          // Only redirect if not already on login page
          if (router.currentRoute.value.name !== "login") {
            router.push("/login");
          }
        } catch (importError) {
          console.warn("Failed to handle auth error:", importError);
          // Fallback: just redirect to login
          if (router.currentRoute.value.name !== "login") {
            router.push("/login");
          }
        } finally {
          // Reset flag after a short delay
          setTimeout(() => {
            isHandlingAuthError = false;
          }, 1000);
        }
      }

      // Return error with message from backend
      const message =
        error.response.data?.message || error.response.data?.error || `Server error (${error.response.status})`;
      error.message = message;
    } else if (error.request) {
      // Network error
      error.message = "Network error - please check your connection";
    } else if (error.code === "ECONNABORTED") {
      // Timeout error
      error.message = "Request timed out - please try again";
    } else {
      // Other error
      error.message = error.message || "An unexpected error occurred";
    }

    return Promise.reject(error);
  },
);

export default apiClient;
