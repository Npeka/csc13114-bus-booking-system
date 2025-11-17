"use client";

import { Button } from "@/components/ui/button";
import { ArrowRightLeft } from "lucide-react";

type SwapLocationsButtonProps = {
  onClick: () => void;
};

export function SwapLocationsButton({ onClick }: SwapLocationsButtonProps) {
  return (
    <Button
      type="button"
      variant="secondary"
      size="icon-sm"
      onClick={onClick}
      className="group hidden lg:flex absolute left-1/2 top-1/2 z-10 h-10 w-10 -translate-x-1/2 -translate-y-1/2 rounded-full border bg-background shadow-elevated transition-transform duration-300 hover:scale-105"
      aria-label="Đổi điểm đi đến"
    >
      <ArrowRightLeft className="h-3.5 w-3.5 transition-transform duration-300 group-hover:rotate-180" />
    </Button>
  );
}

