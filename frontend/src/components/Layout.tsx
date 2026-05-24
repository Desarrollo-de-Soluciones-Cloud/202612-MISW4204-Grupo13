import type { ReactNode } from "react";
import type { User } from "../api/types";

type LayoutProps = Readonly<{
  title: string;
  description?: string;
  user: User;
  onLogout: () => void;
  children: ReactNode;
}>;

function toRoleLabel(role: User["global_role"]): string {
  const map = {
    admin: "Administrador",
    professor: "Profesor",
    monitor: "Monitor",
    assistant: "Asistente graduado",
  } as const;

  return map[role];
}

export default function Layout({ title, description, user, onLogout, children }: LayoutProps) {
  return (
    <div className="app-shell">
      <header className="app-header card">
        <div className="app-brand">
          <p className="app-system-name">Seneprojects</p>
          <h1 className="app-page-title">{title}</h1>
          {description ? <p className="muted page-description">{description}</p> : null}
        </div>

        <div className="app-user-summary">
          <span className="role-chip">{toRoleLabel(user.global_role)}</span>
          <p className="muted">{user.name}</p>
          <p className="muted">{user.email}</p>
          <button className="button-danger button-logout" onClick={onLogout}>
            Cerrar sesión
          </button>
        </div>
      </header>

      <main>{children}</main>
    </div>
  );
}
