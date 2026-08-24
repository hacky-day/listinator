import { openDB, type DBSchema, type IDBPDatabase } from "idb";
import type { Entry, Type } from "@/types";

// OutboxRecord is a queued, not-yet-synced entry mutation. There is at most
// one record per entryId at a time — enqueueing coalesces a new mutation
// into whatever is already queued for that entry instead of appending
// (see offline/outbox.ts) — so the entryId doubles as the record's key.
export interface OutboxRecord {
  entryId: string;
  listID: string;
  type: "create" | "update" | "delete";
  payload?: Entry;
  createdAt: number;
}

interface ListinatorSchema extends DBSchema {
  entries: {
    key: string;
    value: Entry;
    indexes: { "by-listID": string };
  };
  types: {
    key: string;
    value: Type;
  };
  outbox: {
    key: string;
    value: OutboxRecord;
    indexes: { "by-createdAt": number };
  };
}

let dbPromise: Promise<IDBPDatabase<ListinatorSchema>> | undefined;

function getDB() {
  if (dbPromise === undefined) {
    dbPromise = openDB<ListinatorSchema>("listinator", 2, {
      upgrade(db) {
        if (!db.objectStoreNames.contains("entries")) {
          const entries = db.createObjectStore("entries", { keyPath: "ID" });
          entries.createIndex("by-listID", "ListID");
        }
        if (!db.objectStoreNames.contains("types")) {
          db.createObjectStore("types", { keyPath: "ID" });
        }
        if (!db.objectStoreNames.contains("outbox")) {
          const outbox = db.createObjectStore("outbox", {
            keyPath: "entryId",
          });
          outbox.createIndex("by-createdAt", "createdAt");
        }
      },
    });
  }
  return dbPromise;
}

export async function getEntriesForList(listID: string): Promise<Entry[]> {
  const db = await getDB();
  return db.getAllFromIndex("entries", "by-listID", listID);
}

export async function getEntryLocal(id: string): Promise<Entry | undefined> {
  const db = await getDB();
  return db.get("entries", id);
}

export async function putEntry(entry: Entry): Promise<void> {
  const db = await getDB();
  await db.put("entries", entry);
}

export async function deleteEntryLocal(id: string): Promise<void> {
  const db = await getDB();
  await db.delete("entries", id);
}

// replaceEntriesForList overwrites the local cache for a list with a fresh
// server response, pruning any locally-cached entries that are no longer
// present (e.g. deleted by another client while this one was offline).
export async function replaceEntriesForList(
  listID: string,
  entries: Entry[],
): Promise<void> {
  const db = await getDB();
  const tx = db.transaction("entries", "readwrite");
  const existingIDs = await tx.store.index("by-listID").getAllKeys(listID);
  const freshIDs = new Set(entries.map((entry) => entry.ID));
  await Promise.all([
    ...existingIDs
      .filter((id) => !freshIDs.has(id))
      .map((id) => tx.store.delete(id)),
    ...entries.map((entry) => tx.store.put(entry)),
  ]);
  await tx.done;
}

export async function getTypesLocal(): Promise<Type[]> {
  const db = await getDB();
  return db.getAll("types");
}

export async function replaceTypes(types: Type[]): Promise<void> {
  const db = await getDB();
  const tx = db.transaction("types", "readwrite");
  await tx.store.clear();
  await Promise.all(types.map((type) => tx.store.put(type)));
  await tx.done;
}

// getOutboxRecords returns every queued mutation, oldest first, so callers
// can flush them in the order they were made.
export async function getOutboxRecords(): Promise<OutboxRecord[]> {
  const db = await getDB();
  return db.getAllFromIndex("outbox", "by-createdAt");
}

export async function getOutboxRecord(
  entryId: string,
): Promise<OutboxRecord | undefined> {
  const db = await getDB();
  return db.get("outbox", entryId);
}

export async function putOutboxRecord(record: OutboxRecord): Promise<void> {
  const db = await getDB();
  await db.put("outbox", record);
}

export async function deleteOutboxRecord(entryId: string): Promise<void> {
  const db = await getDB();
  await db.delete("outbox", entryId);
}
