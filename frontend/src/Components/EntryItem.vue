<script setup lang="ts">
import { type Entry, type Type } from "@/types.ts";
import { apiUpdateEntry } from "@/api/api";
import Button from "@/Components/Button.vue";

const props = defineProps<{
  entry: Entry;
  types: Type[];
}>();

async function becomesDirty() {
  props.entry._dirty = true;
}

async function update() {
  await apiUpdateEntry(props.entry).catch((error) =>
    alert("unable to update entry" + error),
  );
  props.entry._dirty = false;
}

async function onClickButton() {
  props.entry.Bought = !props.entry.Bought;
  update();
}
</script>

<template>
  <li>
    <div class="entryAttributes">
      <select @focus="becomesDirty" @change="update" v-model="entry.TypeID">
        <option v-for="type in types" :key="type.Name" :value="type.Name">
          {{ type.Icon }}
        </option>
      </select>
      <div>{{ entry.Name }}</div>
    </div>
    <div class="entryAttributes">
      <input
        @focus="becomesDirty"
        @blur="update"
        v-model="entry.Number"
        type="text"
      />
      <Button @click="onClickButton">{{ entry.Bought ? "+" : "✓" }}</Button>
    </div>
  </li>
</template>

<style scoped>
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

select {
  width: 3em;
  height: 2.5em;
  background: var(--color-surface);
  border: 0.1em solid var(--color-border);
  border-radius: var(--border-radius);
}
</style>
