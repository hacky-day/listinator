export interface List {
  ID: string;
}

export interface Entry {
  ID: string;
  Name: string;
  Bought: boolean;
  Number: string;
  ListID: string;
  TypeID: string;
}
