import { z } from "zod";
import type { CaseState } from "./api/client";

const stateImportSchema = z
  .object({
    schema_version: z.literal("civitas.case_state.v1"),
    app_version: z.string(),
    exported_at: z.string(),
    case: z
      .object({
        id: z.string(),
        title: z.string(),
        description: z.string(),
        created_at: z.string(),
        updated_at: z.string(),
        document_ids: z.array(z.string()),
      })
      .passthrough(),
    documents: z.array(
      z
        .object({
          document: z
            .object({
              id: z.string(),
              filename: z.string(),
              sha256: z.string(),
            })
            .passthrough(),
          content_base64: z.string(),
          content_sha256: z.string(),
        })
        .passthrough(),
    ),
  })
  .passthrough();

export function parseCaseState(raw: string): CaseState {
  return stateImportSchema.parse(JSON.parse(raw)) as CaseState;
}
