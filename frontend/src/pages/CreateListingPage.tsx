import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api, type ApiResponse } from '../services/api';
import { ImageUpload } from '../components/ImageUpload';

interface CreateListingPayload {
  title: string;
  description: string;
  price: number;
  category: string;
  // For Sprint 1 we keep image as optional string; backend can evolve later.
  imageUrl?: string;
}

/**
 * Create listing page - protected route.
 * FE-10: Form UI
 * FE-11: Image upload component
 * FE-12: Connect to backend
 */
export function CreateListingPage() {
  const navigate = useNavigate();
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [price, setPrice] = useState('');
  const [category, setCategory] = useState('Books');
  const [imageFile, setImageFile] = useState<File | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    setError(null);

    if (!title.trim() || !description.trim() || !price.trim()) {
      setError('Please fill in all required fields.');
      return;
    }

    const numericPrice = Number(price);
    if (Number.isNaN(numericPrice) || numericPrice < 0) {
      setError('Please enter a valid price.');
      return;
    }

    // For Sprint 1, we send a simple JSON payload and ignore the image file,
    // or you can coordinate with backend to upload image separately.
    const payload: CreateListingPayload = {
      title: title.trim(),
      description: description.trim(),
      price: numericPrice,
      category,
      imageUrl: '', // backend can derive from separate upload in future sprints
    };

    setSubmitting(true);
    try {
      const response = await api.post<ApiResponse<unknown>>('/listings', payload);
      if (!response.data.success) {
        throw new Error(response.data.error || 'Failed to create listing');
      }
      navigate('/');
    } catch (err) {
      console.error('Failed to create listing', err);
      setError('Could not create listing right now. Please try again.');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div style={{ maxWidth: '640px' }}>
      <h1>Create Listing</h1>
      <p style={{ color: '#4b5563', marginBottom: '1rem' }}>
        Share an item with your campus community. Fill in the details below.
      </p>
      <form
        onSubmit={handleSubmit}
        style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}
      >
        <label style={{ display: 'flex', flexDirection: 'column', gap: '0.25rem' }}>
          <span>Title</span>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
            placeholder="e.g. CS101 textbook, hardly used"
            style={{
              padding: '0.5rem 0.75rem',
              borderRadius: '0.5rem',
              border: '1px solid #e5e7eb',
            }}
          />
        </label>

        <label style={{ display: 'flex', flexDirection: 'column', gap: '0.25rem' }}>
          <span>Description</span>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            required
            rows={4}
            placeholder="Add details like condition, why you are selling, pickup location, etc."
            style={{
              padding: '0.5rem 0.75rem',
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
              flex: '1 1 160px',
            }}
          >
            <span>Price ($)</span>
            <input
              type="number"
              min={0}
              step="1"
              value={price}
              onChange={(e) => setPrice(e.target.value)}
              required
              placeholder="e.g. 1500"
              style={{
                padding: '0.5rem 0.75rem',
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
              flex: '1 1 160px',
            }}
          >
            <span>Category</span>
            <select
              value={category}
              onChange={(e) => setCategory(e.target.value)}
              style={{
                padding: '0.5rem 0.75rem',
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

        <div style={{ marginTop: '0.5rem' }}>
          <ImageUpload file={imageFile} onFileChange={setImageFile} />
        </div>

        {error && (
          <p style={{ margin: 0, color: '#b91c1c', fontSize: '0.9rem' }}>
            {error}
          </p>
        )}

        <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '0.75rem' }}>
          <button
            type="submit"
            disabled={submitting}
            style={{
              padding: '0.55rem 1.4rem',
              borderRadius: '999px',
              border: 'none',
              backgroundColor: '#16a34a',
              color: '#fff',
              fontWeight: 600,
              cursor: 'pointer',
              opacity: submitting ? 0.8 : 1,
            }}
          >
            {submitting ? 'Creating…' : 'Create listing'}
          </button>
        </div>
      </form>
    </div>
  );
}

