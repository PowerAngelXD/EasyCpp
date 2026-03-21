/**
 * Execute user-provided JavaScript in a Worker and stream logs back to UI.
 *
 * Message protocol:
 * - from main: { type: 'run', payload: { codeText: string } }
 * - to main:   { type: 'log'|'error'|'done', payload: { text: string } }
 */

/**
 * @param {unknown} value Any value logged by user code.
 * @returns {string} A readable string for terminal output.
 */
function formatTerminalValue(value) {
    if (typeof value === 'string') {
        return value
    }

    try {
        return JSON.stringify(value, null, 2)
    } catch {
        return String(value)
    }
}

/**
 * @param {unknown[]} args Console arguments.
 * @returns {string} Joined string.
 */
function joinConsoleArgs(args) {
    return args.map(formatTerminalValue).join(' ')
}

self.onmessage = (event) => {
    const messageData = event?.data
    if (!messageData || messageData.type !== 'run') {
        return
    }

    const codeText = String(messageData?.payload?.codeText ?? '')

    const originalConsoleLog = console.log
    const originalConsoleError = console.error

    console.log = (...args) => {
        self.postMessage({ type: 'log', payload: { text: joinConsoleArgs(args) } })
        originalConsoleLog(...args)
    }

    console.error = (...args) => {
        self.postMessage({ type: 'error', payload: { text: joinConsoleArgs(args) } })
        originalConsoleError(...args)
    }

    try {
        const runnable = new Function(codeText)
        runnable()
        self.postMessage({ type: 'done', payload: { text: '' } })
    } catch (error) {
        const errorText = error instanceof Error ? `${error.name}: ${error.message}` : String(error)
        self.postMessage({ type: 'error', payload: { text: errorText } })
        self.postMessage({ type: 'done', payload: { text: '' } })
    } finally {
        console.log = originalConsoleLog
        console.error = originalConsoleError
    }
}
