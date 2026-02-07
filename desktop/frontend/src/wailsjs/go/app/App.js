// Stub bindings. Run "wails generate module" or "wails dev" to generate real bindings.
import { Invoke } from '../../runtime/runtime.js'

export function StartPriceStream(serverAddr) {
  return Invoke('StartPriceStream', serverAddr)
}

export function StopPriceStream() {
  return Invoke('StopPriceStream')
}
