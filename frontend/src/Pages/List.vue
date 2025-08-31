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
import { type Entry, type Type, type ContextmenuAction } from "@/types.ts";
import DefaultLayout from "@/Layouts/DefaultLayout.vue";
import Button from "@/Components/Button.vue";
import EntryItem from "@/Components/EntryItem.vue";
import ShareButton from "@/Components/ShareButton.vue";
import Contextmenu from "@/Components/Contextmenu.vue";
import { useNotificationManager } from "@/composables/useNotificationManager";

const { show, clear } = useNotificationManager();

let failureMessage: string | undefined = undefined;

const route = useRoute();
const listID = route.params.id as string;
const entries = ref<Entry[]>([]);
const types = ref<Type[]>([]);
const searchInput = ref<string>("");

const contextmenuVisible = ref(false);
const contextmenuX = ref(0);
const contextmenuY = ref(0);
const contextmenuActions = [
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
  }
}

const activeEntries = computed(() => {
  return entries.value.filter((entry: Entry) =>
    entry.Name.toLowerCase().includes(searchInput.value.toLowerCase()),
  );
});

const activeSortedNotBoughtEntries = computed(() => {
  return activeEntries.value
    .filter((entry) => !entry.Bought)
    .sort((a, b) => {
      // Sort by type
      const typeCompare = a.TypeID.localeCompare(b.TypeID);
      if (typeCompare !== 0) return typeCompare;

      // Sort by name
      return a.Name.localeCompare(b.Name);
    });
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

  try {
    // If only one entry matches the filter, reactivate it if bought, otherwise do nothing.
    if (activeEntries.value.length === 1) {
      const entry = activeEntries.value[0];
      if (entry.Bought === false) {
        return;
      }
      entry.Bought = false;
      await updateEntry(entry);
      return;
    }

    // Create new Entry
    try {
      await apiCreateEntry(searchInput.value, listID);
    } catch (error) {
      show("error", "Unable to create new entry", { logMessage: error });
    }
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
  }
}

// eventSourceTimeout is to prevent a connection without any traffic being keeps open. The server should send a ping every 5 seconds
let eventSourceTimeout = <number | null>null;
function createEventSourceTimeout() {
  eventSourceTimeout = setTimeout(restart, 15_000);
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
    console.log(event);
    restart();
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

async function start() {
  // Clear Event Source
  deleteEventSource();

  // Get all entries initially
  try {
    const newEntries = await apiGetEntries(listID);
    entries.value = newEntries;
  } catch (error) {
    // if it's the first failure, show a message and do not remove it automatically
    if (failureMessage === undefined) {
      failureMessage = show("error", "Unable to load entries", {
        logMessage: error,
        autoHide: false,
      }).id;
    }

    restart();
    return;
  }

  // remove the failure and show a short notice that loading is working again
  if (failureMessage !== undefined) {
    clear(failureMessage);
    failureMessage = undefined;
    show("info", "Entries loaded again");
  }

  // Create Event Source to get updates
  createEventSource();
}

function restart() {
  setTimeout(start, 1000);
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

onMounted(() => {
  getTypes();
  start();
});
onUnmounted(() => {
  deleteEventSource();
});
</script>

<template>
  <DefaultLayout>
    <template v-slot:header>
      <ShareButton></ShareButton>
      <input
        v-model="searchInput"
        @keydown.enter="ensureEntryOnNotBoughtList"
        type="search"
        autocomplete="off"
        placeholder="Search"
      />
      <Button @click="ensureEntryOnNotBoughtList" class="inverted">+</Button>
    </template>
    <template v-slot:main>
      <TransitionGroup name="list" tag="ul" @before-leave="beforeLeave">
        <EntryItem
          v-for="(entry, i) in activeSortedNotBoughtEntries"
          :key="entry.ID"
          v-model="activeSortedNotBoughtEntries[i]"
          :types="types"
          @contextmenu="contextmenuShow($event, entry)"
          @update="updateEntry(entry)"
        >
        </EntryItem>
        <hr v-if="activeBoughtEntries.length > 0" />
        <EntryItem
          v-for="(entry, i) in activeBoughtEntries"
          :key="entry.ID"
          v-model="activeBoughtEntries[i]"
          :types="types"
          @contextmenu="contextmenuShow($event, entry)"
          @update="updateEntry(entry)"
        >
        </EntryItem>
      </TransitionGroup>
      <Contextmenu
        v-if="contextmenuVisible"
        id="contextmenu"
        :x="contextmenuX"
        :y="contextmenuY"
        :actions="contextmenuActions"
        @action="contextmenuHandle"
      />
    </template>
  </DefaultLayout>
</template>

<style>
hr {
  color: var(--color-primary);
  margin: 0.5em 1em;
}

.list-move, /* apply transition to moving elements */
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

ul {
  /* no bullet points */
  list-style: none;
  padding: 0em;
  margin: 0em;
}
</style>
