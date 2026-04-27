import type { User } from "../api/types";

const AUTH_TOKEN_KEY = "auth.access_token";
const AUTH_USER_KEY = "auth.user";

export interface AuthSession {
  accessToken: string;
  user: User;
}

export function saveAuthSession(session: AuthSession): void {
  localStorage.setItem(AUTH_TOKEN_KEY, session.accessToken);
  localStorage.setItem(AUTH_USER_KEY, JSON.stringify(session.user));
}

export function getAuthSession(): AuthSession | null {
  const accessToken = localStorage.getItem(AUTH_TOKEN_KEY);
  const userRaw = localStorage.getItem(AUTH_USER_KEY);

  if (!accessToken || !userRaw) {
    return null;
  }

  try {
    const user = JSON.parse(userRaw) as User;
    return { accessToken, user };
  } catch {
    clearAuthSession();
    return null;
  }
}

export function clearAuthSession(): void {
  localStorage.removeItem(AUTH_TOKEN_KEY);
  localStorage.removeItem(AUTH_USER_KEY);
}

export function getAccessToken(): string | null {
  return localStorage.getItem(AUTH_TOKEN_KEY);
}
