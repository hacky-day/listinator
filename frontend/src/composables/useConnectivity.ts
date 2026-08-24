import { ref } from "vue";

// Module-scope so every composable/component observes the same connectivity
// state instead of each list view tracking its own.
const isOnline = ref(navigator.onLine);

window.addEventListener("online", () => {
  isOnline.value = true;
});
window.addEventListener("offline", () => {
  isOnline.value = false;
});

export function useConnectivity() {
  return { isOnline };
}
