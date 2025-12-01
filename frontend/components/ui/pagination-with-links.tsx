import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";

interface PaginationWithLinksProps {
  page: number;
  totalPages: number;
  createPageURL: (pageNumber: number) => string;
}

export function PaginationWithLinks({
  page,
  totalPages,
  createPageURL,
}: PaginationWithLinksProps) {
  return (
    <Pagination>
      <PaginationContent>
        <PaginationItem>
          <PaginationPrevious
            href={createPageURL(page - 1)}
            aria-disabled={page <= 1}
            className={page <= 1 ? "pointer-events-none opacity-50" : undefined}
          />
        </PaginationItem>

        {/* First Page */}
        {page > 3 && (
          <PaginationItem>
            <PaginationLink href={createPageURL(1)}>1</PaginationLink>
          </PaginationItem>
        )}

        {/* Ellipsis Start */}
        {page > 4 && (
          <PaginationItem>
            <PaginationEllipsis />
          </PaginationItem>
        )}

        {/* Page Numbers */}
        {Array.from({ length: totalPages }, (_, i) => i + 1)
          .filter((p) => {
            if (totalPages <= 7) return true;
            return Math.abs(page - p) <= 2; // Show current +/- 2
          })
          .map((p) => (
            <PaginationItem key={p}>
              <PaginationLink href={createPageURL(p)} isActive={page === p}>
                {p}
              </PaginationLink>
            </PaginationItem>
          ))}

        {/* Ellipsis End */}
        {page < totalPages - 3 && (
          <PaginationItem>
            <PaginationEllipsis />
          </PaginationItem>
        )}

        {/* Last Page */}
        {page < totalPages - 2 && (
          <PaginationItem>
            <PaginationLink href={createPageURL(totalPages)}>
              {totalPages}
            </PaginationLink>
          </PaginationItem>
        )}

        <PaginationItem>
          <PaginationNext
            href={createPageURL(page + 1)}
            aria-disabled={page >= totalPages}
            className={
              page >= totalPages ? "pointer-events-none opacity-50" : undefined
            }
          />
        </PaginationItem>
      </PaginationContent>
    </Pagination>
  );
}
