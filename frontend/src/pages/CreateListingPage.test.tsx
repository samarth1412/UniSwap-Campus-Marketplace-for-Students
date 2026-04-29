import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Route, Routes, useLocation } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { CreateListingPage } from './CreateListingPage';
import { listingsApi } from '../services/api';

vi.mock('../services/api', async () => {
  const actual = await vi.importActual('../services/api');
  return {
    ...actual,
    listingsApi: {
      create: vi.fn(),
      uploadImages: vi.fn(),
    },
  };
});

function DetailStateProbe() {
  const location = useLocation();
  const state = location.state as { flowMessage?: string } | null;
  return <p>{state?.flowMessage ?? 'no-message'}</p>;
}

describe('CreateListingPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listingsApi.create).mockResolvedValue({
      data: {
        success: true,
        data: {
          id: 42,
          title: 'Created listing',
          description: 'Created description',
          price: 30,
          category: 'Books',
          imageUrl: 'https://example.com/listing.jpg',
        },
      },
    } as never);
    vi.mocked(listingsApi.uploadImages).mockResolvedValue({
      data: { success: true },
    } as never);
  });

  it('creates a listing with image upload and forwards success state', async () => {
    const user = userEvent.setup();

    render(
      <MemoryRouter initialEntries={['/create']}>
        <Routes>
          <Route path="/create" element={<CreateListingPage />} />
          <Route path="/listing/:id" element={<DetailStateProbe />} />
        </Routes>
      </MemoryRouter>
    );

    fireEvent.change(screen.getByPlaceholderText('e.g. CS101 textbook, hardly used'), {
      target: { value: 'Graphing Calculator' },
    });
    fireEvent.change(
      screen.getByPlaceholderText('Add details like condition, why you are selling, pickup location, etc.'),
      { target: { value: 'Good condition with case' } }
    );
    fireEvent.change(screen.getByPlaceholderText('e.g. 1500'), { target: { value: '2000' } });

    const input = screen.getByLabelText('Listing photo (optional)');
    const image = new File(['image-content'], 'calculator.png', { type: 'image/png' });
    await user.upload(input, image);
    expect(screen.getByText('Selected: calculator.png')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'Create listing' }));

    await waitFor(() => {
      expect(listingsApi.create).toHaveBeenCalledWith({
        title: 'Graphing Calculator',
        description: 'Good condition with case',
        price: 2000,
        category: 'Books',
      });
    });
    expect(listingsApi.uploadImages).toHaveBeenCalledWith(42, [image]);
    expect(await screen.findByText('Listing created successfully.')).toBeInTheDocument();
  });
});
