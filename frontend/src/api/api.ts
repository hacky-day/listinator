import { type Entry, type List, type User, type Type } from "@/types.ts";

// ApiError carries the HTTP status so callers (namely the write outbox) can
// tell a non-retryable rejection (4xx) apart from a transient failure (5xx,
// or no status at all when fetch itself throws e.g. because we're offline).
export class ApiError extends Error {
  status?: number;

  constructor(message: string, status?: number) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

export async function apiFetchJSON(url: string, options = {}) {
  const response = await fetch(url, options);
  if (!response.ok) {
    throw new ApiError(`API Error, ${response.status}`, response.status);
  }
  return response.json();
}

export async function apiCreateList(list: List): Promise<List> {
  const json = await apiFetchJSON("/api/v1/lists", {
    method: "POST",
    body: JSON.stringify(list),
    headers: {
      "Content-Type": "application/json",
    },
  });
  return json as List;
}

export async function apiGetEntries(listID: string): Promise<Entry[]> {
  const json = await apiFetchJSON(`/api/v1/entries?ListID=${listID}`);
  return json as Entry[];
}

export async function apiCreateEntry(entry: Entry): Promise<Entry> {
  const json = await apiFetchJSON("/api/v1/entries", {
    method: "POST",
    body: JSON.stringify({
      ID: entry.ID,
      Name: entry.Name,
      Number: entry.Number,
      Bought: entry.Bought,
      TypeID: entry.TypeID,
      ListID: entry.ListID,
      ClientUpdatedAt: entry.ClientUpdatedAt,
    }),
    headers: {
      "Content-Type": "application/json",
    },
  });
  return json as Entry;
}

export async function apiGetEntry(entryID: string): Promise<Entry> {
  const json = await apiFetchJSON(`/api/v1/entries/${entryID}`);
  return json as Entry;
}

export async function apiUpdateEntry(entry: Entry): Promise<Entry> {
  const json = await apiFetchJSON(`/api/v1/entries/${entry.ID}`, {
    method: "PUT",
    body: JSON.stringify({
      Name: entry.Name,
      Number: entry.Number,
      Bought: entry.Bought,
      TypeID: entry.TypeID,
      ListID: entry.ListID,
      ClientUpdatedAt: entry.ClientUpdatedAt,
    }),
    headers: {
      "Content-Type": "application/json",
    },
  });
  return json as Entry;
}

// apiDeleteEntry reports deleted: false (instead of throwing) when the
// server rejects the delete as stale (409) because a newer edit landed in
// the meantime — the caller decides how to reconcile, this is not a failure.
export async function apiDeleteEntry(
  entry: Entry,
): Promise<{ deleted: boolean; entry: Entry }> {
  const response = await fetch(`/api/v1/entries/${entry.ID}`, {
    method: "DELETE",
    body: JSON.stringify({ ClientUpdatedAt: entry.ClientUpdatedAt }),
    headers: {
      "Content-Type": "application/json",
    },
  });
  if (response.status === 409) {
    return { deleted: false, entry: (await response.json()) as Entry };
  }
  if (!response.ok) {
    throw new ApiError(`API Error, ${response.status}`, response.status);
  }
  return { deleted: true, entry: (await response.json()) as Entry };
}

export async function apiGetTypes(): Promise<Type[]> {
  const json = await apiFetchJSON(`/api/v1/types`);
  return json as Type[];
}

export async function apiCreateSession(username: string, password: string) {
  const response = await fetch("api/v1/session", {
    method: "POST",
    body: JSON.stringify({
      name: username,
      password: password,
    }),
    headers: {
      "Content-Type": "application/json",
    },
  });
  if (!response.ok) {
    throw new Error(`API Error, ${response.status}`);
  }
}

export async function apiGetSession(): Promise<User> {
  const json = await apiFetchJSON("api/v1/session");
  return json as User;
}

export async function apiDeleteSession() {
  const response = await fetch("api/v1/session", {
    method: "DELETE",
  });
  if (!response.ok) {
    throw new Error(`API Error, ${response.status}`);
  }
}
