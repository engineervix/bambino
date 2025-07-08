import { format, parseISO, isToday, isYesterday } from 'date-fns'

/**
 * Combines date and time strings into a proper Date object
 * @param {string} dateStr - Date in YYYY-MM-DD format
 * @param {string} timeStr - Time in HH:MM format
 * @returns {Date}
 */
export function combineDateAndTime(dateStr, timeStr) {
  return new Date(`${dateStr}T${timeStr}:00`)
}

/**
 * Gets current date in YYYY-MM-DD format
 * @returns {string}
 */
export function getCurrentDate() {
  return new Date().toISOString().split('T')[0]
}

/**
 * Gets current time in HH:MM format
 * @returns {string}
 */
export function getCurrentTime() {
  return new Date().toTimeString().slice(0, 5)
}

/**
 * Extracts date from Date object in YYYY-MM-DD format
 * @param {Date} date 
 * @returns {string}
 */
export function getDateString(date) {
  return date.toISOString().split('T')[0]
}

/**
 * Extracts time from Date object in HH:MM format
 * @param {Date} date 
 * @returns {string}
 */
export function getTimeString(date) {
  return date.toTimeString().slice(0, 5)
}

/**
 * Formats a date for display
 * @param {string|Date} date 
 * @returns {string}
 */
export function formatActivityDate(date) {
  const dateObj = typeof date === 'string' ? parseISO(date) : date
  
  if (isToday(dateObj)) {
    return `Today ${format(dateObj, 'h:mm a')}`
  } else if (isYesterday(dateObj)) {
    return `Yesterday ${format(dateObj, 'h:mm a')}`
  } else {
    return format(dateObj, 'MMM d, h:mm a')
  }
}

/**
 * Formats duration in minutes to readable string
 * @param {number} minutes 
 * @returns {string}
 */
export function formatDuration(minutes) {
  if (minutes < 60) {
    return `${minutes} min`
  }
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  return mins > 0 ? `${hours}h ${mins}m` : `${hours}h`
}

/**
 * Formats relative time (e.g., "2 hours ago")
 * @param {string|Date} date 
 * @returns {string}
 */
export function formatTimeAgo(date) {
  const dateObj = typeof date === 'string' ? parseISO(date) : date
  const now = new Date()
  const diffInMinutes = Math.floor((now - dateObj) / (1000 * 60))
  
  if (diffInMinutes < 1) {
    return 'Just now'
  } else if (diffInMinutes < 60) {
    return `${diffInMinutes} min ago`
  } else if (diffInMinutes < 1440) { // Less than 24 hours
    const hours = Math.floor(diffInMinutes / 60)
    return `${hours}h ago`
  } else {
    const days = Math.floor(diffInMinutes / 1440)
    return `${days}d ago`
  }
}