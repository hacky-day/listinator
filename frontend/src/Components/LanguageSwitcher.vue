<script setup lang="ts">
import { computed } from "vue";

import { availableLocales, locale, setLocale, t, type Locale } from "@/i18n";

const selected = computed({
  get: () => locale.value,
  set: (value: Locale) => setLocale(value),
});
</script>

<template>
  <select v-model="selected" :aria-label="t('common.language')" :title="t('common.language')">
    <option v-for="entry in availableLocales" :key="entry.locale" :value="entry.locale" :title="entry.label">
      {{ entry.short }}
    </option>
  </select>
</template>

<style scoped>
/* Sits in the header on every page, so it stays as small and quiet as
   possible: only the language code, no background of its own. */
select {
  height: 2em;
  padding: 0em 0.2em;
  border: 0em;
  border-radius: var(--border-radius);
  background: transparent;
  color: inherit;
  font-weight: 300;
}

/* The dropdown itself is drawn on the page background, not on the header, so
   it needs the regular text colors instead of the inherited inverted ones. */
option {
  background: var(--color-surface);
  color: var(--color-text);
}

@media (hover: hover) and (pointer: fine) {
  select:hover {
    cursor: pointer;
  }
}
</style>
