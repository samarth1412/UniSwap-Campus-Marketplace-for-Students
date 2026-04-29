import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { HomePage } from './HomePage';
import { listingsApi } from '../services/api';

vi.mock('../context/useWishlist', () => ({
  useWishlist: () => ({
    isWishlisted: () => false,
    toggleWishlist: vi.fn(),
    togglingListingId: null,
  }),
}));

vi.mock('../services/api', async () => {
  const actual = await vi.importActual('../services/api');
  return {
    ...actual,
    listingsApi: {
      getAll: vi.fn(),
    },
  };
});

describe('HomePage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('renders listings returned by the paginated API', async () => {
    vi.mocked(listingsApi.getAll).mockResolvedValue({
      data: {
        success: true,
        data: {
          items: [
            {
              id: 1,
              title: 'Algorithms Textbook',
              description: 'Clean copy',
              price: 25,
              category: 'Books',
              imageUrl: 'https://example.com/book.jpg',
            },
          ],
          page: 1,
          limit: 10,
          total: 1,
          total_pages: 1,
        },
      },
    } as never);

    render(
      <MemoryRouter>
        <HomePage />
      </MemoryRouter>
    );

    await screen.findByText('Algorithms Textbook');

    expect(listingsApi.getAll).toHaveBeenCalledWith({
      search: undefined,
      category: undefined,
      min_price: undefined,
      max_price: undefined,
      page: 1,
      limit: 10,
    });
    expect(screen.getByText('1 listing found')).toBeInTheDocument();
  });

  it('sends the debounced search term to the listings endpoint', async () => {
    vi.mocked(listingsApi.getAll).mockResolvedValue({
      data: {
        success: true,
        data: {
          items: [],
          page: 1,
          limit: 10,
          total: 0,
          total_pages: 1,
        },
      },
    } as never);

    render(
      <MemoryRouter>
        <HomePage />
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(listingsApi.getAll).toHaveBeenCalledTimes(1);
    });

    fireEvent.change(screen.getByLabelText('Search listings'), {
      target: { value: 'laptop' },
    });

    await act(async () => {
      await new Promise((resolve) => window.setTimeout(resolve, 350));
    });

    await waitFor(() => {
      expect(listingsApi.getAll).toHaveBeenLastCalledWith({
        search: 'laptop',
        category: undefined,
        min_price: undefined,
        max_price: undefined,
        page: 1,
        limit: 10,
      });
    });
  });

  it('requests the next backend page when pagination advances', async () => {
    vi.mocked(listingsApi.getAll)
      .mockResolvedValueOnce({
        data: {
          success: true,
          data: {
            items: [
              {
                id: 1,
                title: 'Desk Lamp',
                description: 'Warm light',
                price: 15,
                category: 'Other',
                imageUrl: 'https://example.com/lamp.jpg',
              },
            ],
            page: 1,
            limit: 10,
            total: 12,
            total_pages: 2,
          },
        },
      } as never)
      .mockResolvedValueOnce({
        data: {
          success: true,
          data: {
            items: [
              {
                id: 2,
                title: 'Monitor Stand',
                description: 'Wood stand',
                price: 20,
                category: 'Furniture',
                imageUrl: 'https://example.com/stand.jpg',
              },
            ],
            page: 2,
            limit: 10,
            total: 12,
            total_pages: 2,
          },
        },
      } as never);

    render(
      <MemoryRouter>
        <HomePage />
      </MemoryRouter>
    );

    await screen.findByText('Desk Lamp');

    fireEvent.click(screen.getByRole('button', { name: 'Next' }));

    await waitFor(() => {
      expect(listingsApi.getAll).toHaveBeenLastCalledWith({
        search: undefined,
        category: undefined,
        min_price: undefined,
        max_price: undefined,
        page: 2,
        limit: 10,
      });
    });

    expect(await screen.findByText('Monitor Stand')).toBeInTheDocument();
    expect(screen.getByText('Page 2 of 2')).toBeInTheDocument();
  });
});
