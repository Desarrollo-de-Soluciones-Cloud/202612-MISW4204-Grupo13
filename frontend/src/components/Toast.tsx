import type { ToastType } from "./useToast";

type ToastProps = Readonly<{
  type: ToastType;
  message: string;
  onClose: () => void;
}>;

export default function Toast({ type, message, onClose }: ToastProps) {
  return (
    <output className="toast-container" aria-live="polite">
      <span className={`toast toast-${type}`}>
        <span className="toast-message">{message}</span>
        <button
          type="button"
          className="toast-close"
          aria-label="Cerrar notificación"
          onClick={onClose}
        >
          x
        </button>
      </span>
    </output>
  );
}
