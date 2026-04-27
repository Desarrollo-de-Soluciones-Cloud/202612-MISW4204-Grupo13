import { FormEvent, useEffect, useMemo, useState } from "react";
import {
  createAssignment,
  createPeriod,
  createUser,
  createWorkspace,
  downloadReport,
  getMe,
  listAssignmentsByUser,
  listPeriods,
  listReports,
  listTasks,
  listUsers,
  listWeeksByPeriod,
  listWorkspaces,
  toErrorMessage,
} from "../api/client";
import type { Assignment, GlobalRole, Period, Report, Task, User, Week, Workspace } from "../api/types";
import EmptyState from "../components/EmptyState";
import HelpText from "../components/HelpText";
import Layout from "../components/Layout";
import Loading from "../components/Loading";
import Toast from "../components/Toast";
import useToast from "../components/useToast";

interface AdminDashboardProps {
  user: User;
  onLogout: () => void;
}

export default function AdminDashboard({ user, onLogout }: AdminDashboardProps) {
  const [me, setMe] = useState<User | null>(null);
  const [users, setUsers] = useState<User[]>([]);
  const [periods, setPeriods] = useState<Period[]>([]);
  const [workspaces, setWorkspaces] = useState<Workspace[]>([]);
  const [assignments, setAssignments] = useState<Assignment[]>([]);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [reports, setReports] = useState<Report[]>([]);
  const [reportWeeks, setReportWeeks] = useState<Week[]>([]);
  const [reportFilters, setReportFilters] = useState({
    workspace_id: "",
    week_id: "",
  });
  const [loading, setLoading] = useState(false);
  const { toast, showToast, clearToast } = useToast();

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

  const assignableUsers = useMemo(
    () => users.filter((item) => item.global_role === assignmentForm.role),
    [users, assignmentForm.role],
  );

  const getWorkspaceLabel = (workspace: Workspace): string =>
    `${workspace.name} - ${workspace.type} - ${workspace.state} (ID ${workspace.id})`;

  const getPeriodLabel = (period: Period): string =>
    `${period.name} - ${period.period_state} (ID ${period.id})`;

  const getProfessorLabel = (professor: User): string =>
    `${professor.name} - ${professor.email} - ${professor.global_role} (ID ${professor.id})`;

  const getAssignableUserLabel = (account: User): string =>
    `${account.name} - ${account.email} - ${account.global_role} (ID ${account.id})`;

  const getWeekLabel = (week: Week): string =>
    `Semana ${week.number}: ${week.initial_date} a ${week.final_date} (ID ${week.id})`;

  const loadAll = async () => {
    setLoading(true);

    try {
      const [meResult, usersResult, periodsResult, workspacesResult, tasksResult] =
        await Promise.all([
          getMe(),
          listUsers(),
          listPeriods(),
          listWorkspaces(),
          listTasks(),
        ]);

      const defaultWorkspaceId = workspacesResult.workspaces[0]?.id;
      const selectedWorkspaceId = Number(reportFilters.workspace_id || String(defaultWorkspaceId ?? ""));
      const selectedWeekId = reportFilters.week_id ? Number(reportFilters.week_id) : undefined;
      const reportsResult = selectedWorkspaceId
        ? await listReports({
            workspace_id: selectedWorkspaceId,
            week_id: selectedWeekId,
          })
        : { reports: [] };

      const usersForAssignments = usersResult.users.filter(
        (item) => item.global_role === "monitor" || item.global_role === "assistant",
      );
      const assignmentResponses = await Promise.all(
        usersForAssignments.map((item) => listAssignmentsByUser(item.id)),
      );

      const uniqueAssignments = new Map<number, Assignment>();
      assignmentResponses
        .flatMap((item) => item.assignments)
        .forEach((item) => uniqueAssignments.set(item.id, item));

      const localProfessorUsers = usersResult.users.filter((item) => item.global_role === "professor");

      setMe(meResult);
      setUsers(usersResult.users);
      setPeriods(periodsResult.periods);
      setWorkspaces(workspacesResult.workspaces);
      setAssignments(Array.from(uniqueAssignments.values()));
      setTasks(tasksResult.tasks);
      setReports(reportsResult.reports);
      setReportFilters((previous) => ({
        ...previous,
        workspace_id: previous.workspace_id || String(defaultWorkspaceId ?? ""),
      }));

      setWorkspaceForm((previous) => ({
        ...previous,
        period_id: previous.period_id || String(periodsResult.periods[0]?.id ?? ""),
        user_id: previous.user_id || String(localProfessorUsers[0]?.id ?? ""),
      }));

      setAssignmentForm((previous) => ({
        ...previous,
        workspace_id: previous.workspace_id || String(workspacesResult.workspaces[0]?.id ?? ""),
      }));
    } catch (err) {
      showToast(toErrorMessage(err), "error");
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
      user_id: previous.user_id || String(professorUsers[0]?.id ?? ""),
    }));
  }, [professorUsers]);

  useEffect(() => {
    setAssignmentForm((previous) => ({
      ...previous,
      user_id: String(assignableUsers[0]?.id ?? ""),
    }));
  }, [assignableUsers]);

  useEffect(() => {
    const selectedWorkspaceId = Number(reportFilters.workspace_id);
    if (!selectedWorkspaceId) {
      setReports([]);
      setReportWeeks([]);
      setReportFilters((previous) => ({ ...previous, week_id: "" }));
      return;
    }

    const loadReports = async () => {
      try {
        const selectedWorkspace = workspaces.find((item) => item.id === selectedWorkspaceId);
        if (selectedWorkspace) {
          const weeksResponse = await listWeeksByPeriod(selectedWorkspace.period_id);
          setReportWeeks(weeksResponse.weeks);
          if (!weeksResponse.weeks.some((item) => String(item.id) === reportFilters.week_id)) {
            setReportFilters((previous) => ({ ...previous, week_id: "" }));
          }
        } else {
          setReportWeeks([]);
        }

        const response = await listReports({
          workspace_id: selectedWorkspaceId,
          week_id: reportFilters.week_id ? Number(reportFilters.week_id) : undefined,
        });
        setReports(response.reports);
      } catch (err) {
        showToast(toErrorMessage(err), "error");
      }
    };

    void loadReports();
  }, [reportFilters.workspace_id, workspaces]);

  const handleFilterReports = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    clearToast();

    try {
      if (!reportFilters.workspace_id) {
        showToast("El filtro por curso/proyecto es obligatorio para consultar reportes.", "error");
        return;
      }

      const response = await listReports({
        workspace_id: Number(reportFilters.workspace_id),
        week_id: reportFilters.week_id ? Number(reportFilters.week_id) : undefined,
      });
      setReports(response.reports);

      if (response.reports.length === 0) {
        showToast("No hay reportes generados para los filtros seleccionados.", "error");
      }
    } catch (err) {
      showToast(toErrorMessage(err), "error");
    }
  };

  const handleDownloadReport = async (reportId: number) => {
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
      showToast("Reporte descargado correctamente.", "success");
    } catch (err) {
      showToast(toErrorMessage(err), "error");
    }
  };

  const handleCreateUser = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    clearToast();

    try {
      await createUser(userForm);
      showToast("Usuario registrado correctamente.", "success");
      setUserForm({
        name: "",
        email: "",
        password: "",
        global_role: "professor",
      });
      await loadAll();
    } catch (err) {
      showToast(toErrorMessage(err), "error");
    }
  };

  const handleCreatePeriod = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    clearToast();

    try {
      await createPeriod({
        name: periodForm.name,
        initial_date: periodForm.initial_date,
        weeks_count: Number(periodForm.weeks_count) as 8 | 16,
        period_state: periodForm.period_state,
      });
      showToast("Periodo académico creado correctamente.", "success");
      setPeriodForm({
        name: "",
        initial_date: "",
        weeks_count: "16",
        period_state: "active",
      });
      await loadAll();
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
        user_id: Number(workspaceForm.user_id),
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
      await loadAll();
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
      showToast("Vinculación creada correctamente.", "success");
      setAssignmentForm((previous) => ({
        ...previous,
        weekly_hours: "1",
      }));
      await loadAll();
    } catch (err) {
      showToast(toErrorMessage(err), "error");
    }
  };

  return (
    <Layout
      title="Panel de administración"
      description="Gestiona usuarios, periodos académicos, cursos/proyectos y vinculaciones."
      user={user}
      onLogout={onLogout}
    >
      <div className="actions-row">
        <button onClick={() => void loadAll()} disabled={loading}>
          Recargar datos
        </button>
      </div>

      {loading && <Loading label="Actualizando información..." />}
      {toast ? <Toast type={toast.type} message={toast.message} onClose={clearToast} /> : null}

      <section className="card">
        <h2>Registrar usuario</h2>
        <p className="section-desc">Crea una cuenta nueva para un profesor, monitor o asistente graduado.</p>
        <form className="form-grid" onSubmit={handleCreateUser}>
          <div className="form-field">
            <label>
              Nombre completo
              <input
                value={userForm.name}
                onChange={(event) =>
                  setUserForm((previous) => ({ ...previous, name: event.target.value }))
                }
                required
              />
            </label>
          </div>

          <div className="form-field">
            <label>
              Correo electrónico
              <input
                type="email"
                value={userForm.email}
                onChange={(event) =>
                  setUserForm((previous) => ({ ...previous, email: event.target.value }))
                }
                required
              />
            </label>
          </div>

          <div className="form-field">
            <label>
              Contraseña
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
            <HelpText>Debe tener entre 8 y 72 caracteres.</HelpText>
          </div>

          <div className="form-field">
            <label>
              Rol global
              <select
                value={userForm.global_role}
                onChange={(event) =>
                  setUserForm((previous) => ({
                    ...previous,
                    global_role: event.target.value as GlobalRole,
                  }))
                }
              >
                <option value="admin">Administrador</option>
                <option value="professor">Profesor</option>
                <option value="monitor">Monitor</option>
                <option value="assistant">Asistente graduado</option>
              </select>
            </label>
            <HelpText>
              El rol define qué información podrá consultar y modificar el usuario.
            </HelpText>
          </div>

          <div className="form-actions">
            <button type="submit">Registrar usuario</button>
          </div>
        </form>
      </section>

      <section className="card">
        <h2>Crear periodo académico</h2>
        <p className="section-desc">Define el rango de fechas y número de semanas del periodo lectivo.</p>
        <form className="form-grid" onSubmit={handleCreatePeriod}>
          <div className="form-field">
            <label>
              Nombre del periodo
              <input
                placeholder="2026-20"
                value={periodForm.name}
                onChange={(event) =>
                  setPeriodForm((previous) => ({ ...previous, name: event.target.value }))
                }
                required
              />
            </label>
          </div>

          <div className="form-field">
            <label>
              Fecha inicial
              <input
                type="date"
                value={periodForm.initial_date}
                onChange={(event) =>
                  setPeriodForm((previous) => ({ ...previous, initial_date: event.target.value }))
                }
                required
              />
            </label>
            <HelpText>La fecha inicial debe ser un lunes y debe ser futura.</HelpText>
          </div>

          <div className="form-field">
            <label>
              Número de semanas
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
            <HelpText>Solo se permiten periodos de 8 o 16 semanas.</HelpText>
          </div>

          <div className="form-field">
            <label>
              Estado del periodo
              <select
                value={periodForm.period_state}
                onChange={(event) =>
                  setPeriodForm((previous) => ({
                    ...previous,
                    period_state: event.target.value as "active" | "closed",
                  }))
                }
              >
                <option value="active">Activo</option>
                <option value="closed">Cerrado</option>
              </select>
            </label>
            <HelpText>Un periodo cerrado no permite nuevas inscripciones.</HelpText>
          </div>

          <div className="form-actions">
            <button type="submit">Crear periodo académico</button>
          </div>
        </form>
      </section>

      <section className="card">
        <h2>Crear curso o proyecto</h2>
        <p className="section-desc">Asocia un curso o proyecto al periodo académico y al profesor responsable.</p>
        <form className="form-grid" onSubmit={handleCreateWorkspace}>
          <div className="form-field">
            <label>
              Periodo
              <select
                value={workspaceForm.period_id}
                onChange={(event) =>
                  setWorkspaceForm((previous) => ({ ...previous, period_id: event.target.value }))
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
            <HelpText>
              Solo se pueden crear cursos/proyectos en periodos con inscripción abierta.
            </HelpText>
          </div>

          <div className="form-field">
            <label>
              Profesor responsable
              <select
                value={workspaceForm.user_id}
                onChange={(event) =>
                  setWorkspaceForm((previous) => ({ ...previous, user_id: event.target.value }))
                }
                required
              >
                <option value="" disabled>
                  Selecciona un profesor
                </option>
                {professorUsers.map((item) => (
                  <option key={item.id} value={item.id}>
                    {getProfessorLabel(item)}
                  </option>
                ))}
              </select>
            </label>
            <HelpText>El usuario seleccionado debe tener rol Profesor.</HelpText>
          </div>

          <div className="form-field">
            <label>
              Nombre del curso o proyecto
              <input
                value={workspaceForm.name}
                onChange={(event) =>
                  setWorkspaceForm((previous) => ({ ...previous, name: event.target.value }))
                }
                required
              />
            </label>
          </div>

          <div className="form-field">
            <label>
              Tipo
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
              Fecha inicial
              <input
                type="date"
                value={workspaceForm.initial_date}
                onChange={(event) =>
                  setWorkspaceForm((previous) => ({ ...previous, initial_date: event.target.value }))
                }
                required
              />
            </label>
          </div>

          <div className="form-field">
            <label>
              Fecha final
              <input
                type="date"
                value={workspaceForm.final_date}
                onChange={(event) =>
                  setWorkspaceForm((previous) => ({ ...previous, final_date: event.target.value }))
                }
                required
              />
            </label>
            <HelpText>La fecha inicial debe ser anterior a la fecha final.</HelpText>
          </div>

          <div className="form-field">
            <label>
              Observaciones
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
            <HelpText>Incluye información breve sobre el curso o proyecto.</HelpText>
          </div>

          <div className="form-field">
            <label>
              Estado
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
            <button type="submit">Crear curso o proyecto</button>
          </div>
        </form>
      </section>

      <section className="card">
        <h2>Crear vinculación</h2>
        <p className="section-desc">Asigna un monitor o asistente a un curso o proyecto con sus horas semanales.</p>
        <form className="form-grid" onSubmit={handleCreateAssignment}>
          <div className="form-field">
            <label>
              Rol de vinculación
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
                <option value="assistant">Asistente graduado</option>
              </select>
            </label>
            <HelpText>Solo se pueden vincular monitores o asistentes graduados.</HelpText>
          </div>

          <div className="form-field">
            <label>
              Usuario
              <select
                value={assignmentForm.user_id}
                onChange={(event) =>
                  setAssignmentForm((previous) => ({ ...previous, user_id: event.target.value }))
                }
                required
              >
                <option value="" disabled>
                  Selecciona un usuario
                </option>
                {assignableUsers.map((item) => (
                  <option key={item.id} value={item.id}>
                    {getAssignableUserLabel(item)}
                  </option>
                ))}
              </select>
            </label>
          </div>

          <div className="form-field">
            <label>
              Curso o proyecto
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
              Horas semanales
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
            <HelpText>
              Monitor: máximo 12 horas. Asistente: máximo 22 horas. Los monitores también están
              sujetos a la regla del 40% respecto a asistentes.
            </HelpText>
          </div>

          <div className="form-actions">
            <button type="submit">Crear vinculación</button>
          </div>
        </form>

        <p className="muted">
          El sistema mostrará errores del backend cuando se incumplan reglas de negocio de horas o
          permisos.
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
        <h2>Usuarios registrados</h2>
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Nombre</th>
                <th>Correo</th>
                <th>Rol</th>
              </tr>
            </thead>
            <tbody>
              {users.length === 0 ? (
                <EmptyState colSpan={4} />
              ) : (
                users.map((item) => (
                  <tr key={item.id}>
                    <td>{item.id}</td>
                    <td>{item.name}</td>
                    <td>{item.email}</td>
                    <td>{item.global_role}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </section>

      <section className="card">
        <h2>Periodos académicos</h2>
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Nombre</th>
                <th>Fecha inicial</th>
                <th>Fecha final</th>
                <th>Estado</th>
              </tr>
            </thead>
            <tbody>
              {periods.length === 0 ? (
                <EmptyState colSpan={5} />
              ) : (
                periods.map((item) => (
                  <tr key={item.id}>
                    <td>{item.id}</td>
                    <td>{item.name}</td>
                    <td>{item.initial_date}</td>
                    <td>{item.final_date}</td>
                    <td>{item.period_state}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </section>

      <section className="card">
        <h2>Cursos y proyectos</h2>
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Nombre</th>
                <th>ID del usuario</th>
                <th>ID del periodo</th>
                <th>Tipo</th>
                <th>Estado</th>
              </tr>
            </thead>
            <tbody>
              {workspaces.length === 0 ? (
                <EmptyState colSpan={6} />
              ) : (
                workspaces.map((item) => (
                  <tr key={item.id}>
                    <td>{item.id}</td>
                    <td>{item.name}</td>
                    <td>{item.user_id}</td>
                    <td>{item.period_id}</td>
                    <td>{item.type}</td>
                    <td>{item.state}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </section>

      <section className="card">
        <h2>Vinculaciones</h2>
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>ID del usuario</th>
                <th>ID del curso/proyecto</th>
                <th>Rol</th>
                <th>Horas semanales</th>
              </tr>
            </thead>
            <tbody>
              {assignments.length === 0 ? (
                <EmptyState colSpan={5} />
              ) : (
                assignments.map((item) => (
                  <tr key={item.id}>
                    <td>{item.id}</td>
                    <td>{item.user_id}</td>
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
        <h2>Tareas reportadas</h2>
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>ID del usuario</th>
                <th>ID de vinculación</th>
                <th>Estado</th>
                <th>Horas dedicadas</th>
                <th>Fecha de inicio de semana</th>
              </tr>
            </thead>
            <tbody>
              {tasks.length === 0 ? (
                <EmptyState colSpan={6} />
              ) : (
                tasks.map((item) => (
                  <tr key={item.id}>
                    <td>{item.id}</td>
                    <td>{item.user_id}</td>
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
        <h2>Reportes generados</h2>
        <form className="form-grid" onSubmit={handleFilterReports}>
          <div className="form-field">
            <label>
              ID del curso/proyecto
              <select
                value={reportFilters.workspace_id}
                onChange={(event) =>
                  setReportFilters((previous) => ({
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
              Semana
              <select
                value={reportFilters.week_id}
                onChange={(event) =>
                  setReportFilters((previous) => ({
                    ...previous,
                    week_id: event.target.value,
                  }))
                }
                disabled={reportFilters.workspace_id === "" || reportWeeks.length === 0}
              >
                <option value="">Todas las semanas</option>
                {reportWeeks.map((item) => (
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
                <th>Ruta del archivo</th>
                <th>Fecha de creación</th>
                <th>Acción</th>
              </tr>
            </thead>
            <tbody>
              {reports.length === 0 ? (
                <EmptyState
                  colSpan={8}
                  message="No hay reportes generados para los filtros seleccionados."
                />
              ) : (
                reports.map((item) => (
                  <tr key={item.id}>
                    <td>{item.id}</td>
                    <td>{item.workspace_id}</td>
                    <td>{item.week_id}</td>
                    <td>{item.assignment_id}</td>
                    <td>{item.user_id}</td>
                    <td>{item.file_path}</td>
                    <td>{item.created_at ?? "-"}</td>
                    <td>
                      <button
                        type="button"
                        onClick={() => void handleDownloadReport(item.id)}
                      >
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
