export type GlobalRole = "admin" | "professor" | "monitor" | "assistant";
export type TaskStatus = "abierto" | "en desarrollo" | "finalizado";

export interface User {
  id: number;
  name: string;
  email: string;
  global_role: GlobalRole;
}

export interface AuthResponse {
  access_token: string;
  token_type: "bearer";
  expires_in: number;
  user: User;
}

export interface Period {
  id: number;
  name: string;
  initial_date: string;
  final_date: string;
  inscription_final_date: string;
  weeks_count: number;
  period_state: "active" | "closed";
}

export interface Week {
  id: number;
  period_id: number;
  number: number;
  initial_date: string;
  final_date: string;
}

export interface Workspace {
  id: number;
  period_id: number;
  user_id: number;
  name: string;
  type: "course" | "project";
  initial_date: string;
  final_date: string;
  observations: string;
  state: "active" | "closed";
}

export interface Assignment {
  id: number;
  user_id: number;
  workspace_id: number;
  role: "monitor" | "assistant";
  weekly_hours: number;
}

export interface Task {
  id: number;
  user_id: number;
  assignment_id: number;
  week_id: number | null;
  title: string;
  description: string;
  status: TaskStatus;
  spent_hours: number;
  observations: string;
  week_start_date: string;
  late: boolean;
}

export interface Report {
  id: number;
  workspace_id: number;
  week_id: number;
  assignment_id: number;
  user_id: number;
  type: "weekly_pdf";
  summary: string;
  file_path: string;
}

export interface ListUsersResponse {
  users: User[];
}

export interface ListPeriodsResponse {
  periods: Period[];
}

export interface ListWeeksResponse {
  weeks: Week[];
}

export interface ListWorkspacesResponse {
  workspaces: Workspace[];
}

export interface ListAssignmentsResponse {
  assignments: Assignment[];
}

export interface ListTasksResponse {
  tasks: Task[];
}

export interface ListReportsResponse {
  reports: Report[];
}

export interface GenerateWeeklyReportPayload {
  workspace_id: number;
  week_id: number;
}

export interface GenerateWeeklyReportResponse {
  reports: Report[];
  generated_count: number;
}

export interface CreateTaskPayload {
  assignment_id: number;
  title: string;
  description: string;
  status: TaskStatus;
  spent_hours: number;
  observations: string;
  week_start_date: string;
}

export type UpdateTaskPayload = CreateTaskPayload;
