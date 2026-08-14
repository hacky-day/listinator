// en is the source of truth for all translatable strings. Every other locale is
// typed against it, so a missing or misspelled key fails the type check.
export const en = {
  common: {
    name: "Name",
    login: "Login",
    logout: "Logout",
    done: "Done",
    edit: "Edit",
    delete: "Delete",
    language: "Language",
  },
  home: {
    newListGuest: "New List as Guest",
    errorCreateList: "Unable to create a new list",
  },
  login: {
    title: "Listinator - Login",
    password: "Password",
    failed: "Login failed. Please check your credentials and try again.",
  },
  list: {
    search: "Search",
    recentlyBought: "Recently bought",
    errorCreateEntry: "Unable to create new entry",
    errorUpdateEntry: "Unable to update entry",
    errorDeleteEntry: "Unable to delete entry",
    errorLoadTypes: "Unable to load entry types",
    connectionLost: "Connection lost, trying to reconnect",
    connectionRestored: "Connection restored",
  },
  entryEditor: {
    title: "Edit Entry",
    number: "Number",
    type: "Type",
    errorUpdateEntry: "Unable to update entry",
    errorLoad: "Unable to connect to server",
  },
  share: {
    copied: "URL copied to clipboard successfully",
    failed: "Unable to share link",
  },
  // Keys match the slug column of the types table, see the
  // add_type_slug migration.
  types: {
    misc: "📦 Miscellaneous",
    "fruit-vegetables-nuts": "🍎 Fruit & Vegetables & Nuts",
    "canned-goods": "🥫 Canned Goods",
    "sauces-spices-dressings": "🌶️ Sauces & Spices & Dressings",
    "drinks-alcohol": "🍹 Drinks & Alcohol",
    bakery: "🥖 Bakery",
    spreads: "🥜 Spreads",
    "coffee-tea": "☕ Coffee & Tea",
    "cereals-muesli": "🥣 Cereals & Muesli",
    "pasta-rice": "🍝 Pasta & Rice",
    "cooking-baking": "🍳 Cooking & Baking",
    "meat-fish": "🍖 Meat & Fish",
    frozen: "❄️ Frozen",
    "dairy-chilled": "🧀 Dairy & Chilled",
    "sweets-snacks": "🍬 Sweets & Snacks",
    "household-baby-pets": "🏠 Household & Baby & Pets",
    "ready-meals-broth-sauce": "🍲 Ready Meals & Broth & Sauce",
  },
};

export type Messages = typeof en;
