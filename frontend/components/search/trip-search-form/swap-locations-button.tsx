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
      className="group absolute top-1/2 left-1/2 z-10 hidden size-9 -translate-x-1/2 -translate-y-1/2 rounded-full shadow-elevated transition-transform duration-300 hover:scale-105 lg:flex"
      aria-label="Đổi điểm đi đến"
    >
      <ArrowRightLeft className="h-3.5 w-3.5 transition-transform duration-300 group-hover:rotate-180" />
    </Button>
  );
}
