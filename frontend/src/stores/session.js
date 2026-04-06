import { reactive } from 'vue'

import { getUserById, loginUser, logoutUser, registerUser } from '../utils/api'

const sessionStorageKey = 'easycpp.session'

/**
 * @param {string} token
 * @returns {{ uid?: number } | null}
 */
function parseJwtPayload(token) {
    if (!token) {
        return null
    }

    try {
        const [, payloadPart] = token.split('.')
        if (!payloadPart) {
            return null
        }

        const normalized = payloadPart.replace(/-/g, '+').replace(/_/g, '/')
        const padded = normalized.padEnd(Math.ceil(normalized.length / 4) * 4, '=')
        return JSON.parse(window.atob(padded))
    } catch {
        return null
    }
}

function readStoredSession() {
    try {
        const rawValue = window.localStorage.getItem(sessionStorageKey)
        if (!rawValue) {
            return null
        }
        return JSON.parse(rawValue)
    } catch {
        return null
    }
}

const storedSession = readStoredSession()

export const sessionState = reactive({
    accessToken: storedSession?.accessToken ?? '',
    tokenType: storedSession?.tokenType ?? 'Bearer',
    expiresIn: storedSession?.expiresIn ?? 0,
    sessionId: storedSession?.sessionId ?? '',
    email: storedSession?.email ?? '',
    user: storedSession?.user ?? null,
    isPending: false,
    errorMessage: '',
})

function persistSession() {
    window.localStorage.setItem(
        sessionStorageKey,
        JSON.stringify({
            accessToken: sessionState.accessToken,
            tokenType: sessionState.tokenType,
            expiresIn: sessionState.expiresIn,
            sessionId: sessionState.sessionId,
            email: sessionState.email,
            user: sessionState.user,
        }),
    )
}

function clearSession() {
    sessionState.accessToken = ''
    sessionState.tokenType = 'Bearer'
    sessionState.expiresIn = 0
    sessionState.sessionId = ''
    sessionState.email = ''
    sessionState.user = null
    sessionState.errorMessage = ''
    window.localStorage.removeItem(sessionStorageKey)
}

export async function hydrateCurrentUser() {
    if (!sessionState.accessToken) {
        return null
    }

    const payload = parseJwtPayload(sessionState.accessToken)
    const userId = Number(payload?.uid ?? 0)
    if (!userId) {
        return null
    }

    try {
        const user = await getUserById(userId, sessionState.accessToken)
        sessionState.user = user
        if (user?.email) {
            sessionState.email = user.email
        }
        persistSession()
        return user
    } catch {
        return null
    }
}

export async function registerAccount(input) {
    sessionState.isPending = true
    sessionState.errorMessage = ''

    try {
        await registerUser(input)
        return await loginWithPassword({
            email: input.email,
            password: input.password,
        })
    } catch (error) {
        sessionState.errorMessage = error instanceof Error ? error.message : 'Failed to register.'
        throw error
    } finally {
        sessionState.isPending = false
    }
}

export async function loginWithPassword(input) {
    sessionState.isPending = true
    sessionState.errorMessage = ''

    try {
        const loginData = await loginUser(input)
        sessionState.accessToken = loginData?.accessToken ?? ''
        sessionState.tokenType = loginData?.tokenType ?? 'Bearer'
        sessionState.expiresIn = Number(loginData?.expiresIn ?? 0)
        sessionState.sessionId = loginData?.sessionId ?? ''
        sessionState.email = input.email
        persistSession()
        await hydrateCurrentUser()
        return loginData
    } catch (error) {
        clearSession()
        sessionState.errorMessage = error instanceof Error ? error.message : 'Failed to sign in.'
        throw error
    } finally {
        sessionState.isPending = false
    }
}

export async function logoutCurrentSession() {
    const activeToken = sessionState.accessToken
    clearSession()

    if (!activeToken) {
        return
    }

    try {
        await logoutUser(activeToken)
    } catch {
        // Ignore logout transport failures after local cleanup.
    }
}
