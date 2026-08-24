<script setup lang="ts">
import { useConnectivity } from "@/composables/useConnectivity";
import { t } from "@/i18n";

const { isOnline } = useConnectivity();
</script>

<template>
  <span v-if="!isOnline" class="offline">{{ t("common.offline") }}</span>
</template>

<style scoped>
/* Sits quietly in the header, like LanguageSwitcher: no alarming color, just
   a muted label so being offline reads as a normal state, not an error. */
.offline {
  padding: 0.2em 0.6em;
  border-radius: var(--border-radius);
  background: rgba(255, 255, 255, 0.15);
  font-weight: 300;
  font-size: 0.85em;
  white-space: nowrap;
}
</style>
