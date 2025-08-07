import { format, parseISO, isToday, isYesterday } from "date-fns";

/**
 * Combines date and time strings into a proper Date object
 * @param {string} dateStr - Date in YYYY-MM-DD format
 * @param {string} timeStr - Time in HH:MM format
 * @returns {Date}
 */
export function combineDateAndTime(dateStr, timeStr) {
  return new Date(`${dateStr}T${timeStr}:00`);
}

/**
 * Gets current date in YYYY-MM-DD format
 * Uses local timezone to ensure correct date display
 * @returns {string}
 */
export function getCurrentDate() {
  const now = new Date();
  return getDateString(now);
}

/**
 * Gets current time in HH:MM format
 * Uses local timezone to ensure correct time display
 * @returns {string}
 */
export function getCurrentTime() {
  const now = new Date();
  return getTimeString(now);
}

/**
 * Extracts date from Date object in YYYY-MM-DD format
 * Uses local timezone to avoid timezone conversion issues
 * @param {Date} date
 * @returns {string}
 */
export function getDateString(date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}

/**
 * Extracts time from Date object in HH:MM format
 * Uses local timezone to avoid timezone conversion issues
 * @param {Date} date
 * @returns {string}
 */
export function getTimeString(date) {
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  return `${hours}:${minutes}`;
}

/**
 * Formats a date for display
 * @param {string|Date} date
 * @returns {string}
 */
export function formatActivityDate(date) {
  const dateObj = typeof date === "string" ? parseISO(date) : date;

  if (isToday(dateObj)) {
    return `Today ${format(dateObj, "h:mm a")}`;
  } else if (isYesterday(dateObj)) {
    return `Yesterday ${format(dateObj, "h:mm a")}`;
  } else {
    return format(dateObj, "MMM d, h:mm a");
  }
}

/**
 * Formats duration in minutes to readable string
 * @param {number} minutes
 * @returns {string}
 */
export function formatDuration(minutes) {
  if (minutes < 60) {
    return `${minutes} min`;
  }
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;
  return mins > 0 ? `${hours}h ${mins}m` : `${hours}h`;
}

/**
 * Formats relative time (e.g., "2 hours ago")
 * @param {string|Date} date
 * @returns {string}
 */
export function formatTimeAgo(date) {
  const dateObj = typeof date === "string" ? parseISO(date) : date;
  const now = new Date();
  const diffInMinutes = Math.floor((now - dateObj) / (1000 * 60));

  if (diffInMinutes < 1) {
    return "Just now";
  } else if (diffInMinutes < 60) {
    return `${diffInMinutes} min ago`;
  } else if (diffInMinutes < 1440) {
    // Less than 24 hours
    const hours = Math.floor(diffInMinutes / 60);
    return `${hours}h ago`;
  } else {
    const days = Math.floor(diffInMinutes / 1440);
    return `${days}d ago`;
  }
}

/**
 * Gets current date and time in a format suitable for datetime-local input.
 * YYYY-MM-DDTHH:mm
 * Uses local timezone to ensure correct datetime display
 * @returns {string}
 */
export function getCurrentDateTimeLocal() {
  const now = new Date();
  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, '0');
  const day = String(now.getDate()).padStart(2, '0');
  const hours = String(now.getHours()).padStart(2, '0');
  const minutes = String(now.getMinutes()).padStart(2, '0');
  return `${year}-${month}-${day}T${hours}:${minutes}`;
}

/**
 * Formats a Date object into a string for datetime-local input.
 * Uses local timezone to avoid timezone conversion issues
 * @param {Date} date
 * @returns {string}
 */
export function formatDateTimeLocal(date) {
  const d = new Date(date);
  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  const hours = String(d.getHours()).padStart(2, '0');
  const minutes = String(d.getMinutes()).padStart(2, '0');
  return `${year}-${month}-${day}T${hours}:${minutes}`;
}
