import { ref } from "vue";

import { type Type } from "@/types.ts";
import {
  availableLocales,
  lookup,
  resolveLocale,
  translate,
  type Locale,
  type TranslationKey,
} from "./core";

export { availableLocales, type Locale, type TranslationKey };

const storageKey = "listinator.locale";

// The globals are accessed defensively so that importing this module outside a
// browser, e.g. in the tests, does not throw. Storage access is wrapped as
// well, because browsers may deny it, for example in private mode.
function storedLocale(): string | null {
  try {
    return globalThis.localStorage?.getItem(storageKey) ?? null;
  } catch {
    return null;
  }
}
function storeLocale(newLocale: Locale) {
  try {
    globalThis.localStorage?.setItem(storageKey, newLocale);
  } catch {
    // Not being able to remember the choice is not worth bothering the user.
  }
}
function browserLocales(): readonly string[] {
  return globalThis.navigator?.languages ?? [];
}

export const locale = ref<Locale>(
  resolveLocale([storedLocale(), ...browserLocales()]),
);

/**
 * setLocale switches the language, remembers the choice for the next visit and
 * keeps the lang attribute of the document in sync.
 */
export function setLocale(newLocale: Locale) {
  locale.value = newLocale;
  storeLocale(newLocale);
  applyDocumentLanguage();
}

function applyDocumentLanguage() {
  if (typeof document !== "undefined") {
    document.documentElement.lang = locale.value;
  }
}

/**
 * t translates the key into the currently selected language. It reads the
 * reactive locale, so templates using it re-render on a language switch.
 */
export function t(
  key: TranslationKey,
  params?: Record<string, string | number>,
): string {
  return translate(locale.value, key, params);
}

/**
 * translateType returns the translated name of an entry type. Types the
 * frontend does not know a translation for fall back to the name stored in the
 * database.
 */
export function translateType(type: Type): string {
  return lookup(locale.value, `types.${type.Slug}`) ?? type.Name;
}

applyDocumentLanguage();
