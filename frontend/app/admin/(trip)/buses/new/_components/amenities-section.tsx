"use client";

import { UseFormReturn } from "react-hook-form";
import { CheckSquare } from "lucide-react";
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Checkbox } from "@/components/ui/checkbox";

interface Amenity {
  value: string;
  display_name: string;
}

interface AmenitiesSectionProps {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  form: UseFormReturn<any>;
  amenities: Amenity[];
  isLoading: boolean;
}

export function AmenitiesSection({
  form,
  amenities,
  isLoading,
}: AmenitiesSectionProps) {
  return (
    <FormField
      control={form.control}
      name="amenities"
      render={() => (
        <FormItem>
          <div className="mb-3">
            <FormLabel className="text-sm font-medium">
              <CheckSquare className="mr-2 inline h-4 w-4" />
              Tiện ích
            </FormLabel>
            <FormDescription className="text-xs">
              Chọn các tiện ích có sẵn trên xe
            </FormDescription>
          </div>
          {isLoading ? (
            <div className="space-y-2">
              <div className="h-5 w-24 animate-pulse rounded bg-muted" />
              <div className="h-5 w-24 animate-pulse rounded bg-muted" />
              <div className="h-5 w-24 animate-pulse rounded bg-muted" />
            </div>
          ) : (
            <div className="grid grid-cols-2 gap-3 md:grid-cols-3">
              {amenities?.map((amenity) => (
                <FormField
                  key={amenity.value}
                  control={form.control}
                  name="amenities"
                  render={({ field }) => {
                    return (
                      <FormItem className="flex flex-row items-start space-y-0 space-x-2">
                        <FormControl>
                          <Checkbox
                            checked={field.value?.includes(amenity.value)}
                            onCheckedChange={(checked) => {
                              return checked
                                ? field.onChange([
                                    ...(field.value || []),
                                    amenity.value,
                                  ])
                                : field.onChange(
                                    field.value?.filter(
                                      (value: string) =>
                                        value !== amenity.value,
                                    ),
                                  );
                            }}
                          />
                        </FormControl>
                        <FormLabel className="text-sm leading-none font-normal">
                          {amenity.display_name}
                        </FormLabel>
                      </FormItem>
                    );
                  }}
                />
              ))}
            </div>
          )}
          <FormMessage />
        </FormItem>
      )}
    />
  );
}
