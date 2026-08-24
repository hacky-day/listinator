import { ApiError, apiCreateEntry, apiDeleteEntry, apiUpdateEntry } from "@/api/api";
import {
  type OutboxRecord,
  deleteOutboxRecord,
  getOutboxRecord,
  getOutboxRecords,
  putEntry,
  putOutboxRecord,
} from "@/offline/db";
import type { Entry } from "@/types";

// Coalescing keeps the queue at most one record per entry: a still-queued
// create hasn't reached the server yet, so a follow-up edit is folded into
// it rather than queued separately.
export async function enqueueCreate(entry: Entry): Promise<void> {
  await putOutboxRecord({
    entryId: entry.ID,
    listID: entry.ListID,
    type: "create",
    payload: entry,
    createdAt: Date.now(),
  });
}

export async function enqueueUpdate(entry: Entry): Promise<void> {
  const existing = await getOutboxRecord(entry.ID);
  await putOutboxRecord({
    entryId: entry.ID,
    listID: entry.ListID,
    type: existing?.type === "create" ? "create" : "update",
    payload: entry,
    createdAt: existing?.createdAt ?? Date.now(),
  });
}

export async function enqueueDelete(entry: Entry): Promise<void> {
  const existing = await getOutboxRecord(entry.ID);
  if (existing?.type === "create") {
    // Neither the create nor the delete ever needs to reach the server.
    await deleteOutboxRecord(entry.ID);
    return;
  }
  await putOutboxRecord({
    entryId: entry.ID,
    listID: entry.ListID,
    type: "delete",
    payload: entry,
    createdAt: existing?.createdAt ?? Date.now(),
  });
}

export async function hasQueuedRecordsForList(
  listID: string,
): Promise<boolean> {
  const records = await getOutboxRecords();
  return records.some((record) => record.listID === listID);
}

// sendRecord reports whether the record is done with — synced, or
// non-retryably rejected — and can be dropped from the queue. A network
// error or 5xx means the server may accept it later, so flushing stops here
// to preserve per-entry ordering instead of skipping ahead.
async function sendRecord(record: OutboxRecord): Promise<boolean> {
  try {
    if (record.type === "create") {
      await apiCreateEntry(record.payload as Entry);
    } else if (record.type === "update") {
      await apiUpdateEntry(record.payload as Entry);
    } else {
      const result = await apiDeleteEntry(record.payload as Entry);
      if (!result.deleted) {
        // The server rejected the delete because a newer edit landed in the
        // meantime — restore it locally instead of retrying a delete that
        // will never be accepted. useEntries's next full refetch (right
        // after the outbox drains) picks this back up into the UI.
        await putEntry(result.entry);
      }
    }
    return true;
  } catch (error) {
    if (error instanceof ApiError && error.status !== undefined) {
      // A 404 on a delete means the entry is already gone server-side,
      // which is the desired end state, not a failure. Any other 4xx is a
      // non-retryable rejection (e.g. a bad request) — drop it and move on
      // rather than getting stuck on it forever.
      return error.status < 500;
    }
    return false;
  }
}

const SAFETY_NET_INTERVAL = 30_000;
let safetyTimer: number | null = null;
let flushing = false;

function ensureSafetyTimer() {
  if (safetyTimer === null) {
    safetyTimer = window.setInterval(() => void flushOutbox(), SAFETY_NET_INTERVAL);
  }
}

function clearSafetyTimer() {
  if (safetyTimer !== null) {
    clearInterval(safetyTimer);
    safetyTimer = null;
  }
}

// flushOutbox processes the queue strictly in createdAt order, sequentially,
// so a create always reaches the server before a later update/delete for the
// same entry. It is safe to call repeatedly/concurrently (e.g. from an
// "online" event, an SSE reconnect, and the safety-net timer all firing
// around the same time) — later callers no-op while a flush is in flight.
export async function flushOutbox(): Promise<void> {
  if (flushing) {
    return;
  }
  flushing = true;
  try {
    for (;;) {
      const [next] = await getOutboxRecords();
      if (next === undefined) {
        break;
      }
      if (!(await sendRecord(next))) {
        break;
      }
      await deleteOutboxRecord(next.entryId);
    }
  } finally {
    flushing = false;
  }
  if ((await getOutboxRecords()).length > 0) {
    ensureSafetyTimer();
  } else {
    clearSafetyTimer();
  }
}

export function scheduleFlush(): void {
  if (navigator.onLine) {
    void flushOutbox();
  }
}

window.addEventListener("online", () => void flushOutbox());
