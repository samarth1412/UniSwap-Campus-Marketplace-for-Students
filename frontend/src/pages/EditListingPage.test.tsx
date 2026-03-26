import { render, screen } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
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
    },
  };
});

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
  });
});
