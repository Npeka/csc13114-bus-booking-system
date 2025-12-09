import { cn } from "@/lib/utils";
import { Check } from "lucide-react";

interface Step {
  title: string;
  description?: string;
}

interface StepperProps {
  steps: Step[];
  currentStep: number;
  className?: string;
}

export function Stepper({ steps, currentStep, className }: StepperProps) {
  return (
    <div className={cn("w-full", className)}>
      {/* Progress bar */}
      <div className="relative mb-8">
        <div className="absolute top-5 right-0 left-0 h-0.5 bg-muted">
          <div
            className="h-full bg-primary transition-all duration-300"
            style={{
              width: `${(currentStep / (steps.length - 1)) * 100}%`,
            }}
          />
        </div>

        {/* Steps */}
        <div className="relative flex justify-between">
          {steps.map((step, index) => {
            const isCompleted = index < currentStep;
            const isCurrent = index === currentStep;
            const isUpcoming = index > currentStep;

            return (
              <div key={index} className="flex flex-col items-center">
                {/* Step circle */}
                <div
                  className={cn(
                    "flex h-10 w-10 items-center justify-center rounded-full border-2 bg-background transition-all duration-300",
                    {
                      "border-primary bg-primary text-primary-foreground":
                        isCompleted || isCurrent,
                      "border-muted text-muted-foreground": isUpcoming,
                    },
                  )}
                >
                  {isCompleted ? (
                    <Check className="h-5 w-5" />
                  ) : (
                    <span className="text-sm font-semibold">{index + 1}</span>
                  )}
                </div>

                {/* Step label */}
                <div className="mt-2 text-center">
                  <p
                    className={cn("text-sm font-medium", {
                      "text-foreground": isCurrent,
                      "text-muted-foreground": !isCurrent,
                    })}
                  >
                    {step.title}
                  </p>
                  {step.description && (
                    <p className="mt-1 text-xs text-muted-foreground">
                      {step.description}
                    </p>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
