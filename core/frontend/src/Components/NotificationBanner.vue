<script setup lang="ts">
import { useNotificationManager } from "@/composables/useNotificationManager";

const { notifications, clear } = useNotificationManager();
</script>

<template>
  <div class="notification">
    <transition-group name="list">
      <div
        v-for="notification in notifications"
        class="entry"
        :class="notification.level"
        :key="notification.id"
      >
        <div class="message">
          {{ notification.userMessage }}
        </div>
        <button @click="clear(notification.id)">×</button>
      </div>
    </transition-group>
  </div>
</template>

<style scoped>
.notification {
  position: absolute;
  width: 100%;
}

.list-enter-active,
.list-leave-active {
  transition: all 0.5s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(2em);
}

.entry {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.5em 0.5em;
  opacity: 0.9;
  border-radius: var(--border-radius);
  margin: 0.25em 0.55em;
  gap: 1em;
}

.entry.error {
  background: red;
}

.entry.warning {
  background: orange;
}

.entry.info {
  background: green;
}

.message {
  flex: 1;
  font-weight: 500;
  text-align: center;
}

button {
  background: none;
  border: none;
  color: var(--color-text-inverted);
  font-size: 1.5em;
  font-weight: bold;
  cursor: pointer;
}
</style>
