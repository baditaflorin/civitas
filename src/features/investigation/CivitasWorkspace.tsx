import {
  Activity,
  CircleDollarSign,
  Database,
  FileArchive,
  Github,
  GitPullRequestArrow,
  Link2,
  Search,
  Server,
  Upload,
} from "lucide-react";
import { useMemo, useState, type FormEvent } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  apiError,
  createCivitasClient,
  type CivitasCase,
  type Document,
  type ExportArtifact,
  type Graph,
  type ProcessorTool,
  type SearchResult,
  type TimelineEvent,
  type VersionInfo,
} from "../../lib/api/client";
import { readStoredEndpoint, writeStoredEndpoint } from "../../lib/config";
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
  const [selectedCaseId, setSelectedCaseId] = useState<string>("");
  const [searchTerm, setSearchTerm] = useState("corruption");
  const client = useMemo(() => createCivitasClient(endpoint), [endpoint]);

  const versionQuery = useQuery({
    queryKey: ["version", endpoint],
    queryFn: async () => {
      const { data, error } = await client.GET("/api/v1/version");
      if (error || !data) throw apiError(error, "Backend version unavailable");
      return data as VersionInfo;
    },
  });

  const processorsQuery = useQuery({
    queryKey: ["processors", endpoint],
    queryFn: async () => {
      const { data, error } = await client.GET("/api/v1/processors");
      if (error || !data)
        throw apiError(error, "Processor registry unavailable");
      return data.processors as ProcessorTool[];
    },
  });

  const casesQuery = useQuery({
    queryKey: ["cases", endpoint],
    queryFn: async () => {
      const { data, error } = await client.GET("/api/v1/cases");
      if (error || !data) throw apiError(error, "Cases unavailable");
      return data.cases as CivitasCase[];
    },
  });

  const cases = casesQuery.data ?? [];
  const activeCaseId = selectedCaseId || cases[0]?.id || "";

  const documentsQuery = useQuery({
    queryKey: ["documents", endpoint, activeCaseId],
    enabled: activeCaseId.length > 0,
    queryFn: async () => {
      const { data, error } = await client.GET(
        "/api/v1/cases/{case_id}/documents",
        {
          params: { path: { case_id: activeCaseId } },
        },
      );
      if (error || !data) throw apiError(error, "Documents unavailable");
      return data.documents as Document[];
    },
  });

  const graphQuery = useQuery({
    queryKey: ["graph", endpoint, activeCaseId],
    enabled: activeCaseId.length > 0,
    queryFn: async () => {
      const { data, error } = await client.GET(
        "/api/v1/cases/{case_id}/graph",
        {
          params: { path: { case_id: activeCaseId } },
        },
      );
      if (error || !data) throw apiError(error, "Graph unavailable");
      return data as Graph;
    },
  });

  const timelineQuery = useQuery({
    queryKey: ["timeline", endpoint, activeCaseId],
    enabled: activeCaseId.length > 0,
    queryFn: async () => {
      const { data, error } = await client.GET(
        "/api/v1/cases/{case_id}/timeline",
        {
          params: { path: { case_id: activeCaseId } },
        },
      );
      if (error || !data) throw apiError(error, "Timeline unavailable");
      return data.events as TimelineEvent[];
    },
  });

  const searchQuery = useQuery({
    queryKey: ["search", endpoint, activeCaseId, searchTerm],
    enabled: activeCaseId.length > 0 && searchTerm.trim().length > 0,
    queryFn: async () => {
      const { data, error } = await client.GET(
        "/api/v1/cases/{case_id}/search",
        {
          params: { path: { case_id: activeCaseId }, query: { q: searchTerm } },
        },
      );
      if (error || !data) throw apiError(error, "Search unavailable");
      return data.results as SearchResult[];
    },
  });

  const createCase = useMutation({
    mutationFn: async (form: FormData) => {
      const title = String(form.get("title") ?? "").trim();
      const description = String(form.get("description") ?? "").trim();
      const { data, error } = await client.POST("/api/v1/cases", {
        body: { title, description },
      });
      if (error || !data) throw apiError(error, "Case creation failed");
      return data as CivitasCase;
    },
    onSuccess: (created) => {
      setSelectedCaseId(created.id);
      void queryClient.invalidateQueries({ queryKey: ["cases", endpoint] });
    },
  });

  const uploadDocument = useMutation({
    mutationFn: async (file: File) => {
      const form = new FormData();
      form.append("file", file);
      const { data, error } = await client.POST(
        "/api/v1/cases/{case_id}/documents",
        {
          params: { path: { case_id: activeCaseId } },
          body: form as never,
        },
      );
      if (error || !data) throw apiError(error, "Upload failed");
      return data as Document;
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({
        queryKey: ["documents", endpoint, activeCaseId],
      });
      void queryClient.invalidateQueries({
        queryKey: ["graph", endpoint, activeCaseId],
      });
      void queryClient.invalidateQueries({
        queryKey: ["timeline", endpoint, activeCaseId],
      });
      void queryClient.invalidateQueries({
        queryKey: ["search", endpoint, activeCaseId],
      });
    },
  });

  const createExport = useMutation({
    mutationFn: async () => {
      const { data, error } = await client.POST(
        "/api/v1/cases/{case_id}/exports",
        {
          params: { path: { case_id: activeCaseId } },
          body: { format: "markdown" },
        },
      );
      if (error || !data) throw apiError(error, "Export failed");
      return data as ExportArtifact;
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
              v{appVersion} · {commit}
            </span>
          </nav>
        </div>
      </header>

      <section className="mx-auto grid max-w-7xl gap-4 px-4 py-5 xl:grid-cols-[320px_1fr_340px]">
        <aside className="space-y-4">
          <Panel title="Backend" icon={<Server size={18} />}>
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
            <StatusLine connected={connected} version={versionQuery.data} />
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
            <input
              className="input"
              type="file"
              disabled={!activeCaseId || uploadDocument.isPending}
              onChange={(event) => {
                const file = event.target.files?.[0];
                if (file) uploadDocument.mutate(file);
              }}
            />
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
                  <span>{doc.content_type}</span>
                  <small>{doc.sha256.slice(0, 12)}</small>
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
              <pre className="export-preview">{createExport.data.body}</pre>
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
