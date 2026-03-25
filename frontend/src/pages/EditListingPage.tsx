import { useMemo, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { MOCK_LISTINGS } from '../data/mockListings';

/**
 * FE-18: Edit listing UI
 * Form-only implementation. API hookup is handled in FE-19.
 */
export function EditListingPage() {
  const { id } = useParams<{ id: string }>();
  const numericId = Number(id);
  const listing = useMemo(
    () => MOCK_LISTINGS.find((item) => item.id === numericId),
    [numericId]
  );

  const [title, setTitle] = useState(listing?.title ?? '');
  const [description, setDescription] = useState(listing?.description ?? '');
  const [price, setPrice] = useState(listing ? String(listing.price) : '');
  const [category, setCategory] = useState(listing?.category ?? 'Books');
  const [message, setMessage] = useState<string | null>(null);

  if (!listing) {
    return (
      <div>
        <p>
          <Link to="/">{'< Back to listings'}</Link>
        </p>
        <h1>Listing not found</h1>
        <p>The listing you want to edit does not exist.</p>
      </div>
    );
  }

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    setMessage('Edit UI saved locally. Backend integration will be added in FE-19.');
  };

  return (
    <div style={{ maxWidth: '680px', display: 'flex', flexDirection: 'column', gap: '1rem' }}>
      <p style={{ margin: 0 }}>
        <Link to={`/listing/${listing.id}`}>{'< Back to listing'}</Link>
      </p>
      <div>
        <h1 style={{ marginBottom: '0.5rem' }}>Edit Listing</h1>
        <p style={{ margin: 0, color: '#4b5563' }}>
          Update your listing details. Changes are local for now until the Sprint 2 edit API is ready.
        </p>
      </div>

      <form
        onSubmit={handleSubmit}
        style={{ display: 'flex', flexDirection: 'column', gap: '0.9rem' }}
      >
        <label style={{ display: 'flex', flexDirection: 'column', gap: '0.25rem' }}>
          <span>Title</span>
          <input
            type="text"
            value={title}
            onChange={(event) => setTitle(event.target.value)}
            required
            style={{
              padding: '0.6rem 0.75rem',
              borderRadius: '0.5rem',
              border: '1px solid #e5e7eb',
            }}
          />
        </label>

        <label style={{ display: 'flex', flexDirection: 'column', gap: '0.25rem' }}>
          <span>Description</span>
          <textarea
            value={description}
            onChange={(event) => setDescription(event.target.value)}
            required
            rows={5}
            style={{
              padding: '0.6rem 0.75rem',
              borderRadius: '0.5rem',
              border: '1px solid #e5e7eb',
              resize: 'vertical',
            }}
          />
        </label>

        <div style={{ display: 'flex', gap: '0.75rem', flexWrap: 'wrap' }}>
          <label
            style={{
              display: 'flex',
              flexDirection: 'column',
              gap: '0.25rem',
              flex: '1 1 180px',
            }}
          >
            <span>Price</span>
            <input
              type="number"
              min={0}
              step="1"
              value={price}
              onChange={(event) => setPrice(event.target.value)}
              required
              style={{
                padding: '0.6rem 0.75rem',
                borderRadius: '0.5rem',
                border: '1px solid #e5e7eb',
              }}
            />
          </label>

          <label
            style={{
              display: 'flex',
              flexDirection: 'column',
              gap: '0.25rem',
              flex: '1 1 180px',
            }}
          >
            <span>Category</span>
            <select
              value={category}
              onChange={(event) => setCategory(event.target.value)}
              style={{
                padding: '0.6rem 0.75rem',
                borderRadius: '0.5rem',
                border: '1px solid #e5e7eb',
              }}
            >
              <option value="Books">Books</option>
              <option value="Electronics">Electronics</option>
              <option value="Furniture">Furniture</option>
              <option value="Other">Other</option>
            </select>
          </label>
        </div>

        {message && (
          <p style={{ margin: 0, color: '#065f46', fontSize: '0.95rem' }}>
            {message}
          </p>
        )}

        <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '0.75rem' }}>
          <Link
            to={`/listing/${listing.id}`}
            style={{
              padding: '0.6rem 1rem',
              borderRadius: '999px',
              border: '1px solid #d1d5db',
              color: '#111827',
              textDecoration: 'none',
            }}
          >
            Cancel
          </Link>
          <button
            type="submit"
            style={{
              padding: '0.6rem 1.1rem',
              borderRadius: '999px',
              border: 'none',
              backgroundColor: '#2563eb',
              color: '#fff',
              fontWeight: 600,
              cursor: 'pointer',
            }}
          >
            Save changes
          </button>
        </div>
      </form>
    </div>
  );
}
