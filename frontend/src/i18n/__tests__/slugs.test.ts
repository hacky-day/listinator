import { readdirSync, readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { describe, expect, it } from "vitest";

import { en } from "@/i18n/locales/en";

// The type slugs are defined twice: in the database migrations and as keys of
// the "types" section of the catalogs. This test reads the migrations and
// compares both sides, so a category added in the backend without a
// translation does not slip through unnoticed.
const migrationsDir = fileURLToPath(
  new URL("../../../../database/migrations", import.meta.url),
);

function slugsFromMigrations(): string[] {
  const slugs = new Set<string>();
  for (const file of readdirSync(migrationsDir)) {
    if (!file.endsWith(".sql")) {
      continue;
    }
    const sql = readFileSync(`${migrationsDir}/${file}`, "utf8");
    for (const match of sql.matchAll(/slug\s*=\s*'([^']+)'/g)) {
      slugs.add(match[1]);
    }
  }
  return [...slugs];
}

describe("type slugs", () => {
  it("are all translated", () => {
    const fromMigrations = slugsFromMigrations().sort();
    expect(fromMigrations.length).toBeGreaterThan(0);
    expect(Object.keys(en.types).sort()).toEqual(fromMigrations);
  });
});
