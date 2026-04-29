import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes, useLocation } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { ListingDetailPage } from './ListingDetailPage';
import { listingsApi } from '../services/api';

vi.mock('../services/api', async () => {
  const actual = await vi.importActual('../services/api');
  return {
    ...actual,
    listingsApi: {
      getById: vi.fn(),
      remove: vi.fn(),
      report: vi.fn(),
    },
  };
});

function MyListingsStateProbe() {
  const location = useLocation();
  const state = location.state as { deletedMessage?: string } | null;
  return <p>{state?.deletedMessage ?? 'no-delete-message'}</p>;
}

describe('ListingDetailPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listingsApi.getById).mockResolvedValue({
      data: {
        success: true,
        data: {
          id: 1,
          title: 'Desk Lamp',
          description: 'Warm white light',
          price: 15,
          category: 'Other',
          imageUrl: 'https://example.com/lamp.jpg',
          sellerName: 'Bhumi',
        },
      },
    } as never);
    vi.mocked(listingsApi.remove).mockResolvedValue({
      data: { success: true },
    } as never);
  });

  it('shows and dismisses one-time flow message from navigation state', async () => {
    render(
      <MemoryRouter initialEntries={[{ pathname: '/listing/1', state: { flowMessage: 'Listing updated successfully.' } }]}>
        <Routes>
          <Route path="/listing/:id" element={<ListingDetailPage />} />
        </Routes>
      </MemoryRouter>
    );

    await screen.findByText('Desk Lamp');
    expect(screen.getByText('Listing updated successfully.')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'Dismiss' }));
    expect(screen.queryByText('Listing updated successfully.')).not.toBeInTheDocument();
  });

  it('shows the contact seller action on listing details', async () => {
    render(
      <MemoryRouter initialEntries={['/listing/1']}>
        <Routes>
          <Route path="/listing/:id" element={<ListingDetailPage />} />
        </Routes>
      </MemoryRouter>
    );

    await screen.findByText('Desk Lamp');
    expect(screen.getByRole('button', { name: 'Contact Seller' })).toBeInTheDocument();
  });

  it('navigates to my listings with deleted-state message after delete confirm', async () => {
    render(
      <MemoryRouter initialEntries={['/listing/1']}>
        <Routes>
          <Route path="/listing/:id" element={<ListingDetailPage />} />
          <Route path="/my-listings" element={<MyListingsStateProbe />} />
        </Routes>
      </MemoryRouter>
    );

    await screen.findByText('Desk Lamp');

    fireEvent.click(screen.getByRole('button', { name: 'Delete listing' }));
    fireEvent.click(screen.getByRole('button', { name: 'Delete' }));

    await waitFor(() => {
      expect(listingsApi.remove).toHaveBeenCalledWith(1);
    });
    expect(await screen.findByText('Listing deleted successfully.')).toBeInTheDocument();
  });
});
