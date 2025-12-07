"use client";

import { Fragment } from "react";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  ChevronLeft,
  ChevronRight,
  ChevronsLeft,
  ChevronsRight,
} from "lucide-react";

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  pageSize?: number;
  onPageChange: (page: number) => void;
  onPageSizeChange?: (pageSize: number) => void;
}

export function Pagination({
  currentPage,
  totalPages,
  pageSize = 5,
  onPageChange,
  onPageSizeChange,
}: PaginationProps) {
  // Temporarily show pagination even with 1 page for debugging
  // if (totalPages <= 1) return null;

  const pageSizeOptions = [5, 10, 20, 50, 100];
  const maxVisiblePages = 3; // Reduced from 5 to 3

  // Calculate which pages to show
  const getVisiblePages = () => {
    const pages: number[] = [];

    // Always show first page
    pages.push(1);

    // Calculate middle pages
    if (totalPages <= 5) {
      // If total pages is 5 or less, show all
      for (let i = 2; i < totalPages; i++) {
        pages.push(i);
      }
    } else {
      // Show pages around current page
      if (currentPage <= 3) {
        // Near start: show 2, 3, 4
        pages.push(2, 3, 4);
      } else if (currentPage >= totalPages - 2) {
        // Near end: show last 3 before final
        pages.push(totalPages - 3, totalPages - 2, totalPages - 1);
      } else {
        // Middle: show current and neighbors
        pages.push(currentPage - 1, currentPage, currentPage + 1);
      }
    }

    // Always show last page if more than 1 page
    if (totalPages > 1) {
      pages.push(totalPages);
    }

    return [...new Set(pages)].sort((a, b) => a - b);
  };

  const visiblePages = getVisiblePages();

  return (
    <div className="mt-6 flex flex-col items-center gap-4 sm:flex-row sm:justify-between">
      {/* Page size selector */}
      {onPageSizeChange && (
        <div className="flex items-center gap-2">
          <span className="text-sm text-muted-foreground">Hiển thị:</span>
          <Select
            value={pageSize.toString()}
            onValueChange={(value) => {
              onPageSizeChange(Number(value));
              onPageChange(1); // Reset to first page when changing page size
            }}
          >
            <SelectTrigger className="h-9 w-20">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {pageSizeOptions.map((size) => (
                <SelectItem key={size} value={size.toString()}>
                  {size}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <span className="text-sm text-muted-foreground">mục</span>
        </div>
      )}

      {/* Pagination controls */}
      <div className="flex items-center gap-2">
        {/* First page */}
        <Button
          variant="outline"
          size="icon"
          onClick={() => onPageChange(1)}
          disabled={currentPage === 1}
          className="h-9 w-9"
        >
          <ChevronsLeft className="h-4 w-4" />
        </Button>

        {/* Previous page */}
        <Button
          variant="outline"
          size="icon"
          onClick={() => onPageChange(currentPage - 1)}
          disabled={currentPage === 1}
          className="h-9 w-9"
        >
          <ChevronLeft className="h-4 w-4" />
        </Button>

        {/* Page numbers */}
        {visiblePages.map((page, index) => (
          <Fragment key={page}>
            {/* Show ... if there's a gap */}
            {index > 0 && visiblePages[index - 1] !== page - 1 && (
              <span className="px-2 text-muted-foreground">...</span>
            )}
            <Button
              variant={currentPage === page ? "default" : "outline"}
              size="icon"
              onClick={() => onPageChange(page)}
              className="h-9 w-9"
            >
              {page}
            </Button>
          </Fragment>
        ))}

        {/* Next page */}
        <Button
          variant="outline"
          size="icon"
          onClick={() => onPageChange(currentPage + 1)}
          disabled={currentPage === totalPages}
          className="h-9 w-9"
        >
          <ChevronRight className="h-4 w-4" />
        </Button>

        {/* Last page */}
        <Button
          variant="outline"
          size="icon"
          onClick={() => onPageChange(totalPages)}
          disabled={currentPage === totalPages}
          className="h-9 w-9"
        >
          <ChevronsRight className="h-4 w-4" />
        </Button>

        {/* Page info */}
        <span className="ml-4 text-sm text-muted-foreground">
          Trang {currentPage} / {totalPages}
        </span>
      </div>
    </div>
  );
}
