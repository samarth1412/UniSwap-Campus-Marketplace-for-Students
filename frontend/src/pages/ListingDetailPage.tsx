import { useEffect, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { listingsApi } from '../services/api';
import type { Listing } from '../types/listing';

/**
 * Listing detail - full page for /listing/:id.
 * FE-9: Detail UI
 * FE-15: Report listing button UI + API call
 * FE-17/21: API-backed detail and delete flow
 */
export function ListingDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const numericId = Number(id);

  const [listing, setListing] = useState<Listing | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showReportModal, setShowReportModal] = useState(false);
  const [reportReason, setReportReason] = useState('');
  const [reportSubmitting, setReportSubmitting] = useState(false);
  const [reportMessage, setReportMessage] = useState<string | null>(null);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deleteSubmitting, setDeleteSubmitting] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    async function loadListing() {
      if (!Number.isFinite(numericId) || numericId <= 0) {
        if (!cancelled) {
          setError('Listing not found.');
          setLoading(false);
        }
        return;
      }

      setLoading(true);
      setError(null);

      try {
        const response = await listingsApi.getById(numericId);
        if (!response.data.success || !response.data.data) {
          throw new Error(response.data.error || 'Failed to load listing');
        }

        if (!cancelled) {
          setListing(response.data.data);
        }
      } catch (err) {
        console.error('Failed to load listing detail', err);
        if (!cancelled) {
          setError('Unable to load this listing right now.');
          setListing(null);
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

  const handleOpenReport = () => {
    setReportReason('');
    setReportMessage(null);
    setShowReportModal(true);
  };

  const handleCloseReport = () => {
    if (reportSubmitting) return;
    setShowReportModal(false);
  };

  const handleSubmitReport = async (event: React.FormEvent) => {
    event.preventDefault();
    if (!listing) return;

    if (!reportReason.trim()) {
      setReportMessage('Please provide a reason for reporting this listing.');
      return;
    }

    setReportSubmitting(true);
    setReportMessage(null);

    try {
      const response = await listingsApi.report(listing.id, reportReason.trim());
      if (!response.data.success) {
        throw new Error(response.data.error || 'Failed to submit report');
      }
      setReportMessage('Thank you. Your report has been submitted.');
    } catch (err) {
      console.error('Failed to submit report', err);
      setReportMessage('Could not submit report right now. Please try again later.');
    } finally {
      setReportSubmitting(false);
    }
  };

  const handleDelete = async () => {
    if (!listing) return;

    setDeleteSubmitting(true);
    setDeleteError(null);

    try {
      await listingsApi.remove(listing.id);
      navigate('/my-listings', {
        replace: true,
        state: { deletedMessage: 'Listing deleted successfully.' },
      });
    } catch (err: unknown) {
      console.error('Failed to delete listing', err);
      const ax = err as { response?: { data?: { error?: string } } };
      setDeleteError(ax.response?.data?.error ?? 'Could not delete this listing right now.');
    } finally {
      setDeleteSubmitting(false);
    }
  };

  if (loading) {
    return <p>Loading listing...</p>;
  }

  if (!listing) {
    return (
      <div>
        <p>
          <Link to="/">{'< Back to listings'}</Link>
        </p>
        <h1>Listing not found</h1>
        <p>{error ?? 'The listing you are looking for does not exist.'}</p>
      </div>
    );
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
      <p style={{ margin: 0 }}>
        <Link to="/">{'< Back to listings'}</Link>
      </p>
      {error && <p style={{ margin: 0, color: '#b91c1c' }}>{error}</p>}
      <div
        style={{
          display: 'grid',
          gap: '1.5rem',
          gridTemplateColumns: 'minmax(0, 2fr) minmax(0, 3fr)',
        }}
      >
        <div>
          <img
            src={listing.imageUrl}
            alt={listing.title}
            style={{ width: '100%', maxWidth: '480px', borderRadius: '0.75rem', border: '1px solid #eee' }}
          />
        </div>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
          <span
            style={{
              fontSize: '0.8rem',
              textTransform: 'uppercase',
              letterSpacing: '0.04em',
              color: '#6b7280',
            }}
          >
            {listing.category}
          </span>
          <h1 style={{ margin: 0 }}>{listing.title}</h1>
          <span style={{ fontSize: '1.4rem', fontWeight: 700, color: '#059669' }}>
            Rs. {listing.price.toFixed(0)}
          </span>
          <p style={{ maxWidth: '36rem', lineHeight: 1.6 }}>{listing.description}</p>
          <p style={{ marginTop: '0.5rem', color: '#4b5563' }}>
            <strong>Seller:</strong> {listing.sellerName ?? 'Campus student'}
          </p>
          <div style={{ display: 'flex', gap: '0.75rem', marginTop: '1rem', flexWrap: 'wrap' }}>
            <Link
              to={`/listing/${listing.id}/edit`}
              style={{
                padding: '0.5rem 1rem',
                borderRadius: '999px',
                border: '1px solid #2563eb',
                backgroundColor: '#eff6ff',
                color: '#1d4ed8',
                textDecoration: 'none',
              }}
            >
              Edit listing
            </Link>
            <button
              type="button"
              onClick={() => setShowDeleteModal(true)}
              style={{
                padding: '0.5rem 1rem',
                borderRadius: '999px',
                border: '1px solid #dc2626',
                backgroundColor: '#fef2f2',
                color: '#b91c1c',
                cursor: 'pointer',
              }}
            >
              Delete listing
            </button>
            <button
              type="button"
              onClick={handleOpenReport}
              style={{
                padding: '0.5rem 1rem',
                borderRadius: '999px',
                border: '1px solid #f97316',
                backgroundColor: '#fff7ed',
                color: '#c2410c',
                cursor: 'pointer',
              }}
            >
              Report listing
            </button>
          </div>
        </div>
      </div>

      {showDeleteModal && (
        <div
          style={{
            position: 'fixed',
            inset: 0,
            backgroundColor: 'rgba(0,0,0,0.45)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            zIndex: 50,
          }}
        >
          <div
            style={{
              backgroundColor: '#fff',
              padding: '1.5rem',
              borderRadius: '0.75rem',
              width: '100%',
              maxWidth: '420px',
              boxShadow: '0 10px 25px rgba(0,0,0,0.15)',
            }}
          >
            <h2 style={{ marginTop: 0, marginBottom: '0.75rem' }}>Delete listing?</h2>
            <p style={{ marginTop: 0, marginBottom: '0.75rem', color: '#4b5563' }}>
              This action cannot be undone.
            </p>
            {deleteError && <p style={{ margin: 0, color: '#b91c1c' }}>{deleteError}</p>}
            <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '0.5rem', marginTop: '1rem' }}>
              <button
                type="button"
                onClick={() => {
                  if (!deleteSubmitting) {
                    setShowDeleteModal(false);
                    setDeleteError(null);
                  }
                }}
                disabled={deleteSubmitting}
                style={{
                  padding: '0.45rem 0.9rem',
                  borderRadius: '999px',
                  border: '1px solid #e5e7eb',
                  backgroundColor: '#fff',
                  cursor: 'pointer',
                }}
              >
                Cancel
              </button>
              <button
                type="button"
                onClick={handleDelete}
                disabled={deleteSubmitting}
                style={{
                  padding: '0.45rem 0.9rem',
                  borderRadius: '999px',
                  border: 'none',
                  backgroundColor: '#dc2626',
                  color: '#fff',
                  cursor: 'pointer',
                  opacity: deleteSubmitting ? 0.8 : 1,
                }}
              >
                {deleteSubmitting ? 'Deleting...' : 'Delete'}
              </button>
            </div>
          </div>
        </div>
      )}

      {showReportModal && (
        <div
          style={{
            position: 'fixed',
            inset: 0,
            backgroundColor: 'rgba(0,0,0,0.45)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            zIndex: 50,
          }}
        >
          <div
            style={{
              backgroundColor: '#fff',
              padding: '1.5rem',
              borderRadius: '0.75rem',
              width: '100%',
              maxWidth: '420px',
              boxShadow: '0 10px 25px rgba(0,0,0,0.15)',
            }}
          >
            <h2 style={{ marginTop: 0, marginBottom: '0.75rem' }}>Report listing</h2>
            <p style={{ marginTop: 0, marginBottom: '0.75rem', fontSize: '0.9rem', color: '#4b5563' }}>
              Help keep UniSwap safe by telling us what&apos;s wrong with this listing.
            </p>
            <form onSubmit={handleSubmitReport} style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
              <label style={{ fontSize: '0.9rem', fontWeight: 500 }}>
                Reason
                <select
                  value={reportReason}
                  onChange={(e) => setReportReason(e.target.value)}
                  style={{
                    marginTop: '0.25rem',
                    width: '100%',
                    padding: '0.5rem 0.75rem',
                    borderRadius: '0.5rem',
                    border: '1px solid #e5e7eb',
                  }}
                >
                  <option value="">Select a reason</option>
                  <option value="Spam or misleading">Spam or misleading</option>
                  <option value="Inappropriate or offensive">Inappropriate or offensive</option>
                  <option value="Suspicious or scam">Suspicious or scam</option>
                  <option value="Other">Other</option>
                </select>
              </label>
              {reportMessage && (
                <p
                  style={{
                    margin: 0,
                    fontSize: '0.85rem',
                    color: reportMessage.startsWith('Thank') ? '#047857' : '#b91c1c',
                  }}
                >
                  {reportMessage}
                </p>
              )}
              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '0.5rem', marginTop: '0.75rem' }}>
                <button
                  type="button"
                  onClick={handleCloseReport}
                  disabled={reportSubmitting}
                  style={{
                    padding: '0.45rem 0.9rem',
                    borderRadius: '999px',
                    border: '1px solid #e5e7eb',
                    backgroundColor: '#fff',
                    cursor: 'pointer',
                  }}
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={reportSubmitting}
                  style={{
                    padding: '0.45rem 0.9rem',
                    borderRadius: '999px',
                    border: 'none',
                    backgroundColor: '#dc2626',
                    color: '#fff',
                    cursor: 'pointer',
                    opacity: reportSubmitting ? 0.8 : 1,
                  }}
                >
                  {reportSubmitting ? 'Submitting...' : 'Submit report'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
