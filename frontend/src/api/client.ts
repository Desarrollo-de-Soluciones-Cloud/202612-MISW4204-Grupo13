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

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:80/api";

interface ApiErrorBody {
  error?: string;
  errors?: string[];
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

export function toErrorMessage(error: unknown): string {
  if (error instanceof ApiError) {
    if (error.details.length > 0) {
      return error.details.join(" | ");
    }
    return error.message;
  }

  if (error instanceof Error) {
    return error.message;
  }

  return "unexpected error";
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
    const body = await parseErrorBody(response);
    const details = body?.errors ?? [];
    const message = body?.error ?? details[0] ?? `request failed with status ${response.status}`;
    throw new ApiError(response.status, message, details);
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
    const body = await parseErrorBody(response);
    const details = body?.errors ?? [];
    const message = body?.error ?? details[0] ?? `request failed with status ${response.status}`;
    throw new ApiError(response.status, message, details);
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

export function listReports(filters?: {
  workspace_id?: number;
  week_id?: number;
}): Promise<ListReportsResponse> {
  const params = new URLSearchParams();
  if (filters?.workspace_id) {
    params.set("workspace_id", String(filters.workspace_id));
  }
  if (filters?.week_id) {
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
