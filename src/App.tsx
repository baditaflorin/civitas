import { CivitasWorkspace } from "./features/investigation/CivitasWorkspace";

export function App() {
  return <CivitasWorkspace appVersion={__APP_VERSION__} commit={__COMMIT_SHA__} />;
}
