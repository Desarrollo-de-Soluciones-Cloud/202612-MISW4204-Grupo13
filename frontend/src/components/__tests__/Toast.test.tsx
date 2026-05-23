import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import Toast from "../Toast";

describe("Toast", () => {
  it("renderiza mensaje, usa output y ejecuta onClose", async () => {
    const user = userEvent.setup();
    const onClose = vi.fn();
    const { container } = render(
      <Toast type="error" message="No fue posible guardar." onClose={onClose} />,
    );

    expect(screen.getByText("No fue posible guardar.")).toBeInTheDocument();
    const output = container.querySelector("output.toast-container");
    expect(output).toBeInTheDocument();
    expect(output).toHaveAttribute("aria-live", "polite");

    await user.click(screen.getByRole("button", { name: "Cerrar notificación" }));
    expect(onClose).toHaveBeenCalledTimes(1);
  });
});
