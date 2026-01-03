"use client";

import { PaginationWithLinks } from "@/components/ui/pagination-with-links";

interface TripPaginationProps {
  currentPage: number;
  totalPages: number;
  onPageCreateURL: (pageNumber: number) => string;
}

export function TripPagination({
  currentPage,
  totalPages,
  onPageCreateURL,
}: TripPaginationProps) {
  return (
    <div className="mt-8 flex justify-center">
      <PaginationWithLinks
        page={currentPage}
        totalPages={totalPages}
        createPageURL={onPageCreateURL}
      />
    </div>
  );
}
