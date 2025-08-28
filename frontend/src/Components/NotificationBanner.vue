<script setup lang="ts">
import { useNotificationManager } from "@/composables/useNotificationManager";
import { computed } from "vue";

const { currentNotification, clearError } = useNotificationManager();

// Computed property to check if notification should be visible
const isVisible = computed(() => {
  return currentNotification.value !== null;
});

// Computed property for banner class based on type
const bannerClass = computed(() => {
  if (!currentNotification.value) return "";
  return currentNotification.value.type === "success" ? "success" : "error";
});

function dismissError() {
  clearError();
}
</script>

<template>
  <div v-if="isVisible" class="notification-banner" :class="bannerClass">
    <div class="notification-content">
      <div class="notification-message">
        {{ currentNotification?.userMessage }}
      </div>
      <button
        @click="dismissError"
        class="notification-dismiss"
        aria-label="Dismiss notification"
      >
        ×
      </button>
    </div>
  </div>
</template>

<style scoped>
.notification-banner {
  width: 100%;
  color: var(--color-text-inverted);
  padding: 0.75em 0; /* Remove horizontal padding to prevent overflow */
  animation: slideDown 0.3s ease-out;
  position: relative; /* Add position context for absolute positioning */
}

.notification-banner.error {
  background: var(--color-error);
}

.notification-banner.success {
  background: var(--color-accent);
}

.notification-content {
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  padding: 0 1em; /* Move padding here to maintain spacing inside content */
  gap: 1em;
  width: 100%;
  box-sizing: border-box; /* Ensure padding is included in width calculation */
}

.notification-message {
  flex: 1;
  font-weight: 500;
  font-size: 0.95rem;
  text-align: center;
}

.notification-dismiss {
  background: none;
  border: none;
  color: var(--color-text-inverted);
  font-size: 1.5em;
  font-weight: bold;
  cursor: pointer;
  padding: 0;
  line-height: 1;
  width: 1em;
  height: 1em;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background-color 0.2s ease;
  position: absolute;
  right: 1em; /* Position relative to the banner, with consistent padding */
  top: 50%;
  transform: translateY(-50%); /* Center vertically */
}

@keyframes slideDown {
  from {
    transform: translateY(-100%);
  }
  to {
    transform: translateY(0);
  }
}
</style>
