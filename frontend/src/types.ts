export interface List {
  ID: string;
}

export interface Type {
  ID: string;
  Name: string;
  Slug: string;
  Immutable: boolean;
  Color: string;
  Priority: number;
}

export interface Entry {
  ID: string;
  Name: string;
  Bought: boolean;
  Number: string;
  ListID: string;
  TypeID: string;
  // When the user actually made the edit (captured at the moment of the
  // optimistic local write), as opposed to when the request reached the
  // server. Drives edit-time last-write-wins for offline writes.
  ClientUpdatedAt?: string;
}

export interface User {
  ID: string;
  Name: string;
  IsAdmin: boolean;
}

export type ContextmenuAction = {
  label: string;
  action: string;
};
