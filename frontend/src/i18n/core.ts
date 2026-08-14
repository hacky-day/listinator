import { en, type Messages } from "./locales/en";
import { de } from "./locales/de";

// This module is deliberately free of framework and browser dependencies: only
// pure functions taking everything they need as an argument. The reactive glue
// and the access to localStorage/navigator/document live in index.ts.

export const messages = { en, de };

export type Locale = keyof typeof messages;

// availableLocales lists the locales in the order they are offered in the UI,
// each with its name written in that language itself. short is what the
// switcher in the header displays, label is used as its tooltip.
export const availableLocales: {
  locale: Locale;
  label: string;
  short: string;
}[] = [
  { locale: "en", label: "English", short: "EN" },
  { locale: "de", label: "Deutsch", short: "DE" },
];

export const fallbackLocale: Locale = "en";

// TranslationKey is the union of all dot separated paths to a string in the
// message catalog, e.g. "list.search" or "types.bakery".
type Paths<T> = {
  [K in keyof T & string]: T[K] extends string ? K : `${K}.${Paths<T[K]>}`;
}[keyof T & string];
export type TranslationKey = Paths<Messages>;

function isLocale(value: string): value is Locale {
  return value in messages;
}

/**
 * resolveLocale picks the first supported locale from the candidates. A
 * candidate matches if it equals a supported locale or starts with it followed
 * by a separator, so "de-DE" resolves to "de". Falls back to fallbackLocale.
 */
export function resolveLocale(candidates: (string | null | undefined)[]): Locale {
  for (const candidate of candidates) {
    if (!candidate) {
      continue;
    }
    const normalized = candidate.toLowerCase().split(/[-_]/)[0];
    if (isLocale(normalized)) {
      return normalized;
    }
  }
  return fallbackLocale;
}

/**
 * lookup resolves a dot separated key in the catalog of the given locale and
 * returns undefined if the key does not exist or does not point to a string.
 */
export function lookup(locale: Locale, key: string): string | undefined {
  let current: unknown = messages[locale];
  for (const part of key.split(".")) {
    if (typeof current !== "object" || current === null) {
      return undefined;
    }
    current = (current as Record<string, unknown>)[part];
  }
  return typeof current === "string" ? current : undefined;
}

/**
 * format replaces {placeholder} occurrences with the matching parameter.
 * Placeholders without a parameter are left untouched.
 */
export function format(
  message: string,
  params?: Record<string, string | number>,
): string {
  if (params === undefined) {
    return message;
  }
  return message.replace(/\{(\w+)\}/g, (match, name: string) =>
    name in params ? String(params[name]) : match,
  );
}

/**
 * translate looks the key up in the given locale, falls back to the fallback
 * locale and finally to the key itself, so a missing translation is visible
 * instead of rendering an empty string.
 */
export function translate(
  locale: Locale,
  key: string,
  params?: Record<string, string | number>,
): string {
  const message = lookup(locale, key) ?? lookup(fallbackLocale, key);
  if (message === undefined) {
    return key;
  }
  return format(message, params);
}
