import { render, screen } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { test, expect } from "vitest";
import { CivitasWorkspace } from "./CivitasWorkspace";

test("renders public project links, build metadata, and completed controls", () => {
  const client = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  render(
    <QueryClientProvider client={client}>
      <CivitasWorkspace appVersion="0.3.0" commit="abc1234" />
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
  expect(screen.getByText(/v0.3.0 · abc1234/i)).toBeInTheDocument();
  expect(screen.getByRole("button", { name: /start fresh/i })).toBeVisible();
  expect(screen.getByRole("button", { name: /upload paste/i })).toBeDisabled();
  expect(screen.getByRole("button", { name: /load sample/i })).toBeDisabled();
  expect(screen.getByText(/import state/i)).toBeVisible();
  expect(
    screen.getByRole("button", { name: /download state/i }),
  ).toBeDisabled();
});
