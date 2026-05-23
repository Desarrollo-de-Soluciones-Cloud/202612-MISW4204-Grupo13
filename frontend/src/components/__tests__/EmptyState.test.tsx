import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import EmptyState from "../EmptyState";

describe("EmptyState", () => {
  it("renderiza el mensaje recibido", () => {
    render(
      <table>
        <tbody>
          <EmptyState colSpan={3} message="Sin datos por ahora" />
        </tbody>
      </table>,
    );

    expect(screen.getByText("Sin datos por ahora")).toBeInTheDocument();
  });
});
