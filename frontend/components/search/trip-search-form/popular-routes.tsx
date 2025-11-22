"use client";

import { Button } from "@/components/ui/button";
import { POPULAR_ROUTES } from "./constants";

type PopularRoutesProps = {
  onSelectRoute: (route: { from: string; to: string }) => void;
};

export function PopularRoutes({ onSelectRoute }: PopularRoutesProps) {
  return (
    <div className="border-t pt-6">
      <p className="mb-3 text-sm font-medium text-muted-foreground">
        Tuyến đường phổ biến:
      </p>
      <div className="flex flex-wrap gap-2">
        {POPULAR_ROUTES.map((route) => (
          <Button
            key={route.id}
            type="button"
            variant="outline"
            size="sm"
            onClick={() => onSelectRoute(route)}
            className="text-xs"
          >
            {route.from} → {route.to}
          </Button>
        ))}
      </div>
    </div>
  );
}
