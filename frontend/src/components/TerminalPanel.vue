<script setup>
    import { computed, nextTick, onMounted, ref, watch } from 'vue'

    const props = defineProps({
        lines: { type: Array, required: true },
        isRunning: { type: Boolean, required: true },
    })

    const emit = defineEmits(['clear'])

    const terminalRef = ref(null)

    const hasLines = computed(() => props.lines.length > 0)

    /**
     * Scroll terminal to bottom after output changes.
     * @returns {Promise<void>}
     */
    async function scrollToBottom() {
        await nextTick()
        const terminalElement = /** @type {HTMLDivElement|null} */ (terminalRef.value)
        if (!terminalElement) {
            return
        }
        terminalElement.scrollTop = terminalElement.scrollHeight
    }

    watch(
        () => props.lines.length,
        () => {
            scrollToBottom()
        },
    )

    onMounted(() => {
        scrollToBottom()
    })
</script>

<template>
    <div class="terminalShell">
        <div class="terminalHeader">
            <div class="terminalTitle">Terminal</div>
            <div class="terminalRight">
                <span class="status" :class="{ running: isRunning }">
                    {{ isRunning ? 'Running' : 'Idle' }}
                </span>
                <button
                    class="smallButton"
                    type="button"
                    :disabled="!hasLines"
                    @click="emit('clear')"
                >
                    Clear
                </button>
            </div>
        </div>

        <div ref="terminalRef" class="terminalBody" role="log" aria-label="Program output">
            <div v-if="!hasLines" class="terminalEmpty">No output yet.</div>
            <div
                v-for="line in lines"
                v-else
                :key="line.id"
                class="terminalLine"
                :class="line.type"
            >
                <span class="prefix" aria-hidden="true">
                    {{ line.type === 'error' ? '!' : line.type === 'system' ? '>' : '' }}
                </span>
                <span class="text">{{ line.text }}</span>
            </div>
        </div>
    </div>
</template>

<style scoped>
    .terminalShell {
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
        border-radius: 12px;
        overflow: hidden;
        box-shadow: var(--shadowSm);
        display: flex;
        flex-direction: column;
        min-height: 240px;
        max-height: 78vh;
    }

    .terminalHeader {
        padding: 10px 12px;
        border-bottom: 1px solid var(--borderColor);
        display: flex;
        align-items: baseline;
        justify-content: space-between;
        gap: 12px;
        background: color-mix(in srgb, var(--surfaceColor) 88%, var(--pageBgColor));
    }

    .terminalTitle {
        font-size: 13px;
        font-weight: 650;
        color: var(--textColor);
    }

    .terminalRight {
        display: inline-flex;
        gap: 10px;
        align-items: center;
    }

    .status {
        font-size: 12px;
        color: var(--mutedTextColor);
        padding: 4px 8px;
        border-radius: 999px;
        border: 1px solid var(--borderColor);
        background: var(--insetSurfaceColor);
    }

    .status.running {
        color: var(--textColor);
        border-color: var(--accentBorderColor);
        background: var(--accentSurfaceColor);
    }

    .smallButton {
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
        color: var(--mutedTextColor);
        font-size: 12px;
        padding: 4px 8px;
        border-radius: 8px;
        cursor: pointer;
    }

    .smallButton:hover:enabled {
        border-color: var(--borderColorHover);
        color: var(--textColor);
    }

    .smallButton:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .terminalBody {
        flex: 1 1 auto;
        padding: 12px;
        background: var(--insetSurfaceColor);
        color: var(--textColor);
        font-family: var(--fontMono);
        font-size: 12.5px;
        line-height: 1.55;
        overflow: auto;
    }

    .terminalEmpty {
        color: var(--mutedTextColor);
    }

    .terminalLine {
        display: flex;
        gap: 10px;
        white-space: pre-wrap;
        word-break: break-word;
    }

    .prefix {
        width: 14px;
        color: var(--mutedTextColor);
        flex: 0 0 auto;
    }

    .terminalLine.error .text {
        color: color-mix(in srgb, #ff4d4f 70%, var(--textColor));
    }

    .terminalLine.system .text {
        color: var(--mutedTextColor);
    }
</style>
