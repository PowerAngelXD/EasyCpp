/**
 * @typedef {'log'|'error'|'system'} TerminalLineType
 */

/**
 * @typedef {object} TerminalLine
 * @property {string} id Unique id for rendering.
 * @property {TerminalLineType} type Line type.
 * @property {string} text Text content.
 */

/**
 * Create a JS runner backed by Web Worker.
 * @param {object} args Runner configuration.
 * @param {(line: TerminalLine) => void} args.onLine Callback for each output line.
 * @param {(isRunning: boolean) => void} args.onRunningChange Callback for running state.
 * @returns {{ run: (codeText: string) => void, stop: () => void }} Runner API.
 */
export function createJsRunner({ onLine, onRunningChange }) {
    /** @type {Worker|null} */
    let activeWorker = null

    /** @type {number|null} */
    let timeoutId = null

    /**
     * @param {TerminalLineType} type Line type.
     * @param {string} text Output text.
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
        if (timeoutId !== null) {
            window.clearTimeout(timeoutId)
            timeoutId = null
        }
        if (activeWorker) {
            activeWorker.terminate()
            activeWorker = null
        }
        onRunningChange(false)
    }

    /**
     * Start execution of the provided JS code text.
     * @param {string} codeText The JavaScript code to execute.
     * @returns {void}
     */
    function run(codeText) {
        cleanup()
        onRunningChange(true)

        pushLine('system', 'Running JavaScript...')

        activeWorker = new Worker(new URL('../workers/jsRunnerWorker.js', import.meta.url), {
            type: 'module',
        })

        activeWorker.onmessage = (event) => {
            const messageData = event?.data
            if (!messageData || typeof messageData.type !== 'string') {
                return
            }

            if (messageData.type === 'log') {
                pushLine('log', String(messageData?.payload?.text ?? ''))
                return
            }

            if (messageData.type === 'error') {
                pushLine('error', String(messageData?.payload?.text ?? ''))
                return
            }

            if (messageData.type === 'done') {
                cleanup()
            }
        }

        activeWorker.onerror = () => {
            pushLine('error', 'Worker crashed.')
            cleanup()
        }

        timeoutId = window.setTimeout(() => {
            pushLine('error', 'Timeout: execution took too long.')
            cleanup()
        }, 2000)

        activeWorker.postMessage({ type: 'run', payload: { codeText } })
    }

    /**
     * Stop currently running job.
     * @returns {void}
     */
    function stop() {
        pushLine('system', 'Stopped.')
        cleanup()
    }

    return { run, stop }
}
