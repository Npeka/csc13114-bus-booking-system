"use client";

import * as React from "react";
import { ChevronsUpDown, Loader2, Search, X } from "lucide-react";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";

export interface InfiniteSelectOption {
  value: string;
  label: string;
}

interface InfiniteSelectProps {
  value?: string;
  onValueChange: (value: string) => void;
  options: InfiniteSelectOption[];
  placeholder?: string;
  searchPlaceholder?: string;
  emptyText?: string;
  disabled?: boolean;
  isLoading?: boolean;
  hasNextPage?: boolean;
  isFetchingNextPage?: boolean;
  fetchNextPage?: () => void;
}

export function InfiniteSelect({
  value,
  onValueChange,
  options,
  placeholder = "Chọn...",
  searchPlaceholder = "Tìm kiếm...",
  emptyText = "Không tìm thấy kết quả",
  disabled = false,
  isLoading = false,
  hasNextPage = false,
  isFetchingNextPage = false,
  fetchNextPage,
}: InfiniteSelectProps) {
  const [open, setOpen] = React.useState(false);
  const [searchQuery, setSearchQuery] = React.useState("");
  const listRef = React.useRef<HTMLDivElement>(null);

  const selectedOption = options.find((opt) => opt.value === value);

  const filteredOptions = React.useMemo(() => {
    if (!searchQuery) return options;
    return options.filter((opt) =>
      opt.label.toLowerCase().includes(searchQuery.toLowerCase()),
    );
  }, [options, searchQuery]);

  // Handle scroll for infinite loading
  const handleScroll = React.useCallback(
    (e: React.UIEvent<HTMLDivElement>) => {
      const target = e.currentTarget;
      const bottom =
        target.scrollHeight - target.scrollTop <= target.clientHeight + 50;

      if (bottom && hasNextPage && !isFetchingNextPage && fetchNextPage) {
        fetchNextPage();
      }
    },
    [hasNextPage, isFetchingNextPage, fetchNextPage],
  );

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className="w-full justify-between"
          disabled={disabled}
        >
          <span className="truncate">
            {selectedOption ? selectedOption.label : placeholder}
          </span>
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent
        className="w-[var(--radix-popover-trigger-width)] p-0"
        align="start"
      >
        <div className="flex flex-col">
          {/* Search input */}
          <div className="border-b p-2">
            <div className="relative">
              <Search className="absolute top-2.5 left-2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder={searchPlaceholder}
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-8"
              />
              {searchQuery && (
                <button
                  onClick={() => setSearchQuery("")}
                  className="absolute top-2.5 right-2 text-muted-foreground hover:text-foreground"
                >
                  <X className="h-4 w-4" />
                </button>
              )}
            </div>
          </div>

          {/* Options list */}
          <div
            ref={listRef}
            onScroll={handleScroll}
            className="max-h-[300px] overflow-y-auto"
          >
            {isLoading && filteredOptions.length === 0 ? (
              <div className="flex items-center justify-center py-6">
                <Loader2 className="h-4 w-4 animate-spin" />
              </div>
            ) : filteredOptions.length === 0 ? (
              <div className="py-6 text-center text-sm text-muted-foreground">
                {emptyText}
              </div>
            ) : (
              <>
                {filteredOptions.map((option) => (
                  <button
                    key={option.value}
                    onClick={() => {
                      onValueChange(option.value);
                      setOpen(false);
                      setSearchQuery("");
                    }}
                    className={cn(
                      "flex w-full cursor-pointer items-center px-3 py-2 text-sm hover:bg-accent",
                      value === option.value && "bg-accent",
                    )}
                  >
                    {option.label}
                  </button>
                ))}

                {/* Load more indicator */}
                {hasNextPage && (
                  <div className="flex items-center justify-center py-2">
                    {isFetchingNextPage && (
                      <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
                    )}
                  </div>
                )}
              </>
            )}
          </div>
        </div>
      </PopoverContent>
    </Popover>
  );
}
