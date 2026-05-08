import createClient from "openapi-fetch";
import type { components, paths } from "./schema";

export type CivitasCase = components["schemas"]["Case"];
export type Document = components["schemas"]["Document"];
export type ExportArtifact = components["schemas"]["Export"];
export type Graph = components["schemas"]["Graph"];
export type ProcessorTool = components["schemas"]["ProcessorTool"];
export type SearchResult = components["schemas"]["SearchResult"];
export type TimelineEvent = components["schemas"]["TimelineEvent"];
export type VersionInfo = components["schemas"]["VersionInfo"];

export function createCivitasClient(baseUrl: string) {
  return createClient<paths>({ baseUrl: baseUrl.replace(/\/$/, "") });
}

export function apiError(error: unknown, fallback: string): Error {
  if (error && typeof error === "object" && "message" in error) {
    return new Error(String((error as { message: unknown }).message));
  }
  return new Error(fallback);
}
