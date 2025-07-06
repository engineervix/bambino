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
          feed: '#FFC107',
          pump: '#E91E63',
          diaper: '#795548',
          sleep: '#3F51B5',
          growth: '#4CAF50',
          health: '#9C27B0',
          milestone: '#FF5722'
        }
      }
    }
  }
})
