<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { StartPriceStream, StopPriceStream } from './wailsjs/go/app/App'
import { EventsOn } from './wailsjs/runtime/runtime'

const price = ref<number | null>(null)
const symbol = ref('')
const atTs = ref<number | null>(null)
const error = ref('')
const connected = ref(false)

function connect() {
  error.value = ''
  connected.value = true
  StartPriceStream('').catch((err: unknown) => {
    error.value = err instanceof Error ? err.message : String(err)
    connected.value = false
  })
}

function disconnect() {
  StopPriceStream()
  connected.value = false
}

onMounted(() => {
  EventsOn('price', (payload: { symbol?: string; price?: number; atTs?: number }) => {
    if (payload?.price != null) {
      price.value = payload.price
      symbol.value = payload.symbol ?? 'BTC'
      atTs.value = payload.atTs ?? null
    }
  })
  EventsOn('priceError', (msg: string) => {
    error.value = msg
    connected.value = false
  })
})

onUnmounted(() => {
  disconnect()
})
</script>

<template>
  <div class="app">
    <h1>BTC Price Tracker</h1>
    <p v-if="connected" class="status">Connected to gRPC stream</p>
    <p v-else class="status idle">Not connected</p>
    <div v-if="error" class="error">{{ error }}</div>
    <div class="price-block">
      <span class="symbol">{{ symbol || '—' }}</span>
      <span class="price">{{ price != null ? price.toLocaleString(undefined, { minimumFractionDigits: 2 }) : '—' }}</span>
      <span v-if="atTs != null" class="at-ts">{{ new Date(atTs * 1000).toLocaleTimeString() }}</span>
    </div>
    <div class="actions">
      <button v-if="!connected" @click="connect">Connect</button>
      <button v-else @click="disconnect">Disconnect</button>
    </div>
  </div>
</template>

<style scoped>
.app {
  font-family: system-ui, sans-serif;
  padding: 1.5rem;
  max-width: 360px;
}
h1 {
  font-size: 1.25rem;
  margin: 0 0 0.5rem 0;
}
.status {
  font-size: 0.875rem;
  color: #0a0;
  margin: 0 0 1rem 0;
}
.status.idle {
  color: #666;
}
.error {
  color: #c00;
  font-size: 0.875rem;
  margin-bottom: 0.5rem;
}
.price-block {
  background: #f0f0f0;
  padding: 1rem;
  border-radius: 8px;
  margin-bottom: 1rem;
}
.symbol {
  font-weight: 600;
  margin-right: 0.5rem;
}
.price {
  font-size: 1.5rem;
}
.at-ts {
  display: block;
  font-size: 0.75rem;
  color: #666;
  margin-top: 0.25rem;
}
.actions button {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  cursor: pointer;
}
</style>
