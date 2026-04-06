<script setup>
    import { computed, onMounted, reactive, ref, watch } from 'vue'
    import { RouterLink } from 'vue-router'

    import { createPost, listPosts } from '../utils/api'
    import { sessionState } from '../stores/session'

    function formatDateTime(value) {
        if (!value) {
            return 'Unknown time'
        }

        const dateValue = new Date(value)
        if (Number.isNaN(dateValue.getTime())) {
            return value
        }

        return new Intl.DateTimeFormat('en-US', {
            month: 'short',
            day: 'numeric',
            year: 'numeric',
        }).format(dateValue)
    }

    function mapPostToCard(post) {
        return {
            id: String(post.id),
            title: post.title,
            excerpt: post.summary || post.content?.slice(0, 180) || 'No summary yet.',
            createdAt: formatDateTime(post.createdAt),
            authorName: post.authorName || `User #${post.authorId}`,
            stats: {
                commentsCount: Number(post.commentCount ?? 0),
                likesCount: Number(post.likeCount ?? 0),
                favoritesCount: Number(post.viewCount ?? 0),
            },
            tags: Array.isArray(post.tags) ? post.tags : [],
            difficulty: post.difficulty || 'unknown',
            language: post.language || 'cpp',
        }
    }

    const posts = ref([])
    const isLoading = ref(false)
    const loadError = ref('')
    const createError = ref('')
    const isSubmitting = ref(false)

    const createForm = reactive({
        title: '',
        summary: '',
        content: '#include <iostream>\n\nint main() {\n    std::cout << "Hello, EasyCpp!" << std::endl;\n    return 0;\n}\n',
        language: 'cpp',
        difficulty: 'beginner',
        tagsText: 'cpp, starter',
    })

    const isAuthenticated = computed(() => Boolean(sessionState.accessToken))

    async function loadFeed() {
        isLoading.value = true
        loadError.value = ''

        try {
            const data = await listPosts(sessionState.accessToken)
            posts.value = Array.isArray(data) ? data.map(mapPostToCard) : []
        } catch (error) {
            loadError.value = error instanceof Error ? error.message : 'Failed to load posts.'
            posts.value = []
        } finally {
            isLoading.value = false
        }
    }

    async function submitPost() {
        if (!isAuthenticated.value) {
            createError.value = 'Sign in before creating a post.'
            return
        }

        isSubmitting.value = true
        createError.value = ''

        try {
            await createPost(
                {
                    title: createForm.title,
                    summary: createForm.summary,
                    content: createForm.content,
                    language: createForm.language,
                    difficulty: createForm.difficulty,
                    tags: createForm.tagsText
                        .split(',')
                        .map((tag) => tag.trim())
                        .filter(Boolean),
                },
                sessionState.accessToken,
            )

            createForm.title = ''
            createForm.summary = ''
            createForm.tagsText = 'cpp, starter'
            await loadFeed()
        } catch (error) {
            createError.value = error instanceof Error ? error.message : 'Failed to create post.'
        } finally {
            isSubmitting.value = false
        }
    }

    onMounted(() => {
        loadFeed()
    })

    watch(
        () => sessionState.accessToken,
        () => {
            loadFeed()
        },
    )
</script>

<template>
    <div class="page">
        <div class="pageHeader">
            <div>
                <h1 class="pageTitle">Feed</h1>
                <div class="pageHint">Live posts from the backend API. Empty state is expected on a fresh database.</div>
            </div>
            <button class="refreshButton" type="button" :disabled="isLoading" @click="loadFeed">
                {{ isLoading ? 'Refreshing...' : 'Refresh' }}
            </button>
        </div>

        <section class="composerCard">
            <div class="composerHeader">
                <div>
                    <h2 class="composerTitle">Create Post</h2>
                    <div class="composerHint">
                        {{ isAuthenticated ? 'Publish directly to the backend feed.' : 'Sign in from the sidebar to publish posts.' }}
                    </div>
                </div>
            </div>

            <form class="composerForm" @submit.prevent="submitPost">
                <label class="field fieldWide">
                    <span class="fieldLabel">Title</span>
                    <input
                        v-model="createForm.title"
                        class="fieldInput"
                        type="text"
                        minlength="5"
                        maxlength="120"
                        :disabled="!isAuthenticated || isSubmitting"
                        required
                    />
                </label>

                <label class="field fieldWide">
                    <span class="fieldLabel">Summary</span>
                    <textarea
                        v-model="createForm.summary"
                        class="fieldInput fieldTextarea"
                        rows="3"
                        maxlength="280"
                        :disabled="!isAuthenticated || isSubmitting"
                    ></textarea>
                </label>

                <label class="field fieldWide">
                    <span class="fieldLabel">Content</span>
                    <textarea
                        v-model="createForm.content"
                        class="fieldInput fieldTextarea codeField"
                        rows="10"
                        :disabled="!isAuthenticated || isSubmitting"
                        required
                    ></textarea>
                </label>

                <label class="field">
                    <span class="fieldLabel">Language</span>
                    <select
                        v-model="createForm.language"
                        class="fieldInput"
                        :disabled="!isAuthenticated || isSubmitting"
                    >
                        <option value="cpp">C++</option>
                        <option value="c">C</option>
                    </select>
                </label>

                <label class="field">
                    <span class="fieldLabel">Difficulty</span>
                    <select
                        v-model="createForm.difficulty"
                        class="fieldInput"
                        :disabled="!isAuthenticated || isSubmitting"
                    >
                        <option value="beginner">Beginner</option>
                        <option value="intermediate">Intermediate</option>
                        <option value="advanced">Advanced</option>
                    </select>
                </label>

                <label class="field fieldWide">
                    <span class="fieldLabel">Tags</span>
                    <input
                        v-model="createForm.tagsText"
                        class="fieldInput"
                        type="text"
                        placeholder="cpp, arrays, beginner"
                        :disabled="!isAuthenticated || isSubmitting"
                    />
                </label>

                <div class="composerActions">
                    <div v-if="createError" class="errorText">{{ createError }}</div>
                    <button
                        class="primaryButton"
                        type="submit"
                        :disabled="!isAuthenticated || isSubmitting"
                    >
                        {{ isSubmitting ? 'Publishing...' : 'Publish Post' }}
                    </button>
                </div>
            </form>
        </section>

        <div v-if="loadError" class="messageCard errorText">{{ loadError }}</div>
        <div v-else-if="isLoading" class="messageCard">Loading posts...</div>
        <div v-else-if="!posts.length" class="messageCard">
            No posts yet. Create the first one from the form above.
        </div>

        <div v-else class="stack">
            <article v-for="post in posts" :key="post.id" class="card">
                <div class="cardTop">
                    <div class="meta">
                        <span class="author">{{ post.authorName }}</span>
                        <span class="dot" aria-hidden="true">路</span>
                        <span class="time">{{ post.createdAt }}</span>
                    </div>
                    <div class="stats">
                        <span class="stat">{{ post.stats.commentsCount }} comments</span>
                        <span class="stat">{{ post.stats.likesCount }} likes</span>
                        <span class="stat">{{ post.stats.favoritesCount }} views</span>
                    </div>
                </div>

                <h2 class="cardTitle">
                    <RouterLink class="cardLink" :to="`/posts/${post.id}`">
                        {{ post.title }}
                    </RouterLink>
                </h2>

                <p class="cardExcerpt">{{ post.excerpt }}</p>

                <div class="inlineMeta">
                    <span class="pill">{{ post.language }}</span>
                    <span class="pill">{{ post.difficulty }}</span>
                </div>

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

    .pageHeader,
    .composerHeader,
    .composerActions,
    .cardTop {
        display: flex;
        align-items: flex-start;
        justify-content: space-between;
        gap: 12px;
    }

    .pageHeader {
        margin-bottom: 14px;
    }

    .pageTitle,
    .composerTitle {
        margin: 0;
        color: var(--textColor);
    }

    .pageTitle {
        font-size: 22px;
        font-weight: 650;
    }

    .composerTitle {
        font-size: 17px;
        font-weight: 650;
    }

    .pageHint,
    .composerHint {
        margin-top: 6px;
        font-size: 13px;
        line-height: 1.5;
        color: var(--mutedTextColor);
    }

    .composerCard,
    .card,
    .messageCard {
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
        border-radius: 12px;
        box-shadow: var(--shadowSm);
    }

    .composerCard {
        padding: 16px;
        margin-bottom: 14px;
    }

    .messageCard {
        padding: 14px 16px;
        color: var(--mutedTextColor);
    }

    .composerForm {
        display: grid;
        grid-template-columns: repeat(2, minmax(0, 1fr));
        gap: 12px;
        margin-top: 12px;
    }

    .field {
        display: flex;
        flex-direction: column;
        gap: 6px;
    }

    .fieldWide {
        grid-column: 1 / -1;
    }

    .fieldLabel {
        font-size: 12px;
        color: var(--mutedTextColor);
    }

    .fieldInput {
        width: 100%;
        box-sizing: border-box;
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
        color: var(--textColor);
        border-radius: 10px;
        padding: 10px 12px;
        font: inherit;
    }

    .fieldTextarea {
        resize: vertical;
        min-height: 80px;
    }

    .codeField {
        font-family: var(--fontMono);
        line-height: 1.6;
    }

    .composerActions {
        grid-column: 1 / -1;
        align-items: center;
    }

    .refreshButton,
    .primaryButton {
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
        color: var(--textColor);
        font-size: 13px;
        padding: 9px 12px;
        border-radius: 10px;
        cursor: pointer;
    }

    .primaryButton {
        border-color: var(--accentBorderColor);
        background: var(--accentSurfaceColor);
    }

    .refreshButton:disabled,
    .primaryButton:disabled {
        cursor: not-allowed;
        opacity: 0.65;
    }

    .stack {
        display: flex;
        flex-direction: column;
        gap: 12px;
    }

    .card {
        padding: 14px 14px 12px;
    }

    .meta,
    .stats,
    .inlineMeta,
    .tags {
        display: inline-flex;
        flex-wrap: wrap;
        gap: 8px;
    }

    .meta {
        align-items: center;
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
        color: var(--mutedTextColor);
        font-size: 12px;
    }

    .stat,
    .pill,
    .tag {
        padding: 4px 8px;
        border-radius: 999px;
        border: 1px solid var(--borderColor);
        font-size: 12px;
    }

    .stat,
    .pill {
        background: var(--insetSurfaceColor);
        color: var(--mutedTextColor);
    }

    .tag {
        background: var(--surfaceColor);
        color: var(--mutedTextColor);
    }

    .cardTitle {
        margin: 10px 0 8px;
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

    .inlineMeta {
        margin-bottom: 10px;
    }

    .errorText {
        color: color-mix(in srgb, #ff4d4f 70%, var(--textColor));
    }

    @media (max-width: 760px) {
        .composerForm {
            grid-template-columns: 1fr;
        }

        .cardTop,
        .pageHeader,
        .composerHeader,
        .composerActions {
            flex-direction: column;
            align-items: stretch;
        }
    }
</style>
