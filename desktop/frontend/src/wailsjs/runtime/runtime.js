// Stub runtime used when bindings are not yet generated. Run "wails generate module" or "wails dev" to generate real bindings.
const bridge = typeof window !== 'undefined' ? window : undefined

export function EventsOn(eventName, callback) {
  if (bridge?.runtime?.EventsOn) {
    return bridge.runtime.EventsOn(eventName, callback)
  }
  const key = '__wails_events_' + eventName
  if (!bridge[key]) bridge[key] = []
  bridge[key].push(callback)
  return () => {
    bridge[key] = (bridge[key] || []).filter((cb) => cb !== callback)
  }
}

export function Invoke(method, ...args) {
  if (bridge?.go?.main?.App?.[method]) {
    return bridge.go.main.App[method](...args)
  }
  return Promise.reject(new Error('Go backend not available (run in Wails)'))
}
