interface PaginationProps {
  meta: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
  };
  onPageChange: (page: number) => void;
}

export function Pagination({ meta, onPageChange }: PaginationProps) {
  const { page, page_size, total, total_pages } = meta;
  const startItem = (page - 1) * page_size + 1;
  const endItem = Math.min(page * page_size, total);

  const canGoPrevious = page > 1;
  const canGoNext = page < total_pages;

  return (
    <div className="flex items-center justify-between border-t pt-4">
      <p className="text-sm text-gray-500">
        Hiển thị {startItem} - {endItem} / {total}
      </p>
      <div className="flex gap-2">
        <button
          onClick={() => onPageChange(page - 1)}
          disabled={!canGoPrevious}
          className="rounded border px-3 py-1 text-sm hover:bg-gray-50 disabled:opacity-50"
        >
          Trước
        </button>
        <button
          onClick={() => onPageChange(page + 1)}
          disabled={!canGoNext}
          className="rounded border px-3 py-1 text-sm hover:bg-gray-50 disabled:opacity-50"
        >
          Sau
        </button>
      </div>
    </div>
  );
}
