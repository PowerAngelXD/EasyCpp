<script setup>
    import { onBeforeUnmount, ref, watch } from 'vue'

    import CodeEditor from '../components/CodeEditor.vue'
    import TerminalPanel from '../components/TerminalPanel.vue'
    import { createCppRunner } from '../utils/cppRunner'
    import { createJsRunner } from '../utils/jsRunner'

    const defaultCodeByLanguage = {
        javascript: `function main() {
    console.log('Hello.')
}

main()
`,
        cpp: `#include <iostream>

int main() {
    std::cout << "Hello from C++." << std::endl;
    return 0;
}
`,
        python: `print('Python support is not connected yet.')`,
    }

    const languageId = ref('javascript')
    const codeDrafts = ref({ ...defaultCodeByLanguage })
    const codeText = ref(codeDrafts.value[languageId.value])

    const terminalLines = ref([])
    const isRunning = ref(false)

    const jsRunner = createJsRunner({
        onLine(line) {
            terminalLines.value = [...terminalLines.value, line]
        },
        onRunningChange(nextIsRunning) {
            isRunning.value = nextIsRunning
        },
    })

    const cppRunner = createCppRunner({
        onLine(line) {
            terminalLines.value = [...terminalLines.value, line]
        },
        onRunningChange(nextIsRunning) {
            isRunning.value = nextIsRunning
        },
    })

    watch(languageId, (nextLanguageId, previousLanguageId) => {
        codeDrafts.value[previousLanguageId] = codeText.value
        codeText.value = codeDrafts.value[nextLanguageId] ?? defaultCodeByLanguage[nextLanguageId] ?? ''
    })

    watch(codeText, (nextCodeText) => {
        codeDrafts.value[languageId.value] = nextCodeText
    })

    /**
     * Clear terminal output.
     * @returns {void}
     */
    function clearTerminal() {
        terminalLines.value = []
    }

    /**
     * Run code in the selected language.
     * @returns {void}
     */
    function runCode() {
        if (languageId.value === 'javascript') {
            jsRunner.run(codeText.value)
            return
        }

        if (languageId.value === 'cpp') {
            cppRunner.run(codeText.value)
            return
        }

        terminalLines.value = [
            ...terminalLines.value,
            {
                id: `${Date.now()}-unsupported`,
                type: 'error',
                text: 'This prototype only runs JavaScript in-browser. Other languages require backend compilation.',
            },
        ]
    }

    /**
     * Stop current run.
     * @returns {void}
     */
    function stopRun() {
        if (languageId.value === 'cpp') {
            cppRunner.stop()
            return
        }

        jsRunner.stop()
    }

    onBeforeUnmount(() => {
        jsRunner.stop()
        cppRunner.stop()
    })
</script>

<template>
    <div class="page">
        <div class="pageHeader">
            <div class="headerLeft">
                <h1 class="title">Playground</h1>
                <div class="subtitle">JavaScript runs locally. C++ runs through the backend API.</div>
            </div>
            <div class="headerRight">
                <label class="selectLabel">
                    <span class="selectText">Language</span>
                    <select v-model="languageId" class="select">
                        <option value="javascript">JavaScript</option>
                        <option value="cpp">C++</option>
                        <option value="python">Python (later)</option>
                    </select>
                </label>

                <button class="primaryButton" type="button" :disabled="isRunning" @click="runCode">
                    Run
                </button>
                <button class="button" type="button" :disabled="!isRunning" @click="stopRun">
                    Stop
                </button>
            </div>
        </div>

        <div class="grid">
            <CodeEditor v-model="codeText" :language-id="languageId" />
            <TerminalPanel :lines="terminalLines" :is-running="isRunning" @clear="clearTerminal" />
        </div>
    </div>
</template>

<style scoped>
    .page {
        max-width: 1100px;
    }

    .pageHeader {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 12px;
        margin-bottom: 12px;
    }

    .headerLeft {
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

    .headerRight {
        display: inline-flex;
        align-items: center;
        gap: 10px;
        flex: 0 0 auto;
    }

    .selectLabel {
        display: inline-flex;
        align-items: center;
        gap: 8px;
        padding: 6px 10px;
        border-radius: 10px;
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
    }

    .selectText {
        font-size: 12px;
        color: var(--mutedTextColor);
    }

    .select {
        border: 0;
        background: transparent;
        color: var(--textColor);
        outline: none;
        font: inherit;
        font-size: 13px;
    }

    .select option {
        color: black;
        background: white;
    }

    .button,
    .primaryButton {
        border: 1px solid var(--borderColor);
        background: var(--surfaceColor);
        color: var(--textColor);
        font-size: 13px;
        padding: 8px 12px;
        border-radius: 10px;
        cursor: pointer;
    }

    .primaryButton {
        border-color: var(--accentBorderColor);
        background: var(--accentSurfaceColor);
    }

    .button:hover:enabled,
    .primaryButton:hover:enabled {
        border-color: var(--borderColorHover);
    }

    .button:disabled,
    .primaryButton:disabled {
        opacity: 0.65;
        cursor: not-allowed;
    }

    .grid {
        display: grid;
        grid-template-columns: 1fr;
        gap: 12px;
    }

    @media (min-width: 980px) {
        .grid {
            grid-template-columns: 1fr 1fr;
            align-items: stretch;
        }
    }
</style>
