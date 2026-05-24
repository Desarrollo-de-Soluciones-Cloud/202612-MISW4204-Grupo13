import { render, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
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

    expect(screen.getByRole("heading", { level: 1, name: /Panel del profesor/i })).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByRole("heading", { level: 2, name: /Mis cursos y proyectos/i })).toBeInTheDocument();
      expect(screen.getByRole("heading", { level: 2, name: /Reportes generados/i })).toBeInTheDocument();
    });

    expect(listWorkspacesMock).toHaveBeenCalled();
    expect(listReportsMock).toHaveBeenCalled();
  });

  it("renderiza estado vacio de reportes cuando no hay datos", async () => {
    render(<ProfessorDashboard user={buildUser("professor")} onLogout={vi.fn()} />);

    await waitFor(() => {
      expect(screen.getAllByText(/No hay registros disponibles/i).length).toBeGreaterThan(0);
    });
  });

  it("genera reporte semanal usando el workspace y la semana seleccionados", async () => {
    const user = userEvent.setup();
    listWeeksByPeriodMock.mockResolvedValue({
      weeks: [
        {
          id: 7,
          period_id: 10,
          number: 1,
          initial_date: "2026-01-20",
          final_date: "2026-01-26",
          week_state: "active",
        },
      ],
    });
    generateWeeklyReportMock.mockResolvedValue({ reports: [], generated_count: 1 });
    listReportsMock
      .mockResolvedValueOnce({ reports: [] })
      .mockResolvedValueOnce({ reports: [] })
      .mockResolvedValueOnce({
        reports: [
          {
            id: 99,
            workspace_id: 21,
            week_id: 7,
            assignment_id: 8,
            user_id: 2,
            file_path: "reports/99.pdf",
          },
        ],
      });

    render(<ProfessorDashboard user={buildUser("professor")} onLogout={vi.fn()} />);

    const section = await screen.findByRole("heading", { level: 2, name: /Generar reporte semanal/i });
    const form = section.closest("section");
    if (!form) throw new Error("Expected generate report section");

    const selects = within(form).getAllByRole("combobox");
    await user.selectOptions(selects[0], "21");
    await waitFor(() => {
      expect(listWeeksByPeriodMock).toHaveBeenCalledWith(10);
    });

    await user.selectOptions(selects[1], "7");
    await user.click(within(form).getByRole("button", { name: /Generar reporte PDF/i }));

    await waitFor(() => {
      expect(generateWeeklyReportMock).toHaveBeenCalledWith({
        workspace_id: 21,
        week_id: 7,
      });
    });
  });

  it("crea una vinculacion con los datos diligenciados", async () => {
    const user = userEvent.setup();
    createAssignmentMock.mockResolvedValue({});

    render(<ProfessorDashboard user={buildUser("professor")} onLogout={vi.fn()} />);

    const section = await screen.findByRole("heading", { level: 2, name: /Crear vinculacion/i });
    const form = section.closest("section");
    if (!form) throw new Error("Expected assignment section");

    const selects = within(form).getAllByRole("combobox");
    await user.selectOptions(selects[2], "2");
    await user.clear(within(form).getByLabelText(/Horas semanales/i));
    await user.type(within(form).getByLabelText(/Horas semanales/i), "8");
    await user.click(within(form).getByRole("button", { name: /Crear vinculacion/i }));

    await waitFor(() => {
      expect(createAssignmentMock).toHaveBeenCalledWith({
        user_id: 2,
        workspace_id: 21,
        role: "monitor",
        weekly_hours: 8,
      });
    });
  });
});
