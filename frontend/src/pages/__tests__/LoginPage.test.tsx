import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";
import type { AuthResponse } from "../../api/types";
import LoginPage from "../LoginPage";

const { signInMock, toErrorMessageMock } = vi.hoisted(() => ({
  signInMock: vi.fn(),
  toErrorMessageMock: vi.fn(),
}));

vi.mock("../../api/client", () => ({
  signIn: (...args: unknown[]) => signInMock(...args),
  toErrorMessage: (...args: unknown[]) => toErrorMessageMock(...args),
}));

const successAuth: AuthResponse = {
  access_token: "token-123",
  token_type: "bearer",
  expires_in: 3600,
  user: {
    id: 1,
    name: "Admin",
    email: "admin@example.com",
    global_role: "admin",
  },
};

describe("LoginPage", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renderiza campos de email y password", () => {
    render(<LoginPage onSignedIn={vi.fn()} />);

    expect(screen.getByLabelText("Correo electrónico")).toBeInTheDocument();
    expect(screen.getByLabelText("Contraseña")).toBeInTheDocument();
  });

  it("llama al callback cuando el login es exitoso", async () => {
    const user = userEvent.setup();
    const onSignedIn = vi.fn();
    signInMock.mockResolvedValue(successAuth);

    render(<LoginPage onSignedIn={onSignedIn} />);

    await user.type(screen.getByLabelText("Contraseña"), "secret123");
    await user.click(screen.getByRole("button", { name: "Entrar al sistema" }));

    await waitFor(() => {
      expect(signInMock).toHaveBeenCalledWith("admin@example.com", "secret123");
      expect(onSignedIn).toHaveBeenCalledWith(successAuth);
    });
  });

  it("muestra mensaje de error si login falla", async () => {
    const user = userEvent.setup();
    signInMock.mockRejectedValue(new Error("boom"));
    toErrorMessageMock.mockReturnValue("Credenciales inválidas");

    render(<LoginPage onSignedIn={vi.fn()} />);

    await user.type(screen.getByLabelText("Contraseña"), "bad-password");
    await user.click(screen.getByRole("button", { name: "Entrar al sistema" }));

    expect(await screen.findByText("Credenciales inválidas")).toBeInTheDocument();
  });
});
