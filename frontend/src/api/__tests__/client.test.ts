import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import {
  ApiError,
  createTask,
  createUser,
  deleteTask,
  downloadReport,
  downloadTaskAttachment,
  generateWeeklyReport,
  listReports,
  listUsers,
  listWorkspaces,
  signIn,
  toErrorMessage,
  updateTask,
} from "../client";

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: {
      "content-type": "application/json",
    },
  });
}

describe("api/client", () => {
  const fetchMock = vi.fn();

  beforeEach(() => {
    localStorage.clear();
    fetchMock.mockReset();
    vi.stubGlobal("fetch", fetchMock);
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("agrega Authorization Bearer cuando hay token guardado", async () => {
    localStorage.setItem("auth.access_token", "token-xyz");
    fetchMock.mockResolvedValue(jsonResponse({ users: [] }));

    await listUsers();

    const options = fetchMock.mock.calls[0][1] as RequestInit;
    const headers = options.headers as Headers;
    expect(headers.get("Authorization")).toBe("Bearer token-xyz");
  });

  it("no agrega Authorization cuando no hay sesión", async () => {
    fetchMock.mockResolvedValue(jsonResponse({ users: [] }));

    await listUsers();

    const options = fetchMock.mock.calls[0][1] as RequestInit;
    const headers = options.headers as Headers;
    expect(headers.get("Authorization")).toBeNull();
  });

  it("maneja respuesta JSON exitosa", async () => {
    fetchMock.mockResolvedValue(
      jsonResponse({
        access_token: "token-1",
        token_type: "bearer",
        expires_in: 3600,
        user: {
          id: 1,
          name: "Admin",
          email: "admin@example.com",
          global_role: "admin",
        },
      }),
    );

    const result = await signIn("admin@example.com", "secret");
    expect(result.access_token).toBe("token-1");
    expect(fetchMock.mock.calls[0][0]).toBe("/api/auth/sign-in");
  });

  it("maneja respuesta 204 correctamente", async () => {
    fetchMock.mockResolvedValue(new Response(null, { status: 204 }));

    const result = await deleteTask(10);
    expect(result).toBeUndefined();

    const [url, options] = fetchMock.mock.calls[0] as [string, RequestInit];
    expect(url).toBe("/api/tasks/10");
    expect(options.method).toBe("DELETE");
  });

  it("maneja error HTTP con mensaje JSON", async () => {
    fetchMock.mockResolvedValue(
      jsonResponse({
        message: "title is required",
      }, 400),
    );

    try {
      await listUsers();
      throw new Error("Expected listUsers to fail");
    } catch (error) {
      expect(error).toBeInstanceOf(ApiError);
      expect((error as Error).message).toBe("Completa titulo, descripcion, semana, estado y horas.");
    }
  });

  it("maneja error HTTP con arreglo de detalles", async () => {
    fetchMock.mockResolvedValue(
      jsonResponse(
        {
          details: ["assignment not found", "week is not active"],
        },
        404,
      ),
    );

    try {
      await listUsers();
      throw new Error("Expected listUsers to fail");
    } catch (error) {
      expect(error).toBeInstanceOf(ApiError);
      expect((error as ApiError).details).toEqual(["assignment not found", "week is not active"]);
      expect((error as Error).message).toBe("La vinculacion seleccionada no existe.");
    }
  });

  it("maneja error HTTP con arreglo de errores estructurados", async () => {
    fetchMock.mockResolvedValue(
      jsonResponse(
        {
          errors: [{ field: "week_start_date" }, { message: "spent_hours is required" }],
        },
        400,
      ),
    );

    await expect(listUsers()).rejects.toThrow("Completa titulo, descripcion, semana, estado y horas.");
  });

  it("maneja error HTTP con texto plano", async () => {
    fetchMock.mockResolvedValue(
      new Response("internal   server    error", {
        status: 500,
        headers: {
          "content-type": "text/plain",
        },
      }),
    );

    await expect(listUsers()).rejects.toThrow(
      "Ocurrio un error interno. Revisa que los datos seleccionados existan y esten relacionados correctamente.",
    );
  });

  it("normaliza mensajes de error", () => {
    expect(toErrorMessage(new Error("  mensaje     con   espacios  "))).toBe(
      "mensaje con espacios",
    );
  });

  it("listUsers y listReports usan endpoints esperados", async () => {
    fetchMock.mockResolvedValue(jsonResponse({ users: [] }));
    await listUsers("monitor");

    expect(fetchMock.mock.calls[0][0]).toBe("/api/users?role=monitor");

    fetchMock.mockResolvedValue(jsonResponse({ reports: [] }));
    await listReports({ workspace_id: 7, week_id: 3 });

    expect(fetchMock.mock.calls[1][0]).toBe("/api/reports?workspace_id=7&week_id=3");
  });

  it("listWorkspaces agrega period_id como query param solo cuando existe", async () => {
    fetchMock
      .mockResolvedValueOnce(jsonResponse({ workspaces: [] }))
      .mockResolvedValueOnce(jsonResponse({ workspaces: [] }));

    await listWorkspaces(12);
    await listWorkspaces();

    expect(fetchMock.mock.calls[0][0]).toBe("/api/workspaces?period_id=12");
    expect(fetchMock.mock.calls[1][0]).toBe("/api/workspaces");
  });

  it("createUser envia POST con body JSON y headers esperados", async () => {
    localStorage.setItem("auth.access_token", "token-xyz");
    fetchMock.mockResolvedValue(
      jsonResponse({
        id: 20,
        name: "Neo",
        email: "neo@example.com",
        global_role: "monitor",
      }),
    );

    await createUser({
      name: "Neo",
      email: "neo@example.com",
      password: "abc123",
      global_role: "monitor",
    });

    const [url, options] = fetchMock.mock.calls[0] as [string, RequestInit];
    const headers = options.headers as Headers;

    expect(url).toBe("/api/users");
    expect(options.method).toBe("POST");
    expect(headers.get("Authorization")).toBe("Bearer token-xyz");
    expect(headers.get("Content-Type")).toBe("application/json");
    expect(options.body).toBe(
      JSON.stringify({
        name: "Neo",
        email: "neo@example.com",
        password: "abc123",
        global_role: "monitor",
      }),
    );
  });

  it("updateTask envia PUT con payload JSON cuando no hay adjuntos", async () => {
    fetchMock.mockResolvedValue(
      jsonResponse({
        id: 12,
        user_id: 1,
        assignment_id: 3,
        week_id: null,
        title: "Titulo",
        description: "Desc",
        status: "abierto",
        spent_hours: 2,
        observations: "Obs",
        week_start_date: "2026-05-18",
        late: false,
        attachments: [],
      }),
    );

    await updateTask(12, {
      assignment_id: 3,
      title: "Titulo",
      description: "Desc",
      status: "abierto",
      spent_hours: 2,
      observations: "Obs",
      week_start_date: "2026-05-18",
      attachments: [],
      existing_attachments: [],
    });

    const [url, options] = fetchMock.mock.calls[0] as [string, RequestInit];
    const headers = options.headers as Headers;

    expect(url).toBe("/api/tasks/12");
    expect(options.method).toBe("PUT");
    expect(headers.get("Content-Type")).toBe("application/json");
    expect(options.body).toBe(
      JSON.stringify({
        assignment_id: 3,
        title: "Titulo",
        description: "Desc",
        status: "abierto",
        spent_hours: 2,
        observations: "Obs",
        week_start_date: "2026-05-18",
        attachments: [],
        existing_attachments: [],
      }),
    );
  });

  it("createTask envia FormData cuando hay adjuntos", async () => {
    const file = new File(["contenido"], "soporte.txt", { type: "text/plain" });
    fetchMock.mockResolvedValue(
      jsonResponse({
        id: 40,
        user_id: 1,
        assignment_id: 3,
        week_id: null,
        title: "Titulo",
        description: "Desc",
        status: "abierto",
        spent_hours: 2,
        observations: "Obs",
        week_start_date: "2026-05-18",
        late: false,
        attachments: [],
      }),
    );

    await createTask({
      assignment_id: 3,
      title: "Titulo",
      description: "Desc",
      status: "abierto",
      spent_hours: 2,
      observations: "Obs",
      week_start_date: "2026-05-18",
      attachments: [file],
      existing_attachments: [
        {
          id: "old-1",
          name: "previo.pdf",
          file_path: "tasks/previo.pdf",
          content_type: "application/pdf",
          size: 10,
        },
      ],
    });

    const [url, options] = fetchMock.mock.calls[0] as [string, RequestInit];
    const headers = options.headers as Headers;
    const body = options.body as FormData;

    expect(url).toBe("/api/tasks");
    expect(options.method).toBe("POST");
    expect(headers.get("Content-Type")).toBeNull();
    expect(body.get("assignment_id")).toBe("3");
    expect(body.get("title")).toBe("Titulo");
    expect(body.get("attachments")).toBe(file);
    expect(body.get("existing_attachments")).toBe(
      JSON.stringify([
        {
          id: "old-1",
          name: "previo.pdf",
          file_path: "tasks/previo.pdf",
          content_type: "application/pdf",
          size: 10,
        },
      ]),
    );
  });

  it("generateWeeklyReport envia POST con body JSON y retorna respuesta", async () => {
    fetchMock.mockResolvedValue(jsonResponse({ reports: [], generated_count: 0 }));

    const result = await generateWeeklyReport({ workspace_id: 8, week_id: 2 });

    const [url, options] = fetchMock.mock.calls[0] as [string, RequestInit];
    expect(url).toBe("/api/reports/weekly");
    expect(options.method).toBe("POST");
    expect(options.body).toBe(JSON.stringify({ workspace_id: 8, week_id: 2 }));
    expect(result.generated_count).toBe(0);
  });

  it("downloadReport usa GET y retorna blob", async () => {
    const expectedBlob = new Blob(["pdf"], { type: "application/pdf" });
    localStorage.setItem("auth.access_token", "token-abc");
    fetchMock.mockResolvedValue(
      new Response(expectedBlob, {
        status: 200,
        headers: {
          "content-type": "application/pdf",
        },
      }),
    );

    const result = await downloadReport(55);

    const [url, options] = fetchMock.mock.calls[0] as [string, RequestInit];
    const headers = options.headers as Headers;

    expect(url).toBe("/api/reports/55/download");
    expect(options.method).toBe("GET");
    expect(headers.get("Authorization")).toBe("Bearer token-abc");
    expect(typeof result.arrayBuffer).toBe("function");
    expect(result.size).toBeGreaterThan(0);
    expect(result.type).toBe("application/pdf");
  });

  it("downloadTaskAttachment usa endpoint de descarga de adjuntos", async () => {
    fetchMock.mockResolvedValue(new Response(new Blob(["file"]), { status: 200 }));

    const result = await downloadTaskAttachment(44, "att-9");

    const [url, options] = fetchMock.mock.calls[0] as [string, RequestInit];
    expect(url).toBe("/api/tasks/44/attachments/att-9/download");
    expect(options.method).toBe("GET");
    expect(result.size).toBeGreaterThan(0);
  });
});
