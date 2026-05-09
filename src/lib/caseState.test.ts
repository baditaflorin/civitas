import { describe, expect, test } from "vitest";
import { parseCaseState } from "./caseState";

const validState = {
  schema_version: "civitas.case_state.v1",
  app_version: "0.3.0",
  exported_at: "2026-05-10T00:00:00Z",
  case: {
    id: "case_test",
    title: "Imported case",
    description: "",
    created_at: "2026-05-10T00:00:00Z",
    updated_at: "2026-05-10T00:00:00Z",
    document_ids: ["doc_test"],
  },
  documents: [
    {
      document: {
        id: "doc_test",
        filename: "lead.txt",
        sha256: "abc",
      },
      content_base64: "ZXZpZGVuY2U=",
      content_sha256: "abc",
    },
  ],
};

describe("parseCaseState", () => {
  test("accepts v1 case state JSON", () => {
    expect(parseCaseState(JSON.stringify(validState)).case.title).toBe(
      "Imported case",
    );
  });

  test("rejects unsupported state JSON", () => {
    expect(() =>
      parseCaseState(JSON.stringify({ ...validState, schema_version: "v0" })),
    ).toThrow();
  });
});
