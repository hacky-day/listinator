<script setup lang="ts">
import { onMounted, ref, computed, onUnmounted } from "vue";
import { useRoute } from "vue-router";

import {
  apiGetEntries,
  apiCreateEntry,
  apiGetTypes,
  apiDeleteEntry,
} from "@/api/api.ts";
import { type Entry, type Type, type ContextmenuAction } from "@/types.ts";
import DefaultLayout from "@/Layouts/DefaultLayout.vue";
import Button from "@/Components/Button.vue";
import EntryItem from "@/Components/EntryItem.vue";
import ShareButton from "@/Components/ShareButton.vue";
import Contextmenu from "@/Components/Contextmenu.vue";
import { useNotificationManager } from "@/composables/useNotificationManager";

const { showError } = useNotificationManager();

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
  entry._dirty = true;
  contextmenuX.value = event.clientX;
  contextmenuY.value = event.clientY;
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
    contextmenuTarget.value._dirty = false;
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
      showError("Unable to delete entry", error);
      return;
    }
  }
}

const activeSortedNotBoughtEntries = computed(() => {
  return entries.value
    .filter(
      (entry) =>
        !entry.Bought &&
        entry.Name.toLowerCase().includes(searchInput.value.toLowerCase()),
    )
    .sort((a, b) => {
      // Sort by type
      const typeCompare = a.TypeID.localeCompare(b.TypeID);
      if (typeCompare !== 0) return typeCompare;

      // Sort by name
      return a.Name.localeCompare(b.Name);
    });
});

const activeBoughtEntries = computed(() => {
  return entries.value.filter((entry) => entry.Bought);
});

async function getTypes() {
  try {
    types.value = await apiGetTypes();
  } catch (error) {
    showError("Unable to load entry types", error);
  }
}

async function getEntries() {
  // Get entries from server
  const freshEntries = await apiGetEntries(listID).catch((error) => {
    showError("Unable to load entries", error);
    return [] as Entry[];
  });

  // Just use, if not entries in List
  if (entries.value.length === 0) {
    entries.value = freshEntries;
    return;
  }

  // Merge Entries
  for (const freshEntry of freshEntries) {
    const existingEntry = entries.value.find((e) => e.ID === freshEntry.ID);

    // Add Entries which do not already exist
    if (!existingEntry) {
      entries.value.push(freshEntry);
      continue;
    }

    // Override existing Entries only if there is currently no user interaction on the entry
    if (!existingEntry._dirty) {
      Object.assign(existingEntry, freshEntry);
    }
  }

  // Remove Entries not in the freshEnties Array
  entries.value = entries.value.filter((localEntry) =>
    freshEntries.some((freshEntry) => freshEntry.ID === localEntry.ID),
  );
}

async function createEntry() {
  if (searchInput.value === "") {
    return;
  }
  // Get entries from server
  try {
    const entry = await apiCreateEntry(searchInput.value, listID);
    entries.value.push(entry);
  } catch (error) {
    showError("Unable to create new entry", error);
  }

  // Reset input
  searchInput.value = "";
}

onMounted(() => {
  getTypes();
  getEntries();
  const interval = setInterval(getEntries, 5000);

  onUnmounted(() => {
    clearInterval(interval);
  });
});
</script>

<template>
  <DefaultLayout>
    <template v-slot:header>
      <ShareButton></ShareButton>
      <input
        v-model="searchInput"
        @keydown.enter="createEntry"
        type="search"
        autocomplete="off"
        placeholder="Search"
      />
      <Button @click="createEntry" class="inverted">+</Button>
    </template>
    <template v-slot:main>
      <ul>
        <EntryItem
          v-for="entry in activeSortedNotBoughtEntries"
          :key="entry.ID"
          :entry="entry"
          :types="types"
          @contextmenu="contextmenuShow($event, entry)"
        >
        </EntryItem>
      </ul>
      <hr v-if="activeBoughtEntries.length > 0" />
      <ul>
        <EntryItem
          v-for="entry in activeBoughtEntries"
          :key="entry.ID"
          :entry="entry"
          :types="types"
          @contextmenu="contextmenuShow($event, entry)"
        >
        </EntryItem>
      </ul>
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
