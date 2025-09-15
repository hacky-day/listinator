export interface List {
  ID: string;
}

export interface Type {
  ID: string;
  Name: string;
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
