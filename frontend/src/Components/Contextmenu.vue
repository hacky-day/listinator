<script setup lang="ts">
import { type ContextmenuAction } from "@/types.ts";

defineProps<{
  x: number;
  y: number;
  actions: ContextmenuAction[];
}>();

defineEmits<{
  action: [value: string];
}>();
</script>

<template>
  <div class="contextmenu" :style="{ top: `${y}px`, left: `${x}px` }">
    <div
      v-for="action in actions"
      :key="action.action"
      @click="$emit('action', action.action)"
    >
      {{ action.label }}
    </div>
  </div>
</template>

<style scoped>
.contextmenu {
  position: absolute;
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius);
  padding: 0;
  margin: 0;
  z-index: 9999;
  box-shadow: 0.1em 0.1em 0.1em var(--color-border);
  background: var(--color-surface);
}

.contextmenu div {
  padding: 0.5em 0.8em;
}

.contextmenu div:hover {
  background: var(--color-surface-hover);
  cursor: pointer;
}
</style>
