import { api } from './http';
import type { Listing } from '../types/listing';
import type { ApiResponse, BackendListing, BackendUser } from './types';

function mapBackendListing(listing: BackendListing): Listing {
  return {
    id: listing.id,
    title: listing.title,
    description: listing.description,
    price: listing.price,
    category: listing.category,
    imageUrl:
      listing.primary_image_url?.trim() ||
      'https://via.placeholder.com/640x420?text=Listing+Image',
  };
}

export const profileApi = {
  getMe: () => api.get<ApiResponse<BackendUser>>('/auth/me'),
  getMyListings: async (userId: number | string) => {
    const response = await api.get<ApiResponse<BackendListing[]>>(`/users/${userId}/listings`);
    if (response.data.success && response.data.data) {
      response.data.data = response.data.data.map(mapBackendListing) as never;
    }
    return response as unknown as { data: ApiResponse<Listing[]> };
  },
};
