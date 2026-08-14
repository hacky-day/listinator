import { beforeEach, describe, expect, it } from "vitest";

import { setLocale, translateType } from "@/i18n";
import { type Type } from "@/types";

function makeType(slug: string, name: string): Type {
  return {
    ID: "6c1d3f5e-0000-0000-0000-000000000000",
    Name: name,
    Slug: slug,
    Immutable: false,
    Color: "gray",
    Priority: 0,
  };
}

describe("translateType", () => {
  beforeEach(() => {
    setLocale("en");
  });

  it("translates a known slug into the active language", () => {
    const bakery = makeType("bakery", "🥖 Bakery");
    expect(translateType(bakery)).toBe("🥖 Bakery");
    setLocale("de");
    expect(translateType(bakery)).toBe("🥖 Backwaren");
  });

  it("falls back to the name from the database for an unknown slug", () => {
    const custom = makeType("garden-supplies", "🌱 Garden Supplies");
    expect(translateType(custom)).toBe("🌱 Garden Supplies");
    setLocale("de");
    expect(translateType(custom)).toBe("🌱 Garden Supplies");
  });

  it("falls back to the name for a type without a slug", () => {
    expect(translateType(makeType("", "📦 Miscellaneous"))).toBe(
      "📦 Miscellaneous",
    );
  });
});
