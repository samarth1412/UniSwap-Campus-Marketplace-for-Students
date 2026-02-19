import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { api, type ApiResponse } from '../services/api';
import type { Listing } from '../types/listing';
import { MOCK_LISTINGS } from '../data/mockListings';

/**
 * Listing feed
 * FE-8: UI with mocked data
 * FE-13: Connect to backend
 * FE-14: Basic search/filter UI
 */
export function HomePage() {
  const [listings, setListings] = useState<Listing[]>(MOCK_LISTINGS);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [categoryFilter, setCategoryFilter] = useState<string>('all');

  useEffect(() => {
    let cancelled = false;

    async function loadListings() {
      setLoading(true);
      setError(null);
      try {
        const response = await api.get<ApiResponse<Listing[]>>('/listings');
        if (!response.data.success || !response.data.data) {
          throw new Error(response.data.error || 'Failed to load listings');
        }
        if (!cancelled) {
          // For Sprint 1 we only rely on id, title, price, category.
          setListings(response.data.data);
        }
      } catch (err) {
        // Fallback to mocked data if backend is not ready.
        if (!cancelled) {
          console.error('Failed to load listings, using mock data instead:', err);
          setError('Unable to load listings from server. Showing sample data.');
          setListings(MOCK_LISTINGS);
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
  }, []);

  const filteredListings = listings.filter((listing) => {
    const matchesSearch =
      searchTerm.trim().length === 0 ||
      listing.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
      listing.description.toLowerCase().includes(searchTerm.toLowerCase());

    const matchesCategory =
      categoryFilter === 'all' || listing.category.toLowerCase() === categoryFilter.toLowerCase();

    return matchesSearch && matchesCategory;
  });

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
            style={{
              padding: '0.5rem 0.75rem',
              borderRadius: '999px',
              border: '1px solid #ddd',
              minWidth: '160px',
            }}
          >
            <option value="all">All categories</option>
            <option value="Books">Books</option>
            <option value="Electronics">Electronics</option>
            <option value="Furniture">Furniture</option>
            <option value="Other">Other</option>
          </select>
        </div>
        {loading && <p style={{ margin: 0 }}>Loading listings…</p>}
        {error && (
          <p style={{ margin: 0, color: '#b45309' }}>
            {error}
          </p>
        )}
      </header>

      {filteredListings.length === 0 ? (
        <p>No listings match your search yet.</p>
      ) : (
        <section
          className="listings-grid"
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(4, minmax(0, 1fr))',
            gap: '1.5rem',
          }}
        >
          {filteredListings.map((listing) => (
            <Link
              key={listing.id}
              to={`/listing/${listing.id}`}
              style={{
                textDecoration: 'none',
                color: 'inherit',
              }}
            >
              <article
                style={{
                  border: '1px solid #eee',
                  borderRadius: '0.75rem',
                  overflow: 'hidden',
                  backgroundColor: '#fff',
                  boxShadow: '0 1px 3px rgba(0,0,0,0.04)',
                  display: 'flex',
                  flexDirection: 'column',
                  height: '100%',
                }}
              >
                <div style={{ position: 'relative', paddingTop: '62%' }}>
                  <img
                    src={listing.imageUrl}
                    alt={listing.title}
                    style={{
                      position: 'absolute',
                      inset: 0,
                      width: '100%',
                      height: '100%',
                      objectFit: 'cover',
                    }}
                  />
                </div>
                <div style={{ padding: '0.75rem 0.9rem', display: 'flex', flexDirection: 'column', gap: '0.4rem' }}>
                  <span
                    style={{
                      fontSize: '0.75rem',
                      textTransform: 'uppercase',
                      letterSpacing: '0.04em',
                      color: '#6b7280',
                    }}
                  >
                    {listing.category}
                  </span>
                  <h2
                    style={{
                      margin: 0,
                      fontSize: '1rem',
                      fontWeight: 600,
                      lineHeight: 1.3,
                    }}
                  >
                    {listing.title}
                  </h2>
                  <span style={{ fontWeight: 700, color: '#059669' }}>${listing.price.toFixed(0)}</span>
                </div>
              </article>
            </Link>
          ))}
        </section>
      )}
    </div>
  );
}

