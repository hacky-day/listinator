import { deleteEntryLocal, putEntry } from "@/offline/db";
import { enqueueCreate, enqueueDelete, enqueueUpdate, scheduleFlush } from "@/offline/outbox";
import type { Entry } from "@/types";

// Mirrors database.Entry.BeforeCreate's default so an offline-created entry
// renders under the right group immediately, before it ever reaches the
// server and gets this same default applied there.
const MISC_TYPE_ID = "c29ebd85-812e-4cf6-bfc7-c8368eb83334";

// createEntry, updateEntry and deleteEntry write to the local cache and
// queue the change for sync immediately, without waiting on the network —
// this is what makes writes work offline. useEntries.ts applies the same
// change to its reactive in-memory array around these calls.
export async function createEntry(name: string, listID: string): Promise<Entry> {
  const entry: Entry = {
    ID: crypto.randomUUID(),
    Name: name,
    Number: "",
    Bought: false,
    TypeID: MISC_TYPE_ID,
    ListID: listID,
    ClientUpdatedAt: new Date().toISOString(),
  };
  await putEntry(entry);
  await enqueueCreate(entry);
  scheduleFlush();
  return entry;
}

export async function updateEntry(entry: Entry): Promise<Entry> {
  // entry is typically the live object from useEntries's reactive array
  // (mutated in place by EntryItem.vue's v-model before this is called), a
  // Vue reactive Proxy — IndexedDB's structured clone cannot serialize a
  // Proxy (throws DataCloneError), so store/queue a plain copy instead.
  // ClientUpdatedAt is stamped here, at the moment of the optimistic local
  // write, not at flush time — a delayed flush must still represent when
  // the user actually made the edit, for edit-time last-write-wins to work.
  const plain: Entry = { ...entry, ClientUpdatedAt: new Date().toISOString() };
  await putEntry(plain);
  await enqueueUpdate(plain);
  scheduleFlush();
  return plain;
}

export async function deleteEntry(entry: Entry): Promise<void> {
  // ClientUpdatedAt is stamped at the moment of the optimistic local delete,
  // same reasoning as updateEntry: it must represent when the user acted,
  // not when the queued request eventually flushes.
  const plain: Entry = { ...entry, ClientUpdatedAt: new Date().toISOString() };
  await deleteEntryLocal(entry.ID);
  await enqueueDelete(plain);
  scheduleFlush();
}
