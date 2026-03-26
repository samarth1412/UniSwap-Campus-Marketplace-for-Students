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

import './Pagination.css';

type PaginationProps = {
  currentPage: number; // 1-indexed
  totalPages: number;
  onPrev: () => void;
  onNext: () => void;
};

export function Pagination({ currentPage, totalPages, onPrev, onNext }: PaginationProps) {
  // If there is only one (or zero) page, no need to show pagination controls.
  if (totalPages <= 1) return null;

  return (
    <div className="pagination" aria-label="Pagination">
      <button
        type="button"
        className="pagination__btn"
        onClick={onPrev}
        disabled={currentPage <= 1}
        aria-disabled={currentPage <= 1}
      >
        Previous
      </button>

      <div className="pagination__meta" aria-live="polite">
        Page {currentPage} of {totalPages}
      </div>

      <button
        type="button"
        className="pagination__btn"
        onClick={onNext}
        disabled={currentPage >= totalPages}
        aria-disabled={currentPage >= totalPages}
      >
        Next
      </button>
    </div>
  );
}

