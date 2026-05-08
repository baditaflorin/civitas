type Props = {
  documents: number;
  entities: number;
  events: number;
  processors: number;
};

export function StatStrip({ documents, entities, events, processors }: Props) {
  const items = [
    ["Documents", documents],
    ["Entities", entities],
    ["Timeline", events],
    ["Tools", processors],
  ] as const;

  return (
    <section className="grid grid-cols-2 gap-3 md:grid-cols-4">
      {items.map(([label, value]) => (
        <div key={label} className="rounded-lg border border-line bg-panel p-4">
          <p className="text-sm text-slate-600">{label}</p>
          <strong className="mt-1 block text-3xl">{value}</strong>
        </div>
      ))}
    </section>
  );
}
