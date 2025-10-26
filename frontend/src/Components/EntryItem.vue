<script setup lang="ts">
import { type Entry } from "@/types.ts";

import Button from "@/Components/Button.vue";

const emit = defineEmits<{
  (e: "update"): void;
}>();

const entry = defineModel<Entry>({ required: true });

function onClick() {
  entry.value.Bought = !entry.value.Bought;
  emit("update");
}
</script>

<template>
  <div class="entry">
    <div class="name">
      {{ entry.Name }}
    </div>
    <div class="entryAttributes">
      <input v-model="entry.Number" type="text" @blur="$emit('update')" />
      <Button @click="onClick">{{ entry.Bought ? "+" : "✓" }}</Button>
    </div>
  </div>
</template>

<style scoped>
.entry {
  padding: 0.1em 0em;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--color-surface);
}

@media (hover: hover) and (pointer: fine) {
  .entry:hover {
    background: var(--color-surface-hover);
  }
}

.name {
  padding-left: 0.5em;
}

.entryAttributes {
  display: flex;
  align-items: center;
  gap: 0.5em;
}

input {
  width: 6em;
  border: 0.1em solid var(--color-accent);
  text-align: right;
}
</style>
