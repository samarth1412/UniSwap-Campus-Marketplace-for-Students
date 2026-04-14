import { useEffect, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { listingsApi } from '../services/api';
import { ImageUpload } from '../components/ImageUpload';

export function EditListingPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const numericId = Number(id);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [price, setPrice] = useState('');
  const [category, setCategory] = useState('Books');
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [listingExists, setListingExists] = useState(true);
  const [imageFile, setImageFile] = useState<File | null>(null);
  const [existingImageUrl, setExistingImageUrl] = useState<string | null>(null);

  const isPlaceholderImage = (url: string): boolean => url.includes('placehold.co/');

  useEffect(() => {
    let cancelled = false;

    async function loadListing() {
      if (!Number.isFinite(numericId) || numericId <= 0) {
        if (!cancelled) {
          setListingExists(false);
          setLoadError('Listing not found.');
          setLoading(false);
        }
        return;
      }

      setLoading(true);
      setLoadError(null);

      try {
        const response = await listingsApi.getById(numericId);
        if (!response.data.success || !response.data.data) {
          throw new Error(response.data.error || 'Failed to load listing');
        }

        if (!cancelled) {
          const data = response.data.data;
          setTitle(data.title);
          setDescription(data.description);
          setPrice(String(data.price));
          setCategory(data.category);
          setExistingImageUrl(isPlaceholderImage(data.imageUrl) ? null : data.imageUrl);
          setListingExists(true);
        }
      } catch (err) {
        console.error('Failed to load edit listing data', err);

        if (!cancelled) {
          setListingExists(false);
          setLoadError('Unable to load this listing from the server.');
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    void loadListing();

    return () => {
      cancelled = true;
    };
  }, [numericId]);

  if (loading) {
    return <p>Loading listing...</p>;
  }

  if (!listingExists) {
    return (
      <div>
        <p>
          <Link to="/">{'< Back to listings'}</Link>
        </p>
        <h1>Listing not found</h1>
        <p>{loadError ?? 'The listing you want to edit does not exist.'}</p>
      </div>
    );
  }

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    setSubmitError(null);

    const numericPrice = Number(price);
    if (!title.trim() || !description.trim() || !price.trim()) {
      setSubmitError('Please fill in all required fields.');
      return;
    }
    if (Number.isNaN(numericPrice) || numericPrice < 0) {
      setSubmitError('Please enter a valid price.');
      return;
    }

    setSubmitting(true);
    try {
      const response = await listingsApi.update(numericId, {
        title: title.trim(),
        description: description.trim(),
        price: numericPrice,
        category,
      });
      if (!response.data.success) {
        throw new Error(response.data.error || 'Failed to update listing');
      }

      if (imageFile) {
        const uploadResponse = await listingsApi.uploadImages(numericId, [imageFile]);
        if (!uploadResponse.data.success) {
          throw new Error(uploadResponse.data.error || 'Failed to upload listing image');
        }
      }

      navigate(`/listing/${numericId}`, { replace: true });
    } catch (err: unknown) {
      console.error('Failed to update listing', err);
      const ax = err as { response?: { data?: { error?: string } } };
      setSubmitError(ax.response?.data?.error ?? 'Could not save changes right now. Please try again.');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div style={{ maxWidth: '680px', display: 'flex', flexDirection: 'column', gap: '1rem' }}>
      <p style={{ margin: 0 }}>
        <Link to={`/listing/${numericId}`}>{'< Back to listing'}</Link>
      </p>
      <div>
        <h1 style={{ marginBottom: '0.5rem' }}>Edit Listing</h1>
        <p style={{ margin: 0, color: '#4b5563' }}>
          Update your listing details and save them to the marketplace.
        </p>
      </div>
      {loadError && <p style={{ margin: 0, color: '#b45309' }}>{loadError}</p>}

      <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '0.9rem' }}>
        <label style={{ display: 'flex', flexDirection: 'column', gap: '0.25rem' }}>
          <span>Title</span>
          <input
            type="text"
            value={title}
            onChange={(event) => setTitle(event.target.value)}
            required
            style={{ padding: '0.6rem 0.75rem', borderRadius: '0.5rem', border: '1px solid #e5e7eb' }}
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
          <label style={{ display: 'flex', flexDirection: 'column', gap: '0.25rem', flex: '1 1 180px' }}>
            <span>Price</span>
            <input
              type="number"
              min={0}
              step="1"
              value={price}
              onChange={(event) => setPrice(event.target.value)}
              required
              style={{ padding: '0.6rem 0.75rem', borderRadius: '0.5rem', border: '1px solid #e5e7eb' }}
            />
          </label>

          <label style={{ display: 'flex', flexDirection: 'column', gap: '0.25rem', flex: '1 1 180px' }}>
            <span>Category</span>
            <select
              value={category}
              onChange={(event) => setCategory(event.target.value)}
              style={{ padding: '0.6rem 0.75rem', borderRadius: '0.5rem', border: '1px solid #e5e7eb' }}
            >
              <option value="Books">Books</option>
              <option value="Electronics">Electronics</option>
              <option value="Furniture">Furniture</option>
              <option value="Other">Other</option>
            </select>
          </label>
        </div>

        <div style={{ marginTop: '0.25rem' }}>
          <ImageUpload
            label="Replace listing image"
            file={imageFile}
            onFileChange={setImageFile}
            existingImageUrl={existingImageUrl}
          />
        </div>

        {submitError && <p style={{ margin: 0, color: '#b91c1c', fontSize: '0.95rem' }}>{submitError}</p>}

        <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '0.75rem' }}>
          <Link
            to={`/listing/${numericId}`}
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
            disabled={submitting}
            style={{
              padding: '0.6rem 1.1rem',
              borderRadius: '999px',
              border: 'none',
              backgroundColor: '#2563eb',
              color: '#fff',
              fontWeight: 600,
              cursor: 'pointer',
              opacity: submitting ? 0.8 : 1,
            }}
          >
            {submitting ? 'Saving...' : 'Save changes'}
          </button>
        </div>
      </form>
    </div>
  );
}
