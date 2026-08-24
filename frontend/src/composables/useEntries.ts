import { ref, type Ref } from "vue";
import { apiGetEntries } from "@/api/api";
import { getEntriesForList, putEntry, deleteEntryLocal, replaceEntriesForList } from "@/offline/db";
import {
  createEntry as mutateCreateEntry,
  updateEntry as mutateUpdateEntry,
  deleteEntry as mutateDeleteEntry,
} from "@/offline/entryMutations";
import { flushOutbox, hasQueuedRecordsForList } from "@/offline/outbox";
import type { Entry } from "@/types";

// Per-list reactive caches, shared by every component viewing the same list
// (e.g. List.vue and EntryEditor.vue), so they never diverge and data
// survives navigating away and back without a refetch flash.
const listCaches = new Map<string, Ref<Entry[]>>();

function getListRef(listID: string): Ref<Entry[]> {
  let entries = listCaches.get(listID);
  if (entries === undefined) {
    entries = ref<Entry[]>([]);
    listCaches.set(listID, entries);
  }
  return entries;
}

// upsertLocal replaces the entry with a matching ID in place, or appends it
// if not present. Used both for optimistic local writes and for reconciling
// SSE events, which may echo back this tab's own writes.
function upsertLocal(entries: Ref<Entry[]>, entry: Entry) {
  const index = entries.value.findIndex((item) => item.ID === entry.ID);
  if (index === -1) {
    entries.value.push(entry);
  } else {
    entries.value[index] = entry;
  }
}

function removeLocal(entries: Ref<Entry[]>, id: string) {
  const index = entries.value.findIndex((item) => item.ID === id);
  if (index !== -1) {
    entries.value.splice(index, 1);
  }
}

const BASE_RECONNECT_DELAY = 1_000;
const MAX_RECONNECT_DELAY = 30_000;

export function useEntries(listID: string) {
  const entries = getListRef(listID);

  let eventSource: EventSource | null = null;
  let eventSourceTimeout: number | null = null;
  let restartTimeout: number | null = null;
  let reconnectDelay = BASE_RECONNECT_DELAY;
  let hydrated = false;

  async function hydrate() {
    if (hydrated) {
      return;
    }
    hydrated = true;
    const local = await getEntriesForList(listID);
    if (local.length > 0) {
      entries.value = local;
    }
  }

  function createEventSourceTimeout() {
    eventSourceTimeout = window.setTimeout(() => restart(), 15_000);
  }
  function deleteEventSourceTimeout() {
    if (eventSourceTimeout !== null) {
      clearTimeout(eventSourceTimeout);
      eventSourceTimeout = null;
    }
  }
  function resetEventSourceTimeout() {
    deleteEventSourceTimeout();
    createEventSourceTimeout();
  }

  // restart schedules a reconnect after a growing backoff (capped), so a
  // prolonged outage does not hammer the server/battery with retries every
  // second. It resets whenever a connection attempt succeeds or an external
  // signal (tab becomes visible, browser reports back online) suggests it is
  // worth trying again right away.
  function restart() {
    if (restartTimeout !== null) {
      return;
    }
    restartTimeout = window.setTimeout(() => {
      restartTimeout = null;
      void start();
    }, reconnectDelay);
    reconnectDelay = Math.min(reconnectDelay * 2, MAX_RECONNECT_DELAY);
  }

  // forceReconnect cancels any pending backoff wait and reconnects
  // immediately. Used when we have a good reason to believe the connection
  // can work again (tab became visible, browser fired an "online" event).
  function forceReconnect() {
    if (restartTimeout !== null) {
      clearTimeout(restartTimeout);
      restartTimeout = null;
    }
    reconnectDelay = BASE_RECONNECT_DELAY;
    void start();
  }

  // eventSource keeps the entries updated in the background via
  // Server-Sent-Events. If the connection fails, we do not rely on
  // reconnection from SSE, but create a new one ourselves.
  function createEventSource() {
    eventSource = new EventSource(`/api/v1/entries/events?ListID=${listID}`);
    resetEventSourceTimeout();

    // Do not use the SSE reconnection, because it works differently in
    // every browser.
    eventSource.onerror = () => restart();
    eventSource.addEventListener("ping", () => resetEventSourceTimeout());
    eventSource.addEventListener("create", (event: MessageEvent) => {
      const entry = JSON.parse(event.data) as Entry;
      upsertLocal(entries, entry);
      void putEntry(entry);
    });
    eventSource.addEventListener("update", (event: MessageEvent) => {
      const entry = JSON.parse(event.data) as Entry;
      upsertLocal(entries, entry);
      void putEntry(entry);
    });
    eventSource.addEventListener("delete", (event: MessageEvent) => {
      const entry = JSON.parse(event.data) as Entry;
      removeLocal(entries, entry.ID);
      void deleteEntryLocal(entry.ID);
    });
  }
  function deleteEventSource() {
    deleteEventSourceTimeout();
    if (eventSource !== null) {
      eventSource.close();
      eventSource = null;
    }
  }

  // starting guards against concurrent start() calls racing each other, e.g.
  // when a scheduled restart fires at the same time the tab becomes visible.
  let starting = false;
  async function start() {
    if (starting) {
      return;
    }
    starting = true;
    try {
      deleteEventSource();

      // Load whatever is cached locally first, so the view has something to
      // show immediately (and everything, if offline).
      await hydrate();

      // Drain queued writes before pulling the server's list: overwriting
      // the local cache with a fresh fetch while a create/update/delete is
      // still unsynced would wipe out that optimistic change. If the outbox
      // couldn't fully drain (still offline, or the server is rejecting
      // writes), skip the destructive refetch this cycle and retry later
      // instead of clobbering local state.
      await flushOutbox();
      if (await hasQueuedRecordsForList(listID)) {
        restart();
        return;
      }

      try {
        const fresh = await apiGetEntries(listID);
        entries.value = fresh;
        await replaceEntriesForList(listID, fresh);
      } catch {
        restart();
        return;
      }

      // connection is working again
      reconnectDelay = BASE_RECONNECT_DELAY;

      createEventSource();
    } finally {
      starting = false;
    }
  }

  function stop() {
    deleteEventSource();
    if (restartTimeout !== null) {
      clearTimeout(restartTimeout);
      restartTimeout = null;
    }
  }

  function handleVisibilityChange() {
    if (document.visibilityState === "visible") {
      forceReconnect();
    }
  }
  function handleOnline() {
    forceReconnect();
  }

  function listen() {
    document.addEventListener("visibilitychange", handleVisibilityChange);
    window.addEventListener("online", handleOnline);
  }
  function unlisten() {
    document.removeEventListener("visibilitychange", handleVisibilityChange);
    window.removeEventListener("online", handleOnline);
  }

  async function createEntry(name: string): Promise<void> {
    const entry = await mutateCreateEntry(name, listID);
    upsertLocal(entries, entry);
  }

  async function updateEntry(entry: Entry): Promise<void> {
    const updated = await mutateUpdateEntry(entry);
    upsertLocal(entries, updated);
  }

  async function deleteEntry(entry: Entry): Promise<void> {
    await mutateDeleteEntry(entry);
    removeLocal(entries, entry.ID);
  }

  return {
    entries,
    start,
    stop,
    listen,
    unlisten,
    createEntry,
    updateEntry,
    deleteEntry,
  };
}
