<script setup>
    import { RouterLink } from 'vue-router'

    import { posts } from '../mock/posts'
</script>

<template>
    <div class="page">
        <div class="pageHeader">
            <h1 class="pageTitle">Feed</h1>
            <div class="pageHint">Clean, GitHub-like cards. Click a post to view comments.</div>
        </div>

        <div class="stack">
            <article v-for="post in posts" :key="post.id" class="card">
                <div class="cardTop">
                    <div class="meta">
                        <span class="author">{{ post.author.displayName }}</span>
                        <span class="dot" aria-hidden="true">·</span>
                        <span class="time">{{ post.createdAt }}</span>
                    </div>
                    <div class="stats">
                        <span class="stat">{{ post.stats.commentsCount }} comments</span>
                        <span class="stat">{{ post.stats.likesCount }} likes</span>
                        <span class="stat">{{ post.stats.favoritesCount }} saved</span>
                    </div>
                </div>

                <h2 class="cardTitle">
                    <RouterLink class="cardLink" :to="`/posts/${post.id}`">
                        {{ post.title }}
                    </RouterLink>
                </h2>
                <p class="cardExcerpt">{{ post.excerpt }}</p>

                <div class="tags">
                    <span v-for="tagName in post.tags" :key="tagName" class="tag">
                        {{ tagName }}
                    </span>
                </div>
            </article>
        </div>
    </div>
</template>

<style scoped>
    .page {
        max-width: 980px;
    }

    .pageHeader {
        display: flex;
        flex-direction: column;
        gap: 6px;
        margin-bottom: 14px;
    }

    .pageTitle {
        margin: 0;
        font-size: 22px;
        font-weight: 650;
        color: var(--textColor);
    }

    .pageHint {
        font-size: 13px;
        color: var(--mutedTextColor);
    }

    .card {
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
        border-radius: 12px;
        padding: 14px 14px 12px;
        box-shadow: var(--shadowSm);
    }

    .cardTop {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 12px;
        margin-bottom: 10px;
    }

    .meta {
        display: inline-flex;
        align-items: center;
        gap: 8px;
        color: var(--mutedTextColor);
        font-size: 12px;
    }

    .author {
        color: var(--textColor);
        font-weight: 600;
    }

    .dot {
        opacity: 0.7;
    }

    .stats {
        display: inline-flex;
        gap: 10px;
        color: var(--mutedTextColor);
        font-size: 12px;
        flex: 0 0 auto;
    }

    .stat {
        padding: 4px 8px;
        border-radius: 999px;
        background: var(--insetSurfaceColor);
        border: 1px solid var(--borderColor);
    }

    .cardTitle {
        margin: 0 0 8px;
        font-size: 16px;
        font-weight: 650;
    }

    .cardLink {
        color: var(--textColor);
        text-decoration: none;
    }

    .cardLink:hover {
        color: var(--accentColor);
    }

    .cardExcerpt {
        margin: 0 0 10px;
        color: var(--mutedTextColor);
        font-size: 13px;
        line-height: 1.5;
    }

    .tags {
        display: flex;
        gap: 8px;
        flex-wrap: wrap;
    }

    .tag {
        font-size: 12px;
        color: var(--mutedTextColor);
        padding: 4px 8px;
        border-radius: 999px;
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
    }

    .stack {
        display: flex;
        flex-direction: column;
        gap: 12px;
    }
</style>
