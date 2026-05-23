import { FormEvent, useState } from "react";
import { signIn, toErrorMessage } from "../api/client";
import type { AuthResponse } from "../api/types";
import Loading from "../components/Loading";
import Toast from "../components/Toast";
import useToast from "../components/useToast";

type LoginPageProps = Readonly<{
  onSignedIn: (auth: AuthResponse) => void;
}>;

export default function LoginPage({ onSignedIn }: LoginPageProps) {
  const [email, setEmail] = useState("admin@example.com");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const { toast, showToast, clearToast } = useToast();

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setLoading(true);
    clearToast();

    try {
      const response = await signIn(email, password);
      onSignedIn(response);
    } catch (err) {
      showToast(toErrorMessage(err), "error");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-wrapper">
      <form className="card login-card" onSubmit={handleSubmit}>
        <p className="app-system-name">Seneprojects</p>
        <h1>Iniciar sesión</h1>
        <p className="muted">
          Plataforma para seguimiento semanal de monitores y asistentes graduados.
        </p>
        <p className="hint-note">Usa las credenciales asignadas por el administrador.</p>

        <div className="form-field">
          <label>
            <span>Correo electrónico</span>
            <input
              type="email"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              required
            />
          </label>
        </div>

        <div className="form-field">
          <label>
            <span>Contraseña</span>
            <input
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              required
            />
          </label>
        </div>

        <div className="form-actions login-actions">
          <button type="submit" disabled={loading}>
            Entrar al sistema
          </button>
        </div>

        {loading && <Loading label="Autenticando..." />}
        {toast ? <Toast type={toast.type} message={toast.message} onClose={clearToast} /> : null}
      </form>
    </div>
  );
}

