<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";

import { apiGetEntry } from "@/api/api";
import { getEntryLocal, putEntry } from "@/offline/db";
import type { Entry } from "@/types";
import { router } from "@/router.ts";
import { useNotificationManager } from "@/composables/useNotificationManager";
import { useEntries } from "@/composables/useEntries";
import { useTypes } from "@/composables/useTypes";
import { useConnectivity } from "@/composables/useConnectivity";
import { t, translateType } from "@/i18n";

import DefaultLayout from "@/Layouts/DefaultLayout.vue";
import Button from "@/Components/Button.vue";

const { show } = useNotificationManager();
const { types, refresh: refreshTypes } = useTypes();
const { isOnline } = useConnectivity();
const route = useRoute();
const entryID = route.params.id as string;

const entry = ref<Entry>();
// syncEntry is bound once the entry (and thus its ListID) has loaded, so we
// share the same reactive cache/outbox that List.vue uses for this list.
let syncEntry: ((entry: Entry) => Promise<void>) | undefined;

async function save() {
  if (entry.value === undefined || syncEntry === undefined) {
    return;
  }
  try {
    await syncEntry(entry.value);
  } catch (error) {
    // A failed write while already known offline is expected and gets
    // queued once reconnected, not worth interrupting the user for.
    if (isOnline.value) {
      show("error", t("entryEditor.errorUpdateEntry"), { logMessage: error });
    }
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
  // Show whatever is cached locally first, so an offline visit still has
  // something to render instead of a blank page.
  const cached = await getEntryLocal(entryID);
  if (cached !== undefined) {
    entry.value = cached;
    syncEntry = useEntries(cached.ListID).updateEntry;
  }

  try {
    const fresh = await apiGetEntry(entryID);
    entry.value = fresh;
    syncEntry = useEntries(fresh.ListID).updateEntry;
    await putEntry(fresh);
    await refreshTypes();
  } catch (error) {
    if (entry.value === undefined) {
      show("error", t("entryEditor.errorLoad"), { logMessage: error });
    }
  }
});
</script>
<template>
  <DefaultLayout>
    <template v-slot:header>
      <h1>{{ t("entryEditor.title") }}</h1>
      <Button class="inverted" @click="router.back()">{{
        t("common.done")
      }}</Button>
    </template>
    <template v-slot:main>
      <template v-if="entry">
        <label>
          {{ t("common.name") }}
          <input
            v-model="entry.Name"
            type="text"
            @blur="save"
            @keydown.enter="save"
          />
        </label>
        <label>
          {{ t("entryEditor.number") }}
          <input
            v-model="entry.Number"
            type="text"
            @blur="save"
            @keydown.enter="save"
          />
        </label>
        <div class="typeLabel">{{ t("entryEditor.type") }}</div>
        <ul>
          <li
            v-for="type in types"
            :key="type.ID"
            :class="{ active: type.ID === entry.TypeID }"
            @click="selectType(type.ID)"
            :style="{ borderLeft: `0.3em solid ${type.Color}` }"
          >
            <div>
              {{ translateType(type) }}
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
