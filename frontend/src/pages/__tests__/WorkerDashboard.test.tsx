import { render, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";
import type { Assignment, Task, User } from "../../api/types";
import WorkerDashboard from "../WorkerDashboard";

const {
  createTaskMock,
  deleteTaskMock,
  downloadTaskAttachmentMock,
  getMeMock,
  listAssignmentsByUserMock,
  listTasksMock,
  toErrorMessageMock,
  updateTaskMock,
} = vi.hoisted(() => ({
  createTaskMock: vi.fn(),
  deleteTaskMock: vi.fn(),
  downloadTaskAttachmentMock: vi.fn(),
  getMeMock: vi.fn(),
  listAssignmentsByUserMock: vi.fn(),
  listTasksMock: vi.fn(),
  toErrorMessageMock: vi.fn(),
  updateTaskMock: vi.fn(),
}));

vi.mock("../../api/client", () => ({
  createTask: (...args: unknown[]) => createTaskMock(...args),
  deleteTask: (...args: unknown[]) => deleteTaskMock(...args),
  downloadTaskAttachment: (...args: unknown[]) => downloadTaskAttachmentMock(...args),
  getMe: (...args: unknown[]) => getMeMock(...args),
  listAssignmentsByUser: (...args: unknown[]) => listAssignmentsByUserMock(...args),
  listTasks: (...args: unknown[]) => listTasksMock(...args),
  toErrorMessage: (...args: unknown[]) => toErrorMessageMock(...args),
  updateTask: (...args: unknown[]) => updateTaskMock(...args),
}));

function buildUser(role: User["global_role"]): User {
  return {
    id: 30,
    name: "Worker",
    email: "worker@example.com",
    global_role: role,
  };
}

const assignment: Assignment = {
  id: 99,
  user_id: 30,
  workspace_id: 21,
  role: "monitor",
  weekly_hours: 8,
};

function buildTask(overrides: Partial<Task> = {}): Task {
  return {
    id: 501,
    user_id: 30,
    assignment_id: 99,
    week_id: 7,
    title: "Preparar laboratorio",
    description: "Acompanamiento en practica",
    status: "abierto",
    spent_hours: 2,
    observations: "Sin novedades",
    week_start_date: "2026-05-18",
    late: false,
    attachments: [],
    ...overrides,
  };
}

describe("WorkerDashboard", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    toErrorMessageMock.mockReturnValue("error");
    getMeMock.mockResolvedValue(buildUser("monitor"));
    listAssignmentsByUserMock.mockResolvedValue({ assignments: [assignment] });
    listTasksMock.mockResolvedValue({ tasks: [] });
  });

  it("renderiza panel del monitor con secciones y formulario principal", async () => {
    render(<WorkerDashboard user={buildUser("monitor")} onLogout={vi.fn()} />);

    expect(screen.getByRole("heading", { level: 1, name: "Panel del monitor" })).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByRole("heading", { level: 2, name: "Registrar nueva tarea" })).toBeInTheDocument();
      expect(screen.getByRole("heading", { level: 2, name: "Mis tareas reportadas" })).toBeInTheDocument();
    });

    expect(screen.getByLabelText(/^T.tulo$/)).toBeInTheDocument();
    expect(screen.getByLabelText(/^Fecha/)).toBeInTheDocument();
  });

  it("renderiza panel de asistente y estado vacio de tareas", async () => {
    getMeMock.mockResolvedValue(buildUser("assistant"));

    render(<WorkerDashboard user={buildUser("assistant")} onLogout={vi.fn()} />);

    expect(screen.getByRole("heading", { level: 1, name: "Panel del asistente graduado" })).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByText(/No hay registros disponibles/)).toBeInTheDocument();
    });
  });

  it("muestra estados vacios cuando no hay vinculaciones ni tareas", async () => {
    listAssignmentsByUserMock.mockResolvedValue({ assignments: [] });
    listTasksMock.mockResolvedValue({ tasks: [] });

    render(<WorkerDashboard user={buildUser("monitor")} onLogout={vi.fn()} />);

    await waitFor(() => {
      expect(screen.getAllByText(/No hay registros disponibles/)).toHaveLength(2);
    });

    expect(listAssignmentsByUserMock).toHaveBeenCalledWith(30);
  });

  it("renderiza una vinculacion y una tarea minima", async () => {
    listTasksMock.mockResolvedValue({ tasks: [buildTask()] });

    render(<WorkerDashboard user={buildUser("monitor")} onLogout={vi.fn()} />);

    await waitFor(() => {
      expect(screen.getByText("Preparar laboratorio")).toBeInTheDocument();
    });

    expect(screen.getByText("21")).toBeInTheDocument();
    expect(screen.getByText("Sin archivos")).toBeInTheDocument();
  });

  it("permite diligenciar y enviar el formulario de nueva tarea", async () => {
    const user = userEvent.setup();
    createTaskMock.mockResolvedValue(buildTask({ id: 700, title: "Revision de quiz" }));
    listTasksMock
      .mockResolvedValueOnce({ tasks: [] })
      .mockResolvedValueOnce({ tasks: [buildTask({ id: 700, title: "Revision de quiz" })] });

    render(<WorkerDashboard user={buildUser("monitor")} onLogout={vi.fn()} />);

    await waitFor(() => {
      expect(screen.getByRole("button", { name: "Registrar tarea" })).toBeInTheDocument();
    });

    await user.type(screen.getByLabelText(/^T.tulo$/), "Revision de quiz");
    await user.type(screen.getByLabelText(/^Descripci/), "Apoyo en calificacion");
    await user.clear(screen.getByLabelText(/^Horas/));
    await user.type(screen.getByLabelText(/^Horas/), "3");
    await user.type(screen.getByLabelText(/^Observaciones/), "Entrega completa");
    await user.type(screen.getByLabelText(/^Fecha/), "2026-05-18");
    await user.selectOptions(screen.getByLabelText(/^Estado/), "en desarrollo");
    await user.click(screen.getByRole("button", { name: "Registrar tarea" }));

    await waitFor(() => {
      expect(createTaskMock).toHaveBeenCalledWith({
        assignment_id: 99,
        title: "Revision de quiz",
        description: "Apoyo en calificacion",
        status: "en desarrollo",
        spent_hours: 3,
        observations: "Entrega completa",
        week_start_date: "2026-05-18",
        attachments: [],
        existing_attachments: [],
      });
    });
    expect(screen.getByText("Tarea registrada correctamente.")).toBeInTheDocument();
    expect(screen.getByText("Revision de quiz")).toBeInTheDocument();
  });

  it("permite editar una tarea y quitar un adjunto existente", async () => {
    const user = userEvent.setup();
    const task = buildTask({
      attachments: [
        {
          id: "att-1",
          name: "evidencia.pdf",
          file_path: "tasks/evidencia.pdf",
          content_type: "application/pdf",
          size: 1200,
        },
      ],
    });
    updateTaskMock.mockResolvedValue(buildTask({ title: "Laboratorio actualizado" }));
    listTasksMock
      .mockResolvedValueOnce({ tasks: [task] })
      .mockResolvedValueOnce({ tasks: [buildTask({ title: "Laboratorio actualizado" })] });

    render(<WorkerDashboard user={buildUser("monitor")} onLogout={vi.fn()} />);

    await user.click(await screen.findByRole("button", { name: "Editar" }));
    await user.click(screen.getByRole("button", { name: "Quitar evidencia.pdf" }));

    const editSection = screen.getByRole("heading", { level: 2, name: "Editar tarea" }).closest("section");
    expect(editSection).not.toBeNull();
    const editControls = within(editSection as HTMLElement);

    await user.clear(editControls.getByLabelText(/^T/));
    await user.type(editControls.getByLabelText(/^T/), "Laboratorio actualizado");
    await user.click(editControls.getByRole("button", { name: "Guardar cambios" }));

    await waitFor(() => {
      expect(updateTaskMock).toHaveBeenCalledWith(
        501,
        expect.objectContaining({
          title: "Laboratorio actualizado",
          existing_attachments: [],
        }),
      );
    });
    expect(screen.getByText("Tarea actualizada correctamente.")).toBeInTheDocument();
  });

  it("permite eliminar una tarea con respuesta exitosa", async () => {
    const user = userEvent.setup();
    deleteTaskMock.mockResolvedValue(undefined);
    listTasksMock
      .mockResolvedValueOnce({ tasks: [buildTask()] })
      .mockResolvedValueOnce({ tasks: [] });

    render(<WorkerDashboard user={buildUser("monitor")} onLogout={vi.fn()} />);

    await user.click(await screen.findByRole("button", { name: "Eliminar" }));

    await waitFor(() => {
      expect(deleteTaskMock).toHaveBeenCalledWith(501);
    });
    expect(screen.getByText("Tarea eliminada correctamente.")).toBeInTheDocument();
  });

  it("descarga un adjunto de tarea", async () => {
    const user = userEvent.setup();
    const createObjectURL = vi.fn(() => "blob:task-file");
    const revokeObjectURL = vi.fn();
    const linkClick = vi.spyOn(HTMLAnchorElement.prototype, "click").mockImplementation(() => undefined);
    const originalCreateObjectURL = URL.createObjectURL;
    const originalRevokeObjectURL = URL.revokeObjectURL;
    URL.createObjectURL = createObjectURL;
    URL.revokeObjectURL = revokeObjectURL;
    downloadTaskAttachmentMock.mockResolvedValue(new Blob(["pdf"]));
    listTasksMock.mockResolvedValue({
      tasks: [
        buildTask({
          attachments: [
            {
              id: "att-1",
              name: "evidencia.pdf",
              file_path: "tasks/evidencia.pdf",
              content_type: "application/pdf",
              size: 1200,
            },
          ],
        }),
      ],
    });

    try {
      render(<WorkerDashboard user={buildUser("monitor")} onLogout={vi.fn()} />);

      await user.click(await screen.findByRole("button", { name: "evidencia.pdf" }));

      await waitFor(() => {
        expect(downloadTaskAttachmentMock).toHaveBeenCalledWith(501, "att-1");
      });
      expect(createObjectURL).toHaveBeenCalled();
      expect(revokeObjectURL).toHaveBeenCalledWith("blob:task-file");
    } finally {
      linkClick.mockRestore();
      URL.createObjectURL = originalCreateObjectURL;
      URL.revokeObjectURL = originalRevokeObjectURL;
    }
  });
});
