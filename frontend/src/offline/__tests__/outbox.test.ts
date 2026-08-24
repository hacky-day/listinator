import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { ApiError } from "@/api/api";
import {
  deleteOutboxRecord,
  getEntryLocal,
  getOutboxRecords,
} from "@/offline/db";
import {
  enqueueCreate,
  enqueueDelete,
  enqueueUpdate,
  flushOutbox,
} from "@/offline/outbox";
import type { Entry } from "@/types";

vi.mock("@/api/api", async (importOriginal) => {
  const actual = await importOriginal<typeof import("@/api/api")>();
  return {
    ...actual,
    apiCreateEntry: vi.fn(),
    apiUpdateEntry: vi.fn(),
    apiDeleteEntry: vi.fn(),
  };
});

const { apiCreateEntry, apiUpdateEntry, apiDeleteEntry } = await import("@/api/api");
const mockCreate = vi.mocked(apiCreateEntry);
const mockUpdate = vi.mocked(apiUpdateEntry);
const mockDelete = vi.mocked(apiDeleteEntry);

function makeEntry(overrides: Partial<Entry> = {}): Entry {
  return {
    ID: crypto.randomUUID(),
    Name: "milk",
    Number: "",
    Bought: false,
    ListID: crypto.randomUUID(),
    TypeID: crypto.randomUUID(),
    ClientUpdatedAt: new Date().toISOString(),
    ...overrides,
  };
}

beforeEach(() => {
  mockCreate.mockResolvedValue(makeEntry());
  mockUpdate.mockResolvedValue(makeEntry());
  mockDelete.mockResolvedValue({ deleted: true, entry: makeEntry() });
});

// Every test uses entries with fresh random IDs, so leftover state from a
// prior test can't corrupt a later one by key collision, but a test that
// deliberately leaves a record stuck in the queue (to assert a stalled
// flush) must still drain it here, or getOutboxRecords() in a later test
// would see it too.
afterEach(async () => {
  await flushOutbox().catch(() => undefined);
  await Promise.all(
    (await getOutboxRecords()).map((record) => deleteOutboxRecord(record.entryId)),
  );
  vi.clearAllMocks();
});

describe("enqueue coalescing", () => {
  it("folds a follow-up update into a still-queued create", async () => {
    const entry = makeEntry();
    await enqueueCreate(entry);
    await enqueueUpdate({ ...entry, Name: "eggs" });

    const records = (await getOutboxRecords()).filter((r) => r.entryId === entry.ID);
    expect(records).toHaveLength(1);
    expect(records[0].type).toBe("create");
    expect((records[0].payload as Entry).Name).toBe("eggs");
  });

  it("cancels both create and delete when a queued create is deleted before syncing", async () => {
    const entry = makeEntry();
    await enqueueCreate(entry);
    await enqueueDelete(entry);

    const records = await getOutboxRecords();
    expect(records.some((r) => r.entryId === entry.ID)).toBe(false);
  });

  it("queues a fresh update record after a delete (documents current, not prescribed, behavior)", async () => {
    const entry = makeEntry();
    await enqueueDelete(entry);
    await enqueueUpdate({ ...entry, Name: "eggs" });

    const records = (await getOutboxRecords()).filter((r) => r.entryId === entry.ID);
    expect(records).toHaveLength(1);
    expect(records[0].type).toBe("update");
  });
});

describe("flushOutbox", () => {
  it("sends queued records in createdAt order", async () => {
    const first = makeEntry();
    const second = makeEntry();
    const order: string[] = [];
    mockCreate.mockImplementation(async (e) => {
      order.push(e.ID);
      return e;
    });

    // enqueueCreate stamps createdAt from Date.now(); two calls in the same
    // test can land in the same millisecond, making the order otherwise
    // unspecified. Force them apart.
    const nowSpy = vi.spyOn(Date, "now");
    nowSpy.mockReturnValueOnce(1_000);
    await enqueueCreate(first);
    nowSpy.mockReturnValueOnce(2_000);
    await enqueueCreate(second);
    nowSpy.mockRestore();

    await flushOutbox();

    expect(order).toEqual([first.ID, second.ID]);
    expect(await getOutboxRecords()).toHaveLength(0);
  });

  it("stops at the first record that fails with a network-style error and leaves it queued", async () => {
    const first = makeEntry();
    const second = makeEntry();
    mockCreate.mockRejectedValueOnce(new TypeError("network error"));

    await enqueueCreate(first);
    await enqueueCreate(second);
    await flushOutbox();

    const remaining = await getOutboxRecords();
    expect(remaining.map((r) => r.entryId)).toEqual([first.ID, second.ID]);
  });

  it("stops at the first record that fails with a 5xx and leaves it queued", async () => {
    const entry = makeEntry();
    mockCreate.mockRejectedValueOnce(new ApiError("boom", 503));

    await enqueueCreate(entry);
    await flushOutbox();

    expect(await getOutboxRecords()).toHaveLength(1);
  });

  it("drops a record that fails with a non-409 4xx", async () => {
    const entry = makeEntry();
    mockUpdate.mockRejectedValueOnce(new ApiError("bad request", 400));

    await enqueueUpdate(entry);
    await flushOutbox();

    expect(await getOutboxRecords()).toHaveLength(0);
  });

  it("on a 409 delete conflict, drops the record and restores the entry locally", async () => {
    const entry = makeEntry();
    const resurrected = makeEntry({ ID: entry.ID, Name: "resurrected" });
    mockDelete.mockResolvedValueOnce({ deleted: false, entry: resurrected });

    await enqueueDelete(entry);
    await flushOutbox();

    expect(await getOutboxRecords()).toHaveLength(0);
    const stored = await getEntryLocal(entry.ID);
    expect(stored?.Name).toBe("resurrected");
  });
});
