import { describe, expect, it } from "vitest";

import { format, lookup, resolveLocale, translate } from "@/i18n/core";

describe("resolveLocale", () => {
  it("takes the first supported candidate", () => {
    expect(resolveLocale(["de", "en"])).toBe("de");
    expect(resolveLocale(["en", "de"])).toBe("en");
  });

  it("matches region variants by their language prefix", () => {
    expect(resolveLocale(["de-DE"])).toBe("de");
    expect(resolveLocale(["EN_US"])).toBe("en");
  });

  it("skips unsupported and empty candidates", () => {
    expect(resolveLocale([null, "", "fr-FR", "de"])).toBe("de");
  });

  it("falls back to english without a usable candidate", () => {
    expect(resolveLocale([])).toBe("en");
    expect(resolveLocale(["fr", "es"])).toBe("en");
  });
});

describe("lookup", () => {
  it("resolves a nested key", () => {
    expect(lookup("en", "list.search")).toBe("Search");
    expect(lookup("de", "list.search")).toBe("Suchen");
  });

  it("returns undefined for unknown keys", () => {
    expect(lookup("en", "list.doesNotExist")).toBeUndefined();
    expect(lookup("en", "doesNotExist.search")).toBeUndefined();
  });

  it("returns undefined when the key points at a group instead of a string", () => {
    expect(lookup("en", "list")).toBeUndefined();
  });
});

describe("format", () => {
  it("replaces placeholders with parameters", () => {
    expect(format("Hello {name}!", { name: "Sven" })).toBe("Hello Sven!");
    expect(format("{a} and {b}", { a: 1, b: 2 })).toBe("1 and 2");
  });

  it("leaves placeholders without a parameter untouched", () => {
    expect(format("Hello {name}!", { other: "x" })).toBe("Hello {name}!");
    expect(format("Hello {name}!")).toBe("Hello {name}!");
  });
});

describe("translate", () => {
  it("uses the requested locale", () => {
    expect(translate("de", "common.login")).toBe("Anmelden");
  });

  it("falls back to the key when no locale knows it", () => {
    expect(translate("de", "nope.nope")).toBe("nope.nope");
  });
});
