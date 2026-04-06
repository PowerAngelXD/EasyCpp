<script setup>
    import { computed, onMounted, reactive, ref, watch } from 'vue'
    import { useRoute } from 'vue-router'

    import { createComment, getPostById, getUserById, listCommentsByPostId } from '../utils/api'
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
            hour: '2-digit',
            minute: '2-digit',
        }).format(dateValue)
    }

    const route = useRoute()
    const postId = computed(() => String(route.params.postId ?? ''))
    const isAuthenticated = computed(() => Boolean(sessionState.accessToken))

    const post = ref(null)
    const comments = ref([])
    const isLoading = ref(false)
    const isSubmitting = ref(false)
    const pageError = ref('')
    const commentError = ref('')

    const commentForm = reactive({
        content: '',
    })

    const userNameCache = new Map()

    async function resolveUserName(userId) {
        const cacheKey = String(userId)
        if (userNameCache.has(cacheKey)) {
            return userNameCache.get(cacheKey)
        }

        try {
            const user = await getUserById(userId, sessionState.accessToken)
            const name = user?.username || `User #${userId}`
            userNameCache.set(cacheKey, name)
            return name
        } catch {
            const fallbackName = `User #${userId}`
            userNameCache.set(cacheKey, fallbackName)
            return fallbackName
        }
    }

    async function loadPostPage() {
        isLoading.value = true
        pageError.value = ''

        try {
            const [postData, commentData] = await Promise.all([
                getPostById(postId.value, sessionState.accessToken),
                listCommentsByPostId(postId.value, sessionState.accessToken),
            ])

            post.value = postData

            const rawComments = Array.isArray(commentData) ? commentData : []
            comments.value = await Promise.all(
                rawComments.map(async (comment) => ({
                    id: String(comment.id),
                    authorName: await resolveUserName(comment.authorId),
                    createdAt: formatDateTime(comment.createdAt),
                    content: comment.content,
                })),
            )
        } catch (error) {
            post.value = null
            comments.value = []
            pageError.value = error instanceof Error ? error.message : 'Failed to load post.'
        } finally {
            isLoading.value = false
        }
    }

    async function submitComment() {
        if (!isAuthenticated.value) {
            commentError.value = 'Sign in before publishing a comment.'
            return
        }

        isSubmitting.value = true
        commentError.value = ''

        try {
            await createComment(
                postId.value,
                { content: commentForm.content },
                sessionState.accessToken,
            )
            commentForm.content = ''
            await loadPostPage()
        } catch (error) {
            commentError.value = error instanceof Error ? error.message : 'Failed to publish comment.'
        } finally {
            isSubmitting.value = false
        }
    }

    onMounted(() => {
        loadPostPage()
    })

    watch(
        [postId, () => sessionState.accessToken],
        () => {
            loadPostPage()
        },
    )
</script>

<template>
    <div class="page">
        <div v-if="isLoading" class="messageCard">Loading post...</div>
        <div v-else-if="pageError" class="messageCard errorText">{{ pageError }}</div>

        <template v-else-if="post">
            <div class="header">
                <div class="meta">
                    <span class="author">{{ post.authorName || `User #${post.authorId}` }}</span>
                    <span class="dot" aria-hidden="true">路</span>
                    <span class="time">{{ formatDateTime(post.createdAt) }}</span>
                    <span class="dot" aria-hidden="true">路</span>
                    <span class="time">{{ post.language }}</span>
                </div>
                <h1 class="title">{{ post.title }}</h1>
                <p class="excerpt">{{ post.summary || 'No summary provided.' }}</p>
                <div class="headerBadges">
                    <span class="badge">{{ post.difficulty }}</span>
                    <span class="badge">{{ post.commentCount }} comments</span>
                    <span class="badge">{{ post.likeCount }} likes</span>
                    <span class="badge">{{ post.viewCount }} views</span>
                </div>
            </div>

            <div class="article">
                <h2>Content</h2>
                <pre class="articleBody">{{ post.content }}</pre>
                <div class="tags">
                    <span v-for="tagName in post.tags" :key="tagName" class="tag">
                        {{ tagName }}
                    </span>
                </div>
            </div>

            <section class="comments">
                <div class="commentsHeader">
                    <h2 class="commentsTitle">Comments</h2>
                    <div class="commentsHint">Live data from the backend API.</div>
                </div>

                <form class="composer" @submit.prevent="submitComment">
                    <div class="composerAvatar" aria-hidden="true"></div>
                    <div class="composerMain">
                        <textarea
                            v-model="commentForm.content"
                            class="composerInput"
                            rows="4"
                            maxlength="1000"
                            :disabled="!isAuthenticated || isSubmitting"
                            :placeholder="
                                isAuthenticated
                                    ? 'Write a comment...'
                                    : 'Sign in from the sidebar to comment.'
                            "
                        ></textarea>
                        <div class="composerActions">
                            <div v-if="commentError" class="errorText">{{ commentError }}</div>
                            <button
                                class="primaryButton"
                                type="submit"
                                :disabled="!isAuthenticated || isSubmitting || !commentForm.content.trim()"
                            >
                                {{ isSubmitting ? 'Publishing...' : 'Publish Comment' }}
                            </button>
                        </div>
                    </div>
                </form>

                <div v-if="!comments.length" class="messageCard">No comments yet.</div>

                <div v-else class="threadList">
                    <article v-for="comment in comments" :key="comment.id" class="comment">
                        <div class="commentHeader">
                            <div class="commentMeta">
                                <span class="commentAuthor">{{ comment.authorName }}</span>
                                <span class="dot" aria-hidden="true">路</span>
                                <span class="commentTime">{{ comment.createdAt }}</span>
                            </div>
                        </div>
                        <div class="commentBody">{{ comment.content }}</div>
                    </article>
                </div>
            </section>
        </template>

        <div v-else class="messageCard">Post not found.</div>
    </div>
</template>

<style scoped>
    .page {
        max-width: 980px;
    }

    .header,
    .article,
    .comment,
    .messageCard {
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
        border-radius: 12px;
        box-shadow: var(--shadowSm);
    }

    .header,
    .article {
        padding: 16px;
    }

    .article,
    .comments {
        margin-top: 12px;
    }

    .messageCard {
        padding: 14px 16px;
        color: var(--mutedTextColor);
    }

    .meta,
    .commentMeta,
    .headerBadges,
    .tags {
        display: inline-flex;
        align-items: center;
        gap: 8px;
        flex-wrap: wrap;
    }

    .meta,
    .commentMeta {
        color: var(--mutedTextColor);
        font-size: 12px;
    }

    .author,
    .commentAuthor {
        color: var(--textColor);
        font-weight: 650;
    }

    .dot {
        opacity: 0.7;
    }

    .title {
        margin: 10px 0 6px;
        font-size: 24px;
        font-weight: 700;
        color: var(--textColor);
    }

    .excerpt {
        margin: 0;
        color: var(--mutedTextColor);
        font-size: 13px;
        line-height: 1.55;
    }

    .headerBadges {
        margin-top: 12px;
    }

    .badge,
    .tag {
        padding: 4px 8px;
        border-radius: 999px;
        border: 1px solid var(--borderColor);
        font-size: 12px;
        color: var(--mutedTextColor);
    }

    .badge {
        background: var(--insetSurfaceColor);
    }

    .article h2,
    .commentsTitle {
        margin: 0;
        color: var(--textColor);
        font-size: 16px;
    }

    .articleBody {
        margin: 10px 0 0;
        padding: 14px;
        white-space: pre-wrap;
        word-break: break-word;
        background: var(--insetSurfaceColor);
        border: 1px solid var(--borderColor);
        border-radius: 10px;
        color: var(--textColor);
        font-family: var(--fontMono);
        font-size: 13px;
        line-height: 1.65;
    }

    .tags {
        margin-top: 12px;
    }

    .commentsHeader,
    .composerActions,
    .commentHeader {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 12px;
    }

    .commentsHeader {
        margin: 16px 0 10px;
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

    .composerMain {
        flex: 1 1 auto;
        display: flex;
        flex-direction: column;
        gap: 10px;
    }

    .composerInput {
        width: 100%;
        box-sizing: border-box;
        resize: vertical;
        min-height: 96px;
        border-radius: 12px;
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
        padding: 10px 12px;
        color: var(--textColor);
        font: inherit;
    }

    .primaryButton {
        border: 1px solid var(--accentBorderColor);
        background: var(--accentSurfaceColor);
        color: var(--textColor);
        font-size: 13px;
        padding: 8px 12px;
        border-radius: 10px;
        cursor: pointer;
    }

    .primaryButton:disabled {
        cursor: not-allowed;
        opacity: 0.65;
    }

    .threadList {
        display: flex;
        flex-direction: column;
        gap: 12px;
    }

    .comment {
        padding: 12px;
    }

    .commentHeader {
        margin-bottom: 8px;
    }

    .commentBody {
        font-size: 13px;
        color: var(--textColor);
        line-height: 1.6;
        white-space: pre-wrap;
    }

    .errorText {
        color: color-mix(in srgb, #ff4d4f 70%, var(--textColor));
    }

    @media (max-width: 760px) {
        .commentsHeader,
        .composerActions {
            align-items: stretch;
            flex-direction: column;
        }
    }
</style>
