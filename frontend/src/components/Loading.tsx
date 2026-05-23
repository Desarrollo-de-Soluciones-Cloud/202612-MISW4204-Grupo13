type LoadingProps = Readonly<{
  label?: string;
}>;

export default function Loading({ label = "Cargando..." }: LoadingProps) {
  return <p className="muted">{label}</p>;
}
