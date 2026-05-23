type EmptyStateProps = Readonly<{
  colSpan: number;
  message?: string;
}>;

export default function EmptyState({
  colSpan,
  message = "No hay registros disponibles todavía.",
}: EmptyStateProps) {
  return (
    <tr>
      <td colSpan={colSpan} className="empty-state-cell">
        {message}
      </td>
    </tr>
  );
}
