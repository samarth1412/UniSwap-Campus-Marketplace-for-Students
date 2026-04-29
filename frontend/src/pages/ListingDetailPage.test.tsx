import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes, useLocation } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { ListingDetailPage } from './ListingDetailPage';
import { listingsApi, profileApi } from '../services/api';

vi.mock('../services/api', async () => {
  const actual = await vi.importActual('../services/api');
  return {
    ...actual,
    listingsApi: {
      getById: vi.fn(),
      remove: vi.fn(),
      report: vi.fn(),
    },
    profileApi: {
      getMe: vi.fn(),
      getMyListings: vi.fn(),
    },
  };
});

function MyListingsStateProbe() {
  const location = useLocation();
  const state = location.state as { deletedMessage?: string } | null;
  return <p>{state?.deletedMessage ?? 'no-delete-message'}</p>;
}

function LoginStateProbe() {
  const location = useLocation();
  const state = location.state as { from?: { pathname?: string } } | null;
  return <p>login-from:{state?.from?.pathname ?? 'none'}</p>;
}

describe('ListingDetailPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
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
    vi.mocked(profileApi.getMe).mockResolvedValue({
      data: {
        success: true,
        data: {
          id: 7,
          full_name: 'Samarth',
          email: 'samarth@example.com',
        },
      },
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

  it('hides the contact seller action for the listing owner', async () => {
    localStorage.setItem('token', 'seller-token');
    vi.mocked(listingsApi.getById).mockResolvedValueOnce({
      data: {
        success: true,
        data: {
          id: 1,
          userId: 7,
          title: 'Desk Lamp',
          description: 'Warm white light',
          price: 15,
          category: 'Other',
          imageUrl: 'https://example.com/lamp.jpg',
          sellerName: 'Samarth',
        },
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/listing/1']}>
        <Routes>
          <Route path="/listing/:id" element={<ListingDetailPage />} />
        </Routes>
      </MemoryRouter>
    );

    await screen.findByText('Desk Lamp');
    await waitFor(() => {
      expect(screen.queryByRole('button', { name: 'Contact Seller' })).not.toBeInTheDocument();
    });
  });

  it('sends logged-out users to login when contacting a seller', async () => {
    render(
      <MemoryRouter initialEntries={['/listing/1']}>
        <Routes>
          <Route path="/listing/:id" element={<ListingDetailPage />} />
          <Route path="/login" element={<LoginStateProbe />} />
        </Routes>
      </MemoryRouter>
    );

    await screen.findByText('Desk Lamp');
    fireEvent.click(screen.getByRole('button', { name: 'Contact Seller' }));

    expect(await screen.findByText('login-from:/listing/1')).toBeInTheDocument();
  });

  it('opens the contact seller form for logged-in buyers', async () => {
    localStorage.setItem('token', 'buyer-token');
    vi.mocked(listingsApi.getById).mockResolvedValueOnce({
      data: {
        success: true,
        data: {
          id: 1,
          userId: 12,
          title: 'Desk Lamp',
          description: 'Warm white light',
          price: 15,
          category: 'Other',
          imageUrl: 'https://example.com/lamp.jpg',
          sellerName: 'Bhumi',
        },
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/listing/1']}>
        <Routes>
          <Route path="/listing/:id" element={<ListingDetailPage />} />
        </Routes>
      </MemoryRouter>
    );

    await screen.findByText('Desk Lamp');
    fireEvent.click(await screen.findByRole('button', { name: 'Contact Seller' }));

    expect(screen.getByRole('heading', { name: 'Contact Seller' })).toBeInTheDocument();
    expect(screen.getByLabelText('Message')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Send message' })).toBeInTheDocument();
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
