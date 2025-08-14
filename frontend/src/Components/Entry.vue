<script setup lang="ts">
import { type Entry } from "@/types.ts";

defineProps<{
  entry: Entry;
}>();
</script>

<template>
  <li :id="entry.ID">
    <div>{{ entry.Name }}</div>
  </li>
</template>
