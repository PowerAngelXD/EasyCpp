<script setup>
    import { RouterView } from 'vue-router'

    import SidebarNav from './SidebarNav.vue'
    import ThemeToggle from './ThemeToggle.vue'
</script>

<template>
    <div class="appShell">
        <SidebarNav />
        <div class="main">
            <header class="topBar">
                <div class="topBarLeft">
                    <div class="topBarTitle">Online Learning</div>
                    <div class="topBarSubtitle">Prototype UI</div>
                </div>
                <div class="topBarRight">
                    <ThemeToggle />
                </div>
            </header>

            <main class="content">
                <RouterView />
            </main>
        </div>
    </div>
</template>

<style scoped>
    .appShell {
        min-height: 100vh;
        display: flex;
        background: var(--pageBgColor);
    }

    .main {
        flex: 1 1 auto;
        min-width: 0;
        display: flex;
        flex-direction: column;
    }

    .topBar {
        position: sticky;
        top: 0;
        z-index: 10;
        background: color-mix(in srgb, var(--pageBgColor) 92%, transparent);
        backdrop-filter: blur(10px);
        border-bottom: 1px solid var(--borderColor);
        padding: 12px 16px;
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 12px;
    }

    .topBarLeft {
        display: flex;
        flex-direction: column;
        gap: 2px;
        min-width: 0;
    }

    .topBarTitle {
        font-weight: 600;
        color: var(--textColor);
        line-height: 1.2;
    }

    .topBarSubtitle {
        font-size: 12px;
        color: var(--mutedTextColor);
    }

    .content {
        padding: 18px 18px 40px;
    }

    @media (max-width: 980px) {
        .appShell {
            flex-direction: column;
        }
    }
</style>
