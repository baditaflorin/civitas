import { z } from "zod";

const endpointSchema = z.string().url();

export const defaultApiBaseUrl =
  import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

export function readStoredEndpoint(): string {
  const stored = localStorage.getItem("civitas_api_base_url");
  const parsed = endpointSchema.safeParse(stored);
  return parsed.success ? parsed.data : defaultApiBaseUrl;
}

export function writeStoredEndpoint(value: string) {
  localStorage.setItem("civitas_api_base_url", value);
}
