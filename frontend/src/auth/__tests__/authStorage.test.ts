import { describe, expect, it, beforeEach } from "vitest";
import {
  clearAuthSession,
  getAccessToken,
  getAuthSession,
  saveAuthSession,
} from "../authStorage";

const AUTH_TOKEN_KEY = "auth.access_token";
const AUTH_USER_KEY = "auth.user";

const user = {
  id: 7,
  name: "Ana",
  email: "ana@example.com",
  global_role: "admin" as const,
};

describe("authStorage", () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it("guarda y lee una sesión válida", () => {
    saveAuthSession({ accessToken: "token-123", user });

    expect(getAuthSession()).toEqual({ accessToken: "token-123", user });
  });

  it("limpia la sesión", () => {
    saveAuthSession({ accessToken: "token-123", user });
    clearAuthSession();

    expect(localStorage.getItem(AUTH_TOKEN_KEY)).toBeNull();
    expect(localStorage.getItem(AUTH_USER_KEY)).toBeNull();
  });

  it("retorna null cuando no hay sesión", () => {
    expect(getAuthSession()).toBeNull();
  });

  it("retorna null si localStorage tiene JSON inválido", () => {
    localStorage.setItem(AUTH_TOKEN_KEY, "token-123");
    localStorage.setItem(AUTH_USER_KEY, "{invalid json}");

    expect(getAuthSession()).toBeNull();
    expect(localStorage.getItem(AUTH_TOKEN_KEY)).toBeNull();
    expect(localStorage.getItem(AUTH_USER_KEY)).toBeNull();
  });

  it("mantiene token y usuario correctamente", () => {
    saveAuthSession({ accessToken: "token-abc", user });

    expect(getAccessToken()).toBe("token-abc");
    expect(getAuthSession()?.user).toEqual(user);
  });
});
