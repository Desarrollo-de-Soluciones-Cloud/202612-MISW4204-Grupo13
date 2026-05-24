import { FormEvent, useEffect, useState } from "react";
import {
  createAssignment,
  createWorkspace,
  downloadReport,
  generateWeeklyReport,
  getMe,
  listPeriods,
  listReports,
  listTasks,
  listUsers,
  listWeeksByPeriod,
  listWorkspaces,
  toErrorMessage,
} from "../api/client";
import type { Period, Report, Task, User, Week, Workspace } from "../api/types";
import EmptyState from "../components/EmptyState";
import HelpText from "../components/HelpText";
import Layout from "../components/Layout";
import Loading from "../components/Loading";
import Toast from "../components/Toast";
import useToast from "../components/useToast";

type ProfessorDashboardProps = Readonly<{
  user: User;
  onLogout: () => void;
}>;

const reportPollingIntervalMs = 4000;
const reportPollingAttempts = 15;

function delay(ms: number): Promise<void> {
  return new Promise((resolve) => {
    globalThis.setTimeout(resolve, ms);
  });
}

export default function ProfessorDashboard({ user, onLogout }: ProfessorDashboardProps) {
  const [me, setMe] = useState<User | null>(null);
  const [workspaces, setWorkspaces] = useState<Workspace[]>([]);
  const [periods, setPeriods] = useState<Period[]>([]);
  const [assignableUsers, setAssignableUsers] = useState<User[]>([]);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [reports, setReports] = useState<Report[]>([]);
  const [weeksByPeriod, setWeeksByPeriod] = useState<Record<number, Week[]>>({});
  const [workspaceId, setWorkspaceId] = useState("");
  const [weekId, setWeekId] = useState("");
  const [filterWorkspaceId, setFilterWorkspaceId] = useState("");
  const [filterWeekId, setFilterWeekId] = useState("");
  const [workspaceForm, setWorkspaceForm] = useState({
    period_id: "",
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
  const [loading, setLoading] = useState(false);
  const { toast, showToast, clearToast } = useToast();

  const getReportsWorkspaceId = (fallbackWorkspaceId?: number): number | undefined => {
    if (filterWorkspaceId) {
      return Number(filterWorkspaceId);
    }

    if (workspaceId) {
      return Number(workspaceId);
    }

    return fallbackWorkspaceId;
  };

  const getWorkspaceLabel = (workspace: Workspace): string =>
    `${workspace.name} - ${workspace.type} (ID ${workspace.id})`;

  const getPeriodLabel = (period: Period): string =>
    `${period.name} - ${period.period_state} (ID ${period.id})`;

  const getUserLabel = (account: User): string =>
    `${account.name} - ${account.email} - ${account.global_role} (ID ${account.id})`;

  const getWeekLabel = (week: Week): string =>
    `Semana ${week.number}: ${week.initial_date} a ${week.final_date} (ID ${week.id})`;

  const getReportGenerationContext = (selectedWorkspaceId: number, selectedWeekId: number) => {
    const workspace = workspaces.find((item) => item.id === selectedWorkspaceId);
    const week =
      Object.values(weeksByPeriod)
        .flat()
        .find((item) => item.id === selectedWeekId) ?? null;

    return {
      workspaceLabel: workspace?.name ?? `workspace ${selectedWorkspaceId}`,
      weekLabel: week ? `semana ${week.number} (${week.initial_date} a ${week.final_date})` : `semana ${selectedWeekId}`,
    };
  };

  const loadReportsForFilters = async (selectedWorkspaceId: number, selectedWeekId?: number) => {
    const response = await listReports({
      workspace_id: selectedWorkspaceId,
      week_id: selectedWeekId,
    });
    setReports(response.reports);
    return response.reports;
  };

  const waitForQueuedReports = async (
    selectedWorkspaceId: number,
    selectedWeekId: number,
    baselineCount: number,
    queuedCount: number,
  ) => {
    const expectedTotal = baselineCount + queuedCount;
    let latestCount = baselineCount;
    const { workspaceLabel, weekLabel } = getReportGenerationContext(
      selectedWorkspaceId,
      selectedWeekId,
    );

    for (let attempt = 0; attempt < reportPollingAttempts; attempt += 1) {
      if (attempt > 0) {
        await delay(reportPollingIntervalMs);
      }

      const nextReports = await loadReportsForFilters(selectedWorkspaceId, selectedWeekId);
      latestCount = nextReports.filter((item) => item.week_id === selectedWeekId).length;

      if (latestCount >= expectedTotal) {
        showToast(
          `La generacion de reportes para ${workspaceLabel} en la ${weekLabel} finalizo con exito. Los resultados ya estan disponibles en la lista para descarga.`,
          "success",
        );
        return;
      }
    }

    const generatedCount = Math.max(latestCount - baselineCount, 0);
    if (generatedCount === 0) {
      showToast(
        `La generacion de reportes para ${workspaceLabel} en la ${weekLabel} se esta demorando mas de lo esperado. Revisa mas tarde la lista de reportes.`,
        "info",
      );
      return;
    }

    showToast(
      `La generacion de reportes para ${workspaceLabel} en la ${weekLabel} se esta demorando mas de lo esperado. Ya hay resultados parciales visibles en la lista.`,
      "info",
    );
  };

  const loadWeeksForPeriod = async (periodId: number) => {
    if (weeksByPeriod[periodId]) {
      return;
    }

    const response = await listWeeksByPeriod(periodId);
    setWeeksByPeriod((previous) => ({
      ...previous,
      [periodId]: response.weeks,
    }));
  };

  const selectedWorkspace = workspaces.find((item) => String(item.id) === workspaceId);
  const selectedFilterWorkspace = workspaces.find((item) => String(item.id) === filterWorkspaceId);
  const generateWeeks = selectedWorkspace ? weeksByPeriod[selectedWorkspace.period_id] ?? [] : [];
  const filterWeeks = selectedFilterWorkspace
    ? weeksByPeriod[selectedFilterWorkspace.period_id] ?? []
    : [];

  const loadBase = async () => {
    setLoading(true);

    try {
      const [meResult, workspaceResult, taskResult, periodsResult, monitorResult, assistantResult] = await Promise.all([
        getMe(),
        listWorkspaces(),
        listTasks(),
        listPeriods(),
        listUsers("monitor"),
        listUsers("assistant"),
      ]);

      const availableUsers = [...monitorResult.users, ...assistantResult.users];

      const defaultWorkspaceId = workspaceResult.workspaces[0]?.id;
      const reportResult = defaultWorkspaceId
        ? await listReports({ workspace_id: defaultWorkspaceId })
        : { reports: [] };

      setMe(meResult);
      setWorkspaces(workspaceResult.workspaces);
      setPeriods(periodsResult.periods);
      setAssignableUsers(availableUsers);
      setTasks(taskResult.tasks);
      setReports(reportResult.reports);
      setWorkspaceId((previous) => previous || String(defaultWorkspaceId ?? ""));
      setFilterWorkspaceId((previous) => previous || String(defaultWorkspaceId ?? ""));
      setWorkspaceForm((previous) => ({
        ...previous,
        period_id: previous.period_id || String(periodsResult.periods[0]?.id ?? ""),
      }));
      setAssignmentForm((previous) => ({
        ...previous,
        workspace_id: previous.workspace_id || String(defaultWorkspaceId ?? ""),
      }));
      setAssignmentForm((previous) => ({
        ...previous,
        user_id:
          previous.user_id ||
          String(availableUsers.find((item) => item.global_role === previous.role)?.id ?? ""),
      }));
    } catch (err) {
      showToast(toErrorMessage(err), "error");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    setAssignmentForm((previous) => ({
      ...previous,
      user_id:
        previous.user_id &&
        assignableUsers.some((item) => String(item.id) === previous.user_id && item.global_role === previous.role)
          ? previous.user_id
          : String(assignableUsers.find((item) => item.global_role === previous.role)?.id ?? ""),
    }));
  }, [assignmentForm.role, assignableUsers]);

  useEffect(() => {
    void loadBase();
  }, []);

  useEffect(() => {
    const loadGenerateWeeks = async () => {
      if (!selectedWorkspace) {
        setWeekId("");
        return;
      }

      try {
        await loadWeeksForPeriod(selectedWorkspace.period_id);
      } catch (err) {
        showToast(toErrorMessage(err), "error");
      }
    };

    void loadGenerateWeeks();
  }, [workspaceId, workspaces]);

  useEffect(() => {
    if (generateWeeks.length === 0) {
      setWeekId("");
      return;
    }

    setWeekId((previous) => {
      if (previous && generateWeeks.some((item) => String(item.id) === previous)) {
        return previous;
      }

      return String(generateWeeks[0].id);
    });
  }, [generateWeeks]);

  useEffect(() => {
    const loadFilterWeeks = async () => {
      if (!selectedFilterWorkspace) {
        setFilterWeekId("");
        return;
      }

      try {
        await loadWeeksForPeriod(selectedFilterWorkspace.period_id);
      } catch (err) {
        showToast(toErrorMessage(err), "error");
      }
    };

    void loadFilterWeeks();
  }, [filterWorkspaceId, workspaces]);

  useEffect(() => {
    if (filterWeekId && !filterWeeks.some((item) => String(item.id) === filterWeekId)) {
      setFilterWeekId("");
    }
  }, [filterWeeks]);

  useEffect(() => {
    const selectedWorkspaceId = Number(filterWorkspaceId);
    if (!selectedWorkspaceId) {
      setReports([]);
      return;
    }

    const loadReports = async () => {
      try {
        await loadReportsForFilters(
          selectedWorkspaceId,
          filterWeekId ? Number(filterWeekId) : undefined,
        );
      } catch (err) {
        showToast(toErrorMessage(err), "error");
      }
    };

    void loadReports();
  }, [filterWorkspaceId]);

  const handleGenerateReport = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    clearToast();

    try {
      const selectedWorkspaceId = Number(workspaceId);
      const selectedWeekId = Number(weekId);
      const baselineCount = reports.filter((item) => item.week_id === selectedWeekId).length;
      const { workspaceLabel, weekLabel } = getReportGenerationContext(
        selectedWorkspaceId,
        selectedWeekId,
      );

      const response = await generateWeeklyReport({
        workspace_id: selectedWorkspaceId,
        week_id: selectedWeekId,
      });
      if (response.reports.length === 0) {
        setFilterWorkspaceId(String(selectedWorkspaceId));
        setFilterWeekId(String(selectedWeekId));
        showToast(
          `Se inicio la generacion asincrona de reportes para ${workspaceLabel} en la ${weekLabel}. Te avisaremos cuando los resultados empiecen a quedar disponibles.`,
          "info",
        );
        setWeekId("");
        await waitForQueuedReports(
          selectedWorkspaceId,
          selectedWeekId,
          baselineCount,
          response.generated_count,
        );
        return;
      }

      showToast(
        `La generacion de reportes para ${workspaceLabel} en la ${weekLabel} finalizo con exito. Los resultados ya estan disponibles para descarga.`,
        "success",
      );
      setWeekId("");

      const reportsWorkspaceId = getReportsWorkspaceId(selectedWorkspaceId);
      if (reportsWorkspaceId) {
        await loadReportsForFilters(
          reportsWorkspaceId,
          filterWeekId ? Number(filterWeekId) : undefined,
        );
      } else {
        setReports([]);
      }
    } catch (err) {
      showToast(toErrorMessage(err), "error");
    }
  };

  const handleFilterReports = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    clearToast();

    try {
      if (!filterWorkspaceId) {
        showToast("El filtro por curso/proyecto es obligatorio para consultar reportes.", "error");
        return;
      }

      await loadReportsForFilters(
        Number(filterWorkspaceId),
        filterWeekId ? Number(filterWeekId) : undefined,
      );
    } catch (err) {
      showToast(toErrorMessage(err), "error");
    }
  };

  const handleDownload = async (reportId: number) => {
    clearToast();

    try {
      const blob = await downloadReport(reportId);
      const fileUrl = globalThis.URL.createObjectURL(blob);
      const anchor = document.createElement("a");
      anchor.href = fileUrl;
      anchor.download = `reporte_${reportId}.pdf`;
      document.body.appendChild(anchor);
      anchor.click();
      anchor.remove();
      globalThis.URL.revokeObjectURL(fileUrl);
      showToast(`Reporte ${reportId} descargado correctamente.`, "success");
    } catch (err) {
      showToast(toErrorMessage(err), "error");
    }
  };

  const handleCreateWorkspace = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    clearToast();

    try {
      await createWorkspace({
        period_id: Number(workspaceForm.period_id),
        user_id: me?.id ?? user.id,
        name: workspaceForm.name,
        type: workspaceForm.type,
        initial_date: workspaceForm.initial_date,
        final_date: workspaceForm.final_date,
        observations: workspaceForm.observations,
        state: workspaceForm.state,
      });

      showToast("Curso o proyecto creado correctamente.", "success");
      setWorkspaceForm((previous) => ({
        ...previous,
        name: "",
        initial_date: "",
        final_date: "",
        observations: "",
      }));
      await loadBase();
    } catch (err) {
      showToast(toErrorMessage(err), "error");
    }
  };

  const handleCreateAssignment = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    clearToast();

    try {
      await createAssignment({
        user_id: Number(assignmentForm.user_id),
        workspace_id: Number(assignmentForm.workspace_id),
        role: assignmentForm.role,
        weekly_hours: Number(assignmentForm.weekly_hours),
      });

      showToast("Vinculacion creada correctamente.", "success");
      setAssignmentForm((previous) => ({
        ...previous,
        weekly_hours: "1",
      }));
      await loadBase();
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
          El sistema genera reportes PDF semanales con apoyo de IA usando las tareas registradas,
          sus descripciones, observaciones y las horas reportadas por monitores y asistentes.
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
        <h2>Crear curso o proyecto</h2>
        <p className="section-desc">Registra un nuevo curso/proyecto asociado a tu cuenta de profesor.</p>
        <form className="form-grid" onSubmit={handleCreateWorkspace}>
          <div className="form-field">
            <label>
              <span>Periodo academico</span>
              <select
                value={workspaceForm.period_id}
                onChange={(event) =>
                  setWorkspaceForm((previous) => ({
                    ...previous,
                    period_id: event.target.value,
                  }))
                }
                required
              >
                <option value="" disabled>
                  Selecciona un periodo
                </option>
                {periods.map((item) => (
                  <option key={item.id} value={item.id}>
                    {getPeriodLabel(item)}
                  </option>
                ))}
              </select>
            </label>
          </div>

          <div className="form-field">
            <label>
              <span>Nombre</span>
              <input
                value={workspaceForm.name}
                onChange={(event) =>
                  setWorkspaceForm((previous) => ({
                    ...previous,
                    name: event.target.value,
                  }))
                }
                required
              />
            </label>
          </div>

          <div className="form-field">
            <label>
              <span>Tipo</span>
              <select
                value={workspaceForm.type}
                onChange={(event) =>
                  setWorkspaceForm((previous) => ({
                    ...previous,
                    type: event.target.value as "course" | "project",
                  }))
                }
              >
                <option value="course">Curso</option>
                <option value="project">Proyecto</option>
              </select>
            </label>
          </div>

          <div className="form-field">
            <label>
              <span>Fecha inicial</span>
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
          </div>

          <div className="form-field">
            <label>
              <span>Fecha final</span>
              <input
                type="date"
                value={workspaceForm.final_date}
                onChange={(event) =>
                  setWorkspaceForm((previous) => ({
                    ...previous,
                    final_date: event.target.value,
                  }))
                }
                required
              />
            </label>
          </div>

          <div className="form-field">
            <label>
              <span>Observaciones</span>
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
          </div>

          <div className="form-field">
            <label>
              <span>Estado</span>
              <select
                value={workspaceForm.state}
                onChange={(event) =>
                  setWorkspaceForm((previous) => ({
                    ...previous,
                    state: event.target.value as "active" | "closed",
                  }))
                }
              >
                <option value="active">Activo</option>
                <option value="closed">Cerrado</option>
              </select>
            </label>
          </div>

          <div className="form-actions">
            <button type="submit">Crear curso/proyecto</button>
          </div>
        </form>
      </section>

      <section className="card">
        <h2>Crear vinculacion</h2>
        <p className="section-desc">Vincula monitores o asistentes a uno de tus cursos/proyectos.</p>
        <form className="form-grid" onSubmit={handleCreateAssignment}>
          <div className="form-field">
            <label>
              <span>Curso/proyecto propio</span>
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
                  Selecciona un curso/proyecto
                </option>
                {workspaces.map((item) => (
                  <option key={item.id} value={item.id}>
                    {getWorkspaceLabel(item)}
                  </option>
                ))}
              </select>
            </label>
          </div>

          <div className="form-field">
            <label>
              <span>Rol de vinculacion</span>
              <select
                value={assignmentForm.role}
                onChange={(event) =>
                  setAssignmentForm((previous) => ({
                    ...previous,
                    role: event.target.value as "monitor" | "assistant",
                  }))
                }
              >
                <option value="monitor">Monitor</option>
                <option value="assistant">Asistente</option>
              </select>
            </label>
          </div>

          <div className="form-field">
            <label>
              <span>Usuario</span>
              <select
                value={assignmentForm.user_id}
                onChange={(event) =>
                  setAssignmentForm((previous) => ({
                    ...previous,
                    user_id: event.target.value,
                  }))
                }
                required
              >
                <option value="" disabled>
                  Selecciona un usuario
                </option>
                {assignableUsers
                  .filter((item) => item.global_role === assignmentForm.role)
                  .map((item) => (
                    <option key={item.id} value={item.id}>
                      {getUserLabel(item)}
                    </option>
                  ))}
              </select>
            </label>
          </div>

          <div className="form-field">
            <label>
              <span>Horas semanales</span>
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
          </div>

          <div className="form-actions">
            <button type="submit">Crear vinculacion</button>
          </div>
        </form>
      </section>

      <section className="card">
        <h2>Generar reporte semanal</h2>
        <p className="section-desc">Genera el reporte PDF de actividad para la semana académica seleccionada.</p>
        <form className="form-grid" onSubmit={handleGenerateReport}>
          <div className="form-field">
            <label>
              <span>Curso/proyecto</span>
              <select
                value={workspaceId}
                onChange={(event) => setWorkspaceId(event.target.value)}
                required
              >
                <option value="" disabled>
                  Selecciona un curso/proyecto
                </option>
                {workspaces.map((item) => (
                  <option key={item.id} value={item.id}>
                    {getWorkspaceLabel(item)}
                  </option>
                ))}
              </select>
            </label>
            <HelpText>Selecciona un curso o proyecto propio.</HelpText>
          </div>

          <div className="form-field">
            <label>
              <span>Semana</span>
              <select
                value={weekId}
                onChange={(event) => setWeekId(event.target.value)}
                required
                disabled={generateWeeks.length === 0}
              >
                <option value="" disabled>
                  {generateWeeks.length === 0
                    ? "No hay semanas disponibles"
                    : "Selecciona una semana"}
                </option>
                {generateWeeks.map((item) => (
                  <option key={item.id} value={item.id}>
                    {getWeekLabel(item)}
                  </option>
                ))}
              </select>
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
              <span>Curso/proyecto</span>
              <select
                value={filterWorkspaceId}
                onChange={(event) => setFilterWorkspaceId(event.target.value)}
                required
              >
                <option value="" disabled>
                  Selecciona un curso/proyecto
                </option>
                {workspaces.map((item) => (
                  <option key={item.id} value={item.id}>
                    {getWorkspaceLabel(item)}
                  </option>
                ))}
              </select>
            </label>
          </div>

          <div className="form-field">
            <label>
              <span>Semana</span>
              <select
                value={filterWeekId}
                onChange={(event) => setFilterWeekId(event.target.value)}
                disabled={filterWorkspaceId === "" || filterWeeks.length === 0}
              >
                <option value="">Todas las semanas</option>
                {filterWeeks.map((item) => (
                  <option key={item.id} value={item.id}>
                    {getWeekLabel(item)}
                  </option>
                ))}
              </select>
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

