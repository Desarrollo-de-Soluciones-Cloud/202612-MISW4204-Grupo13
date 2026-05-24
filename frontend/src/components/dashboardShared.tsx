import type { ReactNode } from "react";
import type { User, Week, Workspace } from "../api/types";
import HelpText from "./HelpText";

type SummaryTableProps = Readonly<{
  user: User | null;
}>;

type WorkspaceDetailsFieldsProps = Readonly<{
  form: {
    name: string;
    type: "course" | "project";
    initial_date: string;
    final_date: string;
    observations: string;
    state: "active" | "closed";
  };
  onChange: (updates: Partial<WorkspaceDetailsFieldsProps["form"]>) => void;
  observationsHelpText?: ReactNode;
  finalDateHelpText?: ReactNode;
}>;

type WorkspaceWeekFiltersProps = Readonly<{
  workspaceId: string;
  weekId: string;
  workspaces: Workspace[];
  weeks: Week[];
  workspaceLabel: string;
  weekLabel: string;
  workspacePlaceholder: string;
  weekPlaceholder: string;
  allWeeksLabel?: string;
  disabled?: boolean;
  onWorkspaceChange: (value: string) => void;
  onWeekChange: (value: string) => void;
  getWorkspaceLabel: (workspace: Workspace) => string;
  getWeekLabel: (week: Week) => string;
}>;

export function SessionSummaryTable({ user }: SummaryTableProps) {
  if (!user) {
    return <p className="muted">Sin datos</p>;
  }

  return (
    <table>
      <tbody>
        <tr>
          <th>ID</th>
          <td>{user.id}</td>
        </tr>
        <tr>
          <th>Nombre</th>
          <td>{user.name}</td>
        </tr>
        <tr>
          <th>Correo</th>
          <td>{user.email}</td>
        </tr>
        <tr>
          <th>Rol</th>
          <td>{user.global_role}</td>
        </tr>
      </tbody>
    </table>
  );
}

export function WorkspaceDetailsFields({
  form,
  onChange,
  observationsHelpText,
  finalDateHelpText,
}: WorkspaceDetailsFieldsProps) {
  return (
    <>
      <div className="form-field">
        <label>
          <span>Nombre</span>
          <input
            value={form.name}
            onChange={(event) => onChange({ name: event.target.value })}
            required
          />
        </label>
      </div>

      <div className="form-field">
        <label>
          <span>Tipo</span>
          <select
            value={form.type}
            onChange={(event) =>
              onChange({ type: event.target.value as "course" | "project" })
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
            value={form.initial_date}
            onChange={(event) => onChange({ initial_date: event.target.value })}
            required
          />
        </label>
      </div>

      <div className="form-field">
        <label>
          <span>Fecha final</span>
          <input
            type="date"
            value={form.final_date}
            onChange={(event) => onChange({ final_date: event.target.value })}
            required
          />
        </label>
        {finalDateHelpText ? <HelpText>{finalDateHelpText}</HelpText> : null}
      </div>

      <div className="form-field">
        <label>
          <span>Observaciones</span>
          <input
            value={form.observations}
            onChange={(event) => onChange({ observations: event.target.value })}
            required
          />
        </label>
        {observationsHelpText ? <HelpText>{observationsHelpText}</HelpText> : null}
      </div>

      <div className="form-field">
        <label>
          <span>Estado</span>
          <select
            value={form.state}
            onChange={(event) =>
              onChange({ state: event.target.value as "active" | "closed" })
            }
          >
            <option value="active">Activo</option>
            <option value="closed">Cerrado</option>
          </select>
        </label>
      </div>
    </>
  );
}

export function WorkspaceWeekFilters({
  workspaceId,
  weekId,
  workspaces,
  weeks,
  workspaceLabel,
  weekLabel,
  workspacePlaceholder,
  weekPlaceholder,
  allWeeksLabel,
  disabled,
  onWorkspaceChange,
  onWeekChange,
  getWorkspaceLabel,
  getWeekLabel,
}: WorkspaceWeekFiltersProps) {
  return (
    <>
      <div className="form-field">
        <label>
          <span>{workspaceLabel}</span>
          <select
            value={workspaceId}
            onChange={(event) => onWorkspaceChange(event.target.value)}
            required
          >
            <option value="" disabled>
              {workspacePlaceholder}
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
          <span>{weekLabel}</span>
          <select
            value={weekId}
            onChange={(event) => onWeekChange(event.target.value)}
            disabled={disabled}
          >
            {allWeeksLabel ? <option value="">{allWeeksLabel}</option> : null}
            <option value="" disabled={!allWeeksLabel}>
              {weekPlaceholder}
            </option>
            {weeks.map((item) => (
              <option key={item.id} value={item.id}>
                {getWeekLabel(item)}
              </option>
            ))}
          </select>
        </label>
      </div>
    </>
  );
}
