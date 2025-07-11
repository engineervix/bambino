import { defineStore } from "pinia";
import { format, subDays, addDays, isToday } from "date-fns";
import apiClient from "@/api/client";

export const useStatsStore = defineStore("stats", {
  state: () => ({
    recent: null,
    daily: null,
    weekly: null,
    babyBirthDate: null,
    loading: false,
    dailyLoading: false,
    lastUpdated: null,
    error: null,
    dailyDate: new Date(),
  }),

  getters: {
    isViewingToday: (state) => isToday(state.dailyDate),
  },

  actions: {
    async fetchBabyBirthDate() {
      try {
        const response = await apiClient.get("/babies");
        // Assuming user has at least one baby, take the first one
        if (response.data && response.data.length > 0) {
          this.babyBirthDate = response.data[0].birth_date.split('T')[0]; // Convert to YYYY-MM-DD format
        }
      } catch (err) {
        console.error("Failed to fetch baby birth date:", err);
        this.babyBirthDate = null;
      }
    },    async setDailyDate(date) {
      if (date instanceof Date) {
        this.dailyDate = date;
      } else {
        this.dailyDate = new Date(date);
      }

      this.dailyLoading = true;
      try {
        const dateStr = format(this.dailyDate, "yyyy-MM-dd");
        const timezoneOffset = this.dailyDate.getTimezoneOffset();

        // Fetch daily and weekly stats for the new date
        const [dailyRes, weeklyRes] = await Promise.all([
          apiClient.get(`/stats/daily?date=${dateStr}&tz_offset=${timezoneOffset}`),
          apiClient.get(`/stats/weekly?date=${dateStr}&tz_offset=${timezoneOffset}`),
        ]);

        this.daily = dailyRes.data;
        this.weekly = weeklyRes.data;
        this.error = null; // Clear any previous errors
      } catch (err) {
        this.error = err.response?.data?.message || err.message || "Failed to load daily statistics.";

        // If the error is about querying before birth date, reset to today
        if (err.response?.status === 400 && err.response?.data?.message?.includes("birth date")) {
          this.dailyDate = new Date();
          // Retry with today's date
          try {
            const todayStr = format(new Date(), "yyyy-MM-dd");
            const todayOffset = new Date().getTimezoneOffset();
            const [dailyRes, weeklyRes] = await Promise.all([
              apiClient.get(`/stats/daily?date=${todayStr}&tz_offset=${todayOffset}`),
              apiClient.get(`/stats/weekly?date=${todayStr}&tz_offset=${todayOffset}`),
            ]);
            this.daily = dailyRes.data;
            this.weekly = weeklyRes.data;
            this.error = "Cannot view dates before baby's birth date. Showing today instead.";
          } catch (retryErr) {
            console.error("Failed to load today's data after birth date error:", retryErr);
          }
        }
      } finally {
        this.dailyLoading = false;
      }
    },

    async fetchInitialStats() {
      this.loading = true;
      this.error = null;
      try {
        const today = new Date();
        this.dailyDate = today;
        const dateStr = format(today, "yyyy-MM-dd");
        const timezoneOffset = today.getTimezoneOffset();

        const [recentRes, dailyRes, weeklyRes] = await Promise.all([
          apiClient.get("/stats/recent"),
          apiClient.get(`/stats/daily?date=${dateStr}&tz_offset=${timezoneOffset}`),
          apiClient.get(`/stats/weekly?date=${dateStr}&tz_offset=${timezoneOffset}`),
        ]);

        this.recent = recentRes.data;
        this.daily = dailyRes.data;
        this.weekly = weeklyRes.data;
        this.lastUpdated = new Date();

        // Also fetch baby birth date for date picker constraints
        await this.fetchBabyBirthDate();
      } catch (err) {
        this.error = err.message || "Failed to load statistics.";
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async changeDay(direction) {
      if (direction === "next" && this.isViewingToday) {
        return;
      }

      const newDate = direction === "next" ? addDays(this.dailyDate, 1) : subDays(this.dailyDate, 1);
      await this.setDailyDate(newDate);
    },
  },
});
