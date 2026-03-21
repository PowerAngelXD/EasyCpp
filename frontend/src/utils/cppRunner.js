const defaultApiBaseURL = import.meta.env.VITE_API_BASE_URL ?? 'http://127.0.0.1:8080'

/**
 * @typedef {'log'|'error'|'system'} TerminalLineType
 */

/**
 * @typedef {object} TerminalLine
 * @property {string} id
 * @property {TerminalLineType} type
 * @property {string} text
 */

/**
 * Create a C++ runner backed by backend API.
 * @param {object} args
 * @param {(line: TerminalLine) => void} args.onLine
 * @param {(isRunning: boolean) => void} args.onRunningChange
 * @returns {{ run: (codeText: string, stdinText?: string) => void, stop: () => void }}
 */
export function createCppRunner({ onLine, onRunningChange }) {
    /** @type {AbortController|null} */
    let activeController = null

    /**
     * @param {TerminalLineType} type
     * @param {string} text
     * @returns {void}
     */
    function pushLine(type, text) {
        const lineId = `${Date.now()}-${Math.random().toString(16).slice(2)}`
        onLine({ id: lineId, type, text })
    }

    /**
     * @returns {void}
     */
    function cleanup() {
        activeController = null
        onRunningChange(false)
    }

    /**
     * @param {string} codeText
     * @param {string} [stdinText]
     * @returns {Promise<void>}
     */
    async function run(codeText, stdinText = '') {
        if (activeController) {
            activeController.abort()
        }

        activeController = new AbortController()
        onRunningChange(true)
        pushLine('system', 'Compiling C++...')

        try {
            const response = await fetch(`${defaultApiBaseURL}/api/v1/ide/cpp/run`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    code: codeText,
                    stdin: stdinText,
                    timeLimitMs: 2000,
                }),
                signal: activeController.signal,
            })

            if (!response.ok) {
                pushLine('error', `Request failed: HTTP ${response.status}`)
                return
            }

            const payload = await response.json()
            const data = payload?.data
            if (!data) {
                pushLine('error', 'Invalid response format.')
                return
            }

            if (!data?.compile?.succeeded) {
                pushLine('error', data?.compile?.stderr || 'Compilation failed.')
                return
            }

            pushLine('system', 'Compilation succeeded. Running...')

            const runData = data?.run
            if (!runData) {
                pushLine('error', 'Run phase is missing.')
                return
            }

            if (runData.stdout) {
                pushLine('log', String(runData.stdout))
            }
            if (runData.stderr) {
                pushLine('error', String(runData.stderr))
            }

            if (runData.timedOut) {
                pushLine('error', 'Execution timed out.')
            }

            pushLine('system', `Done in ${data.durationMs} ms.`)
        } catch (error) {
            if (error instanceof DOMException && error.name === 'AbortError') {
                pushLine('system', 'Execution cancelled.')
                return
            }
            pushLine('error', 'Failed to run code. Check backend service.')
        } finally {
            cleanup()
        }
    }

    /**
     * @returns {void}
     */
    function stop() {
        if (activeController) {
            activeController.abort()
        }
    }

    return { run, stop }
}
