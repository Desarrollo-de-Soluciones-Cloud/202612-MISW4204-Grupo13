import { FormEvent, useEffect, useMemo, useState } from "react";
import {
  createTask,
  deleteTask,
  downloadTaskAttachment,
  getMe,
  listAssignmentsByUser,
  listTasks,
  toErrorMessage,
  updateTask,
} from "../api/client";
import type { Assignment, Task, TaskAttachment, TaskStatus, UpdateTaskPayload, User } from "../api/types";
import EmptyState from "../components/EmptyState";
import HelpText from "../components/HelpText";
import Layout from "../components/Layout";
import Loading from "../components/Loading";
import Toast from "../components/Toast";
import useToast from "../components/useToast";

interface WorkerDashboardProps {
  user: User;
  onLogout: () => void;
}

interface TaskFormState {
  assignment_id: string;
  title: string;
  description: string;
  status: TaskStatus;
  spent_hours: string;
  observations: string;
  week_start_date: string;
  attachments: File[];
  existing_attachments: TaskAttachment[];
}

function buildDefaultFormState(): TaskFormState {
  return {
    assignment_id: "",
    title: "",
    description: "",
    status: "abierto",
    spent_hours: "1",
    observations: "",
    week_start_date: "",
    attachments: [],
    existing_attachments: [],
  };
}

const statusOptions: TaskStatus[] = ["abierto", "en desarrollo", "finalizado"];

function normalizeTaskAttachments(task: Task): Task {
  return {
    ...task,
    attachments: Array.isArray(task.attachments) ? task.attachments : [],
  };
}

export default function WorkerDashboard({ user, onLogout }: WorkerDashboardProps) {
  const [me, setMe] = useState<User | null>(null);
  const [assignments, setAssignments] = useState<Assignment[]>([]);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [createForm, setCreateForm] = useState<TaskFormState>(buildDefaultFormState());
  const [editTaskId, setEditTaskId] = useState<number | null>(null);
  const [editForm, setEditForm] = useState<TaskFormState>(buildDefaultFormState());
  const [loading, setLoading] = useState(false);
  const { toast, showToast, clearToast } = useToast();

  const assignmentIds = useMemo(() => assignments.map((item) => item.id), [assignments]);

  const isAssistant = user.global_role === "assistant";

  const dashboardTitle = isAssistant ? "Panel del asistente graduado" : "Panel del monitor";

  const loadData = async () => {
    setLoading(true);

    try {
      const meResult = await getMe();
      const [assignmentResult, taskResult] = await Promise.all([
        listAssignmentsByUser(meResult.id),
        listTasks(),
      ]);

      setMe(meResult);
      setAssignments(assignmentResult.assignments);
      setTasks(taskResult.tasks.map(normalizeTaskAttachments));
      setCreateForm((previous): TaskFormState => ({
        ...previous,
        assignment_id: previous.assignment_id || String(assignmentResult.assignments[0]?.id ?? ""),
      }));
    } catch (err) {
      showToast(toErrorMessage(err), "error");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadData();
  }, []);

  const toPayload = (form: TaskFormState): UpdateTaskPayload => ({
    assignment_id: Number(form.assignment_id),
    title: form.title,
    description: form.description,
    status: form.status,
    spent_hours: Number(form.spent_hours),
    observations: form.observations,
    week_start_date: form.week_start_date,
    attachments: form.attachments,
    existing_attachments: form.existing_attachments,
  });

  const handleCreateTask = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    clearToast();

    try {
      await createTask(toPayload(createForm));
      showToast("Tarea registrada correctamente.", "success");
      setCreateForm({
        ...buildDefaultFormState(),
        assignment_id: String(assignments[0]?.id ?? ""),
      });
      const taskResult = await listTasks();
      setTasks(taskResult.tasks.map(normalizeTaskAttachments));
    } catch (err) {
      showToast(toErrorMessage(err), "error");
    }
  };

  const startEdit = (task: Task) => {
    setEditTaskId(task.id);
    setEditForm({
      assignment_id: String(task.assignment_id),
      title: task.title,
      description: task.description,
      status: task.status,
      spent_hours: String(task.spent_hours),
      observations: task.observations,
      week_start_date: task.week_start_date,
      attachments: [],
      existing_attachments: Array.isArray(task.attachments) ? task.attachments : [],
    });
  };

  const handleDownloadAttachment = async (taskId: number, attachmentId: string, fileName: string) => {
    clearToast();

    try {
      const blob = await downloadTaskAttachment(taskId, attachmentId);
      const url = URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.download = fileName;
      link.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      showToast(toErrorMessage(err), "error");
    }
  };

  const handleUpdateTask = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (!editTaskId) {
      return;
    }

    clearToast();

    try {
      await updateTask(editTaskId, toPayload(editForm));
      showToast("Tarea actualizada correctamente.", "success");
      setEditTaskId(null);
      const taskResult = await listTasks();
      setTasks(taskResult.tasks.map(normalizeTaskAttachments));
    } catch (err) {
      showToast(toErrorMessage(err), "error");
    }
  };

  const handleDeleteTask = async (id: number) => {
    clearToast();

    try {
      await deleteTask(id);
      showToast("Tarea eliminada correctamente.", "success");
      const taskResult = await listTasks();
      setTasks(taskResult.tasks.map(normalizeTaskAttachments));
    } catch (err) {
      showToast(toErrorMessage(err), "error");
    }
  };

  return (
    <Layout
      title={dashboardTitle}
      description="Consulta tus vinculaciones y registra tus tareas semanales."
      user={user}
      onLogout={onLogout}
    >
      <div className="actions-row">
        <button onClick={() => void loadData()} disabled={loading}>
          Recargar datos
        </button>
      </div>

      {loading && <Loading label="Actualizando información..." />}
      {toast ? <Toast type={toast.type} message={toast.message} onClose={clearToast} /> : null}

      <section className="card info-card">
        <h2>Errores frecuentes</h2>
        <ul className="list-help">
          <li>task week is not active</li>
          <li>insufficient permissions</li>
          <li>task status is invalid</li>
          <li>period inscription is closed</li>
          <li>assignment already exists</li>
          <li>monitor weekly hours cannot exceed 12</li>
          <li>assistant weekly hours cannot exceed 22</li>
        </ul>
      </section>

      <section className="card">
        <h2>Resumen de sesión</h2>
        {me ? (
          <table>
            <tbody>
              <tr>
                <th>ID</th>
                <td>{me.id}</td>
              </tr>
              <tr>
                <th>Nombre</th>
                <td>{me.name}</td>
              </tr>
              <tr>
                <th>Rol</th>
                <td>{me.global_role}</td>
              </tr>
            </tbody>
          </table>
        ) : (
          <p className="muted">Sin datos</p>
        )}
      </section>

      <section className="card">
        <h2>Mis vinculaciones</h2>
        <p className="muted">Aquí aparecen los cursos o proyectos en los que estás vinculado.</p>
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>ID del curso/proyecto</th>
                <th>Rol</th>
                <th>Horas semanales</th>
              </tr>
            </thead>
            <tbody>
              {assignments.length === 0 ? (
                <EmptyState colSpan={4} />
              ) : (
                assignments.map((item) => (
                  <tr key={item.id}>
                    <td>{item.id}</td>
                    <td>{item.workspace_id}</td>
                    <td>{item.role}</td>
                    <td>{item.weekly_hours}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </section>

      <section className="card">
        <h2>Registrar nueva tarea</h2>
        <p className="section-desc">Reporta las actividades realizadas durante la semana activa.</p>
        <form className="form-grid" onSubmit={handleCreateTask}>
          <div className="form-field">
            <label>
              Vinculación
              <select
                value={createForm.assignment_id}
                onChange={(event) =>
                  setCreateForm((previous): TaskFormState => ({
                    ...previous,
                    assignment_id: event.target.value,
                  }))
                }
                required
              >
                <option value="" disabled>
                  Selecciona una vinculación
                </option>
                {assignmentIds.map((id) => (
                  <option key={id} value={id}>
                    {id}
                  </option>
                ))}
              </select>
            </label>
            <HelpText>Selecciona la vinculación correspondiente al curso o proyecto.</HelpText>
          </div>

          <div className="form-field">
            <label>
              Título
              <input
                value={createForm.title}
                onChange={(event) =>
                  setCreateForm((previous): TaskFormState => ({ ...previous, title: event.target.value }))
                }
                required
              />
            </label>
            <HelpText>Puedes adjuntar uno o varios archivos de soporte.</HelpText>
          </div>

          <div className="form-field">
            <label>
              Adjuntos
              <input
                type="file"
                multiple
                onChange={(event) =>
                  setCreateForm((previous): TaskFormState => ({
                    ...previous,
                    attachments: Array.from(event.target.files ?? []),
                  }))
                }
              />
            </label>
            <HelpText>{createForm.attachments.length} archivo(s) seleccionados.</HelpText>
          </div>

          <div className="form-field">
            <label>
              Descripción
              <input
                value={createForm.description}
                onChange={(event) =>
                  setCreateForm((previous): TaskFormState => ({
                    ...previous,
                    description: event.target.value,
                  }))
                }
                required
              />
            </label>
          </div>

          <div className="form-field">
            <label>
              Estado
              <select
                value={createForm.status}
                onChange={(event) =>
                  setCreateForm((previous): TaskFormState => ({
                    ...previous,
                    status: event.target.value as TaskStatus,
                  }))
                }
              >
                {statusOptions.map((status) => (
                  <option key={status} value={status}>
                    {status}
                  </option>
                ))}
              </select>
            </label>
            <HelpText>Usa solo: abierto, en desarrollo o finalizado.</HelpText>
          </div>

          <div className="form-field">
            <label>
              Horas dedicadas
              <input
                type="number"
                min={1}
                value={createForm.spent_hours}
                onChange={(event) =>
                  setCreateForm((previous): TaskFormState => ({
                    ...previous,
                    spent_hours: event.target.value,
                  }))
                }
                required
              />
            </label>
            <HelpText>Debe ser un número mayor o igual a 1.</HelpText>
          </div>

          <div className="form-field">
            <label>
              Observaciones
              <input
                value={createForm.observations}
                onChange={(event) =>
                  setCreateForm((previous): TaskFormState => ({
                    ...previous,
                    observations: event.target.value,
                  }))
                }
              />
            </label>
          </div>

          <div className="form-field">
            <label>
              Fecha de inicio de semana
              <input
                type="date"
                value={createForm.week_start_date}
                onChange={(event) =>
                  setCreateForm((previous): TaskFormState => ({
                    ...previous,
                    week_start_date: event.target.value,
                  }))
                }
                required
              />
            </label>
            <HelpText>
              Debe corresponder al lunes de la semana activa. Si la semana no está activa, el
              backend rechazará la tarea.
            </HelpText>
          </div>

          <div className="form-actions">
            <button type="submit">Registrar tarea</button>
          </div>
        </form>
      </section>

      <section className="card">
        <h2>Mis tareas reportadas</h2>
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>ID de vinculación</th>
                <th>Título</th>
                <th>Estado</th>
                <th>Horas dedicadas</th>
                <th>Fecha de inicio de semana</th>
                <th>Adjuntos</th>
                <th>Acciones</th>
              </tr>
            </thead>
            <tbody>
              {tasks.length === 0 ? (
                <EmptyState colSpan={8} />
              ) : (
                tasks.map((item) => (
                  <tr key={item.id}>
                    <td>{item.id}</td>
                    <td>{item.assignment_id}</td>
                    <td>{item.title}</td>
                    <td>{item.status}</td>
                    <td>{item.spent_hours}</td>
                    <td>{item.week_start_date}</td>
                    <td>
                      {item.attachments.length === 0 ? (
                        <span className="muted">Sin archivos</span>
                      ) : (
                        <div className="actions-row">
                          {item.attachments.map((attachment) => (
                            <button
                              key={attachment.id}
                              type="button"
                              className="button-secondary"
                              onClick={() =>
                                void handleDownloadAttachment(item.id, attachment.id, attachment.name)
                              }
                            >
                              {attachment.name}
                            </button>
                          ))}
                        </div>
                      )}
                    </td>
                    <td>
                      <div className="actions-row">
                        <button type="button" onClick={() => startEdit(item)}>
                          Editar
                        </button>
                        <button
                          type="button"
                          className="button-danger"
                          onClick={() => void handleDeleteTask(item.id)}
                        >
                          Eliminar
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </section>

      {editTaskId ? (
        <section className="card">
          <h2>Editar tarea</h2>
          <p className="section-desc">Modifica los datos de la tarea seleccionada.</p>
          <form className="form-grid" onSubmit={handleUpdateTask}>
            <div className="form-field">
              <label>
                Vinculación
                <select
                  value={editForm.assignment_id}
                  onChange={(event) =>
                    setEditForm((previous): TaskFormState => ({
                      ...previous,
                      assignment_id: event.target.value,
                    }))
                  }
                  required
                >
                  <option value="" disabled>
                    Selecciona una vinculación
                  </option>
                  {assignmentIds.map((id) => (
                    <option key={id} value={id}>
                      {id}
                    </option>
                  ))}
                </select>
              </label>
            </div>

            <div className="form-field">
              <label>
                Adjuntar nuevos archivos
                <input
                  type="file"
                  multiple
                  onChange={(event) =>
                    setEditForm((previous): TaskFormState => ({
                      ...previous,
                      attachments: Array.from(event.target.files ?? []),
                    }))
                  }
                />
              </label>
              <HelpText>{editForm.attachments.length} archivo(s) nuevos seleccionados.</HelpText>
            </div>

            <div className="form-field">
              <label>Archivos actuales</label>
              {editForm.existing_attachments.length === 0 ? (
                <p className="muted">Sin archivos adjuntos.</p>
              ) : (
                <div className="actions-row">
                  {editForm.existing_attachments.map((attachment) => (
                    <button
                      key={attachment.id}
                      type="button"
                      className="button-secondary"
                      onClick={() =>
                        setEditForm((previous): TaskFormState => ({
                          ...previous,
                          existing_attachments: previous.existing_attachments.filter(
                            (item) => item.id !== attachment.id,
                          ),
                        }))
                      }
                    >
                      Quitar {attachment.name}
                    </button>
                  ))}
                </div>
              )}
            </div>

            <div className="form-field">
              <label>
                Título
                <input
                  value={editForm.title}
                  onChange={(event) =>
                    setEditForm((previous): TaskFormState => ({ ...previous, title: event.target.value }))
                  }
                  required
                />
              </label>
            </div>

            <div className="form-field">
              <label>
                Descripción
                <input
                  value={editForm.description}
                  onChange={(event) =>
                    setEditForm((previous): TaskFormState => ({
                      ...previous,
                      description: event.target.value,
                    }))
                  }
                  required
                />
              </label>
            </div>

            <div className="form-field">
              <label>
                Estado
                <select
                  value={editForm.status}
                  onChange={(event) =>
                    setEditForm((previous): TaskFormState => ({
                      ...previous,
                      status: event.target.value as TaskStatus,
                    }))
                  }
                >
                  {statusOptions.map((status) => (
                    <option key={status} value={status}>
                      {status}
                    </option>
                  ))}
                </select>
              </label>
            </div>

            <div className="form-field">
              <label>
                Horas dedicadas
                <input
                  type="number"
                  min={1}
                  value={editForm.spent_hours}
                  onChange={(event) =>
                    setEditForm((previous): TaskFormState => ({
                      ...previous,
                      spent_hours: event.target.value,
                    }))
                  }
                  required
                />
              </label>
            </div>

            <div className="form-field">
              <label>
                Observaciones
                <input
                  value={editForm.observations}
                  onChange={(event) =>
                    setEditForm((previous): TaskFormState => ({
                      ...previous,
                      observations: event.target.value,
                    }))
                  }
                />
              </label>
            </div>

            <div className="form-field">
              <label>
                Fecha de inicio de semana
                <input
                  type="date"
                  value={editForm.week_start_date}
                  onChange={(event) =>
                    setEditForm((previous): TaskFormState => ({
                      ...previous,
                      week_start_date: event.target.value,
                    }))
                  }
                  required
                />
              </label>
            </div>

            <div className="form-actions">
              <button type="submit">Guardar cambios</button>
              <button type="button" className="button-secondary" onClick={() => setEditTaskId(null)}>
                Cancelar edición
              </button>
            </div>
          </form>
        </section>
      ) : null}
    </Layout>
  );
}
