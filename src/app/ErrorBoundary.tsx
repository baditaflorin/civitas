import { Component, type ErrorInfo, type ReactNode } from "react";

type Props = {
  children: ReactNode;
};

type State = {
  error: Error | null;
};

export class ErrorBoundary extends Component<Props, State> {
  state: State = { error: null };

  static getDerivedStateFromError(error: Error): State {
    return { error };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    if (import.meta.env.DEV) {
      console.error(error, info);
    }
  }

  render() {
    if (this.state.error) {
      return (
        <main className="grid min-h-screen place-items-center bg-paper p-6 text-ink">
          <section className="w-full max-w-xl rounded-lg border border-line bg-panel p-6 shadow-soft">
            <h1 className="text-2xl font-semibold">Civitas failed to start</h1>
            <p className="mt-3 text-sm text-slate-700">{this.state.error.message}</p>
          </section>
        </main>
      );
    }
    return this.props.children;
  }
}
