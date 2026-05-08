import type { Document, Graph } from "../../lib/api/client";

type Props = {
  graph: Graph | undefined;
  documents: Document[];
};

const demoNodes = [
  { id: "dump", type: "document", label: "Leak dump" },
  { id: "mail", type: "email", label: "source@example.org" },
  { id: "addr", type: "address", label: "42 Civic Street" },
  { id: "date", type: "event", label: "2026-05-08" },
];

export function EvidenceMap({ graph, documents }: Props) {
  const nodes = graph?.nodes.length ? graph.nodes.slice(0, 10) : demoNodes;
  const edges = graph?.edges.length
    ? graph.edges.slice(0, 14)
    : [
        { id: "a", source_id: "dump", target_id: "mail" },
        { id: "b", source_id: "dump", target_id: "addr" },
        { id: "c", source_id: "dump", target_id: "date" },
      ];
  const positions = nodes.map((node, index) => {
    const angle =
      (index / Math.max(nodes.length, 1)) * Math.PI * 2 - Math.PI / 2;
    const radius = index === 0 ? 0 : 150;
    return {
      ...node,
      x: 280 + Math.cos(angle) * radius,
      y: 190 + Math.sin(angle) * radius,
    };
  });

  return (
    <div>
      <svg
        className="h-[380px] w-full"
        viewBox="0 0 560 380"
        role="img"
        aria-label="Evidence relationship map"
      >
        <rect width="560" height="380" rx="8" fill="#f6f4ee" />
        {edges.map((edge) => {
          const source =
            positions.find((node) => node.id === edge.source_id) ??
            positions[0];
          const target =
            positions.find((node) => node.id === edge.target_id) ??
            positions[0];
          return (
            <line
              key={edge.id}
              x1={source.x}
              y1={source.y}
              x2={target.x}
              y2={target.y}
              stroke="#8d826f"
              strokeWidth="2"
            />
          );
        })}
        {positions.map((node) => (
          <g key={node.id}>
            <circle
              cx={node.x}
              cy={node.y}
              r={node.type === "document" ? 32 : 24}
              fill={node.type === "document" ? "#0b5cad" : "#fffdfa"}
              stroke={node.type === "document" ? "#0b5cad" : "#bb4d00"}
              strokeWidth="3"
            />
            <text
              x={node.x}
              y={node.y + 48}
              textAnchor="middle"
              fontSize="12"
              fill="#17202a"
              className="select-none"
            >
              {shortLabel(node.label)}
            </text>
          </g>
        ))}
      </svg>
      <p className="muted mt-3">
        {documents.length
          ? `${documents.length} source documents mapped`
          : "Demo map shown until a backend case is selected"}
      </p>
    </div>
  );
}

function shortLabel(value: string) {
  return value.length > 22 ? `${value.slice(0, 21)}...` : value;
}
