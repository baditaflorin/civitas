import createClient from "openapi-fetch";
import type { components, paths } from "./schema";

export type CivitasCase = components["schemas"]["Case"];
export type Document = components["schemas"]["Document"];
export type ExportArtifact = components["schemas"]["Export"];
export type CaseState = components["schemas"]["CaseState"];
export type Graph = components["schemas"]["Graph"];
export type ProcessorTool = components["schemas"]["ProcessorTool"];
export type SearchResult = components["schemas"]["SearchResult"];
export type TimelineEvent = components["schemas"]["TimelineEvent"];
export type VersionInfo = components["schemas"]["VersionInfo"];

export function createCivitasClient(baseUrl: string) {
  return createClient<paths>({ baseUrl: baseUrl.replace(/\/$/, "") });
}

export function apiError(error: unknown, fallback: string): Error {
  if (hasMessage(error)) {
    return new Error(String(error.message));
  }
  return new Error(fallback);
}

export function requireData<T>(
  result: { data?: T; error?: unknown },
  fallback: string,
): T {
  if (result.error || result.data === undefined) {
    throw apiError(result.error, fallback);
  }
  return result.data;
}

function hasMessage(value: unknown): value is { message: unknown } {
  return value !== null && typeof value === "object" && "message" in value;
}
