import { useEffect, useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import type { Listing } from '../types/listing';
import { profileApi } from '../services/api';
import { formatUsd } from '../utils/currency';

export function MyListingsPage() {
  const [listings, setListings] = useState<Listing[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const location = useLocation();
  const deletedMessage = (location.state as { deletedMessage?: string } | null)?.deletedMessage;
  const fallbackImage = 'https://placehold.co/640x420?text=No+Image';

  useEffect(() => {
    let cancelled = false;

    async function loadMyListings() {
      setLoading(true);
      setError(null);

      try {
        const meResponse = await profileApi.getMe();
        if (!meResponse.data.success || !meResponse.data.data) {
          throw new Error(meResponse.data.error || 'Failed to load user profile');
        }

        const listingsResponse = await profileApi.getMyListings(meResponse.data.data.id);
        if (!listingsResponse.data.success || !listingsResponse.data.data) {
          throw new Error(listingsResponse.data.error || 'Failed to load your listings');
        }

        if (!cancelled) {
          setListings(listingsResponse.data.data);
        }
      } catch (err) {
        console.error('Failed to load my listings', err);
        if (!cancelled) {
          setError('Unable to load your listings right now. Please try again later.');
          setListings([]);
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    void loadMyListings();

    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
      <header style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
        <h1 style={{ margin: 0 }}>My Listings</h1>
        <p style={{ margin: 0, color: '#4b5563' }}>
          Manage the items you have posted to the marketplace.
        </p>
      </header>
      {deletedMessage && <p style={{ margin: 0, color: '#065f46' }}>{deletedMessage}</p>}

      {loading ? (
        <p style={{ margin: 0 }}>Loading your listings...</p>
      ) : error ? (
        <div style={{ border: '1px solid #fecaca', borderRadius: '0.75rem', padding: '1.25rem', background: '#fff1f2' }}>
          <h2 style={{ marginTop: 0, marginBottom: '0.5rem' }}>Could not load listings</h2>
          <p style={{ marginTop: 0, color: '#7f1d1d' }}>
            Please refresh the page or check your connection, then try again.
          </p>
          <Link
            to="/"
            style={{
              display: 'inline-block',
              padding: '0.5rem 0.95rem',
              borderRadius: '999px',
              border: '1px solid #d1d5db',
              color: '#111827',
              textDecoration: 'none',
              backgroundColor: '#fff',
            }}
          >
            Back to all listings
          </Link>
        </div>
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
                src={listing.imageUrl?.trim() ? listing.imageUrl : fallbackImage}
                alt={listing.title}
                style={{ width: '100%', height: '180px', objectFit: 'cover' }}
                onError={(event) => {
                  const image = event.currentTarget;
                  if (image.src !== fallbackImage) {
                    image.src = fallbackImage;
                  }
                }}
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
                <p style={{ margin: 0, color: '#059669', fontWeight: 700 }}>{formatUsd(listing.price)}</p>
                <p style={{ margin: 0, color: '#4b5563', lineHeight: 1.5 }}>{listing.description}</p>
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
