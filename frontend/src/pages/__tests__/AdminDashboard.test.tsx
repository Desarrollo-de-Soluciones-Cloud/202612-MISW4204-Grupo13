import { render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import type { User } from "../../api/types";
import AdminDashboard from "../AdminDashboard";

const {
  createAssignmentMock,
  createPeriodMock,
  createUserMock,
  createWorkspaceMock,
  downloadReportMock,
  getMeMock,
  listAssignmentsByUserMock,
  listPeriodsMock,
  listReportsMock,
  listTasksMock,
  listUsersMock,
  listWeeksByPeriodMock,
  listWorkspacesMock,
  toErrorMessageMock,
} = vi.hoisted(() => ({
  createAssignmentMock: vi.fn(),
  createPeriodMock: vi.fn(),
  createUserMock: vi.fn(),
  createWorkspaceMock: vi.fn(),
  downloadReportMock: vi.fn(),
  getMeMock: vi.fn(),
  listAssignmentsByUserMock: vi.fn(),
  listPeriodsMock: vi.fn(),
  listReportsMock: vi.fn(),
  listTasksMock: vi.fn(),
  listUsersMock: vi.fn(),
  listWeeksByPeriodMock: vi.fn(),
  listWorkspacesMock: vi.fn(),
  toErrorMessageMock: vi.fn(),
}));

vi.mock("../../api/client", () => ({
  createAssignment: (...args: unknown[]) => createAssignmentMock(...args),
  createPeriod: (...args: unknown[]) => createPeriodMock(...args),
  createUser: (...args: unknown[]) => createUserMock(...args),
  createWorkspace: (...args: unknown[]) => createWorkspaceMock(...args),
  downloadReport: (...args: unknown[]) => downloadReportMock(...args),
  getMe: (...args: unknown[]) => getMeMock(...args),
  listAssignmentsByUser: (...args: unknown[]) => listAssignmentsByUserMock(...args),
  listPeriods: (...args: unknown[]) => listPeriodsMock(...args),
  listReports: (...args: unknown[]) => listReportsMock(...args),
  listTasks: (...args: unknown[]) => listTasksMock(...args),
  listUsers: (...args: unknown[]) => listUsersMock(...args),
  listWeeksByPeriod: (...args: unknown[]) => listWeeksByPeriodMock(...args),
  listWorkspaces: (...args: unknown[]) => listWorkspacesMock(...args),
  toErrorMessage: (...args: unknown[]) => toErrorMessageMock(...args),
}));

function buildUser(role: User["global_role"]): User {
  return {
    id: 1,
    name: "Ada",
    email: "ada@example.com",
    global_role: role,
  };
}

describe("AdminDashboard", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    toErrorMessageMock.mockReturnValue("error");
    getMeMock.mockResolvedValue(buildUser("admin"));
    listUsersMock.mockResolvedValue({
      users: [
        buildUser("professor"),
        { id: 2, name: "Mon", email: "mon@example.com", global_role: "monitor" },
        { id: 3, name: "Asis", email: "asis@example.com", global_role: "assistant" },
      ],
    });
    listPeriodsMock.mockResolvedValue({
      periods: [
        {
          id: 10,
          name: "2026-1",
          initial_date: "2026-01-20",
          final_date: "2026-05-20",
          inscription_final_date: "2026-01-31",
          weeks_count: 16,
          period_state: "active",
        },
      ],
    });
    listWorkspacesMock.mockResolvedValue({
      workspaces: [
        {
          id: 20,
          period_id: 10,
          user_id: 1,
          name: "Cloud",
          type: "course",
          initial_date: "2026-01-20",
          final_date: "2026-05-20",
          observations: "ok",
          state: "active",
        },
      ],
    });
    listTasksMock.mockResolvedValue({ tasks: [] });
    listReportsMock.mockResolvedValue({ reports: [] });
    listAssignmentsByUserMock.mockResolvedValue({ assignments: [] });
    listWeeksByPeriodMock.mockResolvedValue({ weeks: [] });
  });

  it("renderiza panel y secciones principales con datos minimos", async () => {
    render(<AdminDashboard user={buildUser("admin")} onLogout={vi.fn()} />);

    expect(screen.getByRole("heading", { level: 1, name: "Panel de administración" })).toBeInTheDocument();
    expect(screen.getByRole("heading", { level: 2, name: "Registrar usuario" })).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByRole("heading", { level: 2, name: "Resumen de sesión" })).toBeInTheDocument();
      expect(screen.getByRole("heading", { level: 2, name: "Reportes generados" })).toBeInTheDocument();
    });

    expect(listUsersMock).toHaveBeenCalled();
    expect(listReportsMock).toHaveBeenCalled();
  });

  it("muestra estado vacio en tablas cuando no hay datos", async () => {
    render(<AdminDashboard user={buildUser("admin")} onLogout={vi.fn()} />);

    await waitFor(() => {
      expect(screen.getAllByText("No hay registros disponibles todavía.").length).toBeGreaterThan(0);
    });
  });
});