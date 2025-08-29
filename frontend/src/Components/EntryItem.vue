<script setup lang="ts">
import { type Entry, type Type } from "@/types.ts";

import Button from "@/Components/Button.vue";

defineProps<{
  types: Type[];
}>();

const emit = defineEmits<{
  (e: "update"): void;
  (e: "changing", value: boolean): void;
}>();

const entry = defineModel<Entry>({ required: true });

function onClick() {
  entry.value.Bought = !entry.value.Bought;
  emit("update");
}
</script>

<template>
  <li>
    <div class="entryAttributes">
      <select
        v-model="entry.TypeID"
        @change="$emit('update')"
        @focus="$emit('changing', true)"
        @blur="$emit('changing', false)"
      >
        <option v-for="type in types" :key="type.Name" :value="type.Name">
          {{ type.Icon }}
        </option>
      </select>
      <div>{{ entry.Name }}</div>
    </div>
    <div class="entryAttributes">
      <input
        v-model="entry.Number"
        type="text"
        @focus="$emit('changing', true)"
        @blur="
          $emit('update');
          $emit('changing', false);
        "
      />
      <Button @click="onClick">{{ entry.Bought ? "+" : "✓" }}</Button>
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
