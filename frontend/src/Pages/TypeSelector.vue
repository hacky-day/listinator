<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";

import { apiGetEntry, apiGetTypes, apiUpdateEntry } from "@/api/api";
import type { Entry, Type } from "@/types";
import { router } from "@/router.ts";
import { useNotificationManager } from "@/composables/useNotificationManager";

import DefaultLayout from "@/Layouts/DefaultLayout.vue";

const { show } = useNotificationManager();
const route = useRoute();
const entryID = route.params.id as string;

const entry = ref<Entry>();
const types = ref<Type[]>();

async function update(typeID: string) {
  if (entry.value === undefined) {
    return;
  }
  entry.value.TypeID = typeID;
  try {
    await apiUpdateEntry(entry.value);
  } catch (error) {
    show("error", "Unable to update entry", { logMessage: error });
    return;
  }
  router.back();
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
      <h1>Listinator - Select Type</h1>
    </template>
    <template v-slot:main>
      <ul>
        <li
          v-for="type in types"
          @click="update(type.Name)"
          :style="{ borderLeft: `0.3em solid ${type.Color}` }"
        >
          <div>
            {{ type.Icon }}
            {{ type.Name }}
          </div>
        </li>
      </ul>
    </template>
  </DefaultLayout>
</template>

<style scoped>
ul {
  /* no bullet points */
  list-style: none;
  padding: 0em;
  margin: 0em;
}

div {
  min-height: 2.5em;
  padding: 0.2em;
  margin: 0.5em 0em;
  border-radius: var(--border-radius);
  display: flex;
  cursor: pointer;
  align-items: center;
  background: var(--color-surface);
}

@media (hover: hover) and (pointer: fine) {
  div:hover {
    background: var(--color-surface-hover);
  }
}
</style>
