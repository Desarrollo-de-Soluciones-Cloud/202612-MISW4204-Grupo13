import { render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import type { User } from "../../api/types";
import ProfessorDashboard from "../ProfessorDashboard";

const {
  createAssignmentMock,
  createWorkspaceMock,
  downloadReportMock,
  generateWeeklyReportMock,
  getMeMock,
  listPeriodsMock,
  listReportsMock,
  listTasksMock,
  listUsersMock,
  listWeeksByPeriodMock,
  listWorkspacesMock,
  toErrorMessageMock,
} = vi.hoisted(() => ({
  createAssignmentMock: vi.fn(),
  createWorkspaceMock: vi.fn(),
  downloadReportMock: vi.fn(),
  generateWeeklyReportMock: vi.fn(),
  getMeMock: vi.fn(),
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
  createWorkspace: (...args: unknown[]) => createWorkspaceMock(...args),
  downloadReport: (...args: unknown[]) => downloadReportMock(...args),
  generateWeeklyReport: (...args: unknown[]) => generateWeeklyReportMock(...args),
  getMe: (...args: unknown[]) => getMeMock(...args),
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
    id: 10,
    name: "Prof",
    email: "prof@example.com",
    global_role: role,
  };
}

describe("ProfessorDashboard", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    toErrorMessageMock.mockReturnValue("error");
    getMeMock.mockResolvedValue(buildUser("professor"));
    listWorkspacesMock.mockResolvedValue({
      workspaces: [
        {
          id: 21,
          period_id: 10,
          user_id: 10,
          name: "Proyecto IA",
          type: "project",
          initial_date: "2026-01-20",
          final_date: "2026-05-20",
          observations: "obs",
          state: "active",
        },
      ],
    });
    listTasksMock.mockResolvedValue({ tasks: [] });
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
    listUsersMock.mockResolvedValue({
      users: [
        { id: 2, name: "Mon", email: "mon@example.com", global_role: "monitor" },
      ],
    });
    listReportsMock.mockResolvedValue({ reports: [] });
    listWeeksByPeriodMock.mockResolvedValue({ weeks: [] });
  });

  it("renderiza panel de profesor y secciones clave", async () => {
    render(<ProfessorDashboard user={buildUser("professor")} onLogout={vi.fn()} />);

    expect(screen.getByRole("heading", { level: 1, name: "Panel del profesor" })).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByRole("heading", { level: 2, name: "Mis cursos y proyectos" })).toBeInTheDocument();
      expect(screen.getByRole("heading", { level: 2, name: "Reportes generados" })).toBeInTheDocument();
    });

    expect(listWorkspacesMock).toHaveBeenCalled();
    expect(listReportsMock).toHaveBeenCalled();
  });

  it("renderiza estado vacio de reportes cuando no hay datos", async () => {
    render(<ProfessorDashboard user={buildUser("professor")} onLogout={vi.fn()} />);

    await waitFor(() => {
      expect(screen.getAllByText("No hay registros disponibles todavía.").length).toBeGreaterThan(0);
    });
  });
});