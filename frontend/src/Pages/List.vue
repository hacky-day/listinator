<script setup lang="ts">
import { onMounted, ref, computed } from "vue";
import { useRoute } from "vue-router";

import { apiGetEntries, apiCreateEntry, apiGetTypes, apiUpdateEntry } from "@/api/api.ts";
import { type Entry, type Type } from "@/types.ts";
import DefaultLayout from "@/Layouts/DefaultLayout.vue";
import Button from "@/Components/Button.vue";

const route = useRoute();
const listID = route.params.id as string;
const entries = ref<Entry[]>([]);
const types = ref<Type[]>([]);
const searchInput = ref<string>("");

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
async function onSelectChange(entry: Entry, event: Event) {
  entry.TypeID = (event.target as HTMLSelectElement).value
  try {
    await apiUpdateEntry(entry)
  } catch (error) {
    alert("unable to update entry" + error)
  }
}

async function onClickButton(entry: Entry) {
  entry.Bought = !entry.Bought
  try {
    await apiUpdateEntry(entry)
  } catch (error) {
    alert("unable to update entry" + error)
  }
}

async function onInputFocusOut(entry: Entry, event: Event) {
  entry.Number = (event.target as HTMLInputElement).value
  try {
    await apiUpdateEntry(entry)
  } catch (error) {
    alert("unable to update entry" + error)
  }
}

onMounted(async () => {
  await getTypes();
  await getEntries();
});
</script>

<template>
  <DefaultLayout>
    <template v-slot:header>
      <input v-model="searchInput" @keydown.enter="createEntry" type="search" autocomplete="off" placeholder="Search" />
      <Button @click="createEntry" class="inverted">+</Button>
    </template>
    <template v-slot:main>
      <ul>
        <li v-for="entry in activeSortedNotBoughtEntries" :key="entry.ID">
          <div class="entryAttributes">
            <select @change="onSelectChange(entry, $event)" v-model="entry.TypeID">
              <option v-for="type in types" :key="type.Name" :value="type.Name">
                {{ type.Icon }}
              </option>
            </select>
            <div>{{ entry.Name }}</div>
          </div>
          <div class="entryAttributes">
            <input @focusout="onInputFocusOut(entry, $event)" v-model="entry.Number" type="text" />
            <Button @click="onClickButton(entry)">✓</Button>
          </div>
        </li>
      </ul>
      <hr v-if="activeBoughtEntries.length > 0" />
      <ul>
        <li v-for="entry in activeBoughtEntries" :key="entry.ID">
          <div class="entryAttributes">
            <select @change="onSelectChange(entry, $event)" v-model="entry.TypeID">
              <option v-for="type in types" :key="type.Name" :value="type.Name">
                {{ type.Icon }}
              </option>
            </select>
            <div>{{ entry.Name }}</div>
          </div>
          <div class="entryAttributes">
            <input @focusout="onInputFocusOut(entry, $event)" v-model="entry.Number" type="text" />
            <Button @click="onClickButton(entry)">+</Button>
          </div>
        </li>
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
  color: var(--color-text);
}

ul {
  /* no bullet points */
  list-style: none;
  padding: 0em;
  margin: 0em;
}

#share-button {
  background: transparent;
  min-width: 0em;
  margin: 0em;
  padding: 0em;
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

select {
  width: 3em;
  height: 2.5em;
  background: var(--color-surface);
  border: 0.1em solid var(--color-border);
  border-radius: var(--border-radius);
}
</style>
