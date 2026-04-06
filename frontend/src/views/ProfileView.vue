<script setup>
    import { computed } from 'vue'

    import { sessionState } from '../stores/session'

    const contributionWeeks = 18
    const contributionDaysPerWeek = 7

    const contributions = Array.from(
        { length: contributionWeeks * contributionDaysPerWeek },
        (_, index) => {
            const dayIndex = index % contributionDaysPerWeek
            const weekIndex = Math.floor(index / contributionDaysPerWeek)

            const intensitySeed = (weekIndex * 13 + dayIndex * 7 + 3) % 10
            const intensityLevel =
                intensitySeed <= 2 ? 0 : intensitySeed <= 5 ? 1 : intensitySeed <= 7 ? 2 : 3

            return {
                id: String(index),
                level: intensityLevel,
            }
        },
    )

    const profile = computed(() => {
        if (sessionState.user) {
            return {
                title: sessionState.user.username,
                subtitle: sessionState.user.email,
                bio: sessionState.user.bio || 'No bio provided yet.',
                hint: 'Live profile data loaded from the backend.',
            }
        }

        return {
            title: 'Guest',
            subtitle: 'Sign in from the sidebar',
            bio: 'This page switches to backend user data after login. Activity heatmap remains a UI prototype for now.',
            hint: 'Authentication-driven profile state',
        }
    })
</script>

<template>
    <div class="page">
        <section class="profileHeader">
            <div class="avatar" aria-hidden="true"></div>
            <div class="profileText">
                <h1 class="title">{{ profile.title }}</h1>
                <div class="subtitle">{{ profile.subtitle }}</div>
                <div class="bio">{{ profile.bio }}</div>
            </div>
        </section>

        <section class="card activityCard">
            <div class="cardHeader">
                <h2 class="cardTitle">Activity</h2>
                <div class="cardHint">Heatmap is still a frontend prototype.</div>
            </div>

            <div class="heatmap" role="img" aria-label="Contribution heatmap">
                <div
                    v-for="cell in contributions"
                    :key="cell.id"
                    class="cell"
                    :class="`level${cell.level}`"
                ></div>
            </div>
        </section>

        <section class="card">
            <div class="cardHeader">
                <h2 class="cardTitle">Backend Status</h2>
                <div class="cardHint">{{ profile.hint }}</div>
            </div>
            <div class="about">
                <p>
                    Current session id:
                    <strong>{{ sessionState.sessionId || 'No active session' }}</strong>
                </p>
                <p>
                    Current auth mode:
                    <strong>{{ sessionState.accessToken ? 'Authenticated' : 'Guest' }}</strong>
                </p>
            </div>
        </section>
    </div>
</template>

<style scoped>
    .page {
        max-width: 980px;
    }

    .profileHeader {
        display: flex;
        gap: 14px;
        align-items: center;
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
        border-radius: 12px;
        padding: 16px;
        box-shadow: var(--shadowSm);
    }

    .avatar {
        width: 64px;
        height: 64px;
        border-radius: 18px;
        background: linear-gradient(135deg, var(--accentColor), var(--accentColorSoft));
        flex: 0 0 auto;
    }

    .profileText {
        display: flex;
        flex-direction: column;
        gap: 6px;
        min-width: 0;
    }

    .title {
        margin: 0;
        font-size: 22px;
        font-weight: 700;
        color: var(--textColor);
    }

    .subtitle {
        font-size: 13px;
        color: var(--mutedTextColor);
    }

    .bio {
        font-size: 13px;
        color: var(--mutedTextColor);
        line-height: 1.6;
    }

    .card {
        margin-top: 12px;
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
        border-radius: 12px;
        padding: 16px;
        box-shadow: var(--shadowSm);
    }

    .cardHeader {
        display: flex;
        align-items: baseline;
        justify-content: space-between;
        gap: 12px;
        margin-bottom: 12px;
    }

    .cardTitle {
        margin: 0;
        color: var(--textColor);
        font-size: 15px;
    }

    .cardHint {
        font-size: 12px;
        color: var(--mutedTextColor);
    }

    .activityCard {
        padding-top: 12px;
        padding-bottom: 12px;
    }

    .heatmap {
        display: grid;
        grid-auto-flow: column;
        grid-template-rows: repeat(7, 12px);
        gap: 4px;
        align-content: center;
        justify-content: flex-start;
        width: 100%;
        min-height: 112px;
        padding: 10px 8px;
        border-radius: 10px;
        background: var(--insetSurfaceColor);
        border: 1px solid var(--borderColor);
        box-sizing: border-box;
    }

    .cell {
        width: 12px;
        height: 12px;
        border-radius: 3px;
        background: var(--heatmap0);
        border: 1px solid color-mix(in srgb, var(--borderColor) 65%, transparent);
    }

    .cell.level0 {
        background: var(--heatmap0);
    }

    .cell.level1 {
        background: var(--heatmap1);
    }

    .cell.level2 {
        background: var(--heatmap2);
    }

    .cell.level3 {
        background: var(--heatmap3);
    }

    .about {
        display: flex;
        flex-direction: column;
        gap: 8px;
    }

    .about p {
        margin: 0;
        font-size: 13px;
        line-height: 1.6;
        color: var(--mutedTextColor);
    }

    .about strong {
        color: var(--textColor);
        font-weight: 650;
    }
</style>
