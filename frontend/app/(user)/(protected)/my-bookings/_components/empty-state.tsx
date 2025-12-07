import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import Link from "next/link";

interface EmptyStateProps {
  title: string;
  description?: string;
  showAction?: boolean;
}

export function EmptyState({
  title,
  description,
  showAction = false,
}: EmptyStateProps) {
  return (
    <Card>
      <CardContent className="py-12 text-center">
        <p className="text-muted-foreground">{title}</p>
        {description && (
          <p className="mt-2 text-sm text-muted-foreground">{description}</p>
        )}
        {showAction && (
          <Button asChild className="mt-4">
            <Link href="/">Đặt vé ngay</Link>
          </Button>
        )}
      </CardContent>
    </Card>
  );
}
