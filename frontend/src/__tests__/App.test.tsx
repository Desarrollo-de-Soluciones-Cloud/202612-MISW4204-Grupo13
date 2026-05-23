import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";
import type { AuthResponse } from "../api/types";
import type { AuthSession } from "../auth/authStorage";
import App from "../App";

const { authStorageMock } = vi.hoisted(() => ({
  authStorageMock: {
    getAuthSession: vi.fn(),
    saveAuthSession: vi.fn(),
    clearAuthSession: vi.fn(),
  },
}));

vi.mock("../auth/authStorage", () => authStorageMock);

vi.mock("../pages/LoginPage", () => ({
  default: ({ onSignedIn }: { onSignedIn: (auth: AuthResponse) => void }) => (
    <div>
      <p>login-page</p>
      <button
        type="button"
        onClick={() =>
          onSignedIn({
            access_token: "token-new",
            token_type: "bearer",
            expires_in: 3600,
            user: {
              id: 9,
              name: "Nuevo",
              email: "nuevo@example.com",
              global_role: "admin",
            },
          })
        }
      >
        do-login
      </button>
    </div>
  ),
}));

vi.mock("../pages/AdminDashboard", () => ({
  default: ({ onLogout }: { onLogout: () => void }) => (
    <div>
      <p>admin-dashboard</p>
      <button type="button" onClick={onLogout}>
        do-logout
      </button>
    </div>
  ),
}));

vi.mock("../pages/ProfessorDashboard", () => ({
  default: () => <p>professor-dashboard</p>,
}));

vi.mock("../pages/WorkerDashboard", () => ({
  default: () => <p>worker-dashboard</p>,
}));

function createSession(role: AuthSession["user"]["global_role"]): AuthSession {
  return {
    accessToken: "token-123",
    user: {
      id: 1,
      name: "User",
      email: "user@example.com",
      global_role: role,
    },
  };
}

describe("App", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("muestra LoginPage cuando no hay sesion", () => {
    authStorageMock.getAuthSession.mockReturnValue(null);

    render(<App />);

    expect(screen.getByText("login-page")).toBeInTheDocument();
  });

  it("muestra dashboard admin cuando hay sesion de administrador", () => {
    authStorageMock.getAuthSession.mockReturnValue(createSession("admin"));

    render(<App />);

    expect(screen.getByText("admin-dashboard")).toBeInTheDocument();
  });

  it("muestra dashboard professor cuando hay sesion de profesor", () => {
    authStorageMock.getAuthSession.mockReturnValue(createSession("professor"));

    render(<App />);

    expect(screen.getByText("professor-dashboard")).toBeInTheDocument();
  });

  it("muestra dashboard worker cuando hay sesion de monitor", () => {
    authStorageMock.getAuthSession.mockReturnValue(createSession("monitor"));

    render(<App />);

    expect(screen.getByText("worker-dashboard")).toBeInTheDocument();
  });

  it("muestra dashboard worker cuando hay sesion de asistente", () => {
    authStorageMock.getAuthSession.mockReturnValue(createSession("assistant"));

    render(<App />);

    expect(screen.getByText("worker-dashboard")).toBeInTheDocument();
  });

  it("guarda la sesion cuando LoginPage informa ingreso exitoso", async () => {
    const user = userEvent.setup();
    authStorageMock.getAuthSession.mockReturnValue(null);

    render(<App />);
    await user.click(screen.getByRole("button", { name: "do-login" }));

    expect(authStorageMock.saveAuthSession).toHaveBeenCalledWith({
      accessToken: "token-new",
      user: {
        id: 9,
        name: "Nuevo",
        email: "nuevo@example.com",
        global_role: "admin",
      },
    });
    expect(screen.getByText("admin-dashboard")).toBeInTheDocument();
  });

  it("logout limpia sesion y vuelve a login", async () => {
    const user = userEvent.setup();
    authStorageMock.getAuthSession.mockReturnValue(createSession("admin"));

    render(<App />);
    await user.click(screen.getByRole("button", { name: "do-logout" }));

    expect(authStorageMock.clearAuthSession).toHaveBeenCalledTimes(1);
    expect(screen.getByText("login-page")).toBeInTheDocument();
  });
});
