const defaultApiBaseURL = import.meta.env.VITE_API_BASE_URL ?? 'http://127.0.0.1:8080'

export class ApiError extends Error {
    constructor(message, status, payload) {
        super(message)
        this.name = 'ApiError'
        this.status = status
        this.payload = payload
    }
}

/**
 * @param {string} path
 * @param {object} [options]
 * @param {string} [options.method]
 * @param {Record<string, string>} [options.headers]
 * @param {string} [options.token]
 * @param {unknown} [options.body]
 * @returns {Promise<any>}
 */
export async function apiRequest(path, options = {}) {
    const { method = 'GET', headers = {}, token = '', body } = options

    const requestHeaders = new Headers(headers)
    if (body !== undefined) {
        requestHeaders.set('Content-Type', 'application/json')
    }
    if (token) {
        requestHeaders.set('Authorization', `Bearer ${token}`)
    }

    const response = await fetch(`${defaultApiBaseURL}${path}`, {
        method,
        headers: requestHeaders,
        body: body === undefined ? undefined : JSON.stringify(body),
    })

    const payload = await response.json().catch(() => null)
    const message = payload?.message || `Request failed: HTTP ${response.status}`
    if (!response.ok) {
        throw new ApiError(message, response.status, payload)
    }

    return payload?.data
}

export function registerUser(input) {
    return apiRequest('/api/v1/auth/register', {
        method: 'POST',
        body: input,
    })
}

export function loginUser(input) {
    return apiRequest('/api/v1/auth/login', {
        method: 'POST',
        body: input,
    })
}

export function logoutUser(token) {
    return apiRequest('/api/v1/auth/logout', {
        method: 'POST',
        token,
    })
}

export function listPosts(token = '') {
    return apiRequest('/api/v1/posts', { token })
}

export function getPostById(postId, token = '') {
    return apiRequest(`/api/v1/posts/${postId}`, { token })
}

export function createPost(input, token) {
    return apiRequest('/api/v1/posts', {
        method: 'POST',
        body: input,
        token,
    })
}

export function listCommentsByPostId(postId, token = '') {
    return apiRequest(`/api/v1/posts/${postId}/comments`, { token })
}

export function createComment(postId, input, token) {
    return apiRequest(`/api/v1/posts/${postId}/comments`, {
        method: 'POST',
        body: input,
        token,
    })
}

export function getUserById(userId, token = '') {
    return apiRequest(`/api/v1/users/${userId}`, { token })
}
