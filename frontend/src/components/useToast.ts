import { useCallback, useEffect, useRef, useState } from "react";

export type ToastType = "success" | "error" | "info";

export interface ToastState {
  message: string;
  type: ToastType;
}

const DEFAULT_DURATION_MS = 3000;

export default function useToast(durationMs: number = DEFAULT_DURATION_MS) {
  const [toast, setToast] = useState<ToastState | null>(null);
  const timeoutRef = useRef<number | null>(null);

  const clearTimeoutRef = useCallback(() => {
    if (timeoutRef.current !== null) {
      window.clearTimeout(timeoutRef.current);
      timeoutRef.current = null;
    }
  }, []);

  const clearToast = useCallback(() => {
    clearTimeoutRef();
    setToast(null);
  }, [clearTimeoutRef]);

  const showToast = useCallback(
    (message: string, type: ToastType = "info") => {
      clearTimeoutRef();
      setToast({ message, type });
      timeoutRef.current = window.setTimeout(() => {
        setToast(null);
        timeoutRef.current = null;
      }, durationMs);
    },
    [clearTimeoutRef, durationMs],
  );

  useEffect(() => {
    return () => {
      clearTimeoutRef();
    };
  }, [clearTimeoutRef]);

  return { toast, showToast, clearToast };
}
