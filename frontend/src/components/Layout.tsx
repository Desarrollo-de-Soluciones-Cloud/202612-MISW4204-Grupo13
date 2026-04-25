import type { ReactNode } from "react";
import type { User } from "../api/types";

interface LayoutProps {
  title: string;
  user: User;
  onLogout: () => void;
  children: ReactNode;
}

export default function Layout({ title, user, onLogout, children }: LayoutProps) {
  return (
    <div className="app-shell">
      <header className="topbar">
        <div>
          <h1>{title}</h1>
          <p className="muted">
            {user.name} ({user.global_role}) - {user.email}
          </p>
        </div>
        <button onClick={onLogout}>Cerrar sesion</button>
      </header>
      <main>{children}</main>
    </div>
  );
}
