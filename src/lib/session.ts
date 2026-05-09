import { z } from "zod";

const endpointSchema = z.string().url();
const textSchema = z.string().max(500);

export const defaultApiBaseUrl =
  import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

const keys = {
  endpoint: "civitas_api_base_url",
  selectedCase: "civitas_selected_case_id",
  search: "civitas_search_term",
} as const;

export function readStoredEndpoint(): string {
  const stored = localStorage.getItem(keys.endpoint);
  const parsed = endpointSchema.safeParse(stored);
  return parsed.success ? parsed.data : defaultApiBaseUrl;
}

export function writeStoredEndpoint(value: string) {
  localStorage.setItem(keys.endpoint, value);
}

export function readStoredSelectedCase(): string {
  return readText(keys.selectedCase, "");
}

export function writeStoredSelectedCase(value: string) {
  writeText(keys.selectedCase, value);
}

export function readStoredSearchTerm(): string {
  return readText(keys.search, "corruption");
}

export function writeStoredSearchTerm(value: string) {
  writeText(keys.search, value);
}

export function clearStoredSession() {
  Object.values(keys).forEach((key) => localStorage.removeItem(key));
}

function readText(key: string, fallback: string): string {
  const parsed = textSchema.safeParse(localStorage.getItem(key));
  return parsed.success && parsed.data ? parsed.data : fallback;
}

function writeText(key: string, value: string) {
  localStorage.setItem(key, value);
}
