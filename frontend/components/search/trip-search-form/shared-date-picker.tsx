"use client";

import { useEffect, useRef, useState } from "react";
import { Calendar } from "@/components/ui/calendar";

type SharedDatePickerProps = {
  isOpen: boolean;
  onClose: () => void;
  departureDate: Date | undefined;
  returnDate: Date | undefined;
  onDepartureDateChange: (date: Date | undefined) => void;
  onReturnDateChange: (date: Date | undefined) => void;
  activeField: "departure" | "return";
  triggerRef: React.RefObject<HTMLDivElement | null>;
};

export function SharedDatePicker({
  isOpen,
  onClose,
  departureDate,
  returnDate,
  onDepartureDateChange,
  onReturnDateChange,
  activeField,
  triggerRef,
}: SharedDatePickerProps) {
  const pickerRef = useRef<HTMLDivElement>(null);
  const [position, setPosition] = useState({ top: 0, left: 0 });

  // Calculate position based on trigger element
  useEffect(() => {
    if (isOpen && triggerRef.current) {
      const rect = triggerRef.current.getBoundingClientRect();
      const pickerWidth = 680; // Approximate width for 2 months
      const centerX = rect.left + rect.width / 2;

      setPosition({
        top: rect.bottom + 8,
        left: Math.max(
          16,
          Math.min(
            centerX - pickerWidth / 2,
            window.innerWidth - pickerWidth - 16,
          ),
        ),
      });
    }
  }, [isOpen, triggerRef]);

  // Handle outside clicks
  useEffect(() => {
    if (!isOpen) return;

    const handleClickOutside = (event: MouseEvent) => {
      const target = event.target as Node;

      // Don't close if clicking on the picker itself or the trigger
      if (
        pickerRef.current?.contains(target) ||
        triggerRef.current?.contains(target)
      ) {
        return;
      }

      onClose();
    };

    document.addEventListener("pointerdown", handleClickOutside);
    return () =>
      document.removeEventListener("pointerdown", handleClickOutside);
  }, [isOpen, onClose, triggerRef]);

  if (!isOpen) return null;

  const normalizeDeparture = departureDate
    ? new Date(
        departureDate.getFullYear(),
        departureDate.getMonth(),
        departureDate.getDate(),
      )
    : undefined;
  const normalizeReturn = returnDate
    ? new Date(
        returnDate.getFullYear(),
        returnDate.getMonth(),
        returnDate.getDate(),
      )
    : undefined;

  const startDate =
    normalizeDeparture && normalizeReturn
      ? normalizeDeparture <= normalizeReturn
        ? normalizeDeparture
        : normalizeReturn
      : normalizeDeparture;
  const endDate =
    normalizeDeparture && normalizeReturn
      ? normalizeDeparture <= normalizeReturn
        ? normalizeReturn
        : normalizeDeparture
      : normalizeReturn;

  const modifiers = {
    range_start: startDate ? [startDate] : [],
    range_end: endDate ? [endDate] : [],
    range_middle:
      startDate && endDate
        ? Array.from(
            {
              length:
                Math.floor(
                  (endDate.getTime() - startDate.getTime()) /
                    (1000 * 60 * 60 * 24),
                ) - 1,
            },
            (_, i) => {
              const date = new Date(startDate);
              date.setDate(date.getDate() + i + 1);
              return date;
            },
          )
        : [],
  };

  const handleSelect = (date: Date | undefined) => {
    if (!date) return;

    if (activeField === "departure") {
      onDepartureDateChange(date);

      if (returnDate && date > returnDate) {
        onReturnDateChange(undefined);
      }
    } else if (activeField === "return") {
      onReturnDateChange(date);

      if (!departureDate) {
        onDepartureDateChange(date);
      } else if (date < departureDate) {
        onDepartureDateChange(date);
        onReturnDateChange(departureDate);
      }
    }
  };

  return (
    <div
      ref={pickerRef}
      className="fixed z-50 animate-in fade-in-0 zoom-in-95 slide-in-from-top-2"
      style={{
        top: `${position.top}px`,
        left: `${position.left}px`,
      }}
    >
      <div className="rounded-lg border border-border bg-card shadow-lg p-4">
        {/* Active field indicator */}
        <div className="mb-3 flex items-center justify-center gap-2 text-sm">
          <span
            className={`px-3 py-1 rounded-full transition-colors ${
              activeField === "departure"
                ? "bg-brand-primary text-white font-medium"
                : "bg-muted text-muted-foreground"
            }`}
          >
            Ngày đi
          </span>
          <span className="text-muted-foreground">→</span>
          <span
            className={`px-3 py-1 rounded-full transition-colors ${
              activeField === "return"
                ? "bg-brand-primary text-white font-medium"
                : "bg-muted text-muted-foreground"
            }`}
          >
            Ngày về
          </span>
        </div>

        <Calendar
          mode="single"
          selected={activeField === "departure" ? departureDate : returnDate}
          onSelect={handleSelect}
          disabled={(date) => date < new Date(new Date().setHours(0, 0, 0, 0))}
          numberOfMonths={2}
          className="rounded-lg"
          modifiers={modifiers}
        />
      </div>
    </div>
  );
}
