import { FormEvent, useEffect, useState } from "react";
import {
  downloadReport,
  generateWeeklyReport,
  getMe,
  listReports,
  listTasks,
  listUsers,
  listWorkspaces,
  toErrorMessage,
} from "../api/client";
import type { Report, Task, User, Workspace } from "../api/types";
import ErrorMessage from "../components/ErrorMessage";
import Layout from "../components/Layout";
import Loading from "../components/Loading";

interface ProfessorDashboardProps {
  user: User;
  onLogout: () => void;
}

export default function ProfessorDashboard({ user, onLogout }: ProfessorDashboardProps) {
  const [me, setMe] = useState<User | null>(null);
  const [workspaces, setWorkspaces] = useState<Workspace[]>([]);
  const [monitors, setMonitors] = useState<User[]>([]);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [reports, setReports] = useState<Report[]>([]);
  const [workspaceId, setWorkspaceId] = useState("");
  const [weekId, setWeekId] = useState("");
  const [filterWorkspaceId, setFilterWorkspaceId] = useState("");
  const [filterWeekId, setFilterWeekId] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const loadBase = async () => {
    setLoading(true);
    setError(null);

    try {
      const [meResult, workspaceResult, monitorResult, taskResult, reportResult] = await Promise.all([
        getMe(),
        listWorkspaces(),
        listUsers("monitor"),
        listTasks(),
        listReports(),
      ]);

      setMe(meResult);
      setWorkspaces(workspaceResult.workspaces);
      setMonitors(monitorResult.users);
      setTasks(taskResult.tasks);
      setReports(reportResult.reports);
    } catch (err) {
      setError(toErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadBase();
  }, []);

  const handleGenerateReport = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    setSuccess(null);

    try {
      const response = await generateWeeklyReport({
        workspace_id: Number(workspaceId),
        week_id: Number(weekId),
      });
      setSuccess(`Reportes generados: ${response.generated_count}`);
      const reportResult = await listReports();
      setReports(reportResult.reports);
    } catch (err) {
      setError(toErrorMessage(err));
    }
  };

  const handleFilterReports = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    setSuccess(null);

    try {
      const response = await listReports({
        workspace_id: filterWorkspaceId ? Number(filterWorkspaceId) : undefined,
        week_id: filterWeekId ? Number(filterWeekId) : undefined,
      });
      setReports(response.reports);
    } catch (err) {
      setError(toErrorMessage(err));
    }
  };

  const handleDownload = async (reportId: number) => {
    setError(null);
    setSuccess(null);

    try {
      const blob = await downloadReport(reportId);
      const fileUrl = window.URL.createObjectURL(blob);
      const anchor = document.createElement("a");
      anchor.href = fileUrl;
      anchor.download = `report_${reportId}.pdf`;
      document.body.appendChild(anchor);
      anchor.click();
      anchor.remove();
      window.URL.revokeObjectURL(fileUrl);
      setSuccess(`Reporte ${reportId} descargado.`);
    } catch (err) {
      setError(toErrorMessage(err));
    }
  };

  return (
    <Layout title="Professor Dashboard" user={user} onLogout={onLogout}>
      <div className="actions-row">
        <button onClick={() => void loadBase()} disabled={loading}>
          Recargar datos
        </button>
      </div>

      {loading && <Loading label="Cargando dashboard..." />}
      <ErrorMessage message={error} />
      {success && <div className="success-box">{success}</div>}

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
        <h2>GET /workspaces</h2>
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Name</th>
              <th>Period ID</th>
              <th>State</th>
            </tr>
          </thead>
          <tbody>
            {workspaces.map((item) => (
              <tr key={item.id}>
                <td>{item.id}</td>
                <td>{item.name}</td>
                <td>{item.period_id}</td>
                <td>{item.state}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>

      <section className="card">
        <h2>GET /users?role=monitor</h2>
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Name</th>
              <th>Email</th>
            </tr>
          </thead>
          <tbody>
            {monitors.map((item) => (
              <tr key={item.id}>
                <td>{item.id}</td>
                <td>{item.name}</td>
                <td>{item.email}</td>
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
        <h2>POST /reports/weekly</h2>
        <form className="grid-form" onSubmit={handleGenerateReport}>
          <label>
            workspace_id
            <input
              type="number"
              value={workspaceId}
              onChange={(event) => setWorkspaceId(event.target.value)}
              required
              min={1}
            />
          </label>

          <label>
            week_id
            <input
              type="number"
              value={weekId}
              onChange={(event) => setWeekId(event.target.value)}
              required
              min={1}
            />
          </label>

          <button type="submit">Generar reporte semanal</button>
        </form>
      </section>

      <section className="card">
        <h2>GET /reports (filtros opcionales)</h2>
        <form className="grid-form" onSubmit={handleFilterReports}>
          <label>
            workspace_id (opcional)
            <input
              type="number"
              value={filterWorkspaceId}
              onChange={(event) => setFilterWorkspaceId(event.target.value)}
              min={1}
            />
          </label>

          <label>
            week_id (opcional)
            <input
              type="number"
              value={filterWeekId}
              onChange={(event) => setFilterWeekId(event.target.value)}
              min={1}
            />
          </label>

          <button type="submit">Aplicar filtros</button>
        </form>

        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Workspace ID</th>
              <th>Week ID</th>
              <th>Assignment ID</th>
              <th>User ID</th>
              <th>PDF</th>
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
                <td>
                  <button type="button" onClick={() => void handleDownload(item.id)}>
                    Descargar
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>
    </Layout>
  );
}
