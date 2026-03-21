<script setup>
    import { onBeforeUnmount, ref } from 'vue'

    import CodeEditor from '../components/CodeEditor.vue'
    import TerminalPanel from '../components/TerminalPanel.vue'
    import { createCppRunner } from '../utils/cppRunner'

    const languageId = ref('cpp')

    const codeText = ref(`#include <iostream>
using namespace std;

int main() {
    cout << "Hello, EasyCpp!" << endl;
    return 0;
}
`)

    const terminalLines = ref([])
    const isRunning = ref(false)

    const cppRunner = createCppRunner({
        onLine(line) {
            terminalLines.value = [...terminalLines.value, line]
        },
        onRunningChange(nextIsRunning) {
            isRunning.value = nextIsRunning
        },
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
        cppRunner.run(codeText.value)
    }

    /**
     * Stop current run.
     * @returns {void}
     */
    function stopRun() {
        cppRunner.stop()
    }

    onBeforeUnmount(() => {
        cppRunner.stop()
    })
</script>

<template>
    <div class="page">
        <div class="pageHeader">
            <div class="headerLeft">
                <h1 class="title">Playground</h1>
                <div class="subtitle">C++ online IDE (compile + run via backend API)</div>
            </div>
            <div class="headerRight">
                <div class="selectLabel">
                    <span class="selectText">Language</span>
                    <span class="selectValue">C++</span>
                </div>

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

    .selectValue {
        color: var(--textColor);
        font-size: 13px;
        font-weight: 600;
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
