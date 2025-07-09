import { defineStore } from "pinia";
import apiClient from "@/api/client";

export const useStatsStore = defineStore("stats", {
  state: () => ({
    recent: null,
    daily: null,
    weekly: null,
    loading: false,
    lastUpdated: null,
    error: null,
  }),

  actions: {
    async fetchStats() {
      this.loading = true;
      this.error = null;
      try {
        const [recentRes, dailyRes, weeklyRes] = await Promise.all([
          apiClient.get("/stats/recent"),
          apiClient.get("/stats/daily"),
          apiClient.get("/stats/weekly"),
        ]);

        this.recent = recentRes.data;
        this.daily = dailyRes.data;
        this.weekly = weeklyRes.data;
        this.lastUpdated = new Date();
      } catch (err) {
        this.error = err.message || "Failed to load statistics.";
        throw err;
      } finally {
        this.loading = false;
      }
    },
  },
});
