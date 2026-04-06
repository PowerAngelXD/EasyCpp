<script setup>
    import { RouterLink } from 'vue-router'

    import AuthPanel from './AuthPanel.vue'
</script>

<template>
    <aside class="sidebar">
        <RouterLink class="sidebarBrand" to="/profile">
            <div class="brandMark" aria-hidden="true"></div>
            <div class="brandText">
                <div class="brandTitle">Study</div>
                <div class="brandSubtitle">Learning Platform</div>
            </div>
        </RouterLink>

        <nav class="nav" aria-label="Primary">
            <RouterLink class="navItem" to="/">
                <span class="navIcon" aria-hidden="true">H</span>
                <span class="navLabel">Home</span>
            </RouterLink>
            <RouterLink class="navItem" to="/playground">
                <span class="navIcon" aria-hidden="true">C</span>
                <span class="navLabel">Playground</span>
            </RouterLink>
            <RouterLink class="navItem" to="/profile">
                <span class="navIcon" aria-hidden="true">P</span>
                <span class="navLabel">Profile</span>
            </RouterLink>
        </nav>

        <div class="sidebarFooter">
            <div class="footerHint">GitHub-like, clean UI</div>
            <AuthPanel />
        </div>
    </aside>
</template>

<style scoped>
    .sidebar {
        width: 260px;
        background: var(--surfaceColor);
        border-right: 1px solid var(--borderColor);
        padding: 16px 12px;
        display: flex;
        flex-direction: column;
        gap: 16px;
    }

    .sidebarBrand {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 10px 10px;
        border-radius: 10px;
        text-decoration: none;
        color: var(--textColor);
    }

    .brandMark {
        width: 28px;
        height: 28px;
        border-radius: 8px;
        background: linear-gradient(135deg, var(--accentColor), var(--accentColorSoft));
    }

    .brandText {
        display: flex;
        flex-direction: column;
        gap: 2px;
        min-width: 0;
    }

    .brandTitle {
        font-weight: 600;
        color: var(--textColor);
        line-height: 1.1;
    }

    .brandSubtitle {
        font-size: 12px;
        color: var(--mutedTextColor);
    }

    .nav {
        display: flex;
        flex-direction: column;
        gap: 6px;
    }

    .navItem {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 10px 10px;
        border-radius: 10px;
        color: var(--textColor);
        text-decoration: none;
        border: 1px solid transparent;
    }

    .navItem:hover {
        background: var(--surfaceColorHover);
        border-color: var(--borderColor);
    }

    .navItem.router-link-active {
        background: var(--accentSurfaceColor);
        border-color: var(--accentBorderColor);
    }

    .navIcon {
        width: 22px;
        height: 22px;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        border-radius: 7px;
        font-size: 12px;
        color: var(--mutedTextColor);
        background: var(--insetSurfaceColor);
        border: 1px solid var(--borderColor);
        flex: 0 0 auto;
    }

    .navLabel {
        font-size: 14px;
    }

    .sidebarFooter {
        margin-top: auto;
        padding: 10px 10px 2px;
        color: var(--mutedTextColor);
        font-size: 12px;
        display: flex;
        flex-direction: column;
        gap: 10px;
    }
</style>
