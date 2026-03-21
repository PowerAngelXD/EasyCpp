<script setup>
    import { computed } from 'vue'

    const props = defineProps({
        modelValue: { type: String, required: true },
        languageId: { type: String, required: true },
    })

    const emit = defineEmits(['update:modelValue'])

    const editorAriaLabel = computed(() => `Code editor (${props.languageId})`)

    /**
     * Update editor content.
     * @param {Event} event Input event.
     * @returns {void}
     */
    function onInput(event) {
        const target = /** @type {HTMLTextAreaElement} */ (event.target)
        emit('update:modelValue', target.value)
    }
</script>

<template>
    <div class="editorShell">
        <div class="editorHeader">
            <div class="editorTitle">Editor</div>
            <div class="editorHint">{{ languageId }}</div>
        </div>
        <textarea
            class="editor"
            :value="modelValue"
            :aria-label="editorAriaLabel"
            spellcheck="false"
            autocapitalize="off"
            autocomplete="off"
            @input="onInput"
        ></textarea>
    </div>
</template>

<style scoped>
    .editorShell {
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
        border-radius: 12px;
        overflow: hidden;
        box-shadow: var(--shadowSm);
        display: flex;
        flex-direction: column;
        min-height: 360px;
    }

    .editorHeader {
        padding: 10px 12px;
        border-bottom: 1px solid var(--borderColor);
        display: flex;
        align-items: baseline;
        justify-content: space-between;
        gap: 12px;
        background: color-mix(in srgb, var(--surfaceColor) 88%, var(--pageBgColor));
    }

    .editorTitle {
        font-size: 13px;
        font-weight: 650;
        color: var(--textColor);
    }

    .editorHint {
        font-size: 12px;
        color: var(--mutedTextColor);
    }

    .editor {
        flex: 1 1 auto;
        width: 100%;
        resize: none;
        border: 0;
        outline: none;
        padding: 12px;
        background: transparent;
        color: var(--textColor);
        font-family: var(--fontMono);
        font-size: 13px;
        line-height: 1.6;
        tab-size: 4;
    }
</style>
