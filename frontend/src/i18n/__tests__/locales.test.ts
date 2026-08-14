import { describe, expect, it } from "vitest";

import { messages, type Locale } from "@/i18n/core";
import { en } from "@/i18n/locales/en";

// flatten turns a nested catalog into dot separated key/value pairs, so the
// catalogs can be compared regardless of their nesting.
function flatten(value: unknown, prefix = ""): Record<string, string> {
  if (typeof value === "string") {
    return { [prefix]: value };
  }
  const result: Record<string, string> = {};
  for (const [key, child] of Object.entries(value as Record<string, unknown>)) {
    Object.assign(result, flatten(child, prefix === "" ? key : `${prefix}.${key}`));
  }
  return result;
}

const reference = flatten(en);
const locales = Object.keys(messages) as Locale[];

describe.each(locales)("catalog %s", (locale) => {
  const catalog = flatten(messages[locale]);

  it("has exactly the keys of the english catalog", () => {
    expect(Object.keys(catalog).sort()).toEqual(Object.keys(reference).sort());
  });

  it("has no empty translations", () => {
    const empty = Object.entries(catalog)
      .filter(([, value]) => value.trim() === "")
      .map(([key]) => key);
    expect(empty).toEqual([]);
  });
});
