import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import Loading from "../Loading";

describe("Loading", () => {
  it("renderiza el texto de carga", () => {
    render(<Loading label="Cargando tareas..." />);

    expect(screen.getByText("Cargando tareas...")).toBeInTheDocument();
  });
});
