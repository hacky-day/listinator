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
      alert("unable to delete entry" + error);
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
    alert("unable to get types:" + error);
  }
}

async function getEntries() {
  // Get entries from server
  try {
    entries.value = await apiGetEntries(listID);
  } catch (error) {
    alert("unable to get entries: " + error);
  }
}

async function createEntry() {
  if (searchInput.value === "") {
    return;
  }
  // Get entries from server
  try {
    await apiCreateEntry(searchInput.value, listID);
  } catch (error) {
    alert("unable to get entries:" + error);
  }

  // Reset input
  searchInput.value = "";
}

async function updateEntry(entry: Entry) {
  try {
    await apiUpdateEntry(entry);
  } catch (error) {
    alert("unable to update entry: " + error);
    return;
  }
}

onMounted(() => {
  getTypes();
  getEntries();

  const evtSource = new EventSource(`/api/v1/entries/events?ListID=${listID}`);
  evtSource.onerror = (event) => {
    console.log(event);
  };
  evtSource.addEventListener("create", (event: MessageEvent) => {
    const entry = JSON.parse(event.data) as Entry;
    entries.value.push(entry);
  });
  evtSource.addEventListener("update", async (event: MessageEvent) => {
    const entry = JSON.parse(event.data) as Entry;
    const index = entries.value.findIndex((item) => item.ID === entry.ID);
    if (index === -1) {
      return;
    }
    entries.value[index] = entry;
  });
  onUnmounted(() => {
    evtSource.close();
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
          v-for="(entry, i) in activeSortedNotBoughtEntries"
          :key="entry.ID"
          v-model="activeSortedNotBoughtEntries[i]"
          :types="types"
          @contextmenu="contextmenuShow($event, entry)"
          @update="updateEntry(entry)"
        >
        </EntryItem>
      </ul>
      <hr v-if="activeBoughtEntries.length > 0" />
      <ul>
        <EntryItem
          v-for="(entry, i) in activeBoughtEntries"
          :key="entry.ID"
          v-model="activeBoughtEntries[i]"
          :types="types"
          @contextmenu="contextmenuShow($event, entry)"
          @update="updateEntry(entry)"
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
