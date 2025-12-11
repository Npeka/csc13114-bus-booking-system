import { Clock } from "lucide-react";

interface CountdownTimerProps {
  seconds: number;
}

export function CountdownTimer({ seconds }: CountdownTimerProps) {
  const minutes = Math.floor(seconds / 60);
  const secs = seconds % 60;

  // Color coding: green > 2min, yellow 1-2min, red < 1min
  const isUrgent = minutes < 1;
  const isWarning = minutes >= 1 && minutes < 2;

  if (seconds <= 0) {
    return (
      <div className="flex items-center gap-2 rounded-md bg-destructive/10 px-3 py-2 text-sm text-destructive">
        <Clock className="h-4 w-4" />
        <span className="font-medium">
          Hết thời gian giữ chỗ - Vui lòng chọn lại ghế
        </span>
      </div>
    );
  }

  return (
    <div
      className={`flex items-center gap-2 rounded-md px-3 py-2 text-sm font-medium transition-colors ${
        isUrgent
          ? "bg-destructive/10 text-destructive"
          : isWarning
            ? "bg-orange-500/10 text-orange-600 dark:text-orange-500"
            : "bg-blue-500/10 text-blue-600 dark:text-blue-500"
      }`}
    >
      <Clock className="h-4 w-4 animate-pulse" />
      <span>
        Ghế được giữ trong{" "}
        <span className="font-bold tabular-nums">
          {minutes}:{secs.toString().padStart(2, "0")}
        </span>
      </span>
    </div>
  );
}
