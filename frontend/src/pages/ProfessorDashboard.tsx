import { FormEvent, useEffect, useState } from "react";
import {
  downloadReport,
  generateWeeklyReport,
  getMe,
  listReports,
  listTasks,
  listWorkspaces,
  toErrorMessage,
} from "../api/client";
import type { Report, Task, User, Workspace } from "../api/types";
import EmptyState from "../components/EmptyState";
import HelpText from "../components/HelpText";
import Layout from "../components/Layout";
import Loading from "../components/Loading";
import Toast from "../components/Toast";
import useToast from "../components/useToast";

interface ProfessorDashboardProps {
  user: User;
  onLogout: () => void;
}

export default function ProfessorDashboard({ user, onLogout }: ProfessorDashboardProps) {
  const [me, setMe] = useState<User | null>(null);
  const [workspaces, setWorkspaces] = useState<Workspace[]>([]);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [reports, setReports] = useState<Report[]>([]);
  const [workspaceId, setWorkspaceId] = useState("");
  const [weekId, setWeekId] = useState("");
  const [filterWorkspaceId, setFilterWorkspaceId] = useState("");
  const [filterWeekId, setFilterWeekId] = useState("");
  const [loading, setLoading] = useState(false);
  const { toast, showToast, clearToast } = useToast();

  const loadBase = async () => {
    setLoading(true);

    try {
      const [meResult, workspaceResult, taskResult, reportResult] = await Promise.all([
        getMe(),
        listWorkspaces(),
        listTasks(),
        listReports(),
      ]);

      setMe(meResult);
      setWorkspaces(workspaceResult.workspaces);
      setTasks(taskResult.tasks);
      setReports(reportResult.reports);
      setWorkspaceId((previous) => previous || String(workspaceResult.workspaces[0]?.id ?? ""));
    } catch (err) {
      showToast(toErrorMessage(err), "error");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadBase();
  }, []);

  const handleGenerateReport = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    clearToast();

    try {
      const response = await generateWeeklyReport({
        workspace_id: Number(workspaceId),
        week_id: Number(weekId),
      });
      showToast(`Se generaron ${response.generated_count} reportes semanales.`, "success");
      setWeekId("");
      const reportResult = await listReports();
      setReports(reportResult.reports);
    } catch (err) {
      showToast(toErrorMessage(err), "error");
    }
  };

  const handleFilterReports = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    clearToast();

    try {
      const response = await listReports({
        workspace_id: filterWorkspaceId ? Number(filterWorkspaceId) : undefined,
        week_id: filterWeekId ? Number(filterWeekId) : undefined,
      });
      setReports(response.reports);
    } catch (err) {
      showToast(toErrorMessage(err), "error");
    }
  };

  const handleDownload = async (reportId: number) => {
    clearToast();

    try {
      const blob = await downloadReport(reportId);
      const fileUrl = window.URL.createObjectURL(blob);
      const anchor = document.createElement("a");
      anchor.href = fileUrl;
      anchor.download = `reporte_${reportId}.pdf`;
      document.body.appendChild(anchor);
      anchor.click();
      anchor.remove();
      window.URL.revokeObjectURL(fileUrl);
      showToast(`Reporte ${reportId} descargado correctamente.`, "success");
    } catch (err) {
      showToast(toErrorMessage(err), "error");
    }
  };

  return (
    <Layout
      title="Panel del profesor"
      description="Consulta tus cursos/proyectos, revisa tareas reportadas y genera reportes semanales."
      user={user}
      onLogout={onLogout}
    >
      <div className="actions-row">
        <button onClick={() => void loadBase()} disabled={loading}>
          Recargar datos
        </button>
      </div>

      {loading && <Loading label="Actualizando información..." />}
      {toast ? <Toast type={toast.type} message={toast.message} onClose={clearToast} /> : null}

      <section className="card info-card">
        <h2>Alcance MVP actual</h2>
        <p>
          Actualmente el sistema genera reportes PDF semanales básicos. La integración con IA
          externa, los adjuntos de tareas y Cloud Storage quedan registrados como deuda técnica
          para la siguiente fase.
        </p>
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
                <th>Correo</th>
                <td>{me.email}</td>
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
        <h2>Mis cursos y proyectos</h2>
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Nombre</th>
                <th>ID del periodo</th>
                <th>Estado</th>
              </tr>
            </thead>
            <tbody>
              {workspaces.length === 0 ? (
                <EmptyState colSpan={4} />
              ) : (
                workspaces.map((item) => (
                  <tr key={item.id}>
                    <td>{item.id}</td>
                    <td>{item.name}</td>
                    <td>{item.period_id}</td>
                    <td>{item.state}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </section>

      <section className="card">
        <h2>Tareas reportadas por monitores y asistentes</h2>
        <p className="muted">
          Estas son las tareas registradas por los monitores y asistentes vinculados a tus
          cursos/proyectos.
        </p>
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>ID de vinculación</th>
                <th>Estado</th>
                <th>Horas dedicadas</th>
                <th>Fecha de inicio de semana</th>
              </tr>
            </thead>
            <tbody>
              {tasks.length === 0 ? (
                <EmptyState colSpan={5} />
              ) : (
                tasks.map((item) => (
                  <tr key={item.id}>
                    <td>{item.id}</td>
                    <td>{item.assignment_id}</td>
                    <td>{item.status}</td>
                    <td>{item.spent_hours}</td>
                    <td>{item.week_start_date}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </section>

      <section className="card">
        <h2>Generar reporte semanal</h2>
        <p className="section-desc">Genera el reporte PDF de actividad para la semana académica seleccionada.</p>
        <form className="form-grid" onSubmit={handleGenerateReport}>
          <div className="form-field">
            <label>
              ID del curso/proyecto
              <input
                type="number"
                value={workspaceId}
                onChange={(event) => setWorkspaceId(event.target.value)}
                required
                min={1}
              />
            </label>
            <HelpText>Selecciona un curso o proyecto propio.</HelpText>
          </div>

          <div className="form-field">
            <label>
              ID de semana
              <input
                type="number"
                value={weekId}
                onChange={(event) => setWeekId(event.target.value)}
                required
                min={1}
              />
            </label>
            <HelpText>
              El reporte se genera para la semana académica seleccionada.
            </HelpText>
          </div>

          <div className="form-actions">
            <button type="submit">Generar reporte PDF</button>
          </div>
        </form>
      </section>

      <section className="card">
        <h2>Reportes generados</h2>
        <p className="section-desc">Filtra por curso/proyecto o semana para consultar los reportes disponibles.</p>
        <form className="form-grid" onSubmit={handleFilterReports}>
          <div className="form-field">
            <label>
              ID del curso/proyecto (opcional)
              <input
                type="number"
                value={filterWorkspaceId}
                onChange={(event) => setFilterWorkspaceId(event.target.value)}
                min={1}
              />
            </label>
          </div>

          <div className="form-field">
            <label>
              ID de semana (opcional)
              <input
                type="number"
                value={filterWeekId}
                onChange={(event) => setFilterWeekId(event.target.value)}
                min={1}
              />
            </label>
          </div>

          <div className="form-actions">
            <button type="submit">Aplicar filtros</button>
          </div>
        </form>

        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>ID del curso/proyecto</th>
                <th>ID de semana</th>
                <th>ID de vinculación</th>
                <th>ID del usuario</th>
                <th>Acción</th>
              </tr>
            </thead>
            <tbody>
              {reports.length === 0 ? (
                <EmptyState colSpan={6} />
              ) : (
                reports.map((item) => (
                  <tr key={item.id}>
                    <td>{item.id}</td>
                    <td>{item.workspace_id}</td>
                    <td>{item.week_id}</td>
                    <td>{item.assignment_id}</td>
                    <td>{item.user_id}</td>
                    <td>
                      <button type="button" onClick={() => void handleDownload(item.id)}>
                        Descargar PDF
                      </button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </section>
    </Layout>
  );
}
