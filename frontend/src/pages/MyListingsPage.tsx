import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import type { Listing } from '../types/listing';
import { MOCK_LISTINGS } from '../data/mockListings';

/**
 * FE-22: My listings page UI
 * Uses mock data until the backend user listings endpoint is ready.
 */
export function MyListingsPage() {
  const [listings, setListings] = useState<Listing[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const timer = window.setTimeout(() => {
      setListings(MOCK_LISTINGS.slice(0, 2));
      setLoading(false);
    }, 250);

    return () => window.clearTimeout(timer);
  }, []);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
      <header style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
        <h1 style={{ margin: 0 }}>My Listings</h1>
        <p style={{ margin: 0, color: '#4b5563' }}>
          Manage the items you have posted. This page is using placeholder data until the backend endpoint is available.
        </p>
      </header>

      {loading ? (
        <p style={{ margin: 0 }}>Loading your listings...</p>
      ) : listings.length === 0 ? (
        <div style={{ border: '1px solid #e5e7eb', borderRadius: '0.75rem', padding: '1.5rem' }}>
          <h2 style={{ marginTop: 0 }}>No listings yet</h2>
          <p style={{ color: '#4b5563' }}>
            Start by creating your first item listing for the marketplace.
          </p>
          <Link
            to="/create"
            style={{
              display: 'inline-block',
              padding: '0.55rem 1rem',
              borderRadius: '999px',
              backgroundColor: '#16a34a',
              color: '#fff',
              textDecoration: 'none',
              fontWeight: 600,
            }}
          >
            Create listing
          </Link>
        </div>
      ) : (
        <section
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fill, minmax(240px, 1fr))',
            gap: '1.25rem',
          }}
        >
          {listings.map((listing) => (
            <article
              key={listing.id}
              style={{
                border: '1px solid #e5e7eb',
                borderRadius: '0.75rem',
                overflow: 'hidden',
                backgroundColor: '#fff',
              }}
            >
              <img
                src={listing.imageUrl}
                alt={listing.title}
                style={{ width: '100%', height: '180px', objectFit: 'cover' }}
              />
              <div style={{ padding: '1rem', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
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
                <h2 style={{ margin: 0, fontSize: '1rem' }}>{listing.title}</h2>
                <p style={{ margin: 0, color: '#059669', fontWeight: 700 }}>
                  Rs. {listing.price.toFixed(0)}
                </p>
                <p style={{ margin: 0, color: '#4b5563', lineHeight: 1.5 }}>
                  {listing.description}
                </p>
                <div style={{ display: 'flex', gap: '0.5rem', marginTop: '0.5rem', flexWrap: 'wrap' }}>
                  <Link
                    to={`/listing/${listing.id}`}
                    style={{
                      padding: '0.45rem 0.8rem',
                      borderRadius: '999px',
                      border: '1px solid #d1d5db',
                      color: '#111827',
                      textDecoration: 'none',
                    }}
                  >
                    View
                  </Link>
                  <Link
                    to={`/listing/${listing.id}/edit`}
                    style={{
                      padding: '0.45rem 0.8rem',
                      borderRadius: '999px',
                      border: '1px solid #2563eb',
                      backgroundColor: '#eff6ff',
                      color: '#1d4ed8',
                      textDecoration: 'none',
                    }}
                  >
                    Edit
                  </Link>
                </div>
              </div>
            </article>
          ))}
        </section>
      )}
    </div>
  );
}
