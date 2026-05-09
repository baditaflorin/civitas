import { render, screen } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { test, expect } from "vitest";
import { CivitasWorkspace } from "./CivitasWorkspace";

test("renders public project links and build metadata", () => {
  const client = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  render(
    <QueryClientProvider client={client}>
      <CivitasWorkspace appVersion="0.2.0" commit="abc1234" />
    </QueryClientProvider>,
  );

  expect(screen.getByRole("link", { name: /star on github/i })).toHaveAttribute(
    "href",
    "https://github.com/baditaflorin/civitas",
  );
  expect(screen.getByRole("link", { name: /paypal/i })).toHaveAttribute(
    "href",
    "https://www.paypal.com/paypalme/florinbadita",
  );
  expect(screen.getByText(/v0.2.0 · abc1234/i)).toBeInTheDocument();
});
