import { useEffect } from "react";

export default function Notification({
  type = "info",
  message,
  onClose,
  duration = 5000,
}) {
  useEffect(() => {
    if (!message || !duration) {
      return undefined;
    }
    const timer = setTimeout(() => {
      onClose?.();
    }, duration);
    return () => clearTimeout(timer);
  }, [message, duration, onClose]);

  if (!message) {
    return null;
  }

  return (
    <div className={`notification notification--${type}`}>
      <span>{message}</span>
      <button type="button" onClick={onClose} aria-label="Dismiss notification">
        ×
      </button>
    </div>
  );
}

