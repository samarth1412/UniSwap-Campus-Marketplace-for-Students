import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { MyListingsPage } from './MyListingsPage';
import { profileApi } from '../services/api';

vi.mock('../services/api', async () => {
  const actual = await vi.importActual('../services/api');
  return {
    ...actual,
    profileApi: {
      getMe: vi.fn(),
      getMyListings: vi.fn(),
    },
  };
});

describe('MyListingsPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows the empty state when the current user has no listings', async () => {
    vi.mocked(profileApi.getMe).mockResolvedValue({
      data: {
        success: true,
        data: {
          id: 1,
          full_name: 'Samarth',
          email: 'samarth@example.com',
        },
      },
    } as never);
    vi.mocked(profileApi.getMyListings).mockResolvedValue({
      data: {
        success: true,
        data: [],
      },
    } as never);

    render(
      <MemoryRouter>
        <MyListingsPage />
      </MemoryRouter>
    );

    expect(screen.getByText('Loading your listings...')).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByText('No listings yet')).toBeInTheDocument();
    });
  });
});
