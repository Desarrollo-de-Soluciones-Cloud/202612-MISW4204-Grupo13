import type { ReactNode } from "react";

interface HelpTextProps {
  children: ReactNode;
}

export default function HelpText({ children }: HelpTextProps) {
  return (
    <p className="help-text">
      <span className="help-icon">ⓘ</span>
      {children}
    </p>
  );
}
