import { ReactNode } from "react";

interface InfoLayoutProps {
  children: ReactNode;
}

export default function InfoLayout({ children }: InfoLayoutProps) {
  return (
    <div className="min-h-screen bg-secondary/30">
      <div className="container py-8">{children}</div>
    </div>
  );
}
