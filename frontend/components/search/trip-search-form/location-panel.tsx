"use client";

import { RefObject } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

export type LocationPanelProps = {
  searchInputRef: RefObject<HTMLInputElement | null>;
  searchValue: string;
  onSearchChange: (value: string) => void;
  locations: string[];
  onSelect: (city: string) => void;
  recentLocations: string[];
};

export function LocationPanel({
  searchInputRef,
  searchValue,
  onSearchChange,
  locations,
  onSelect,
  recentLocations,
}: LocationPanelProps) {
  return (
    <div
      data-location-panel
      className="absolute left-1/2 top-0 z-40 w-[min(28rem,calc(100vw-2rem))] -translate-x-1/2 -translate-y-4 rounded-2xl border bg-white p-4 shadow-elevated animate-in fade-in-0 zoom-in-95"
      style={{ transformOrigin: "top" }}
    >
      <div className="space-y-5">
        <Input
          ref={searchInputRef}
          placeholder="Nhập tên tỉnh, thành phố..."
          value={searchValue}
          onChange={(event) => onSearchChange(event.target.value)}
        />
        <div>
          <p className="mb-2 text-xs font-semibold tracking-wide text-muted-foreground">
            TỈNH/THÀNH PHỐ
          </p>
          <div className="custom-scroll max-h-64 overflow-y-auto pr-1">
            {locations.length > 0 ? (
              locations.map((city) => (
                <button
                  key={city}
                  type="button"
                  className="flex w-full items-center justify-between border-b py-3 text-left text-sm last:border-none hover:bg-muted"
                  onClick={() => onSelect(city)}
                >
                  <span className="px-2">{city}</span>
                </button>
              ))
            ) : (
              <div className="px-4 py-3 text-sm text-muted-foreground">
                Không tìm thấy địa điểm
              </div>
            )}
          </div>
        </div>
        {recentLocations.length > 0 && (
          <div>
            <p className="mb-2 text-xs font-semibold tracking-wide text-muted-foreground">
              TÌM KIẾM GẦN ĐÂY
            </p>
            <div className="flex flex-wrap gap-2">
              {recentLocations.map((city) => (
                <Button
                  key={city}
                  type="button"
                  variant="outline"
                  size="sm"
                  className="rounded-full"
                  onClick={() => onSelect(city)}
                >
                  {city}
                </Button>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
