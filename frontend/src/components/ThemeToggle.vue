<script setup>
    import { computed, onMounted, ref } from 'vue'

    import {
        applyTheme,
        getStoredThemePreference,
        getSystemThemePreference,
        setStoredThemePreference,
    } from '../utils/theme'

    const themeName = ref('light')

    const isDarkTheme = computed(() => themeName.value === 'dark')

    /**
     * Toggle between light and dark theme.
     * @returns {void}
     */
    function toggleTheme() {
        themeName.value = themeName.value === 'dark' ? 'light' : 'dark'
        applyTheme(themeName.value)
        setStoredThemePreference(themeName.value)
    }

    onMounted(() => {
        const storedThemeName = getStoredThemePreference()
        themeName.value = storedThemeName ?? getSystemThemePreference()
        applyTheme(themeName.value)
    })
</script>

<template>
    <button class="themeToggle" type="button" :aria-pressed="isDarkTheme" @click="toggleTheme">
        <svg
            v-if="!isDarkTheme"
            class="themeToggleIcon"
            role="presentation"
            aria-hidden="true"
            viewBox="0 0 24 24"
        >
            <circle cx="12" cy="12" r="5" fill="none" stroke="currentColor" stroke-width="1.5" />
            <path
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
                d="M12 1v3M12 20v3M3.41 3.41l2.83 2.83M17.76 17.76l2.83 2.83M1 12h3M20 12h3M3.41 20.59l2.83-2.83M17.76 6.24l2.83-2.83"
            />
        </svg>
        <svg
            v-else
            class="themeToggleIcon"
            role="presentation"
            aria-hidden="true"
            viewBox="0 0 24 24"
        >
            <path fill="currentColor" d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
        </svg>
        <span class="themeToggleLabel">{{ isDarkTheme ? 'Dark theme' : 'Light theme' }}</span>
    </button>
</template>

<style scoped>
    .themeToggle {
        display: inline-flex;
        align-items: center;
        gap: 8px;
        padding: 6px 10px;
        border-radius: 8px;
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
        color: var(--textColor);
        cursor: pointer;
        font: inherit;
    }

    .themeToggle:hover {
        border-color: var(--borderColorHover);
    }

    .themeToggle:focus-visible {
        outline: 2px solid var(--accentColor);
        outline-offset: 2px;
    }

    .themeToggleIcon {
        width: 16px;
        height: 16px;
        color: var(--mutedTextColor);
    }

    .themeToggleLabel {
        font-size: 13px;
        color: var(--mutedTextColor);
    }
</style>
