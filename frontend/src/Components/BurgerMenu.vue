<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useNotificationManager } from "@/composables/useNotificationManager";
import { apiGetSession, apiDeleteSession } from "@/api/api.ts";

const { show } = useNotificationManager();

const menuVisible = ref(false);
const loggedIn = ref(false);
const loaded = ref(false);

async function checkSession() {
  try {
    await apiGetSession();
    loggedIn.value = true;
  } catch (error) {
    loggedIn.value = false;
  }
}

async function logout() {
  try {
    await apiDeleteSession();
    loggedIn.value = false;
    show("info", "Logged out successfully");
  } catch (error) {
    show("error", "Unable to logout", { logMessage: error });
  }
}

/**
 * share shares the current URL. Either via share or via copy to clipboard.
 * Support for different Browser and mobile are hard to test.
 */
async function share() {
  try {
    if (navigator.share) {
      await navigator.share({
        url: window.location.href,
      });
      return;
    }
    await navigator.clipboard.writeText(window.location.href);
    show("info", "URL copied to clipboard successfully");
  } catch (error) {
    show("error", "Unable to share link", { logMessage: error });
  }
}

function toggleMenu() {
  menuVisible.value = !menuVisible.value;
  if (menuVisible.value) {
    document.addEventListener("click", handleOutsideClick);
  } else {
    document.removeEventListener("click", handleOutsideClick);
  }
}

function handleOutsideClick(event: Event) {
  const menu = document.getElementById("burger-menu");
  const button = document.getElementById("burger-button");
  
  if (
    menu &&
    button &&
    !menu.contains(event.target as HTMLElement) &&
    !button.contains(event.target as HTMLElement)
  ) {
    closeMenu();
  }
}

function closeMenu() {
  menuVisible.value = false;
  document.removeEventListener("click", handleOutsideClick);
}

function handleAction(action: string) {
  closeMenu();
  
  switch (action) {
    case "share":
      share();
      break;
    case "login":
      // Navigate to login page
      window.location.href = "/login";
      break;
    case "logout":
      logout();
      break;
  }
}

onMounted(async () => {
  await checkSession();
  loaded.value = true;
});
</script>

<template>
  <div class="burger-menu-container">
    <button
      id="burger-button"
      class="burger-button"
      @click="toggleMenu"
      aria-label="Menu"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="20"
        height="20"
        fill="white"
        viewBox="0 0 24 24"
      >
        <path d="M3 6h18v2H3V6zm0 5h18v2H3v-2zm0 5h18v2H3v-2z" />
      </svg>
    </button>
    
    <div
      v-if="menuVisible && loaded"
      id="burger-menu"
      class="burger-menu"
    >
      <div class="menu-item" @click="handleAction('share')">
        Share
      </div>
      <div v-if="!loggedIn" class="menu-item" @click="handleAction('login')">
        Login
      </div>
      <div v-if="loggedIn" class="menu-item" @click="handleAction('logout')">
        Logout
      </div>
    </div>
  </div>
</template>

<style scoped>
.burger-menu-container {
  position: relative;
}

.burger-button {
  border: 0;
  background: transparent;
  cursor: pointer;
  padding: 0.25em;
  display: flex;
  align-items: center;
  justify-content: center;
}

.burger-button:hover {
  background: rgba(255, 255, 255, 0.1);
  border-radius: var(--border-radius);
}

.burger-menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 0.25em;
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius);
  padding: 0;
  z-index: 9999;
  box-shadow: 0.1em 0.1em 0.1em var(--color-border);
  background: var(--color-surface);
  min-width: 8em;
}

.menu-item {
  padding: 0.5em 0.8em;
  cursor: pointer;
  color: var(--color-text);
}

.menu-item:hover {
  background: var(--color-surface-hover);
}

.menu-item:first-child {
  border-top-left-radius: var(--border-radius);
  border-top-right-radius: var(--border-radius);
}

.menu-item:last-child {
  border-bottom-left-radius: var(--border-radius);
  border-bottom-right-radius: var(--border-radius);
}
</style>