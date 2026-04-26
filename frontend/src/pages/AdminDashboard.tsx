import { FormEvent, useEffect, useMemo, useState } from "react";
import {
  createAssignment,
  createPeriod,
  createUser,
  createWorkspace,
  getMe,
  listPeriods,
  listReports,
  listTasks,
  listUsers,
  listWorkspaces,
  toErrorMessage,
} from "../api/client";
import type { GlobalRole, Period, Report, Task, User, Workspace } from "../api/types";
import ErrorMessage from "../components/ErrorMessage";
import Layout from "../components/Layout";
import Loading from "../components/Loading";

interface AdminDashboardProps {
  user: User;
  onLogout: () => void;
}

export default function AdminDashboard({ user, onLogout }: AdminDashboardProps) {
  const [me, setMe] = useState<User | null>(null);
  const [users, setUsers] = useState<User[]>([]);
  const [periods, setPeriods] = useState<Period[]>([]);
  const [workspaces, setWorkspaces] = useState<Workspace[]>([]);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [reports, setReports] = useState<Report[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const [userForm, setUserForm] = useState({
    name: "",
    email: "",
    password: "",
    global_role: "professor" as GlobalRole,
  });
  const [periodForm, setPeriodForm] = useState({
    name: "",
    initial_date: "",
    weeks_count: "16",
    period_state: "active" as "active" | "closed",
  });
  const [workspaceForm, setWorkspaceForm] = useState({
    period_id: "",
    user_id: "",
    name: "",
    type: "project" as "course" | "project",
    initial_date: "",
    final_date: "",
    observations: "",
    state: "active" as "active" | "closed",
  });
  const [assignmentForm, setAssignmentForm] = useState({
    user_id: "",
    workspace_id: "",
    role: "monitor" as "monitor" | "assistant",
    weekly_hours: "1",
  });

  const professorUsers = useMemo(
    () => users.filter((item) => item.global_role === "professor"),
    [users],
  );
  const workerUsers = useMemo(() => {
    return users.filter((item) => item.global_role === assignmentForm.role);
  }, [users, assignmentForm.role]);

  const loadAll = async () => {
    setLoading(true);
    setError(null);

    try {
      const [meResult, usersResult, periodsResult, workspacesResult, tasksResult, reportsResult] =
        await Promise.all([
          getMe(),
          listUsers(),
          listPeriods(),
          listWorkspaces(),
          listTasks(),
          listReports(),
        ]);

      setMe(meResult);
      setUsers(usersResult.users);
      setPeriods(periodsResult.periods);
      setWorkspaces(workspacesResult.workspaces);
      setTasks(tasksResult.tasks);
      setReports(reportsResult.reports);

      setWorkspaceForm((previous) => ({
        ...previous,
        period_id: previous.period_id || String(periodsResult.periods[0]?.id ?? ""),
        user_id: previous.user_id || String(professorUsers[0]?.id ?? usersResult.users.find((u) => u.global_role === "professor")?.id ?? ""),
      }));
      setAssignmentForm((previous) => ({
        ...previous,
        workspace_id: previous.workspace_id || String(workspacesResult.workspaces[0]?.id ?? ""),
      }));
    } catch (err) {
      setError(toErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadAll();
  }, []);

  useEffect(() => {
    setWorkspaceForm((previous) => ({
      ...previous,
      user_id:
        previous.user_id || String(professorUsers[0]?.id ?? ""),
    }));
  }, [professorUsers]);

  useEffect(() => {
    setAssignmentForm((previous) => ({
      ...previous,
      user_id: String(workerUsers[0]?.id ?? ""),
    }));
  }, [workerUsers]);

  const handleCreateUser = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    setSuccess(null);

    try {
      await createUser(userForm);
      setSuccess("Usuario creado correctamente.");
      setUserForm({
        name: "",
        email: "",
        password: "",
        global_role: "professor",
      });
      await loadAll();
    } catch (err) {
      setError(toErrorMessage(err));
    }
  };

  const handleCreatePeriod = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    setSuccess(null);

    try {
      await createPeriod({
        name: periodForm.name,
        initial_date: periodForm.initial_date,
        weeks_count: Number(periodForm.weeks_count) as 8 | 16,
        period_state: periodForm.period_state,
      });
      setSuccess("Periodo creado correctamente.");
      setPeriodForm({
        name: "",
        initial_date: "",
        weeks_count: "16",
        period_state: "active",
      });
      await loadAll();
    } catch (err) {
      setError(toErrorMessage(err));
    }
  };

  const handleCreateWorkspace = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    setSuccess(null);

    try {
      await createWorkspace({
        period_id: Number(workspaceForm.period_id),
        user_id: Number(workspaceForm.user_id),
        name: workspaceForm.name,
        type: workspaceForm.type,
        initial_date: workspaceForm.initial_date,
        final_date: workspaceForm.final_date,
        observations: workspaceForm.observations,
        state: workspaceForm.state,
      });
      setSuccess("Workspace creado correctamente.");
      setWorkspaceForm((previous) => ({
        ...previous,
        name: "",
        initial_date: "",
        final_date: "",
        observations: "",
      }));
      await loadAll();
    } catch (err) {
      setError(toErrorMessage(err));
    }
  };

  const handleCreateAssignment = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    setSuccess(null);

    try {
      await createAssignment({
        user_id: Number(assignmentForm.user_id),
        workspace_id: Number(assignmentForm.workspace_id),
        role: assignmentForm.role,
        weekly_hours: Number(assignmentForm.weekly_hours),
      });
      setSuccess("Assignment creado correctamente.");
      setAssignmentForm((previous) => ({
        ...previous,
        weekly_hours: "1",
      }));
      await loadAll();
    } catch (err) {
      setError(toErrorMessage(err));
    }
  };

  return (
    <Layout title="Admin Dashboard" user={user} onLogout={onLogout}>
      <div className="actions-row">
        <button onClick={() => void loadAll()} disabled={loading}>
          Recargar datos
        </button>
      </div>

      {loading && <Loading label="Cargando dashboard..." />}
      <ErrorMessage message={error} />
      {success && <div className="success-box">{success}</div>}

      <section className="card">
        <h2>POST /users</h2>
        <form className="grid-form" onSubmit={handleCreateUser}>
          <label>
            name
            <input
              value={userForm.name}
              onChange={(event) =>
                setUserForm((previous) => ({ ...previous, name: event.target.value }))
              }
              required
            />
          </label>

          <label>
            email
            <input
              type="email"
              value={userForm.email}
              onChange={(event) =>
                setUserForm((previous) => ({ ...previous, email: event.target.value }))
              }
              required
            />
          </label>

          <label>
            password
            <input
              type="password"
              value={userForm.password}
              onChange={(event) =>
                setUserForm((previous) => ({ ...previous, password: event.target.value }))
              }
              required
              minLength={8}
            />
          </label>

          <label>
            global_role
            <select
              value={userForm.global_role}
              onChange={(event) =>
                setUserForm((previous) => ({
                  ...previous,
                  global_role: event.target.value as GlobalRole,
                }))
              }
            >
              <option value="admin">admin</option>
              <option value="professor">professor</option>
              <option value="monitor">monitor</option>
              <option value="assistant">assistant</option>
            </select>
          </label>

          <button type="submit">Crear usuario</button>
        </form>
      </section>

      <section className="card">
        <h2>POST /periods</h2>
        <form className="grid-form" onSubmit={handleCreatePeriod}>
          <label>
            name
            <input
              placeholder="2026-20"
              value={periodForm.name}
              onChange={(event) =>
                setPeriodForm((previous) => ({ ...previous, name: event.target.value }))
              }
              required
            />
          </label>

          <label>
            initial_date
            <input
              type="date"
              value={periodForm.initial_date}
              onChange={(event) =>
                setPeriodForm((previous) => ({ ...previous, initial_date: event.target.value }))
              }
              required
            />
          </label>

          <label>
            weeks_count
            <select
              value={periodForm.weeks_count}
              onChange={(event) =>
                setPeriodForm((previous) => ({ ...previous, weeks_count: event.target.value }))
              }
            >
              <option value="8">8</option>
              <option value="16">16</option>
            </select>
          </label>

          <label>
            period_state
            <select
              value={periodForm.period_state}
              onChange={(event) =>
                setPeriodForm((previous) => ({
                  ...previous,
                  period_state: event.target.value as "active" | "closed",
                }))
              }
            >
              <option value="active">active</option>
              <option value="closed">closed</option>
            </select>
          </label>

          <button type="submit">Crear periodo</button>
        </form>
      </section>

      <section className="card">
        <h2>POST /workspaces</h2>
        <form className="grid-form" onSubmit={handleCreateWorkspace}>
          <label>
            period_id
            <select
              value={workspaceForm.period_id}
              onChange={(event) =>
                setWorkspaceForm((previous) => ({ ...previous, period_id: event.target.value }))
              }
              required
            >
              <option value="" disabled>
                Seleccione periodo
              </option>
              {periods.map((item) => (
                <option key={item.id} value={item.id}>
                  {item.id} - {item.name}
                </option>
              ))}
            </select>
          </label>

          <label>
            user_id (professor)
            <select
              value={workspaceForm.user_id}
              onChange={(event) =>
                setWorkspaceForm((previous) => ({ ...previous, user_id: event.target.value }))
              }
              required
            >
              <option value="" disabled>
                Seleccione professor
              </option>
              {professorUsers.map((item) => (
                <option key={item.id} value={item.id}>
                  {item.id} - {item.name}
                </option>
              ))}
            </select>
          </label>

          <label>
            name
            <input
              value={workspaceForm.name}
              onChange={(event) =>
                setWorkspaceForm((previous) => ({ ...previous, name: event.target.value }))
              }
              required
            />
          </label>

          <label>
            type
            <select
              value={workspaceForm.type}
              onChange={(event) =>
                setWorkspaceForm((previous) => ({
                  ...previous,
                  type: event.target.value as "course" | "project",
                }))
              }
            >
              <option value="course">course</option>
              <option value="project">project</option>
            </select>
          </label>

          <label>
            initial_date
            <input
              type="date"
              value={workspaceForm.initial_date}
              onChange={(event) =>
                setWorkspaceForm((previous) => ({
                  ...previous,
                  initial_date: event.target.value,
                }))
              }
              required
            />
          </label>

          <label>
            final_date
            <input
              type="date"
              value={workspaceForm.final_date}
              onChange={(event) =>
                setWorkspaceForm((previous) => ({ ...previous, final_date: event.target.value }))
              }
              required
            />
          </label>

          <label>
            observations
            <input
              value={workspaceForm.observations}
              onChange={(event) =>
                setWorkspaceForm((previous) => ({
                  ...previous,
                  observations: event.target.value,
                }))
              }
              required
            />
          </label>

          <label>
            state
            <select
              value={workspaceForm.state}
              onChange={(event) =>
                setWorkspaceForm((previous) => ({
                  ...previous,
                  state: event.target.value as "active" | "closed",
                }))
              }
            >
              <option value="active">active</option>
              <option value="closed">closed</option>
            </select>
          </label>

          <button type="submit">Crear workspace</button>
        </form>
      </section>

      <section className="card">
        <h2>POST /assignments</h2>
        <form className="grid-form" onSubmit={handleCreateAssignment}>
          <label>
            role
            <select
              value={assignmentForm.role}
              onChange={(event) =>
                setAssignmentForm((previous) => ({
                  ...previous,
                  role: event.target.value as "monitor" | "assistant",
                }))
              }
            >
              <option value="monitor">monitor</option>
              <option value="assistant">assistant</option>
            </select>
          </label>

          <label>
            user_id
            <select
              value={assignmentForm.user_id}
              onChange={(event) =>
                setAssignmentForm((previous) => ({ ...previous, user_id: event.target.value }))
              }
              required
            >
              <option value="" disabled>
                Seleccione usuario
              </option>
              {workerUsers.map((item) => (
                <option key={item.id} value={item.id}>
                  {item.id} - {item.name}
                </option>
              ))}
            </select>
          </label>

          <label>
            workspace_id
            <select
              value={assignmentForm.workspace_id}
              onChange={(event) =>
                setAssignmentForm((previous) => ({
                  ...previous,
                  workspace_id: event.target.value,
                }))
              }
              required
            >
              <option value="" disabled>
                Seleccione workspace
              </option>
              {workspaces.map((item) => (
                <option key={item.id} value={item.id}>
                  {item.id} - {item.name}
                </option>
              ))}
            </select>
          </label>

          <label>
            weekly_hours
            <input
              type="number"
              min={1}
              value={assignmentForm.weekly_hours}
              onChange={(event) =>
                setAssignmentForm((previous) => ({
                  ...previous,
                  weekly_hours: event.target.value,
                }))
              }
              required
            />
          </label>

          <button type="submit">Crear assignment</button>
        </form>

        <p className="muted">
          Errores de negocio esperados: assignment already exists, period inscription is closed,
          monitor weekly hours cannot exceed 12, assistant weekly hours cannot exceed 22.
        </p>
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
                <th>Email</th>
                <td>{me.email}</td>
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
        <h2>GET /users</h2>
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Name</th>
              <th>Email</th>
              <th>Role</th>
            </tr>
          </thead>
          <tbody>
            {users.map((item) => (
              <tr key={item.id}>
                <td>{item.id}</td>
                <td>{item.name}</td>
                <td>{item.email}</td>
                <td>{item.global_role}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>

      <section className="card">
        <h2>GET /periods</h2>
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Name</th>
              <th>Initial Date</th>
              <th>Final Date</th>
              <th>State</th>
            </tr>
          </thead>
          <tbody>
            {periods.map((item) => (
              <tr key={item.id}>
                <td>{item.id}</td>
                <td>{item.name}</td>
                <td>{item.initial_date}</td>
                <td>{item.final_date}</td>
                <td>{item.period_state}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>

      <section className="card">
        <h2>GET /workspaces</h2>
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Name</th>
              <th>User ID</th>
              <th>Period ID</th>
              <th>Type</th>
              <th>State</th>
            </tr>
          </thead>
          <tbody>
            {workspaces.map((item) => (
              <tr key={item.id}>
                <td>{item.id}</td>
                <td>{item.name}</td>
                <td>{item.user_id}</td>
                <td>{item.period_id}</td>
                <td>{item.type}</td>
                <td>{item.state}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>

      <section className="card">
        <h2>GET /tasks</h2>
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>User ID</th>
              <th>Assignment ID</th>
              <th>Status</th>
              <th>Hours</th>
              <th>Week Start</th>
            </tr>
          </thead>
          <tbody>
            {tasks.map((item) => (
              <tr key={item.id}>
                <td>{item.id}</td>
                <td>{item.user_id}</td>
                <td>{item.assignment_id}</td>
                <td>{item.status}</td>
                <td>{item.spent_hours}</td>
                <td>{item.week_start_date}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>

      <section className="card">
        <h2>GET /reports</h2>
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Workspace ID</th>
              <th>Week ID</th>
              <th>Assignment ID</th>
              <th>User ID</th>
            </tr>
          </thead>
          <tbody>
            {reports.map((item) => (
              <tr key={item.id}>
                <td>{item.id}</td>
                <td>{item.workspace_id}</td>
                <td>{item.week_id}</td>
                <td>{item.assignment_id}</td>
                <td>{item.user_id}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>
    </Layout>
  );
}
