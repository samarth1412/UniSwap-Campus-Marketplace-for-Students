import { useEffect, useState } from 'react';
import { Pagination } from '../components/Pagination';
import { ListingCard } from '../components/ListingCard';
import { listingsApi } from '../services/api';
import type { Listing } from '../types/listing';

/**
 * Listing feed
 * FE-8: UI with mocked data
 * FE-13: Connect to backend
 * FE-14: Basic search/filter UI
 */
export function HomePage() {
  const [items, setItems] = useState<Listing[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [debouncedSearchTerm, setDebouncedSearchTerm] = useState('');
  const [categoryFilter, setCategoryFilter] = useState('');
  const [minPrice, setMinPrice] = useState('');
  const [maxPrice, setMaxPrice] = useState('');
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const [totalPages, setTotalPages] = useState(1);

  const parsedMinPrice = parsePriceFilter(minPrice);
  const parsedMaxPrice = parsePriceFilter(maxPrice);

  useEffect(() => {
    const nextValue = searchTerm.trim();
    const timeoutId = window.setTimeout(() => {
      setDebouncedSearchTerm(nextValue);
    }, 300);

    return () => {
      window.clearTimeout(timeoutId);
    };
  }, [searchTerm]);

  useEffect(() => {
    // New filter values always start from page 1.
    setPage(1);
  }, [debouncedSearchTerm, categoryFilter, minPrice, maxPrice]);

  useEffect(() => {
    let cancelled = false;

    async function loadListings() {
      setLoading(true);
      setError(null);

      if (parsedMinPrice === null || parsedMaxPrice === null) {
        setError('Price filters must be valid non-negative numbers.');
        setItems([]);
        setTotalPages(1);
        setLoading(false);
        return;
      }

      if (
        parsedMinPrice !== undefined &&
        parsedMaxPrice !== undefined &&
        parsedMinPrice > parsedMaxPrice
      ) {
        setError('Min price cannot be greater than max price.');
        setItems([]);
        setTotalPages(1);
        setLoading(false);
        return;
      }

      try {
        const response = await listingsApi.getAll({
          search: debouncedSearchTerm || undefined,
          category: categoryFilter || undefined,
          min_price: parsedMinPrice,
          max_price: parsedMaxPrice,
          page,
          limit,
        });
        if (!response.data.success || !response.data.data) {
          throw new Error(response.data.error || 'Failed to load listings');
        }

        if (!cancelled) {
          setItems(response.data.data.items);
          setPage(response.data.data.page);
          setTotalPages(Math.max(1, response.data.data.total_pages));
        }
      } catch (err) {
        if (!cancelled) {
          console.error('Failed to load listings', err);
          setError('Unable to load listings from server.');
          setItems([]);
          setTotalPages(1);
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    void loadListings();

    return () => {
      cancelled = true;
    };
  }, [page, limit, debouncedSearchTerm, categoryFilter, parsedMinPrice, parsedMaxPrice]);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
      <header style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <h1 style={{ margin: 0 }}>Listings</h1>
          <span style={{ fontSize: '0.9rem', color: '#666' }}>
            Browse items from your campus community.
          </span>
        </div>
        <div
          style={{
            display: 'flex',
            gap: '0.75rem',
            flexWrap: 'wrap',
            alignItems: 'center',
          }}
        >
          <input
            type="text"
            placeholder="Search by title or description"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            aria-label="Search listings"
            style={{
              flex: '1 1 220px',
              padding: '0.5rem 0.75rem',
              borderRadius: '999px',
              border: '1px solid #ddd',
            }}
          />
          <select
            value={categoryFilter}
            onChange={(e) => setCategoryFilter(e.target.value)}
            aria-label="Filter by category"
            style={{
              padding: '0.5rem 0.75rem',
              borderRadius: '999px',
              border: '1px solid #ddd',
              minWidth: '160px',
            }}
          >
            <option value="">All categories</option>
            <option value="Books">Books</option>
            <option value="Electronics">Electronics</option>
            <option value="Furniture">Furniture</option>
            <option value="Other">Other</option>
          </select>
          <input
            type="number"
            min="0"
            placeholder="Min price"
            value={minPrice}
            onChange={(e) => setMinPrice(e.target.value)}
            aria-label="Minimum price"
            style={{
              flex: '0 1 140px',
              padding: '0.5rem 0.75rem',
              borderRadius: '999px',
              border: '1px solid #ddd',
            }}
          />
          <input
            type="number"
            min="0"
            placeholder="Max price"
            value={maxPrice}
            onChange={(e) => setMaxPrice(e.target.value)}
            aria-label="Maximum price"
            style={{
              flex: '0 1 140px',
              padding: '0.5rem 0.75rem',
              borderRadius: '999px',
              border: '1px solid #ddd',
            }}
          />
        </div>
        {loading && <p style={{ margin: 0 }}>Loading listings...</p>}
        {error && <p style={{ margin: 0, color: '#b45309' }}>{error}</p>}
      </header>

      {items.length === 0 && !loading ? (
        <p>No listings match your search yet.</p>
      ) : (
        <section
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fill, minmax(220px, 1fr))',
            gap: '1.5rem',
          }}
        >
          {items.map((listing) => (
            <ListingCard key={listing.id} listing={listing} />
          ))}
        </section>
      )}

      <Pagination
        currentPage={page}
        totalPages={totalPages}
        disabled={loading}
        onPageChange={(nextPage) => setPage(nextPage)}
      />
    </div>
  );
}

function parsePriceFilter(value: string): number | undefined | null {
  const trimmed = value.trim();
  if (trimmed === '') return undefined;

  const parsed = Number(trimmed);
  if (!Number.isFinite(parsed) || parsed < 0) {
    return null;
  }

  return parsed;
}
