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
  const normalized = value.trim();
  const withReplaceAll = normalized as string & {
    replaceAll?: (searchValue: RegExp | string, replaceValue: string) => string;
  };

  return withReplaceAll.replaceAll
    ? withReplaceAll.replaceAll(/\s+/g, " ")
    : normalized.replace(/\s+/g, " ");
}

function includesAny(value: string, terms: readonly string[]): boolean {
  return terms.some((term) => value.includes(term));
}

function includesEntityNotFound(value: string, entity: string): boolean {
  return value.includes(entity) && includesAny(value, ["not found", "does not exist"]);
}

function includesTaskValidationFields(value: string): boolean {
  return includesAny(value, [
    "title",
    "description",
    "week_start_date",
    "status",
    "spent_hours",
    "binding",
    "required",
  ]);
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

type DomainMessageRule = {
  matches: (normalized: string, status: number) => boolean;
  message: string | ((message: string, status: number) => string);
};

const DOMAIN_MESSAGE_RULES: DomainMessageRule[] = [
  {
    matches: (normalized) => includesEntityNotFound(normalized, "week"),
    message: "La semana seleccionada no existe.",
  },
  {
    matches: (normalized) => includesEntityNotFound(normalized, "assignment"),
    message: "La vinculacion seleccionada no existe.",
  },
  {
    matches: (normalized) => normalized.includes("workspace") && normalized.includes("closed"),
    message: "No se pueden crear tareas en un curso o proyecto cerrado.",
  },
  {
    matches: (normalized) =>
      includesAny(normalized, [
        "week is not active",
        "late update forbidden",
        "cannot be updated",
        "cannot be modified",
      ]),
    message: "La tarea no se puede modificar porque la semana ya no esta activa.",
  },
  {
    matches: (normalized, status) => status === 400 && includesTaskValidationFields(normalized),
    message: "Completa titulo, descripcion, semana, estado y horas.",
  },
  {
    matches: (normalized) => includesAny(normalized, ["40%", "40 percent", "forty"]),
    message: "La vinculacion no cumple la regla del 40% entre horas de monitor y asistente.",
  },
  {
    matches: (normalized) => normalized.includes("already exists") && normalized.includes("assignment"),
    message: "Ya existe una vinculacion equivalente para ese usuario y curso/proyecto.",
  },
  {
    matches: (normalized) =>
      includesAny(normalized, [
        "monitor weekly hours cannot exceed",
        "assistant weekly hours cannot exceed",
      ]),
    message: (message) => message,
  },
  {
    matches: (normalized) => normalized.includes("period") && normalized.includes("closed"),
    message: "No se puede crear un curso/proyecto en un periodo cerrado.",
  },
  {
    matches: (normalized) => normalized.includes("professor can only create workspaces for themselves"),
    message: "Como profesor, solo puedes crear cursos/proyectos asociados a tu propia cuenta.",
  },
  {
    matches: (normalized) => normalized.includes("professor can only create assignments in their own workspaces"),
    message: "Como profesor, solo puedes crear vinculaciones en cursos/proyectos que te pertenecen.",
  },
  {
    matches: (normalized) => normalized.includes("professor") && includesAny(normalized, ["not found", "invalid"]),
    message: "El profesor seleccionado no existe o no tiene rol de profesor.",
  },
  {
    matches: (normalized) => normalized.includes("date") && normalized.includes("period"),
    message: "Revisa las fechas del periodo academico.",
  },
  {
    matches: (normalized) => includesAny(normalized, ["weeks_count", "weeks count"]),
    message: "Revisa la cantidad de semanas del periodo academico.",
  },
  {
    matches: (normalized) => normalized.includes("no tasks") && normalized.includes("week"),
    message: "No hay tareas reportadas para ese curso/proyecto y semana.",
  },
  {
    matches: (normalized) => includesAny(normalized, ["report workspace not found", "workspace not found"]),
    message: "No se encontro el curso/proyecto seleccionado para generar el reporte.",
  },
  {
    matches: (normalized) => normalized.includes("report week not found") || (normalized.includes("week") && normalized.includes("report")),
    message: "No se encontro la semana seleccionada para generar el reporte.",
  },
  {
    matches: (normalized) => includesAny(normalized, ["ai report generation failed", "ai generation failed"]),
    message: "No fue posible generar el resumen con IA. Revisa la configuracion del servicio.",
  },
  {
    matches: (normalized) => includesAny(normalized, ["pdf report generation failed", "pdf generation failed"]),
    message: "No fue posible generar el PDF del reporte.",
  },
  {
    matches: (normalized) => normalized.includes("report file not found") || normalized.includes("download"),
    message: "No se pudo descargar el PDF. Verifica que el archivo exista.",
  },
  {
    matches: (normalized) => normalized === "internal server error",
    message: (_message, status) => formatStatusMessage(status),
  },
];

function mapDomainErrorMessage(rawMessage: string, status: number): string {
  const message = normalizeText(rawMessage);
  const normalized = message.toLowerCase();

  for (const rule of DOMAIN_MESSAGE_RULES) {
    if (!rule.matches(normalized, status)) {
      continue;
    }

    return typeof rule.message === "string" ? rule.message : rule.message(message, status);
  }

  return message;
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null;
}

function getStringValue(source: Record<string, unknown>, key: string): string | null {
  const value = source[key];
  if (typeof value !== "string") {
    return null;
  }

  const normalized = normalizeText(value);
  return normalized === "" ? null : normalized;
}

function extractMessageFromObject(value: unknown): string | null {
  if (!isRecord(value)) {
    return null;
  }

  return getStringValue(value, "message") ?? getStringValue(value, "error") ?? getStringValue(value, "field");
}

function extractMessagesFromDetails(details: ApiErrorBody["details"]): string[] {
  if (typeof details === "string") {
    const normalized = normalizeText(details);
    return normalized ? [normalized] : [];
  }

  if (!Array.isArray(details)) {
    return [];
  }

  return details
    .filter((detail): detail is string => typeof detail === "string")
    .map((detail) => normalizeText(detail))
    .filter((detail) => detail !== "");
}

function extractMessagesFromErrors(errors: ApiErrorBody["errors"]): string[] {
  if (!Array.isArray(errors)) {
    return [];
  }

  return errors
    .map((item) => {
      if (typeof item === "string") {
        return normalizeText(item);
      }

      return extractMessageFromObject(item);
    })
    .filter((message): message is string => Boolean(message));
}

function extractApiMessages(body: ApiErrorBody | null): string[] {
  if (!body) {
    return [];
  }

  const topLevelMessages = [body.error, body.message]
    .filter((value): value is string => typeof value === "string")
    .map((value) => normalizeText(value))
    .filter((value) => value !== "");

  return [
    ...topLevelMessages,
    ...extractMessagesFromDetails(body.details),
    ...extractMessagesFromErrors(body.errors),
  ];
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

  if (init.body && !(init.body instanceof FormData) && !headers.has("Content-Type")) {
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
  if (payload.attachments?.length || payload.existing_attachments?.length) {
    const formData = buildTaskFormData(payload);
    return request<Task>("/tasks", {
      method: "POST",
      body: formData,
    });
  }

  return request<Task>("/tasks", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function updateTask(id: number, payload: UpdateTaskPayload): Promise<Task> {
  if (payload.attachments?.length || payload.existing_attachments?.length) {
    const formData = buildTaskFormData(payload);
    return request<Task>(`/tasks/${id}`, {
      method: "PUT",
      body: formData,
    });
  }

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
  const suffix = query ? `?${query}` : "";
  return request<ListReportsResponse>(`/reports${suffix}`);
}

export function downloadReport(id: number): Promise<Blob> {
  return requestBlob(`/reports/${id}/download`);
}

export function downloadTaskAttachment(taskId: number, attachmentId: string): Promise<Blob> {
  return requestBlob(`/tasks/${taskId}/attachments/${attachmentId}/download`);
}

function buildTaskFormData(payload: CreateTaskPayload): FormData {
  const formData = new FormData();
  formData.append("assignment_id", String(payload.assignment_id));
  formData.append("title", payload.title);
  formData.append("description", payload.description);
  formData.append("status", payload.status);
  formData.append("spent_hours", String(payload.spent_hours));
  formData.append("observations", payload.observations);
  formData.append("week_start_date", payload.week_start_date);

  for (const file of payload.attachments ?? []) {
    formData.append("attachments", file);
  }

  if (payload.existing_attachments?.length) {
    formData.append("existing_attachments", JSON.stringify(payload.existing_attachments));
  }

  return formData;
}

export type {
  Assignment,
  Period,
  Task,
  TaskAttachment,
  User,
  Week,
  Workspace,
  CreateTaskPayload,
  UpdateTaskPayload,
} from "./types";
