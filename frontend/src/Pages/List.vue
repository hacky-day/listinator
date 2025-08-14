<script setup lang="ts">
import { onMounted, ref, computed } from "vue";
import { useRoute } from "vue-router";

import { router } from "@/router.ts";
import { apiGetEntries, apiCreateEntry } from "@/api/api.ts";
import { type Entry } from "@/type.ts";
import DefaultLayout from "@/Layouts/DefaultLayout.vue";
import Button from "@/Components/Button.vue";
import BuyEntry from "@/Components/BuyEntry.vue";

const route = useRoute();
const listID = route.params.id;

const entries = ref<Entry[]>([]);
const searchInput = ref<string>("");

const activeSortedNotBoughtEntries = computed(() => {
  return entries.value
    .filter(
      (entry) =>
        !entry.bought &&
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
  return entries.value.filter((entry) => entry.bought);
});

async function getEntries() {
  // Get entries from server
  try {
    entries.value = await apiGetEntries(listID);
  } catch (error) {
    alert("unable to get entries:" + error);
  }
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
    alert("unable to get entries:" + error);
  }

  // Reset input
  searchInput.value = "";
}

onMounted(async () => {
  await getEntries();
});
</script>

<template>
  <DefaultLayout>
    <template v-slot:header>
      <button id="share-button">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="20"
          height="20"
          fill="white"
          viewBox="0 0 640 640"
        >
          <!--!Font Awesome Free v7.0.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2025 Fonticons, Inc.-->
          <path
            d="M448 256C501 256 544 213 544 160C544 107 501 64 448 64C395 64 352 107 352 160C352 165.4 352.5 170.8 353.3 176L223.6 248.1C206.7 233.1 184.4 224 160 224C107 224 64 267 64 320C64 373 107 416 160 416C184.4 416 206.6 406.9 223.6 391.9L353.3 464C352.4 469.2 352 474.5 352 480C352 533 395 576 448 576C501 576 544 533 544 480C544 427 501 384 448 384C423.6 384 401.4 393.1 384.4 408.1L254.7 336C255.6 330.8 256 325.5 256 320C256 314.5 255.5 309.2 254.7 304L384.4 231.9C401.3 246.9 423.6 256 448 256z"
          />
        </svg>
      </button>
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
        <BuyEntry
          v-for="entry in activeSortedNotBoughtEntries"
          :entry="entry"
          :key="entry.id"
        ></BuyEntry>
      </ul>
      <hr v-if="activeBoughtEntries.len > 0" />
      <ul>
        <BuyEntry
          v-for="entry in activeBoughtEntries"
          :entry="entry"
          :key="entry.id"
        ></BuyEntry>
      </ul>
    </template>
  </DefaultLayout>
</template>

<style>
hr {
  color: var(--color-primary);
  margin: 0.5em 1em;
}

header input {
  flex-grow: 1; /* max width possible */
}

input {
  height: 2.5em;
  box-sizing: border-box;
  padding: 0.1em 0.2em;
  border-radius: var(--border-radius);
  color: var(--color-text);
  border: 0.1em solid var(--color-border);
  background-color: var(--color-surface);
  color: var(--color-text);
}

ul {
  /* no bullet points */
  list-style: none;
  padding: 0em;
  margin: 0em;
}

li {
  padding: 0.1em 0em;
  margin: 0.2em 0em;
  border-radius: var(--border-radius);
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--color-surface);
}

@media (hover: hover) and (pointer: fine) {
  li:hover {
    background: var(--color-surface-hover);
  }
}

.entryAttributes {
  display: flex;
  align-items: center;
  gap: 0.5em;
}

select {
  width: 3em;
  height: 2.5em;
  background: var(--color-surface);
  border: 0.1em solid var(--color-border);
  border-radius: var(--border-radius);
}

li input {
  width: 6em;
  border: 0.1em solid var(--color-accent);
  text-align: right;
}

.context-menu {
  position: absolute;
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius);
  padding: 0;
  margin: 0;
  z-index: 9999;
  box-shadow: 0.1em 0.1em 0.1em var(--color-border);
  background: var(--color-surface);
}

.context-menu div {
  padding: 0.5em 0.8em;
}

.context-menu div:hover {
  background: var(--color-surface-hover);
  cursor: pointer;
}

#share-button {
  background: transparent;
  min-width: 0em;
  margin: 0em;
  padding: 0em;
}
</style>
