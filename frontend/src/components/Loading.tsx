interface LoadingProps {
  label?: string;
}

export default function Loading({ label = "Cargando..." }: LoadingProps) {
  return <p className="muted">{label}</p>;
}
