import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import ErrorMessage from "../ErrorMessage";

describe("ErrorMessage", () => {
  it("renderiza el mensaje recibido con la clase de error", () => {
    const { container } = render(<ErrorMessage message="Error de validacion" />);

    expect(screen.getByText("Error de validacion")).toBeInTheDocument();
    expect(container.firstChild).toHaveClass("error-box");
  });

  it("no renderiza contenido cuando recibe null, undefined o string vacio", () => {
    const { rerender, container } = render(<ErrorMessage message={null} />);
    expect(container.firstChild).toBeNull();

    rerender(<ErrorMessage />);
    expect(container.firstChild).toBeNull();

    rerender(<ErrorMessage message="" />);
    expect(container.firstChild).toBeNull();
  });
});
