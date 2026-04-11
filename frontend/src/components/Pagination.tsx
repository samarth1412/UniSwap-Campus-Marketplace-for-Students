import { useMemo } from 'react';
import './Pagination.css';

type PaginationProps = {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  disabled?: boolean;
};

export function Pagination({
  currentPage,
  totalPages,
  onPageChange,
  disabled = false,
}: PaginationProps) {
  const safeTotalPages = useMemo(() => Math.max(1, totalPages), [totalPages]);
  const safeCurrentPage = useMemo(
    () => Math.min(Math.max(1, currentPage), safeTotalPages),
    [currentPage, safeTotalPages]
  );

  if (safeTotalPages <= 1) return null;

  return (
    <nav className="pagination" aria-label="Pagination">
      <button
        type="button"
        className="pagination__btn"
        onClick={() => onPageChange(safeCurrentPage - 1)}
        disabled={disabled || safeCurrentPage <= 1}
      >
        Previous
      </button>

      <span className="pagination__indicator" aria-live="polite">
        Page {safeCurrentPage} of {safeTotalPages}
      </span>

      <button
        type="button"
        className="pagination__btn"
        onClick={() => onPageChange(safeCurrentPage + 1)}
        disabled={disabled || safeCurrentPage >= safeTotalPages}
      >
        Next
      </button>
    </nav>
  );
}
