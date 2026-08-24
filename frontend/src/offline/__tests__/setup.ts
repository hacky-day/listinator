import "fake-indexeddb/auto";

// The "node" test environment has neither global, and offline/outbox.ts
// registers a "window online" listener as a module-level side effect (so it
// fires even if no component ever calls useEntries/useConnectivity) —
// importing it would otherwise throw at module load time. Node's globalThis
// already implements EventTarget and the timer functions outbox.ts needs,
// so aliasing window to it is enough; the listener is never dispatched in
// these tests.
const globalRecord = globalThis as unknown as Record<string, unknown>;
if (typeof window === "undefined") {
  // Node's globalThis isn't an EventTarget, so a plain stub covers what
  // outbox.ts needs: registering (never dispatching, in these tests) an
  // "online" listener, and the safety-net setInterval/clearInterval.
  globalRecord.window = {
    addEventListener: () => {},
    removeEventListener: () => {},
    setInterval,
    clearInterval,
    setTimeout,
    clearTimeout,
  };
}
if (typeof navigator === "undefined") {
  globalRecord.navigator = { onLine: true };
}
