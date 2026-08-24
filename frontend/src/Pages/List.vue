<script setup lang="ts">
import { onMounted, ref, computed, onUnmounted } from "vue";
import { useRoute } from "vue-router";

import { router } from "@/router.ts";
import { useNotificationManager } from "@/composables/useNotificationManager";
import { useEntries } from "@/composables/useEntries";
import { useTypes } from "@/composables/useTypes";
import { useConnectivity } from "@/composables/useConnectivity";
import { locale, t, translateType } from "@/i18n";
import { type Entry, type Type, type ContextmenuAction } from "@/types.ts";

import DefaultLayout from "@/Layouts/DefaultLayout.vue";
import Button from "@/Components/Button.vue";
import EntryItem from "@/Components/EntryItem.vue";
import ShareButton from "@/Components/ShareButton.vue";
import Contextmenu from "@/Components/Contextmenu.vue";

const { show } = useNotificationManager();
const { isOnline } = useConnectivity();

const route = useRoute();
const listID = route.params.id as string;
const {
  entries,
  start,
  stop,
  listen,
  unlisten,
  createEntry,
  updateEntry: syncEntry,
  deleteEntry,
} = useEntries(listID);
const { types, refresh: refreshTypes } = useTypes();
const searchInput = ref<string>("");

const contextmenuVisible = ref(false);
const contextmenuX = ref(0);
const contextmenuY = ref(0);
// Computed, so the labels follow a language switch instead of being frozen at
// the language active during the first render.
const contextmenuActions = computed<ContextmenuAction[]>(() => [
  { label: t("common.edit"), action: "edit" },
  { label: t("common.delete"), action: "delete" },
]);
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

async function contextmenuHandle(action: string) {
  const targetEntry = contextmenuTarget.value as Entry;
  contextmenuClose();

  if (action === "delete") {
    try {
      await deleteEntry(targetEntry);
    } catch (error) {
      // A failed write while already known offline is expected and gets
      // queued once reconnected, not worth interrupting the user for.
      if (isOnline.value) {
        show("error", t("list.errorDeleteEntry"), { logMessage: error });
      }
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
        // Sort by name, using the collation rules of the active language
        return a.Name.localeCompare(b.Name, locale.value);
      });
  });
  return groups;
});

const activeBoughtEntries = computed(() => {
  return activeEntries.value.filter((entry) => entry.Bought);
});

async function getTypes() {
  try {
    await refreshTypes();
  } catch (error) {
    if (types.value.length === 0) {
      show("error", t("list.errorLoadTypes"), { logMessage: error });
    }
  }
}

async function ensureEntryOnNotBoughtList() {
  if (searchInput.value === "") {
    return;
  }

  // Create new Entry
  try {
    await createEntry(searchInput.value);
  } catch (error) {
    if (isOnline.value) {
      show("error", t("list.errorCreateEntry"), { logMessage: error });
    }
  } finally {
    // Reset input
    searchInput.value = "";
  }
}

async function updateEntry(entry: Entry) {
  try {
    await syncEntry(entry);
  } catch (error) {
    if (isOnline.value) {
      show("error", t("list.errorUpdateEntry"), { logMessage: error });
    }
    return;
  } finally {
    // Reset input
    searchInput.value = "";
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

onMounted(() => {
  getTypes();
  start();
  listen();
});
onUnmounted(() => {
  stop();
  unlisten();
});
</script>

<template>
  <DefaultLayout>
    <template v-slot:header>
      <ShareButton></ShareButton>
      <input v-model="searchInput" @keydown.enter="ensureEntryOnNotBoughtList" type="search" autocomplete="off"
        :placeholder="t('list.search')" />
      <Button @click="ensureEntryOnNotBoughtList" class="inverted">+</Button>
    </template>
    <template v-slot:main>
      <TransitionGroup name="list" @before-leave="beforeLeave" tag="ul">
        <template v-for="type in types">
          <div :key="type.ID" v-if="activeSortedNotBoughtEntriesbyType[type.ID].length > 0" class="divider">
            {{ translateType(type) }}
          </div>
          <li :style="{ borderLeft: `0.3em solid ${type.Color}` }"
            v-for="(entry, i) in activeSortedNotBoughtEntriesbyType[type.ID]" :key="entry.ID">
            <EntryItem v-model="activeSortedNotBoughtEntriesbyType[type.ID][i]"
              @contextmenu="contextmenuShow($event, entry)" @update="updateEntry(entry)">
            </EntryItem>
          </li>
        </template>

        <li key="bought" v-if="activeBoughtEntries.length > 0" class="divider bought">
          {{ t("list.recentlyBought") }}
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

  /* Flex items do not shrink below their intrinsic width by default, which
     makes the header overflow on narrow screens. */
  min-width: 0;
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
