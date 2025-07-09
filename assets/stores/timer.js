import { defineStore } from "pinia";
import { ref, computed, watch } from "vue";
import apiClient from "@/api/client";
import { useIntervalFn } from "@vueuse/core";

export const useTimerStore = defineStore("timer", () => {
  const activeTimers = ref({
    feed: null,
    pump: null,
    sleep: null,
  });

  const isValidating = ref(false); // Add validation state
  const formattedDurations = ref({});
  let intervalControls;

  // Debounced persistence
  let persistenceTimeout = null;
  const STORAGE_KEY = "bambino-timers";
  const PERSISTENCE_DELAY = 300;

  // Enhanced validation with retry logic
  async function validateTimers() {
    if (isValidating.value) return false; // Prevent concurrent validation

    isValidating.value = true;
    const currentTimers = { ...activeTimers.value };
    let hasChanges = false;

    try {
      for (const [type, timer] of Object.entries(currentTimers)) {
        if (!timer?.activityId) continue;

        try {
          const response = await apiClient.get(`/activities/${timer.activityId}`);
          const activity = response.data;

          if (activity.end_time) {
            console.log(`Timer ${type} completed - activity has end_time`);
            activeTimers.value[type] = null;
            hasChanges = true;
          }
        } catch (error) {
          if (error.response?.status === 404) {
            console.log(`Timer ${type} invalid - activity not found`);
            activeTimers.value[type] = null;
            hasChanges = true;
          }
          // For other errors (network, etc.), keep timer and retry later
        }
      }

      if (hasChanges) {
        persistTimers();
      }
    } finally {
      isValidating.value = false;
    }

    return hasChanges;
  }

  // Initialize with better error handling
  async function initializeTimers() {
    console.log("Initializing timers...");

    loadPersistedTimers();

    if (hasActiveTimers.value) {
      try {
        await validateTimers();
      } catch (error) {
        console.warn("Timer validation failed, will retry later:", error);
        // Schedule retry in 30 seconds
        setTimeout(() => {
          if (hasActiveTimers.value) {
            validateTimers();
          }
        }, 30000);
      }
    }
  }

  // Rest of your existing code...
  function startTimer(type, activityId) {
    if (!type || !activityId) {
      console.error("Invalid timer start parameters:", { type, activityId });
      return false;
    }

    if (activeTimers.value[type]) {
      console.warn(`Timer already active for ${type}`);
      return false;
    }

    const timer = {
      startTime: new Date().toISOString(),
      activityId: activityId,
      type: type,
    };

    activeTimers.value[type] = timer;
    debouncedPersist();

    console.log(`Timer started for ${type}:`, timer);
    return true;
  }

  function stopTimer(type) {
    if (!activeTimers.value[type]) {
      console.warn(`No active timer found for ${type}`);
      return false;
    }

    console.log(`Timer stopped for ${type}`);
    activeTimers.value[type] = null;
    debouncedPersist();
    return true;
  }

  function getTimerDuration(type) {
    const timer = activeTimers.value[type];
    if (!timer) return 0;

    const startTime = new Date(timer.startTime);
    if (isNaN(startTime.getTime())) {
      console.warn(`Invalid startTime for ${type} timer:`, timer.startTime);
      return 0;
    }

    return Math.floor((Date.now() - startTime.getTime()) / 1000);
  }

  function _getFormattedDuration(type) {
    const duration = getTimerDuration(type);
    const hours = Math.floor(duration / 3600);
    const minutes = Math.floor((duration % 3600) / 60);
    const seconds = duration % 60;

    if (hours > 0) {
      return `${hours}h ${minutes}m`;
    } else {
      return `${minutes.toString().padStart(2, "0")}:${seconds.toString().padStart(2, "0")}`;
    }
  }

  // Computed properties
  const hasActiveTimers = computed(() => {
    return Object.values(activeTimers.value).some((timer) => timer !== null);
  });

  const activeTimerCount = computed(() => {
    return Object.values(activeTimers.value).filter((timer) => timer !== null).length;
  });

  function updateFormattedDurations() {
    const newDurations = {};
    Object.keys(activeTimers.value).forEach((type) => {
      if (activeTimers.value[type]) {
        newDurations[type] = _getFormattedDuration(type);
      } else {
        newDurations[type] = null;
      }
    });
    formattedDurations.value = newDurations;
  }

  // Timer update interval
  intervalControls = useIntervalFn(updateFormattedDurations, 1000, { immediate: false });

  watch(
    hasActiveTimers,
    (hasActive) => {
      if (hasActive) {
        updateFormattedDurations();
        intervalControls.resume();
      } else {
        if (intervalControls.isActive.value) {
          intervalControls.pause();
        }
        formattedDurations.value = {};
      }
    },
    { immediate: true },
  );

  // Persistence functions (keep your existing implementation)
  function debouncedPersist() {
    if (persistenceTimeout) {
      clearTimeout(persistenceTimeout);
    }

    persistenceTimeout = setTimeout(() => {
      persistTimers();
      persistenceTimeout = null;
    }, PERSISTENCE_DELAY);
  }

  function persistTimers() {
    try {
      const timersToSave = {};
      Object.keys(activeTimers.value).forEach((type) => {
        if (activeTimers.value[type]) {
          timersToSave[type] = activeTimers.value[type];
        }
      });

      if (Object.keys(timersToSave).length > 0) {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(timersToSave));
      } else {
        localStorage.removeItem(STORAGE_KEY);
      }
    } catch (error) {
      console.warn("Failed to persist timers:", error);
    }
  }

  function loadPersistedTimers() {
    try {
      const saved = localStorage.getItem(STORAGE_KEY);
      if (saved) {
        const timers = JSON.parse(saved);

        Object.keys(timers).forEach((type) => {
          if (timers[type]?.activityId && timers[type]?.startTime) {
            const startTime = new Date(timers[type].startTime);
            if (!isNaN(startTime.getTime())) {
              activeTimers.value[type] = timers[type];
            }
          }
        });
      }
    } catch (error) {
      console.warn("Failed to load persisted timers:", error);
      localStorage.removeItem(STORAGE_KEY);
    }
  }

  function clearAllTimers() {
    activeTimers.value = {
      feed: null,
      pump: null,
      sleep: null,
    };

    if (persistenceTimeout) {
      clearTimeout(persistenceTimeout);
      persistenceTimeout = null;
    }

    localStorage.removeItem(STORAGE_KEY);
  }

  return {
    // State
    activeTimers,
    isValidating,
    formattedDurations,

    // Computed
    hasActiveTimers,
    activeTimerCount,

    // Actions
    startTimer,
    stopTimer,
    getActiveTimer: (type) => activeTimers.value[type],
    hasActiveTimer: (type) => !!activeTimers.value[type],
    getTimerDuration,

    // Lifecycle
    initializeTimers,
    validateTimers,
    clearAllTimers,
    persistTimers,
  };
});
