import { act, renderHook } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import useToast from "../useToast";

describe("useToast", () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it("inicia sin toast", () => {
    const { result } = renderHook(() => useToast());

    expect(result.current.toast).toBeNull();
  });

  it("showToast muestra mensaje y tipo", () => {
    const { result } = renderHook(() => useToast());

    act(() => {
      result.current.showToast("Guardado", "success");
    });

    expect(result.current.toast).toEqual({ message: "Guardado", type: "success" });
  });

  it("clearToast limpia el toast", () => {
    const { result } = renderHook(() => useToast());

    act(() => {
      result.current.showToast("Error", "error");
      result.current.clearToast();
    });

    expect(result.current.toast).toBeNull();
  });

  it("auto-cierre funciona con timers", () => {
    const { result } = renderHook(() => useToast(1000));

    act(() => {
      result.current.showToast("Temporal", "info");
    });

    expect(result.current.toast).toEqual({ message: "Temporal", type: "info" });

    act(() => {
      vi.advanceTimersByTime(1000);
    });

    expect(result.current.toast).toBeNull();
  });

  it("si se muestra un nuevo toast, reemplaza el anterior", () => {
    const { result } = renderHook(() => useToast());

    act(() => {
      result.current.showToast("Primero", "info");
      result.current.showToast("Segundo", "error");
    });

    expect(result.current.toast).toEqual({ message: "Segundo", type: "error" });
  });
});
