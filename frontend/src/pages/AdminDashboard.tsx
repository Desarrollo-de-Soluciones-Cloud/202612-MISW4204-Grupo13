import { useEffect, useState } from "react";
import {
  getMe,
  listPeriods,
  listReports,
  listTasks,
  listUsers,
  listWorkspaces,
  toErrorMessage,
} from "../api/client";
import type { Period, Report, Task, User, Workspace } from "../api/types";
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
    } catch (err) {
      setError(toErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadAll();
  }, []);

  return (
    <Layout title="Admin Dashboard" user={user} onLogout={onLogout}>
      <div className="actions-row">
        <button onClick={() => void loadAll()} disabled={loading}>
          Recargar datos
        </button>
      </div>

      {loading && <Loading label="Cargando dashboard..." />}
      <ErrorMessage message={error} />

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
