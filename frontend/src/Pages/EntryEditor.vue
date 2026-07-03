<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";

import { apiGetEntry, apiGetTypes, apiUpdateEntry } from "@/api/api";
import type { Entry, Type } from "@/types";
import { router } from "@/router.ts";
import { useNotificationManager } from "@/composables/useNotificationManager";

import DefaultLayout from "@/Layouts/DefaultLayout.vue";
import Button from "@/Components/Button.vue";

const { show } = useNotificationManager();
const route = useRoute();
const entryID = route.params.id as string;

const entry = ref<Entry>();
const types = ref<Type[]>();

async function save() {
  if (entry.value === undefined) {
    return;
  }
  try {
    await apiUpdateEntry(entry.value);
  } catch (error) {
    show("error", "Unable to update entry", { logMessage: error });
  }
}

async function selectType(typeID: string) {
  if (entry.value === undefined) {
    return;
  }
  entry.value.TypeID = typeID;
  await save();
}

onMounted(async () => {
  try {
    entry.value = await apiGetEntry(entryID);
    types.value = await apiGetTypes();
  } catch (error) {
    show("error", "unable to connect to server", { logMessage: error });
  }
});
</script>
<template>
  <DefaultLayout>
    <template v-slot:header>
      <h1>Edit Entry</h1>
      <Button class="inverted" @click="router.back()">Done</Button>
    </template>
    <template v-slot:main>
      <template v-if="entry">
        <label>
          Name
          <input
            v-model="entry.Name"
            type="text"
            @blur="save"
            @keydown.enter="save"
          />
        </label>
        <label>
          Number
          <input
            v-model="entry.Number"
            type="text"
            @blur="save"
            @keydown.enter="save"
          />
        </label>
        <div class="typeLabel">Type</div>
        <ul>
          <li
            v-for="type in types"
            :key="type.ID"
            :class="{ active: type.ID === entry.TypeID }"
            @click="selectType(type.ID)"
            :style="{ borderLeft: `0.3em solid ${type.Color}` }"
          >
            <div>
              {{ type.Name }}
            </div>
          </li>
        </ul>
      </template>
    </template>
  </DefaultLayout>
</template>

<style scoped>
h1 {
  margin: 0em;
  font-size: 1.1em;
  font-weight: 500;
}

label {
  display: flex;
  flex-direction: column;
  gap: 0.3em;
  margin-bottom: 1em;
  font-weight: 300;
}

label input {
  border: 0.1em solid var(--color-accent);
}

.typeLabel {
  margin-bottom: 0.5em;
  font-weight: 300;
}

ul {
  list-style: none;
  padding: 0em;
  margin: 0em;
}

li {
  min-height: 2.5em;
  padding: 0.2em;
  margin: 0.5em 0em;
  border-radius: var(--border-radius);
  display: flex;
  cursor: pointer;
  align-items: center;
  background: var(--color-surface);
}

li.active {
  outline: 0.15em solid var(--color-primary);
}

@media (hover: hover) and (pointer: fine) {
  li:hover {
    background: var(--color-surface-hover);
  }
}
</style>
