import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Route, Routes, useLocation } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { EditListingPage } from './EditListingPage';
import { listingsApi } from '../services/api';

vi.mock('../services/api', async () => {
  const actual = await vi.importActual('../services/api');
  return {
    ...actual,
    listingsApi: {
      getById: vi.fn(),
      update: vi.fn(),
      uploadImages: vi.fn(),
    },
  };
});

function DetailStateProbe() {
  const location = useLocation();
  const state = location.state as { flowMessage?: string } | null;
  return <p>{state?.flowMessage ?? 'no-message'}</p>;
}

describe('EditListingPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listingsApi.getById).mockResolvedValue({
      data: {
        success: true,
        data: {
          id: 1,
          title: 'Sample listing',
          description: 'Existing description',
          price: 10,
          category: 'Books',
          imageUrl: 'https://via.placeholder.com/640x420?text=Listing+Image',
        },
      },
    });
    vi.mocked(listingsApi.update).mockResolvedValue({
      data: { success: true },
    } as never);
    vi.mocked(listingsApi.uploadImages).mockResolvedValue({
      data: { success: true },
    } as never);
  });

  it('prefills the form with listing data from the API', async () => {
    render(
      <MemoryRouter initialEntries={['/listing/1/edit']}>
        <Routes>
          <Route path="/listing/:id/edit" element={<EditListingPage />} />
        </Routes>
      </MemoryRouter>
    );

    await screen.findByDisplayValue('Sample listing');
    expect(screen.getByDisplayValue('Existing description')).toBeInTheDocument();
    expect(screen.getByDisplayValue('10')).toBeInTheDocument();
    expect(screen.getByDisplayValue('Books')).toBeInTheDocument();
    expect(screen.getByText('Current image selected. Choose a new one to replace it.')).toBeInTheDocument();
  });

  it('uploads a replacement image and redirects with success state', async () => {
    const user = userEvent.setup();
    render(
      <MemoryRouter initialEntries={['/listing/1/edit']}>
        <Routes>
          <Route path="/listing/:id/edit" element={<EditListingPage />} />
          <Route path="/listing/:id" element={<DetailStateProbe />} />
        </Routes>
      </MemoryRouter>
    );

    await screen.findByDisplayValue('Sample listing');

    const fileInput = screen.getByLabelText('Replace listing image');
    const replacement = new File(['replacement-bytes'], 'replacement.jpg', { type: 'image/jpeg' });
    await user.upload(fileInput, replacement);

    expect(screen.getByText('Selected: replacement.jpg')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'Save changes' }));

    await waitFor(() => {
      expect(listingsApi.update).toHaveBeenCalledWith(1, {
        title: 'Sample listing',
        description: 'Existing description',
        price: 10,
        category: 'Books',
      });
    });
    expect(listingsApi.uploadImages).toHaveBeenCalledWith(1, [replacement]);
    expect(await screen.findByText('Listing updated successfully.')).toBeInTheDocument();
  });
});
