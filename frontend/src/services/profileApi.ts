import { api } from './http';
import type { Listing } from '../types/listing';
import type { ApiResponse, BackendListing, BackendUser } from './types';

const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api';
const ASSET_BASE = API_BASE.replace(/\/api\/?$/, '');

function mapBackendListing(listing: BackendListing): Listing {
  return {
    id: listing.id,
    title: listing.title,
    description: listing.description,
    price: listing.price,
    category: listing.category,
    imageUrl: listingImageUrl(listing.primary_image_url),
  };
}

function listingImageUrl(imageUrl?: string): string {
  const trimmed = imageUrl?.trim();
  if (!trimmed) return 'https://placehold.co/640x420?text=Listing+Image';
  if (/^https?:\/\//i.test(trimmed)) return trimmed;
  return `${ASSET_BASE}${trimmed.startsWith('/') ? '' : '/'}${trimmed}`;
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
