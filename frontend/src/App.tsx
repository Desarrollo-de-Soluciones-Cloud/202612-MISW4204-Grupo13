import { useMemo, useState } from "react";
import type { AuthResponse } from "./api/types";
import { clearAuthSession, getAuthSession, saveAuthSession } from "./auth/authStorage";
import AdminDashboard from "./pages/AdminDashboard";
import LoginPage from "./pages/LoginPage";
import ProfessorDashboard from "./pages/ProfessorDashboard";
import WorkerDashboard from "./pages/WorkerDashboard";

export default function App() {
  const initialSession = useMemo(() => getAuthSession(), []);
  const [session, setSession] = useState(initialSession);

  const handleSignedIn = (auth: AuthResponse) => {
    const nextSession = {
      accessToken: auth.access_token,
      user: auth.user,
    };

    saveAuthSession(nextSession);
    setSession(nextSession);
  };

  const handleLogout = () => {
    clearAuthSession();
    setSession(null);
  };

  if (!session) {
    return <LoginPage onSignedIn={handleSignedIn} />;
  }

  if (session.user.global_role === "admin") {
    return <AdminDashboard user={session.user} onLogout={handleLogout} />;
  }

  if (session.user.global_role === "professor") {
    return <ProfessorDashboard user={session.user} onLogout={handleLogout} />;
  }

  return <WorkerDashboard user={session.user} onLogout={handleLogout} />;
}
