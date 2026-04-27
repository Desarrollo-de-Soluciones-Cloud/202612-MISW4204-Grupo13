import type { ToastType } from "./useToast";

interface ToastProps {
  type: ToastType;
  message: string;
  onClose: () => void;
}

export default function Toast({ type, message, onClose }: ToastProps) {
  return (
    <div className="toast-container" role="status" aria-live="polite">
      <div className={`toast toast-${type}`}>
        <p className="toast-message">{message}</p>
        <button
          type="button"
          className="toast-close"
          aria-label="Cerrar notificación"
          onClick={onClose}
        >
          x
        </button>
      </div>
    </div>
  );
}
