/**
 * Perform fuzzy search matching against Vietnamese city names
 * Handles accent-insensitive matching
 */
export function fuzzyMatchCity(city: string, query: string): boolean {
  const normalizedCity = normalizeText(city);
  const normalizedQuery = normalizeText(query);
  const tokens = normalizedQuery.split(/\s+/).filter(Boolean);
  return tokens.every((token) => normalizedCity.includes(token));
}

/**
 * Normalize Vietnamese text for comparison
 * Removes diacritics and converts to lowercase
 */
export function normalizeText(text: string): string {
  return text
    .normalize("NFD")
    .replace(/\p{Diacritic}/gu, "")
    .toLowerCase();
}
