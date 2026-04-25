import { FormEvent, useEffect, useMemo, useState } from "react";
import {
  createTask,
  deleteTask,
  getMe,
  listAssignmentsByUser,
  listTasks,
  toErrorMessage,
  updateTask,
} from "../api/client";
import type { Assignment, Task, TaskStatus, UpdateTaskPayload, User } from "../api/types";
import ErrorMessage from "../components/ErrorMessage";
import Layout from "../components/Layout";
import Loading from "../components/Loading";

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
}

const defaultForm: TaskFormState = {
  assignment_id: "",
  title: "",
  description: "",
  status: "abierto",
  spent_hours: "1",
  observations: "",
  week_start_date: "",
};

const statusOptions: TaskStatus[] = ["abierto", "en desarrollo", "finalizado"];

export default function WorkerDashboard({ user, onLogout }: WorkerDashboardProps) {
  const [me, setMe] = useState<User | null>(null);
  const [assignments, setAssignments] = useState<Assignment[]>([]);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [createForm, setCreateForm] = useState<TaskFormState>(defaultForm);
  const [editTaskId, setEditTaskId] = useState<number | null>(null);
  const [editForm, setEditForm] = useState<TaskFormState>(defaultForm);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const assignmentIds = useMemo(() => assignments.map((item) => item.id), [assignments]);

  const loadData = async () => {
    setLoading(true);
    setError(null);

    try {
      const meResult = await getMe();
      const [assignmentResult, taskResult] = await Promise.all([
        listAssignmentsByUser(meResult.id),
        listTasks(),
      ]);

      setMe(meResult);
      setAssignments(assignmentResult.assignments);
      setTasks(taskResult.tasks);
      setCreateForm((previous) => ({
        ...previous,
        assignment_id:
          previous.assignment_id || String(assignmentResult.assignments[0]?.id ?? ""),
      }));
    } catch (err) {
      setError(toErrorMessage(err));
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
  });

  const handleCreateTask = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    setSuccess(null);

    try {
      await createTask(toPayload(createForm));
      setSuccess("Task creada correctamente.");
      setCreateForm({
        ...defaultForm,
        assignment_id: String(assignments[0]?.id ?? ""),
      });
      const taskResult = await listTasks();
      setTasks(taskResult.tasks);
    } catch (err) {
      setError(toErrorMessage(err));
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
    });
  };

  const handleUpdateTask = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (!editTaskId) {
      return;
    }

    setError(null);
    setSuccess(null);

    try {
      await updateTask(editTaskId, toPayload(editForm));
      setSuccess("Task actualizada correctamente.");
      setEditTaskId(null);
      const taskResult = await listTasks();
      setTasks(taskResult.tasks);
    } catch (err) {
      setError(toErrorMessage(err));
    }
  };

  const handleDeleteTask = async (id: number) => {
    setError(null);
    setSuccess(null);

    try {
      await deleteTask(id);
      setSuccess("Task eliminada correctamente.");
      const taskResult = await listTasks();
      setTasks(taskResult.tasks);
    } catch (err) {
      setError(toErrorMessage(err));
    }
  };

  return (
    <Layout title="Worker Dashboard (monitor/assistant)" user={user} onLogout={onLogout}>
      <div className="actions-row">
        <button onClick={() => void loadData()} disabled={loading}>
          Recargar datos
        </button>
      </div>

      {loading && <Loading label="Cargando dashboard..." />}
      <ErrorMessage message={error} />
      {success && <div className="success-box">{success}</div>}

      <section className="card">
        <h2>Errores backend esperados en tasks</h2>
        <ul>
          <li>task week is not active</li>
          <li>insufficient permissions</li>
          <li>task status is invalid</li>
        </ul>
      </section>

      <section className="card">
        <h2>GET /auth/me</h2>
        {me ? (
          <table>
            <tbody>
              <tr>
                <th>ID</th>
                <td>{me.id}</td>
              </tr>
              <tr>
                <th>Name</th>
                <td>{me.name}</td>
              </tr>
              <tr>
                <th>Role</th>
                <td>{me.global_role}</td>
              </tr>
            </tbody>
          </table>
        ) : (
          <p className="muted">Sin datos</p>
        )}
      </section>

      <section className="card">
        <h2>GET /assignments?user_id=&lt;currentUser.id&gt;</h2>
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Workspace ID</th>
              <th>Role</th>
              <th>Weekly Hours</th>
            </tr>
          </thead>
          <tbody>
            {assignments.map((item) => (
              <tr key={item.id}>
                <td>{item.id}</td>
                <td>{item.workspace_id}</td>
                <td>{item.role}</td>
                <td>{item.weekly_hours}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>

      <section className="card">
        <h2>POST /tasks</h2>
        <form className="grid-form" onSubmit={handleCreateTask}>
          <label>
            assignment_id
            <select
              value={createForm.assignment_id}
              onChange={(event) =>
                setCreateForm((previous) => ({
                  ...previous,
                  assignment_id: event.target.value,
                }))
              }
              required
            >
              <option value="" disabled>
                Seleccione assignment
              </option>
              {assignmentIds.map((id) => (
                <option key={id} value={id}>
                  {id}
                </option>
              ))}
            </select>
          </label>

          <label>
            title
            <input
              value={createForm.title}
              onChange={(event) =>
                setCreateForm((previous) => ({ ...previous, title: event.target.value }))
              }
              required
            />
          </label>

          <label>
            description
            <input
              value={createForm.description}
              onChange={(event) =>
                setCreateForm((previous) => ({
                  ...previous,
                  description: event.target.value,
                }))
              }
              required
            />
          </label>

          <label>
            status
            <select
              value={createForm.status}
              onChange={(event) =>
                setCreateForm((previous) => ({
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

          <label>
            spent_hours
            <input
              type="number"
              min={1}
              value={createForm.spent_hours}
              onChange={(event) =>
                setCreateForm((previous) => ({
                  ...previous,
                  spent_hours: event.target.value,
                }))
              }
              required
            />
          </label>

          <label>
            observations
            <input
              value={createForm.observations}
              onChange={(event) =>
                setCreateForm((previous) => ({
                  ...previous,
                  observations: event.target.value,
                }))
              }
            />
          </label>

          <label>
            week_start_date
            <input
              type="date"
              value={createForm.week_start_date}
              onChange={(event) =>
                setCreateForm((previous) => ({
                  ...previous,
                  week_start_date: event.target.value,
                }))
              }
              required
            />
          </label>

          <button type="submit">Crear task</button>
        </form>
      </section>

      <section className="card">
        <h2>GET /tasks</h2>
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Assignment</th>
              <th>Title</th>
              <th>Status</th>
              <th>Hours</th>
              <th>Week Start</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            {tasks.map((item) => (
              <tr key={item.id}>
                <td>{item.id}</td>
                <td>{item.assignment_id}</td>
                <td>{item.title}</td>
                <td>{item.status}</td>
                <td>{item.spent_hours}</td>
                <td>{item.week_start_date}</td>
                <td>
                  <button type="button" onClick={() => startEdit(item)}>
                    Editar
                  </button>
                  <button type="button" onClick={() => void handleDeleteTask(item.id)}>
                    Eliminar
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>

      {editTaskId ? (
        <section className="card">
          <h2>PUT /tasks/{editTaskId}</h2>
          <form className="grid-form" onSubmit={handleUpdateTask}>
            <label>
              assignment_id
              <select
                value={editForm.assignment_id}
                onChange={(event) =>
                  setEditForm((previous) => ({
                    ...previous,
                    assignment_id: event.target.value,
                  }))
                }
                required
              >
                <option value="" disabled>
                  Seleccione assignment
                </option>
                {assignmentIds.map((id) => (
                  <option key={id} value={id}>
                    {id}
                  </option>
                ))}
              </select>
            </label>

            <label>
              title
              <input
                value={editForm.title}
                onChange={(event) =>
                  setEditForm((previous) => ({ ...previous, title: event.target.value }))
                }
                required
              />
            </label>

            <label>
              description
              <input
                value={editForm.description}
                onChange={(event) =>
                  setEditForm((previous) => ({
                    ...previous,
                    description: event.target.value,
                  }))
                }
                required
              />
            </label>

            <label>
              status
              <select
                value={editForm.status}
                onChange={(event) =>
                  setEditForm((previous) => ({
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

            <label>
              spent_hours
              <input
                type="number"
                min={1}
                value={editForm.spent_hours}
                onChange={(event) =>
                  setEditForm((previous) => ({
                    ...previous,
                    spent_hours: event.target.value,
                  }))
                }
                required
              />
            </label>

            <label>
              observations
              <input
                value={editForm.observations}
                onChange={(event) =>
                  setEditForm((previous) => ({
                    ...previous,
                    observations: event.target.value,
                  }))
                }
              />
            </label>

            <label>
              week_start_date
              <input
                type="date"
                value={editForm.week_start_date}
                onChange={(event) =>
                  setEditForm((previous) => ({
                    ...previous,
                    week_start_date: event.target.value,
                  }))
                }
                required
              />
            </label>

            <div className="actions-row">
              <button type="submit">Guardar cambios</button>
              <button type="button" onClick={() => setEditTaskId(null)}>
                Cancelar
              </button>
            </div>
          </form>
        </section>
      ) : null}
    </Layout>
  );
}
