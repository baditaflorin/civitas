import {
  Activity,
  CircleDollarSign,
  Clipboard,
  Database,
  Download,
  FileArchive,
  FileInput,
  Github,
  GitPullRequestArrow,
  Link2,
  Printer,
  RotateCcw,
  Search,
  Server,
  Upload,
} from "lucide-react";
import { useEffect, useMemo, useState, type FormEvent } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  apiError,
  createCivitasClient,
  requireData,
  type CaseState,
  type VersionInfo,
} from "../../lib/api/client";
import {
  clearStoredSession,
  defaultApiBaseUrl,
  readStoredEndpoint,
  readStoredSearchTerm,
  readStoredSelectedCase,
  writeStoredEndpoint,
  writeStoredSearchTerm,
  writeStoredSelectedCase,
} from "../../lib/config";
import { parseCaseState } from "../../lib/caseState";
import { EvidenceMap } from "./EvidenceMap";
import { StatStrip } from "./StatStrip";

type Props = {
  appVersion: string;
  commit: string;
};

const repoUrl = "https://github.com/baditaflorin/civitas";
const paypalUrl = "https://www.paypal.com/paypalme/florinbadita";

export function CivitasWorkspace({ appVersion, commit }: Props) {
  const queryClient = useQueryClient();
  const [endpoint, setEndpoint] = useState(readStoredEndpoint);
  const [draftEndpoint, setDraftEndpoint] = useState(endpoint);
  const [selectedCaseId, setSelectedCaseId] = useState(readStoredSelectedCase);
  const [searchTerm, setSearchTerm] = useState(readStoredSearchTerm);
  const [pasteText, setPasteText] = useState("");
  const [workflowMessage, setWorkflowMessage] = useState("");
  const client = useMemo(() => createCivitasClient(endpoint), [endpoint]);

  const githubCommitQuery = useQuery({
    queryKey: ["github-main-commit"],
    queryFn: async () => {
      const response = await fetch(
        "https://api.github.com/repos/baditaflorin/civitas/commits/main",
      );
      if (!response.ok) {
        throw new Error("GitHub commit unavailable");
      }
      const payload = (await response.json()) as { sha?: string };
      return payload.sha?.slice(0, 7) ?? commit;
    },
    staleTime: 300_000,
  });

  const versionQuery = useQuery({
    queryKey: ["version", endpoint],
    queryFn: async () => {
      return requireData(
        await client.GET("/api/v1/version"),
        "Backend version unavailable",
      );
    },
  });

  const processorsQuery = useQuery({
    queryKey: ["processors", endpoint],
    queryFn: async () => {
      return requireData(
        await client.GET("/api/v1/processors"),
        "Processor registry unavailable",
      ).processors;
    },
  });

  const casesQuery = useQuery({
    queryKey: ["cases", endpoint],
    queryFn: async () => {
      return requireData(await client.GET("/api/v1/cases"), "Cases unavailable")
        .cases;
    },
  });

  const cases = casesQuery.data ?? [];
  const activeCaseId = selectedCaseId || cases[0]?.id || "";

  useEffect(() => {
    if (activeCaseId) {
      writeStoredSelectedCase(activeCaseId);
    }
  }, [activeCaseId]);

  useEffect(() => {
    writeStoredSearchTerm(searchTerm);
  }, [searchTerm]);

  const documentsQuery = useQuery({
    queryKey: ["documents", endpoint, activeCaseId],
    enabled: activeCaseId.length > 0,
    queryFn: async () => {
      return requireData(
        await client.GET("/api/v1/cases/{case_id}/documents", {
          params: { path: { case_id: activeCaseId } },
        }),
        "Documents unavailable",
      ).documents;
    },
  });

  const graphQuery = useQuery({
    queryKey: ["graph", endpoint, activeCaseId],
    enabled: activeCaseId.length > 0,
    queryFn: async () => {
      return requireData(
        await client.GET("/api/v1/cases/{case_id}/graph", {
          params: { path: { case_id: activeCaseId } },
        }),
        "Graph unavailable",
      );
    },
  });

  const timelineQuery = useQuery({
    queryKey: ["timeline", endpoint, activeCaseId],
    enabled: activeCaseId.length > 0,
    queryFn: async () => {
      return requireData(
        await client.GET("/api/v1/cases/{case_id}/timeline", {
          params: { path: { case_id: activeCaseId } },
        }),
        "Timeline unavailable",
      ).events;
    },
  });

  const searchQuery = useQuery({
    queryKey: ["search", endpoint, activeCaseId, searchTerm],
    enabled: activeCaseId.length > 0 && searchTerm.trim().length > 0,
    queryFn: async () => {
      return requireData(
        await client.GET("/api/v1/cases/{case_id}/search", {
          params: { path: { case_id: activeCaseId }, query: { q: searchTerm } },
        }),
        "Search unavailable",
      ).results;
    },
  });

  const createCase = useMutation({
    mutationFn: async (form: FormData) => {
      const title = String(form.get("title") ?? "").trim();
      const description = String(form.get("description") ?? "").trim();
      return requireData(
        await client.POST("/api/v1/cases", {
          body: { title, description },
        }),
        "Case creation failed",
      );
    },
    onSuccess: (created) => {
      setSelectedCaseId(created.id);
      void queryClient.invalidateQueries({ queryKey: ["cases", endpoint] });
    },
  });

  function invalidateEvidence(caseId = activeCaseId) {
    void queryClient.invalidateQueries({
      queryKey: ["documents", endpoint, caseId],
    });
    void queryClient.invalidateQueries({
      queryKey: ["graph", endpoint, caseId],
    });
    void queryClient.invalidateQueries({
      queryKey: ["timeline", endpoint, caseId],
    });
    void queryClient.invalidateQueries({
      queryKey: ["search", endpoint, caseId],
    });
  }

  const uploadDocument = useMutation({
    mutationFn: async (file: File) => {
      const form = new FormData();
      form.append("file", file);
      return requireData(
        await client.POST("/api/v1/cases/{case_id}/documents", {
          params: { path: { case_id: activeCaseId } },
          body: form as never,
        }),
        "Upload failed",
      );
    },
    onSuccess: () => {
      invalidateEvidence();
    },
  });

  const createExport = useMutation({
    mutationFn: async () => {
      return requireData(
        await client.POST("/api/v1/cases/{case_id}/exports", {
          params: { path: { case_id: activeCaseId } },
          body: { format: "markdown" },
        }),
        "Export failed",
      );
    },
  });

  const exportState = useMutation({
    mutationFn: async () => {
      return requireData(
        await client.GET("/api/v1/cases/{case_id}/state", {
          params: { path: { case_id: activeCaseId } },
        }),
        "State export failed",
      );
    },
    onSuccess: (state) => {
      downloadJSON(state, `${safeFilename(state.case.title)}-state.json`);
      setWorkflowMessage("Case state downloaded.");
    },
  });

  const importState = useMutation({
    mutationFn: async (state: CaseState) => {
      return requireData(
        await client.POST("/api/v1/case-states/import", { body: state }),
        "State import failed",
      );
    },
    onSuccess: (item) => {
      setSelectedCaseId(item.id);
      setWorkflowMessage(`Imported ${item.title}.`);
      void queryClient.invalidateQueries({ queryKey: ["cases", endpoint] });
      invalidateEvidence(item.id);
    },
  });

  const documents = documentsQuery.data ?? [];
  const processors = processorsQuery.data ?? [];
  const timeline = timelineQuery.data ?? [];
  const results = searchQuery.data ?? [];
  const connected = versionQuery.isSuccess;

  function saveEndpoint(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    writeStoredEndpoint(draftEndpoint);
    setEndpoint(draftEndpoint);
    setWorkflowMessage("API endpoint saved.");
  }

  function resetSession() {
    clearStoredSession();
    setEndpoint(defaultApiBaseUrl);
    setDraftEndpoint(defaultApiBaseUrl);
    setSelectedCaseId("");
    setSearchTerm("corruption");
    setWorkflowMessage(
      "Local UI state cleared. Backend evidence was not deleted.",
    );
  }

  async function uploadFiles(files: File[]) {
    if (!activeCaseId || files.length === 0) return;
    let completed = 0;
    const failed: string[] = [];
    for (const file of files) {
      setWorkflowMessage(
        `Uploading ${file.name} (${completed + 1}/${files.length})...`,
      );
      try {
        await uploadDocument.mutateAsync(file);
        completed += 1;
      } catch (error) {
        failed.push(
          `${file.name}: ${apiError(error, "Upload failed").message}`,
        );
      }
    }
    setWorkflowMessage(
      failed.length
        ? `Uploaded ${completed}/${files.length}. ${failed.join(" ")}`
        : `Uploaded ${completed}/${files.length} file${completed === 1 ? "" : "s"}.`,
    );
  }

  async function uploadPaste() {
    const trimmed = pasteText.trim();
    if (!trimmed || !activeCaseId) return;
    const isHTML = /^<!doctype html|^<html|<body|<article|<table/i.test(
      trimmed,
    );
    const extension = isHTML ? "html" : "txt";
    const type = isHTML ? "text/html" : "text/plain";
    const file = new File([trimmed], `pasted-evidence.${extension}`, { type });
    await uploadFiles([file]);
    setPasteText("");
  }

  async function loadSampleEvidence() {
    if (!activeCaseId) return;
    const response = await fetch(
      `${import.meta.env.BASE_URL}data/v1/sample.json`,
    );
    if (!response.ok) {
      setWorkflowMessage("Sample evidence could not be loaded.");
      return;
    }
    const body = await response.text();
    await uploadFiles([
      new File([body], "civitas-sample.json", { type: "application/json" }),
    ]);
  }

  async function importStateFile(file: File) {
    try {
      importState.mutate(parseCaseState(await file.text()));
    } catch (error) {
      setWorkflowMessage(apiError(error, "State import failed").message);
    }
  }

  async function copyExport() {
    if (!createExport.data?.body) return;
    try {
      await navigator.clipboard.writeText(createExport.data.body);
      setWorkflowMessage("Safe export copied.");
    } catch {
      setWorkflowMessage("Clipboard permission was not available.");
    }
  }

  function downloadExport() {
    const body = createExport.data?.body;
    if (!body) return;
    downloadText(
      body,
      `${safeFilename(cases.find((item) => item.id === activeCaseId)?.title ?? "civitas-export")}.md`,
    );
    setWorkflowMessage("Safe export downloaded.");
  }

  function printExport() {
    const body = createExport.data?.body;
    if (!body) return;
    printText(body);
  }

  return (
    <main className="min-h-screen bg-paper text-ink">
      <header className="border-b border-line bg-panel">
        <div className="mx-auto flex max-w-7xl flex-col gap-4 px-4 py-4 lg:flex-row lg:items-center lg:justify-between">
          <div>
            <p className="text-xs font-bold uppercase tracking-wide text-signal">
              Civitas
            </p>
            <h1 className="text-3xl font-semibold">Investigation OS</h1>
          </div>
          <nav
            className="flex flex-wrap items-center gap-2 text-sm"
            aria-label="Project links"
          >
            <a className="link-button" href={repoUrl}>
              <Github size={16} /> Star on GitHub
            </a>
            <a className="link-button" href={paypalUrl}>
              <CircleDollarSign size={16} /> PayPal
            </a>
            <span className="version-pill">
              v{appVersion} · {githubCommitQuery.data ?? commit}
            </span>
          </nav>
        </div>
      </header>

      <section className="mx-auto grid max-w-7xl gap-4 px-4 py-5 xl:grid-cols-[320px_1fr_340px]">
        <aside className="space-y-4">
          <Panel title="Settings" icon={<Server size={18} />}>
            <form className="space-y-3" onSubmit={saveEndpoint}>
              <label className="field-label" htmlFor="endpoint">
                API endpoint
              </label>
              <input
                id="endpoint"
                className="input"
                value={draftEndpoint}
                onChange={(event) => setDraftEndpoint(event.target.value)}
              />
              <button className="primary-button" type="submit">
                <Link2 size={16} /> Connect
              </button>
            </form>
            <button
              className="secondary-button mt-3"
              type="button"
              onClick={resetSession}
            >
              <RotateCcw size={16} /> Start fresh
            </button>
            <StatusLine connected={connected} version={versionQuery.data} />
            {workflowMessage && (
              <p className="muted mt-3" aria-live="polite">
                {workflowMessage}
              </p>
            )}
          </Panel>

          <Panel title="Cases" icon={<Database size={18} />}>
            <form
              className="space-y-3"
              action={(form) => createCase.mutate(form)}
            >
              <input
                className="input"
                name="title"
                placeholder="Case title"
                required
              />
              <textarea
                className="input min-h-20"
                name="description"
                placeholder="Case notes"
              />
              <button
                className="primary-button"
                type="submit"
                disabled={createCase.isPending}
              >
                <GitPullRequestArrow size={16} /> New case
              </button>
            </form>
            <div className="mt-4 space-y-2">
              {cases.map((item) => (
                <button
                  key={item.id}
                  className={
                    item.id === activeCaseId ? "case-row-active" : "case-row"
                  }
                  type="button"
                  onClick={() => setSelectedCaseId(item.id)}
                >
                  <span>{item.title}</span>
                  <small>{item.document_ids.length} docs</small>
                </button>
              ))}
            </div>
          </Panel>
        </aside>

        <section className="space-y-4">
          <StatStrip
            documents={documents.length}
            entities={documents.reduce(
              (count, doc) => count + doc.entities.length,
              0,
            )}
            events={timeline.length}
            processors={processors.filter((item) => item.available).length}
          />
          <Panel title="Evidence Map" icon={<Activity size={18} />}>
            <EvidenceMap graph={graphQuery.data} documents={documents} />
          </Panel>
          <Panel title="Search" icon={<Search size={18} />}>
            <div className="flex flex-col gap-3 sm:flex-row">
              <input
                className="input"
                value={searchTerm}
                onChange={(event) => setSearchTerm(event.target.value)}
                placeholder="Search corpus"
              />
            </div>
            <div className="mt-4 space-y-3">
              {results.map((result) => (
                <article key={result.document_id} className="result-row">
                  <strong>{result.filename}</strong>
                  <p>{result.snippet}</p>
                </article>
              ))}
              {results.length === 0 && (
                <p className="muted">No matching evidence yet.</p>
              )}
            </div>
          </Panel>
        </section>

        <aside className="space-y-4">
          <Panel title="Upload" icon={<Upload size={18} />}>
            <div
              className="drop-zone"
              onDragOver={(event) => event.preventDefault()}
              onDrop={(event) => {
                event.preventDefault();
                void uploadFiles(Array.from(event.dataTransfer.files));
              }}
            >
              <input
                className="input"
                type="file"
                multiple
                disabled={!activeCaseId || uploadDocument.isPending}
                onChange={(event) => {
                  void uploadFiles(Array.from(event.target.files ?? []));
                  event.currentTarget.value = "";
                }}
              />
              <p className="muted mt-3">
                Drop files here, select several files, or use the paste box.
              </p>
            </div>
            <textarea
              className="input mt-3 min-h-20"
              value={pasteText}
              onChange={(event) => setPasteText(event.target.value)}
              placeholder="Paste text or HTML evidence"
              disabled={!activeCaseId || uploadDocument.isPending}
            />
            <div className="mt-3 grid gap-2">
              <button
                className="secondary-button"
                type="button"
                disabled={!activeCaseId || !pasteText.trim()}
                onClick={() => void uploadPaste()}
              >
                <Clipboard size={16} /> Upload paste
              </button>
              <button
                className="secondary-button"
                type="button"
                disabled={!activeCaseId}
                onClick={() => void loadSampleEvidence()}
              >
                <FileInput size={16} /> Load sample
              </button>
              <label className="secondary-button">
                <FileInput size={16} /> Import state
                <input
                  className="sr-only"
                  type="file"
                  accept="application/json"
                  disabled={importState.isPending}
                  onChange={(event) => {
                    const file = event.target.files?.[0];
                    if (file) void importStateFile(file);
                    event.currentTarget.value = "";
                  }}
                />
              </label>
            </div>
            <p className="muted" aria-live="polite">
              {uploadDocument.isPending
                ? "Processing evidence..."
                : uploadDocument.error?.message}
            </p>
          </Panel>

          <Panel title="Documents" icon={<FileArchive size={18} />}>
            <div className="space-y-3">
              {documents.map((doc) => (
                <article key={doc.id} className="document-row">
                  <strong>{doc.filename}</strong>
                  <span>
                    {doc.shape ?? doc.content_type} · {doc.state ?? doc.status}{" "}
                    · {Math.round((doc.confidence ?? 0) * 100)}%
                  </span>
                  <small>{doc.sha256.slice(0, 12)}</small>
                  {(doc.anomalies ?? []).slice(0, 1).map((anomaly) => (
                    <small key={anomaly.id}>{anomaly.message}</small>
                  ))}
                </article>
              ))}
              {documents.length === 0 && (
                <p className="muted">No documents loaded.</p>
              )}
            </div>
          </Panel>

          <Panel title="Processors" icon={<Activity size={18} />}>
            <div className="processor-list">
              {processors.map((tool) => (
                <span
                  key={tool.name}
                  className={tool.available ? "processor-on" : "processor-off"}
                >
                  {tool.name}
                </span>
              ))}
            </div>
          </Panel>

          <Panel title="Timeline" icon={<Activity size={18} />}>
            <ol className="space-y-3">
              {timeline.slice(0, 5).map((event) => (
                <li key={event.id} className="timeline-row">
                  <time>{new Date(event.when).toLocaleDateString()}</time>
                  <span>{event.label}</span>
                </li>
              ))}
            </ol>
            <button
              className="primary-button mt-4"
              type="button"
              disabled={!activeCaseId || createExport.isPending}
              onClick={() => createExport.mutate()}
            >
              <FileArchive size={16} /> Safe export
            </button>
            {createExport.data?.body && (
              <>
                <div className="mt-3 grid gap-2">
                  <button
                    className="secondary-button"
                    type="button"
                    onClick={() => void copyExport()}
                  >
                    <Clipboard size={16} /> Copy export
                  </button>
                  <button
                    className="secondary-button"
                    type="button"
                    onClick={downloadExport}
                  >
                    <Download size={16} /> Download markdown
                  </button>
                  <button
                    className="secondary-button"
                    type="button"
                    onClick={printExport}
                  >
                    <Printer size={16} /> Print export
                  </button>
                </div>
                <pre className="export-preview">{createExport.data.body}</pre>
              </>
            )}
            <button
              className="secondary-button mt-3"
              type="button"
              disabled={!activeCaseId || exportState.isPending}
              onClick={() => exportState.mutate()}
            >
              <Download size={16} /> Download state
            </button>
            {activeCaseId && (
              <pre className="export-preview">{`curl -s ${endpoint.replace(/\/$/, "")}/api/v1/cases/${activeCaseId}/state > civitas-case-state.json`}</pre>
            )}
          </Panel>
        </aside>
      </section>
    </main>
  );
}

function Panel({
  title,
  icon,
  children,
}: {
  title: string;
  icon: React.ReactNode;
  children: React.ReactNode;
}) {
  return (
    <section className="rounded-lg border border-line bg-panel p-4 shadow-sm">
      <h2 className="mb-4 flex items-center gap-2 text-base font-semibold">
        {icon}
        {title}
      </h2>
      {children}
    </section>
  );
}

function StatusLine({
  connected,
  version,
}: {
  connected: boolean;
  version: VersionInfo | undefined;
}) {
  return (
    <div className="mt-4 rounded-md border border-line bg-paper p-3 text-sm">
      <span className={connected ? "status-dot-on" : "status-dot-off"} />
      {connected
        ? `Backend ${version?.version} · ${version?.commit}`
        : "Backend offline"}
    </div>
  );
}

function downloadJSON(value: CaseState, filename: string) {
  downloadText(JSON.stringify(value, null, 2), filename, "application/json");
}

function downloadText(
  body: string,
  filename: string,
  type = "text/markdown;charset=utf-8",
) {
  const url = URL.createObjectURL(new Blob([body], { type }));
  const link = document.createElement("a");
  link.href = url;
  link.download = filename;
  link.click();
  URL.revokeObjectURL(url);
}

function printText(body: string) {
  const popup = window.open("", "_blank", "noopener,noreferrer");
  if (!popup) return;
  popup.document.write(
    `<pre style="white-space:pre-wrap;font:14px/1.5 ui-monospace,Menlo,monospace;">${escapeHTML(body)}</pre>`,
  );
  popup.document.close();
  popup.focus();
  popup.print();
}

function escapeHTML(value: string) {
  return value
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;");
}

function safeFilename(value: string) {
  const cleaned = value
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/(^-|-$)/g, "");
  return cleaned || "civitas";
}
