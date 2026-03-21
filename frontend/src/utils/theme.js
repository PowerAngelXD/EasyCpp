const themeStorageKey = 'olp.theme'

/**
 * Read the preferred theme from localStorage.
 * @returns {('light'|'dark'|null)} The stored theme or null when not set.
 */
export function getStoredThemePreference() {
    try {
        const storedThemeName = window.localStorage.getItem(themeStorageKey)
        if (storedThemeName === 'light' || storedThemeName === 'dark') {
            return storedThemeName
        }
        return null
    } catch {
        return null
    }
}

/**
 * Persist the preferred theme into localStorage.
 * @param {('light'|'dark')} themeName The theme to persist.
 * @returns {void}
 */
export function setStoredThemePreference(themeName) {
    try {
        window.localStorage.setItem(themeStorageKey, themeName)
    } catch {
        // If storage is unavailable, we still allow runtime theme switching.
    }
}

/**
 * Detect the OS/browser preferred color scheme.
 * @returns {('light'|'dark')} The preferred theme from media query.
 */
export function getSystemThemePreference() {
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        return 'dark'
    }
    return 'light'
}

/**
 * Apply theme to the document root.
 * @param {('light'|'dark')} themeName The theme to apply.
 * @returns {void}
 */
export function applyTheme(themeName) {
    document.documentElement.dataset.theme = themeName
    document.documentElement.style.colorScheme = themeName
}
