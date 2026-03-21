import { createApp } from 'vue'
import './style.css'
import App from './App.vue'

import { router } from './router'
import { applyTheme, getStoredThemePreference, getSystemThemePreference } from './utils/theme'

const storedThemeName = getStoredThemePreference()
const initialThemeName = storedThemeName ?? getSystemThemePreference()
applyTheme(initialThemeName)

createApp(App).use(router).mount('#app')
