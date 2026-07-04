<script setup lang="ts">
import { onMounted, ref, computed, onUnmounted } from "vue";
import { useRoute } from "vue-router";

import {
  apiGetEntries,
  apiCreateEntry,
  apiGetTypes,
  apiDeleteEntry,
  apiUpdateEntry,
} from "@/api/api.ts";
import { router } from "@/router.ts";
import { useNotificationManager } from "@/composables/useNotificationManager";
import { type Entry, type Type, type ContextmenuAction } from "@/types.ts";

import DefaultLayout from "@/Layouts/DefaultLayout.vue";
import Button from "@/Components/Button.vue";
import EntryItem from "@/Components/EntryItem.vue";
import ShareButton from "@/Components/ShareButton.vue";
import Contextmenu from "@/Components/Contextmenu.vue";

const { show, clear } = useNotificationManager();

const route = useRoute();
const listID = route.params.id as string;
const entries = ref<Entry[]>([]);
const types = ref<Type[]>([]);
const searchInput = ref<string>("");

const contextmenuVisible = ref(false);
const contextmenuX = ref(0);
const contextmenuY = ref(0);
const contextmenuActions = [
  { label: "Edit", action: "edit" },
  { label: "Delete", action: "delete" },
] as ContextmenuAction[];
const contextmenuTarget = ref<Entry>();

function contextmenuShow(event: MouseEvent, entry: Entry) {
  event.preventDefault();
  if (contextmenuVisible.value === true) {
    contextmenuClose();
  }
  contextmenuVisible.value = true;
  contextmenuTarget.value = entry;
  contextmenuX.value = event.pageX;
  contextmenuY.value = event.pageY;
  document.addEventListener("click", contextmenuHandleOutsideClick);
}

function contextmenuHandleOutsideClick(event: Event) {
  if (
    !document
      .getElementById("contextmenu")
      ?.contains(event?.target as HTMLElement)
  ) {
    contextmenuClose();
  }
}

function contextmenuClose() {
  contextmenuVisible.value = false;
  if (contextmenuTarget.value !== undefined) {
    contextmenuTarget.value = undefined;
  }
  contextmenuX.value = 0;
  contextmenuY.value = 0;
  document.removeEventListener("click", contextmenuHandleOutsideClick);
}

function contextmenuHandle(action: string) {
  const targetEntry = contextmenuTarget.value as Entry;
  contextmenuClose();

  if (action === "delete") {
    try {
      apiDeleteEntry(targetEntry);
      entries.value = entries.value.filter(
        (entry) => targetEntry.ID !== entry?.ID,
      );
    } catch (error) {
      show("error", "Unable to delete entry", { logMessage: error });
      return;
    }
  } else if (action === "edit") {
    router.push({ name: "entryEditor", params: { id: targetEntry.ID } });
  }
}

const activeEntries = computed(() => {
  return entries.value.filter((entry: Entry) =>
    entry.Name.toLowerCase().includes(searchInput.value.toLowerCase()),
  );
});

const activeSortedNotBoughtEntriesbyType = computed(() => {
  const groups: Record<string, Entry[]> = {};
  types.value.forEach((type: Type) => {
    groups[type.ID] = activeEntries.value
      .filter((entry: Entry) => !entry.Bought && entry.TypeID === type.ID)
      .sort((a: Entry, b: Entry) => {
        // Sort by name
        return a.Name.localeCompare(b.Name);
      });
  });
  return groups;
});

const activeBoughtEntries = computed(() => {
  return activeEntries.value.filter((entry) => entry.Bought);
});

async function getTypes() {
  try {
    types.value = await apiGetTypes();
  } catch (error) {
    show("error", "Unable to load entry types", { logMessage: error });
  }
}

async function ensureEntryOnNotBoughtList() {
  if (searchInput.value === "") {
    return;
  }

  // Create new Entry
  try {
    await apiCreateEntry(searchInput.value, listID);
  } catch (error) {
    show("error", "Unable to create new entry", { logMessage: error });
  } finally {
    // Reset input
    searchInput.value = "";
  }
}

async function updateEntry(entry: Entry) {
  try {
    await apiUpdateEntry(entry);
  } catch (error) {
    show("error", "Unable to update entry", { logMessage: error });
    return;
  } finally {
    // Reset input
    searchInput.value = "";
  }
}

// connectionIssue tracks the notification shown while the connection is down,
// so repeated failures do not stack multiple toasts and we can clear it once
// things recover.
let connectionIssue: string | undefined = undefined;
function reportConnectionIssue(error?: unknown) {
  if (connectionIssue === undefined) {
    connectionIssue = show("error", "Connection lost, trying to reconnect", {
      logMessage: error,
      autoHide: false,
    }).id;
  }
}
function clearConnectionIssue() {
  if (connectionIssue !== undefined) {
    clear(connectionIssue);
    connectionIssue = undefined;
    show("info", "Connection restored");
  }
}

// reconnectDelay grows with each consecutive failure (up to a cap) so a
// prolonged outage does not hammer the server/battery with retries every
// second. It resets whenever a connection attempt succeeds or an external
// signal (tab becomes visible, browser reports back online) suggests it is
// worth trying again right away.
const BASE_RECONNECT_DELAY = 1_000;
const MAX_RECONNECT_DELAY = 30_000;
let reconnectDelay = BASE_RECONNECT_DELAY;
let restartTimeout = <number | null>null;

function restart(error?: unknown) {
  reportConnectionIssue(error);
  if (restartTimeout !== null) {
    return;
  }
  restartTimeout = setTimeout(() => {
    restartTimeout = null;
    start();
  }, reconnectDelay);
  reconnectDelay = Math.min(reconnectDelay * 2, MAX_RECONNECT_DELAY);
}

// forceReconnect cancels any pending backoff wait and reconnects immediately.
// Used when we have a good reason to believe the connection can work again
// (tab became visible, browser fired an "online" event).
function forceReconnect() {
  if (restartTimeout !== null) {
    clearTimeout(restartTimeout);
    restartTimeout = null;
  }
  reconnectDelay = BASE_RECONNECT_DELAY;
  start();
}

// eventSourceTimeout is to prevent a connection without any traffic being keeps open. The server should send a ping every 5 seconds
let eventSourceTimeout = <number | null>null;
function createEventSourceTimeout() {
  eventSourceTimeout = setTimeout(
    () => restart("no ping received in time"),
    15_000,
  );
}
function deleteEventSourceTimeout() {
  if (eventSourceTimeout !== null) {
    clearTimeout(eventSourceTimeout);
  }
}
function resetEventSourceTimeout() {
  deleteEventSourceTimeout();
  createEventSourceTimeout();
}

// eventSource keeps the entries updated in the background, with an Sever-Sent Event.
// If the connection fails, we do not rely on reconnection from SSE, but just create a new one ourself.
let eventSource = <EventSource | null>null;
function createEventSource() {
  eventSource = new EventSource(`/api/v1/entries/events?ListID=${listID}`);

  resetEventSourceTimeout();

  // Do not use the SSE Reconnection, because it works different on all browser
  eventSource.onerror = (event) => {
    restart(event);
  };
  eventSource.addEventListener("ping", () => {
    resetEventSourceTimeout();
  });
  eventSource.addEventListener("create", (event: MessageEvent) => {
    const entry = JSON.parse(event.data) as Entry;
    entries.value.push(entry);
  });
  eventSource.addEventListener("update", async (event: MessageEvent) => {
    const entry = JSON.parse(event.data) as Entry;
    const index = entries.value.findIndex((item) => item.ID === entry.ID);
    if (index === -1) {
      return;
    }
    entries.value[index] = entry;
  });
  eventSource.addEventListener("delete", async (event: MessageEvent) => {
    const entry = JSON.parse(event.data) as Entry;
    const index = entries.value.findIndex((item) => item.ID === entry.ID);
    if (index === -1) {
      return;
    }
    entries.value.splice(index, 1);
  });
  return eventSource;
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
    // Clear Event Source
    deleteEventSource();

    // Get all entries initially
    try {
      const newEntries = await apiGetEntries(listID);
      entries.value = newEntries;
    } catch (error) {
      restart(error);
      return;
    }

    // connection is working again, clear any lingering notice
    clearConnectionIssue();
    reconnectDelay = BASE_RECONNECT_DELAY;

    // Create Event Source to get updates
    createEventSource();
  } finally {
    starting = false;
  }
}

// beforeLeave fixes the width of the element before the transition leave is
// applied. This way, when setting the position to absolute the element is not
// shrinked due to the flexbox.
// See https://github.com/vuejs/vue/issues/9713#issuecomment-572153283
function beforeLeave(el: Element) {
  const { width } = window.getComputedStyle(el);
  // This is a silly cast, but the vue hook absolutly wants a Element although it is actually a HTMLElement
  (el as HTMLElement).style.width = width;
}

// handleVisibilityChange reconnects right away when the tab/app becomes
// visible again. Mobile browsers commonly freeze timers and drop network
// connections while backgrounded, so waiting for the regular timeout/backoff
// would otherwise leave the view stuck until the user restarts the app.
function handleVisibilityChange() {
  if (document.visibilityState === "visible") {
    forceReconnect();
  }
}
function handleOnline() {
  forceReconnect();
}
function handleOffline() {
  reportConnectionIssue("browser reported offline");
}

onMounted(() => {
  getTypes();
  start();
  document.addEventListener("visibilitychange", handleVisibilityChange);
  window.addEventListener("online", handleOnline);
  window.addEventListener("offline", handleOffline);
});
onUnmounted(() => {
  deleteEventSource();
  if (restartTimeout !== null) {
    clearTimeout(restartTimeout);
    restartTimeout = null;
  }
  document.removeEventListener("visibilitychange", handleVisibilityChange);
  window.removeEventListener("online", handleOnline);
  window.removeEventListener("offline", handleOffline);
});
</script>

<template>
  <DefaultLayout>
    <template v-slot:header>
      <ShareButton></ShareButton>
      <input v-model="searchInput" @keydown.enter="ensureEntryOnNotBoughtList" type="search" autocomplete="off"
        placeholder="Search" />
      <Button @click="ensureEntryOnNotBoughtList" class="inverted">+</Button>
    </template>
    <template v-slot:main>
      <TransitionGroup name="list" @before-leave="beforeLeave" tag="ul">
        <template v-for="type in types">
          <div :key="type.ID" v-if="activeSortedNotBoughtEntriesbyType[type.ID].length > 0" class="divider">
            {{ type.Name }}
          </div>
          <li :style="{ borderLeft: `0.3em solid ${type.Color}` }"
            v-for="(entry, i) in activeSortedNotBoughtEntriesbyType[type.ID]" :key="entry.ID">
            <EntryItem v-model="activeSortedNotBoughtEntriesbyType[type.ID][i]"
              @contextmenu="contextmenuShow($event, entry)" @update="updateEntry(entry)">
            </EntryItem>
          </li>
        </template>

        <li key="bought" v-if="activeBoughtEntries.length > 0" class="divider bought">
          Recently bought
        </li>
        <li v-for="(entry, i) in activeBoughtEntries" :key="entry.ID">
          <EntryItem v-model="activeBoughtEntries[i]" @contextmenu="contextmenuShow($event, entry)"
            @update="updateEntry(entry)">
          </EntryItem>
        </li>
      </TransitionGroup>
      <Contextmenu v-if="contextmenuVisible" id="contextmenu" :x="contextmenuX" :y="contextmenuY"
        :actions="contextmenuActions" @action="contextmenuHandle" />
    </template>
  </DefaultLayout>
</template>

<style>
hr {
  color: var(--color-primary);
  margin: 0.5em 1em;
}

.list-move,
/* apply transition to moving elements */
.list-enter-active,
.list-leave-active {
  transition: all 0.5s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
}

/* ensure leaving items are taken out of layout flow so that moving
   animations can be calculated correctly. */
.list-leave-active {
  position: absolute;
}

ul {
  /* no bullet points */
  list-style: none;
  padding: 0em;
  margin: 0em;
}

header input {
  flex-grow: 1;
  /* max width possible */
}

input {
  height: 2.5em;
  box-sizing: border-box;
  padding: 0.1em 0.2em;
  border-radius: var(--border-radius);
  color: var(--color-text);
  border: 0.1em solid var(--color-border);
  background-color: var(--color-surface);
}

.divider {
  display: flex;
  align-items: center;
  text-align: center;
  font-weight: 300;
}

.bought {
  margin: 1em 0em;
}

.divider::before,
.divider::after {
  content: "";
  flex: 1;
  border-bottom: 0.1em solid #000;
  margin: 0 0.5em;
}
</style>
