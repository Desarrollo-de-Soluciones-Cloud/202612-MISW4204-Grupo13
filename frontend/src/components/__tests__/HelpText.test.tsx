import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import HelpText from "../HelpText";

describe("HelpText", () => {
  it("renderiza el texto de ayuda", () => {
    render(<HelpText>Selecciona un periodo válido.</HelpText>);

    expect(screen.getByText("Selecciona un periodo válido.")).toBeInTheDocument();
  });
});
