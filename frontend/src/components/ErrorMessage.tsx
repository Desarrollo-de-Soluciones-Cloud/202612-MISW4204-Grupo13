type ErrorMessageProps = Readonly<{
  message?: string | null;
}>;

export default function ErrorMessage({ message }: ErrorMessageProps) {
  if (!message) {
    return null;
  }

  return <div className="error-box">{message}</div>;
}
