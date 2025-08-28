<script setup lang="ts">
import NotificationBanner from "@/Components/NotificationBanner.vue";
</script>

<template>
  <NotificationBanner />
  <header>
    <slot name="header"></slot>
  </header>
  <main>
    <slot name="main"></slot>
  </main>
</template>

<style>
:root {
  /* Primary Colors */
  --color-primary: #4f46e5;
  --color-primary-light: #6366f1;
  --color-primary-dark: #4338ca;

  /* Accent / Secondary */
  --color-accent: #10b981;
  --color-accent-light: #34d399;
  --color-accent-dark: #059669;

  --color-error: #dc2626;

  /* Neutral / Base */
  --color-background: #f9fafb;
  --color-surface: #ffffff;
  --color-surface-hover: #f3f4f6;
  --color-border: #e5e7eb;
  --color-muted: #9ca3af;

  /* Text */
  --color-text: #111827;
  --color-text-muted: #6b7280;
  --color-text-inverted: #ffffff;

  --border-radius: 0.75em;
}

main {
  padding: 0.5em;
}

.invisible {
  display: none;
}

* {
  font-size: 16px;
  font-family: sans;
}

body {
  padding: 0em;
  margin: 0em;
  background: var(--color-background);
  color: var(--color-text);
}

header {
  height: 2.5em;
  padding: 0.5em;
  background: var(--color-primary);
  color: var(--color-text-inverted);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 0.5em;
}
</style>
