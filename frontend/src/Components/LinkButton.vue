<script setup lang="ts">
const props = defineProps({
  to: { type: String, required: true },
});
</script>

<template>
  <router-link class="link-button" :to="to">
    <slot></slot>
  </router-link>
</template>

<style>
.link-button {
  text-align: center;
  text-decoration: None;
  line-height: 2.5em;
  height: 2.5em;
  min-width: 2.5em;
  border-radius: var(--border-radius);
  border: 0em;
  background: var(--color-accent);
  color: var(--color-text-inverted);
}
</style>
