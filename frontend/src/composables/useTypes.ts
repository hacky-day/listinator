import { ref } from "vue";
import { apiGetTypes } from "@/api/api";
import { getTypesLocal, replaceTypes } from "@/offline/db";
import type { Type } from "@/types";

// Module-scope so every page sees the same cached type list instead of each
// fetching it independently.
const types = ref<Type[]>([]);
let hydrated = false;

async function hydrateFromCache() {
  if (hydrated) {
    return;
  }
  hydrated = true;
  const local = await getTypesLocal();
  if (local.length > 0) {
    types.value = local;
  }
}

// refresh hydrates from IndexedDB first (so offline callers already have
// something to render), then tries to fetch fresh data from the server. It
// throws only if there is neither a cache nor a working network, so a caller
// can distinguish "nothing to show at all" from "showing stale cached data".
async function refresh(): Promise<void> {
  await hydrateFromCache();
  const fresh = await apiGetTypes();
  types.value = fresh;
  await replaceTypes(fresh);
}

export function useTypes() {
  return { types, refresh };
}
