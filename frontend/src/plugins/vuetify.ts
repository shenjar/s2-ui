/**
 * plugins/vuetify.ts
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles/main.css'

import { fa, en, vi, zhHans, zhHant, ru } from 'vuetify/locale'

// Composables
import { createVuetify } from 'vuetify'

// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  defaults: {
    VTextField: {
      variant: 'outlined',
      color: 'primary',
    },
    VSelect: {
      variant: 'outlined',
      color: 'primary',
    },
    VCombobox: {
      variant: 'outlined',
      color: 'primary',
    },
    VTextarea: {
      variant: 'outlined',
      color: 'primary',
    },
    VAutocomplete: {
      variant: 'outlined',
      color: 'primary',
    },
    VSwitch: {
      color: 'primary',
    },
    VCheckbox: {
      color: 'primary',
    },
    VTabs: {
      color: 'primary',
    },
    VSlider: {
      color: 'primary',
    },
    VDataTable: {
      hover: true,
    },
    VDataTableServer: {
      hover: true,
    },
  },
  theme: {
    defaultTheme: localStorage.getItem('theme') ?? 'system',
    themes: {
      light: {
        colors: {
          primary: '#4F46E5',
          secondary: '#0EA5E9',
          accent: '#7C3AED',
          background: '#F2F4F8',
          surface: '#FFFFFF',
          error: '#DC2626',
          success: '#16A34A',
          warning: '#D97706',
          info: '#0284C7',
        },
        variables: {
          'border-color': '#475569',
          'border-opacity': 0.14,
        },
      },
      dark: {
        colors: {
          primary: '#818CF8',
          secondary: '#38BDF8',
          accent: '#A78BFA',
          background: '#0E1117',
          surface: '#171B24',
          error: '#F87171',
          success: '#4ADE80',
          warning: '#FBBF24',
          info: '#7DD3FC',
        },
        variables: {
          'border-color': '#94A3B8',
          'border-opacity': 0.16,
        },
      },
    },
  },
  locale: {
    locale: localStorage.getItem("locale") ?? 'en',
    fallback: 'en',
    messages: { en, fa, vi, zhHans, zhHant, ru },
  },
})
