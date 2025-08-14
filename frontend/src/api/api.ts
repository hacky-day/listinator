import { type Entry, type List, type Type } from "@/types.ts";

export async function apiFetchJSON(url: string, options = {}) {
  const response = await fetch(url, options);
  if (!response.ok) {
    throw new Error(`API Error, ${response.status}`);
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

export async function apiCreateEntry(
  name: string,
  listID: string,
): Promise<Entry> {
  const json = await apiFetchJSON("/api/v1/entries", {
    method: "POST",
    body: JSON.stringify({
      Name: name,
      ListID: listID,
    }),
    headers: {
      "Content-Type": "application/json",
    },
  });
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
    }),
    headers: {
      "Content-Type": "application/json",
    },
  });
  return json as Entry;
}

export async function apiDeleteEntry(entry: Entry): Promise<Entry> {
  const json = await apiFetchJSON(`/api/v1/entries/${entry.ID}`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    },
  });
  return json as Entry;
}

export async function apiGetTypes(): Promise<Type[]> {
  const json = await apiFetchJSON(`/api/v1/types`);
  return json as Type[];
}
