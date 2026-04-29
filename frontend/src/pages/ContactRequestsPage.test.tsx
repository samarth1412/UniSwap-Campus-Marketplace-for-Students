import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { ContactRequestsPage } from './ContactRequestsPage';
import { contactRequestsApi } from '../services/api';

vi.mock('../services/api', async () => {
  const actual = await vi.importActual('../services/api');
  return {
    ...actual,
    contactRequestsApi: {
      getReceived: vi.fn(),
    },
  };
});

describe('ContactRequestsPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows loading then empty state', async () => {
    vi.mocked(contactRequestsApi.getReceived).mockResolvedValueOnce({
      data: {
        success: true,
        data: [],
      },
    } as never);

    render(
      <MemoryRouter>
        <ContactRequestsPage />
      </MemoryRouter>
    );

    expect(screen.getByText('Loading contact requests...')).toBeInTheDocument();
    expect(await screen.findByText('No contact requests yet')).toBeInTheDocument();
  });

  it('renders received contact requests with listing context', async () => {
    vi.mocked(contactRequestsApi.getReceived).mockResolvedValueOnce({
      data: {
        success: true,
        data: [
          {
            id: 10,
            listingId: 8,
            listingTitle: 'Office Chair',
            buyerId: 5,
            buyerName: 'Aarav Singh',
            buyerEmail: 'aarav@campus.edu',
            message: 'Hi, can I pick this up tomorrow?',
            status: 'pending',
            createdAt: '2026-04-29T18:00:00Z',
          },
        ],
      },
    } as never);

    render(
      <MemoryRouter>
        <ContactRequestsPage />
      </MemoryRouter>
    );

    expect(await screen.findByText('Office Chair')).toBeInTheDocument();
    expect(screen.getByText('Aarav Singh (aarav@campus.edu)')).toBeInTheDocument();
    expect(screen.getByText('Hi, can I pick this up tomorrow?')).toBeInTheDocument();
  });

  it('shows backend error response when load fails', async () => {
    vi.mocked(contactRequestsApi.getReceived).mockRejectedValueOnce({
      response: {
        data: {
          error: 'failed to fetch received contact requests',
        },
      },
    } as never);

    render(
      <MemoryRouter>
        <ContactRequestsPage />
      </MemoryRouter>
    );

    expect(await screen.findByText('Could not load contact requests')).toBeInTheDocument();
    expect(screen.getByText('failed to fetch received contact requests')).toBeInTheDocument();
  });

  it('retries loading after an error', async () => {
    vi.mocked(contactRequestsApi.getReceived)
      .mockRejectedValueOnce({
        response: {
          data: {
            error: 'temporary failure',
          },
        },
      } as never)
      .mockResolvedValueOnce({
        data: {
          success: true,
          data: [],
        },
      } as never);

    render(
      <MemoryRouter>
        <ContactRequestsPage />
      </MemoryRouter>
    );

    await screen.findByText('temporary failure');
    fireEvent.click(screen.getByRole('button', { name: 'Retry' }));

    await waitFor(() => {
      expect(contactRequestsApi.getReceived).toHaveBeenCalledTimes(2);
    });
    expect(await screen.findByText('No contact requests yet')).toBeInTheDocument();
  });
});
