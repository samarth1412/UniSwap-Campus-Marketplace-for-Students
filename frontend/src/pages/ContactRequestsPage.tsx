import { useCallback, useEffect, useState } from 'react';
import { contactRequestsApi } from '../services/api';
import type { ReceivedContactRequest } from '../services/contactRequestsApi';

function formatTimestamp(value: string): string {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return new Intl.DateTimeFormat(undefined, {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(date);
}

export function ContactRequestsPage() {
  const [requests, setRequests] = useState<ReceivedContactRequest[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadRequests = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await contactRequestsApi.getReceived();
      if (!response.data.success) {
        throw new Error(response.data.error || 'Failed to load contact requests');
      }
      setRequests(response.data.data ?? []);
    } catch (err: unknown) {
      const ax = err as { response?: { data?: { error?: string } } };
      setError(ax.response?.data?.error ?? 'Unable to load contact requests right now.');
      setRequests([]);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadRequests();
  }, [loadRequests]);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
      <header style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
        <h1 style={{ margin: 0 }}>Contact Requests</h1>
        <p style={{ margin: 0, color: '#4b5563' }}>
          Review messages buyers sent about your listings.
        </p>
      </header>

      {loading ? (
        <p style={{ margin: 0 }}>Loading contact requests...</p>
      ) : error ? (
        <div
          style={{
            border: '1px solid #fecaca',
            borderRadius: '0.75rem',
            padding: '1.25rem',
            background: '#fff1f2',
          }}
        >
          <h2 style={{ marginTop: 0, marginBottom: '0.5rem' }}>Could not load contact requests</h2>
          <p style={{ marginTop: 0, color: '#7f1d1d' }}>{error}</p>
          <button
            type="button"
            onClick={() => void loadRequests()}
            style={{
              padding: '0.5rem 0.95rem',
              borderRadius: '999px',
              border: '1px solid #d1d5db',
              color: '#111827',
              backgroundColor: '#fff',
              cursor: 'pointer',
            }}
          >
            Retry
          </button>
        </div>
      ) : requests.length === 0 ? (
        <div style={{ border: '1px solid #e5e7eb', borderRadius: '0.75rem', padding: '1.5rem' }}>
          <h2 style={{ marginTop: 0 }}>No contact requests yet</h2>
          <p style={{ color: '#4b5563', margin: 0 }}>
            When buyers contact you about your listings, their messages will appear here.
          </p>
        </div>
      ) : (
        <section style={{ display: 'grid', gap: '1rem' }} aria-label="Received contact requests">
          {requests.map((request) => (
            <article
              key={request.id}
              style={{
                border: '1px solid #e5e7eb',
                borderRadius: '0.75rem',
                padding: '1rem',
                backgroundColor: '#fff',
                display: 'flex',
                flexDirection: 'column',
                gap: '0.5rem',
              }}
            >
              <p style={{ margin: 0, color: '#1f2937' }}>
                <strong>Listing:</strong> {request.listingTitle}
              </p>
              <p style={{ margin: 0, color: '#1f2937' }}>
                <strong>Buyer:</strong> {request.buyerName} ({request.buyerEmail})
              </p>
              <p style={{ margin: 0, color: '#374151', lineHeight: 1.5 }}>{request.message}</p>
              <p style={{ margin: 0, color: '#6b7280', fontSize: '0.9rem' }}>
                Received {formatTimestamp(request.createdAt)}
              </p>
            </article>
          ))}
        </section>
      )}
    </div>
  );
}
