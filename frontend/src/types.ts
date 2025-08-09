export interface List {
  ID: string;
}

export interface Type {
  Name: string;
  Icon: string;
}

export interface Entry {
  ID: string;
  Name: string;
  Bought: boolean;
  Number: string;
  ListID: string;
  TypeID: string;

  // Helper
  _dirty: boolean;
}

export interface Session {
  uuid: string;
}

export type ContextmenuAction = {
  label: string;
  action: string;
};
