import { FormEvent, useState } from "react";
import { signIn, toErrorMessage } from "../api/client";
import type { AuthResponse } from "../api/types";
import ErrorMessage from "../components/ErrorMessage";
import Loading from "../components/Loading";

interface LoginPageProps {
  onSignedIn: (auth: AuthResponse) => void;
}

export default function LoginPage({ onSignedIn }: LoginPageProps) {
  const [email, setEmail] = useState("admin@example.com");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const response = await signIn(email, password);
      onSignedIn(response);
    } catch (err) {
      setError(toErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-wrapper">
      <form className="card" onSubmit={handleSubmit}>
        <h1>Ingreso</h1>
        <p className="muted">Frontend minimo para Entrega 2</p>

        <label>
          Email
          <input
            type="email"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            required
          />
        </label>

        <label>
          Password
          <input
            type="password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            required
          />
        </label>

        <button type="submit" disabled={loading}>
          Ingresar
        </button>

        {loading && <Loading label="Autenticando..." />}
        <ErrorMessage message={error} />
      </form>
    </div>
  );
}
