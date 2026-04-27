import { getAccessToken } from "../auth/authStorage";
import type {
  Assignment,
  AuthResponse,
  CreateAssignmentPayload,
  CreatePeriodPayload,
  CreateTaskPayload,
  CreateUserPayload,
  CreateWorkspacePayload,
  GenerateWeeklyReportPayload,
  GenerateWeeklyReportResponse,
  ListAssignmentsResponse,
  ListPeriodsResponse,
  ListReportsResponse,
  ListTasksResponse,
  ListUsersResponse,
  ListWeeksResponse,
  ListWorkspacesResponse,
  Period,
  Task,
  User,
  Week,
  Workspace,
  UpdateTaskPayload,
} from "./types";

const API_BASE_URL = "/api";

interface ApiErrorBody {
  error?: string;
  message?: string;
  details?: string | string[];
  errors?: Array<string | { message?: string; error?: string; field?: string }>;
}

export class ApiError extends Error {
  status: number;
  details: string[];

  constructor(status: number, message: string, details: string[] = []) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.details = details;
  }
}

function normalizeText(value: string): string {
  return value.trim().replace(/\s+/g, " ");
}

function formatStatusMessage(status: number): string {
  if (status === 400) {
    return "Revisa los datos ingresados. Hay campos invalidos o incompletos.";
  }
  if (status === 401) {
    return "Tu sesion expiro o no has iniciado sesion.";
  }
  if (status === 403) {
    return "No tienes permisos para realizar esta accion.";
  }
  if (status === 404) {
    return "No se encontro el recurso solicitado.";
  }
  if (status === 409) {
    return "La operacion no se puede completar porque genera conflicto con datos existentes.";
  }
  if (status === 500) {
    return "Ocurrio un error interno. Revisa que los datos seleccionados existan y esten relacionados correctamente.";
  }

  return "Ocurrio un error inesperado en la solicitud.";
}

function mapDomainErrorMessage(rawMessage: string, status: number): string {
  const message = normalizeText(rawMessage);
  const normalized = message.toLowerCase();

  if (
    normalized.includes("week") &&
    (normalized.includes("not found") || normalized.includes("does not exist"))
  ) {
    return "La semana seleccionada no existe.";
  }

  if (
    normalized.includes("assignment") &&
    (normalized.includes("not found") || normalized.includes("does not exist"))
  ) {
    return "La vinculacion seleccionada no existe.";
  }

  if (normalized.includes("workspace") && normalized.includes("closed")) {
    return "No se pueden crear tareas en un curso o proyecto cerrado.";
  }

  if (
    normalized.includes("week is not active") ||
    normalized.includes("late update forbidden") ||
    normalized.includes("cannot be updated") ||
    normalized.includes("cannot be modified")
  ) {
    return "La tarea no se puede modificar porque la semana ya no esta activa.";
  }

  if (
    normalized.includes("title") ||
    normalized.includes("description") ||
    normalized.includes("week_start_date") ||
    normalized.includes("status") ||
    normalized.includes("spent_hours") ||
    normalized.includes("binding") ||
    normalized.includes("required")
  ) {
    if (status === 400) {
      return "Completa titulo, descripcion, semana, estado y horas.";
    }
  }

  if (normalized.includes("40%") || normalized.includes("40 percent") || normalized.includes("forty")) {
    return "La vinculacion no cumple la regla del 40% entre horas de monitor y asistente.";
  }

  if (normalized.includes("already exists") && normalized.includes("assignment")) {
    return "Ya existe una vinculacion equivalente para ese usuario y curso/proyecto.";
  }

  if (
    normalized.includes("monitor weekly hours cannot exceed") ||
    normalized.includes("assistant weekly hours cannot exceed")
  ) {
    return message;
  }

  if (normalized.includes("period") && normalized.includes("closed")) {
    return "No se puede crear un curso/proyecto en un periodo cerrado.";
  }

  if (normalized.includes("professor") && (normalized.includes("not found") || normalized.includes("invalid"))) {
    return "El profesor seleccionado no existe o no tiene rol de profesor.";
  }

  if (normalized.includes("date") && normalized.includes("period")) {
    return "Revisa las fechas del periodo academico.";
  }

  if (normalized.includes("weeks_count") || normalized.includes("weeks count")) {
    return "Revisa la cantidad de semanas del periodo academico.";
  }

  if (normalized.includes("no tasks") && normalized.includes("week")) {
    return "No hay tareas reportadas para ese curso/proyecto y semana.";
  }

  if (normalized.includes("report workspace not found") || normalized.includes("workspace not found")) {
    return "No se encontro el curso/proyecto seleccionado para generar el reporte.";
  }

  if (normalized.includes("report week not found") || (normalized.includes("week") && normalized.includes("report"))) {
    return "No se encontro la semana seleccionada para generar el reporte.";
  }

  if (normalized.includes("ai report generation failed") || normalized.includes("ai generation failed")) {
    return "No fue posible generar el resumen con IA. Revisa la configuracion del servicio.";
  }

  if (normalized.includes("pdf report generation failed") || normalized.includes("pdf generation failed")) {
    return "No fue posible generar el PDF del reporte.";
  }

  if (normalized.includes("report file not found") || normalized.includes("download")) {
    return "No se pudo descargar el PDF. Verifica que el archivo exista.";
  }

  if (normalized === "internal server error") {
    return formatStatusMessage(status);
  }

  return message;
}

function extractApiMessages(body: ApiErrorBody | null): string[] {
  if (!body) {
    return [];
  }

  const messages: string[] = [];

  if (typeof body.error === "string" && body.error.trim() !== "") {
    messages.push(body.error);
  }

  if (typeof body.message === "string" && body.message.trim() !== "") {
    messages.push(body.message);
  }

  if (typeof body.details === "string" && body.details.trim() !== "") {
    messages.push(body.details);
  }

  if (Array.isArray(body.details)) {
    for (const detail of body.details) {
      if (typeof detail === "string" && detail.trim() !== "") {
        messages.push(detail);
      }
    }
  }

  if (Array.isArray(body.errors)) {
    for (const item of body.errors) {
      if (typeof item === "string" && item.trim() !== "") {
        messages.push(item);
        continue;
      }

      if (item && typeof item === "object") {
        const value = item.message || item.error || item.field;
        if (typeof value === "string" && value.trim() !== "") {
          messages.push(value);
        }
      }
    }
  }

  return messages;
}

async function parseErrorResponse(response: Response): Promise<{ message: string; details: string[] }> {
  const contentType = response.headers.get("content-type") || "";
  if (contentType.includes("application/json")) {
    const body = await parseErrorBody(response);
    const messages = extractApiMessages(body);
    if (messages.length > 0) {
      return {
        message: messages[0],
        details: messages,
      };
    }
  } else {
    try {
      const text = normalizeText(await response.text());
      if (text !== "") {
        return { message: text, details: [text] };
      }
    } catch {
      // ignore parse errors and fallback to status message
    }
  }

  return {
    message: formatStatusMessage(response.status),
    details: [],
  };
}

export function toErrorMessage(error: unknown): string {
  if (error instanceof ApiError) {
    const source = error.details[0] || error.message || formatStatusMessage(error.status);
    return mapDomainErrorMessage(source, error.status);
  }

  if (error instanceof Error) {
    const message = normalizeText(error.message || "");
    if (message.toLowerCase() === "internal server error") {
      return formatStatusMessage(500);
    }
    return message || "Ocurrio un error inesperado.";
  }

  return "Ocurrio un error inesperado.";
}

async function parseErrorBody(response: Response): Promise<ApiErrorBody | null> {
  const contentType = response.headers.get("content-type") || "";
  if (!contentType.includes("application/json")) {
    return null;
  }

  try {
    return (await response.json()) as ApiErrorBody;
  } catch {
    return null;
  }
}

async function request<T>(
  path: string,
  init: RequestInit = {},
  requiresAuth = true,
): Promise<T> {
  const headers = new Headers(init.headers ?? {});

  if (init.body && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }

  if (!headers.has("Accept")) {
    headers.set("Accept", "application/json");
  }

  if (requiresAuth) {
    const token = getAccessToken();
    if (token) {
      headers.set("Authorization", `Bearer ${token}`);
    }
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers,
  });

  if (!response.ok) {
    const parsedError = await parseErrorResponse(response);
    const message = mapDomainErrorMessage(parsedError.message, response.status);
    throw new ApiError(response.status, message, parsedError.details);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return (await response.json()) as T;
}

async function requestBlob(path: string): Promise<Blob> {
  const headers = new Headers();
  const token = getAccessToken();

  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: "GET",
    headers,
  });

  if (!response.ok) {
    const parsedError = await parseErrorResponse(response);
    const message = mapDomainErrorMessage(parsedError.message, response.status);
    throw new ApiError(response.status, message, parsedError.details);
  }

  return response.blob();
}

export function signIn(email: string, password: string): Promise<AuthResponse> {
  return request<AuthResponse>(
    "/auth/sign-in",
    {
      method: "POST",
      body: JSON.stringify({ email, password }),
    },
    false,
  );
}

export function getMe(): Promise<User> {
  return request<User>("/auth/me");
}

export function listUsers(role?: string): Promise<ListUsersResponse> {
  const search = role ? `?role=${encodeURIComponent(role)}` : "";
  return request<ListUsersResponse>(`/users${search}`);
}

export function createUser(payload: CreateUserPayload): Promise<User> {
  return request<User>("/users", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function listPeriods(): Promise<ListPeriodsResponse> {
  return request<ListPeriodsResponse>("/periods");
}

export function createPeriod(payload: CreatePeriodPayload): Promise<Period> {
  return request<Period>("/periods", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function listWeeksByPeriod(periodId: number): Promise<ListWeeksResponse> {
  return request<ListWeeksResponse>(`/weeks/periods/${periodId}`);
}

export function listWorkspaces(periodId?: number): Promise<ListWorkspacesResponse> {
  const search = periodId ? `?period_id=${periodId}` : "";
  return request<ListWorkspacesResponse>(`/workspaces${search}`);
}

export function createWorkspace(payload: CreateWorkspacePayload): Promise<Workspace> {
  return request<Workspace>("/workspaces", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function listAssignmentsByUser(userId: number): Promise<ListAssignmentsResponse> {
  return request<ListAssignmentsResponse>(`/assignments?user_id=${userId}`);
}

export function createAssignment(payload: CreateAssignmentPayload): Promise<Assignment> {
  return request<Assignment>("/assignments", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function listTasks(): Promise<ListTasksResponse> {
  return request<ListTasksResponse>("/tasks");
}

export function createTask(payload: CreateTaskPayload): Promise<Task> {
  return request<Task>("/tasks", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function updateTask(id: number, payload: UpdateTaskPayload): Promise<Task> {
  return request<Task>(`/tasks/${id}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
}

export function deleteTask(id: number): Promise<void> {
  return request<void>(`/tasks/${id}`, { method: "DELETE" });
}

export function generateWeeklyReport(
  payload: GenerateWeeklyReportPayload,
): Promise<GenerateWeeklyReportResponse> {
  return request<GenerateWeeklyReportResponse>("/reports/weekly", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function listReports(filters: {
  workspace_id: number;
  week_id?: number;
}): Promise<ListReportsResponse> {
  const params = new URLSearchParams();
  params.set("workspace_id", String(filters.workspace_id));
  if (filters.week_id) {
    params.set("week_id", String(filters.week_id));
  }

  const query = params.toString();
  return request<ListReportsResponse>(`/reports${query ? `?${query}` : ""}`);
}

export function downloadReport(id: number): Promise<Blob> {
  return requestBlob(`/reports/${id}/download`);
}

export type {
  Assignment,
  Period,
  Task,
  User,
  Week,
  Workspace,
  CreateTaskPayload,
  UpdateTaskPayload,
};