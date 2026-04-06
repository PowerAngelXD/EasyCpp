<script setup>
    import { computed, reactive, ref } from 'vue'

    import {
        loginWithPassword,
        logoutCurrentSession,
        registerAccount,
        sessionState,
    } from '../stores/session'

    const authMode = ref('login')
    const loginForm = reactive({
        email: 'test@example.com',
        password: '12345678',
    })
    const registerForm = reactive({
        username: 'testuser',
        email: 'test@example.com',
        password: '12345678',
        bio: 'Hello, I am testuser',
    })

    const isAuthenticated = computed(() => Boolean(sessionState.accessToken))
    const displayName = computed(
        () => sessionState.user?.username || sessionState.email || 'Authenticated user',
    )

    async function submitLogin() {
        await loginWithPassword(loginForm)
    }

    async function submitRegister() {
        await registerAccount(registerForm)
    }
</script>

<template>
    <section class="authPanel">
        <template v-if="isAuthenticated">
            <div class="panelHeader">
                <div>
                    <div class="panelTitle">Signed In</div>
                    <div class="panelHint">Token-based backend session is active.</div>
                </div>
            </div>

            <div class="identityCard">
                <div class="identityName">{{ displayName }}</div>
                <div class="identityMeta">{{ sessionState.user?.email ?? sessionState.email }}</div>
                <div v-if="sessionState.user?.bio" class="identityBio">{{ sessionState.user.bio }}</div>
            </div>

            <button class="authButton secondary" type="button" @click="logoutCurrentSession">
                Logout
            </button>
        </template>

        <template v-else>
            <div class="panelHeader">
                <div>
                    <div class="panelTitle">Connect Backend</div>
                    <div class="panelHint">Sign in to create posts and comments.</div>
                </div>
            </div>

            <div class="modeSwitch" role="tablist" aria-label="Authentication mode">
                <button
                    class="modeButton"
                    :class="{ active: authMode === 'login' }"
                    type="button"
                    @click="authMode = 'login'"
                >
                    Login
                </button>
                <button
                    class="modeButton"
                    :class="{ active: authMode === 'register' }"
                    type="button"
                    @click="authMode = 'register'"
                >
                    Register
                </button>
            </div>

            <form v-if="authMode === 'login'" class="authForm" @submit.prevent="submitLogin">
                <label class="field">
                    <span class="fieldLabel">Email</span>
                    <input v-model="loginForm.email" class="fieldInput" type="email" required />
                </label>
                <label class="field">
                    <span class="fieldLabel">Password</span>
                    <input
                        v-model="loginForm.password"
                        class="fieldInput"
                        type="password"
                        minlength="8"
                        required
                    />
                </label>
                <button class="authButton" type="submit" :disabled="sessionState.isPending">
                    {{ sessionState.isPending ? 'Signing in...' : 'Sign In' }}
                </button>
            </form>

            <form v-else class="authForm" @submit.prevent="submitRegister">
                <label class="field">
                    <span class="fieldLabel">Username</span>
                    <input
                        v-model="registerForm.username"
                        class="fieldInput"
                        type="text"
                        minlength="3"
                        required
                    />
                </label>
                <label class="field">
                    <span class="fieldLabel">Email</span>
                    <input v-model="registerForm.email" class="fieldInput" type="email" required />
                </label>
                <label class="field">
                    <span class="fieldLabel">Password</span>
                    <input
                        v-model="registerForm.password"
                        class="fieldInput"
                        type="password"
                        minlength="8"
                        required
                    />
                </label>
                <label class="field">
                    <span class="fieldLabel">Bio</span>
                    <textarea
                        v-model="registerForm.bio"
                        class="fieldInput fieldTextarea"
                        rows="3"
                        maxlength="200"
                    ></textarea>
                </label>
                <button class="authButton" type="submit" :disabled="sessionState.isPending">
                    {{ sessionState.isPending ? 'Creating account...' : 'Create Account' }}
                </button>
            </form>

            <div v-if="sessionState.errorMessage" class="errorMessage">
                {{ sessionState.errorMessage }}
            </div>
        </template>
    </section>
</template>

<style scoped>
    .authPanel {
        border: 1px solid var(--borderColor);
        background: color-mix(in srgb, var(--surfaceColor) 92%, var(--pageBgColor));
        border-radius: 12px;
        padding: 12px;
        display: flex;
        flex-direction: column;
        gap: 12px;
    }

    .panelHeader {
        display: flex;
        align-items: flex-start;
        justify-content: space-between;
        gap: 10px;
    }

    .panelTitle {
        color: var(--textColor);
        font-size: 14px;
        font-weight: 650;
    }

    .panelHint {
        margin-top: 4px;
        font-size: 12px;
        color: var(--mutedTextColor);
        line-height: 1.45;
    }

    .modeSwitch {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 6px;
    }

    .modeButton,
    .authButton {
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
        color: var(--textColor);
        border-radius: 10px;
        cursor: pointer;
        font-size: 13px;
    }

    .modeButton {
        padding: 8px 10px;
    }

    .modeButton.active {
        background: var(--accentSurfaceColor);
        border-color: var(--accentBorderColor);
    }

    .authButton {
        padding: 10px 12px;
    }

    .authButton.secondary {
        background: transparent;
        color: var(--mutedTextColor);
    }

    .authButton:disabled {
        cursor: not-allowed;
        opacity: 0.65;
    }

    .authForm {
        display: flex;
        flex-direction: column;
        gap: 10px;
    }

    .field {
        display: flex;
        flex-direction: column;
        gap: 6px;
    }

    .fieldLabel {
        font-size: 12px;
        color: var(--mutedTextColor);
    }

    .fieldInput {
        width: 100%;
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
        color: var(--textColor);
        border-radius: 10px;
        padding: 10px 12px;
        font: inherit;
        box-sizing: border-box;
    }

    .fieldTextarea {
        resize: vertical;
        min-height: 78px;
    }

    .identityCard {
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
        border-radius: 10px;
        padding: 10px 12px;
    }

    .identityName {
        color: var(--textColor);
        font-size: 14px;
        font-weight: 650;
    }

    .identityMeta,
    .identityBio,
    .errorMessage {
        margin-top: 6px;
        font-size: 12px;
        line-height: 1.5;
        color: var(--mutedTextColor);
    }

    .errorMessage {
        color: color-mix(in srgb, #ff4d4f 70%, var(--textColor));
    }
</style>
