import { defineStore } from "pinia";
import { format } from "date-fns";
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
        const today = new Date();
        const dateStr = format(today, "yyyy-MM-dd");
        const timezoneOffset = today.getTimezoneOffset();

        const [recentRes, dailyRes, weeklyRes] = await Promise.all([
          apiClient.get("/stats/recent"),
          apiClient.get(`/stats/daily?date=${dateStr}&tz_offset=${timezoneOffset}`),
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
