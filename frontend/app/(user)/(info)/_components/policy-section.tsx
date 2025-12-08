import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ReactNode } from "react";
import { LucideIcon } from "lucide-react";

interface PolicySectionProps {
  title: string;
  children: ReactNode;
  id?: string;
  icon?: LucideIcon;
}

export function PolicySection({
  title,
  children,
  id,
  icon: Icon,
}: PolicySectionProps) {
  return (
    <section id={id} className="scroll-mt-20">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-xl">
            {Icon && <Icon className="h-6 w-6 text-primary" />}
            {title}
          </CardTitle>
        </CardHeader>
        <CardContent className="prose prose-neutral dark:prose-invert max-w-none pt-0">
          {children}
        </CardContent>
      </Card>
    </section>
  );
}
