import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import 'vuetify/styles'

export default createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'dark',
    themes: {
      dark: {
        dark: true,
        colors: {
          primary: '#4CAF50',
          secondary: '#2196F3',
          background: '#121212', // overall app background
          surface: '#1E1E1E',    // cards, sheets, dialogs
          accent1: '#FF9F45',    // warm accent used for highlights/charts
          accent2: '#00BCD4',    // cool accent counterpart
          success: '#4CAF50',    // semantic success
          error: '#FF5252',      // semantic error
          feed: '#FFC107',
          pump: '#E91E63',
          diaper: '#795548',
          sleep: '#3F51B5',
          growth: '#4CAF50',
          health: '#9C27B0',
          milestone: '#FF5722',
          // Background tints (surface-tone variants)
          'feed-bg': '#332700',     // dark amber tint
          'pump-bg': '#3B0617',     // muted rose tint
          'diaper-bg': '#2A231F',   // mocha tint
          'sleep-bg': '#09123B',    // indigo tint
          'growth-bg': '#0E2910',   // deep green tint
          'health-bg': '#1A0D26',   // violet tint
          'milestone-bg': '#331407' // orange tint
        }
      }
    }
  }
})
