interface HelpTextProps {
  children: string;
}

export default function HelpText({ children }: HelpTextProps) {
  return (
    <p className="help-text">
      <span className="help-icon">ⓘ</span>
      {children}
    </p>
  );
}
