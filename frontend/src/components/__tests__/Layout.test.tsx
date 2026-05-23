import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import Layout from "../Layout";

describe("Layout", () => {
  it("renderiza titulo, datos del usuario y children", () => {
    render(
      <Layout
        title="Panel de prueba"
        description="Descripcion corta"
        user={{
          id: 1,
          name: "Ada",
          email: "ada@example.com",
          global_role: "admin",
        }}
        onLogout={vi.fn()}
      >
        <section>contenido interno</section>
      </Layout>,
    );

    expect(screen.getByRole("heading", { level: 1, name: "Panel de prueba" })).toBeInTheDocument();
    expect(screen.getByText("Ada")).toBeInTheDocument();
    expect(screen.getByText("ada@example.com")).toBeInTheDocument();
    expect(screen.getByText("Administrador")).toBeInTheDocument();
    expect(screen.getByText("contenido interno")).toBeInTheDocument();
  });

  it("ejecuta onLogout al presionar cerrar sesion", async () => {
    const user = userEvent.setup();
    const onLogout = vi.fn();

    render(
      <Layout
        title="Panel"
        user={{
          id: 2,
          name: "Lin",
          email: "lin@example.com",
          global_role: "monitor",
        }}
        onLogout={onLogout}
      >
        <div />
      </Layout>,
    );

    await user.click(screen.getByRole("button", { name: "Cerrar sesión" }));
    expect(onLogout).toHaveBeenCalledTimes(1);
  });
});