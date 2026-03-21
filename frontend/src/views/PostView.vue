<script setup>
    import { computed } from 'vue'
    import { useRoute } from 'vue-router'

    import { posts } from '../mock/posts'

    const route = useRoute()

    const postId = computed(() => String(route.params.postId ?? ''))
    const post = computed(() => posts.find((item) => item.id === postId.value) ?? null)

    const commentThreads = [
        {
            id: 'c1',
            author: { displayName: 'Z', handle: 'z' },
            createdAt: '2026-03-19',
            body: 'Nice! The layout feels clean and readable.',
            replies: [
                {
                    id: 'c1-r1',
                    author: { displayName: 'Z', handle: 'z' },
                    createdAt: '2026-03-19',
                    body: 'Thanks. Next step is adding the code editor playground view.',
                },
            ],
        },
        {
            id: 'c2',
            author: { displayName: 'Z', handle: 'z' },
            createdAt: '2026-03-18',
            body: 'Could we make the theme toggle remember the last choice?',
            replies: [],
        },
    ]
</script>

<template>
    <div class="page" v-if="post">
        <div class="header">
            <div class="meta">
                <span class="author">{{ post.author.displayName }}</span>
                <span class="dot" aria-hidden="true">·</span>
                <span class="time">{{ post.createdAt }}</span>
            </div>
            <h1 class="title">{{ post.title }}</h1>
            <p class="excerpt">{{ post.excerpt }}</p>
        </div>

        <div class="article">
            <h2>What you will build</h2>
            <ul>
                <li>Left sidebar navigation</li>
                <li>Feed + post detail with GitHub-like comments</li>
                <li>Profile page with contribution heatmap</li>
                <li>Light/Dark theme toggle (stored + system fallback)</li>
            </ul>

            <h2>Design notes</h2>
            <p>
                This prototype uses a small set of tokens (border, surface, muted text, accent) so
                switching themes is only changing CSS variables.
            </p>
        </div>

        <section class="comments">
            <div class="commentsHeader">
                <h2 class="commentsTitle">Comments</h2>
                <div class="commentsHint">Prototype of GitHub-style thread UI.</div>
            </div>

            <div class="composer">
                <div class="composerAvatar" aria-hidden="true"></div>
                <div class="composerBox" aria-label="Write a comment" role="textbox">
                    Type a comment (UI only)
                </div>
            </div>

            <div class="threadList">
                <article v-for="thread in commentThreads" :key="thread.id" class="thread">
                    <div class="threadLine" aria-hidden="true"></div>
                    <div class="comment">
                        <div class="commentHeader">
                            <div class="commentMeta">
                                <span class="commentAuthor">{{ thread.author.displayName }}</span>
                                <span class="dot" aria-hidden="true">·</span>
                                <span class="commentTime">{{ thread.createdAt }}</span>
                            </div>
                            <div class="commentActions">
                                <button class="ghostButton" type="button">Reply</button>
                                <button class="ghostButton" type="button">Like</button>
                                <button class="ghostButton" type="button">Save</button>
                            </div>
                        </div>
                        <div class="commentBody">{{ thread.body }}</div>
                    </div>

                    <div v-if="thread.replies.length" class="replies">
                        <article
                            v-for="reply in thread.replies"
                            :key="reply.id"
                            class="comment reply"
                        >
                            <div class="commentHeader">
                                <div class="commentMeta">
                                    <span class="commentAuthor">{{
                                        reply.author.displayName
                                    }}</span>
                                    <span class="dot" aria-hidden="true">·</span>
                                    <span class="commentTime">{{ reply.createdAt }}</span>
                                </div>
                                <div class="commentActions">
                                    <button class="ghostButton" type="button">Reply</button>
                                </div>
                            </div>
                            <div class="commentBody">{{ reply.body }}</div>
                        </article>
                    </div>
                </article>
            </div>
        </section>
    </div>

    <div v-else class="page">
        <h1 class="title">Post not found</h1>
        <p class="excerpt">The post id is invalid.</p>
    </div>
</template>

<style scoped>
    .page {
        max-width: 980px;
    }

    .header {
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
        border-radius: 12px;
        padding: 16px;
        box-shadow: var(--shadowSm);
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
        font-weight: 650;
    }

    .dot {
        opacity: 0.7;
    }

    .title {
        margin: 8px 0 6px;
        font-size: 22px;
        font-weight: 700;
        color: var(--textColor);
    }

    .excerpt {
        margin: 0;
        color: var(--mutedTextColor);
        font-size: 13px;
        line-height: 1.55;
    }

    .article {
        margin-top: 12px;
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
        border-radius: 12px;
        padding: 16px;
        box-shadow: var(--shadowSm);
    }

    .article h2 {
        margin: 0 0 8px;
        color: var(--textColor);
        font-size: 15px;
    }

    .article p,
    .article li {
        color: var(--mutedTextColor);
        font-size: 13px;
        line-height: 1.6;
    }

    .comments {
        margin-top: 12px;
    }

    .commentsHeader {
        display: flex;
        align-items: baseline;
        justify-content: space-between;
        gap: 12px;
        margin: 16px 0 10px;
    }

    .commentsTitle {
        margin: 0;
        font-size: 16px;
        color: var(--textColor);
    }

    .commentsHint {
        font-size: 12px;
        color: var(--mutedTextColor);
    }

    .composer {
        display: flex;
        gap: 12px;
        align-items: flex-start;
        margin-bottom: 12px;
    }

    .composerAvatar {
        width: 34px;
        height: 34px;
        border-radius: 999px;
        background: var(--insetSurfaceColor);
        border: 1px solid var(--borderColor);
        flex: 0 0 auto;
    }

    .composerBox {
        flex: 1 1 auto;
        min-height: 42px;
        border-radius: 12px;
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
        padding: 10px 12px;
        color: var(--mutedTextColor);
        font-size: 13px;
    }

    .threadList {
        display: flex;
        flex-direction: column;
        gap: 12px;
    }

    .thread {
        position: relative;
        padding-left: 18px;
    }

    .threadLine {
        position: absolute;
        left: 6px;
        top: 0;
        bottom: 0;
        width: 2px;
        background: var(--borderColor);
        border-radius: 999px;
    }

    .comment {
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
        border-radius: 12px;
        padding: 12px;
        box-shadow: var(--shadowSm);
    }

    .reply {
        margin-top: 10px;
        margin-left: 22px;
    }

    .commentHeader {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 12px;
        margin-bottom: 8px;
    }

    .commentMeta {
        display: inline-flex;
        gap: 8px;
        align-items: center;
        color: var(--mutedTextColor);
        font-size: 12px;
    }

    .commentAuthor {
        color: var(--textColor);
        font-weight: 650;
    }

    .commentActions {
        display: inline-flex;
        gap: 6px;
    }

    .ghostButton {
        border: 1px solid transparent;
        background: transparent;
        color: var(--mutedTextColor);
        font-size: 12px;
        padding: 4px 8px;
        border-radius: 8px;
        cursor: pointer;
    }

    .ghostButton:hover {
        background: var(--surfaceColorHover);
        border-color: var(--borderColor);
        color: var(--textColor);
    }

    .commentBody {
        font-size: 13px;
        color: var(--textColor);
        line-height: 1.6;
    }
</style>
