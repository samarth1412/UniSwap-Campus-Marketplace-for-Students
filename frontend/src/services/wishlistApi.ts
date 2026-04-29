import { api } from './http';
import type { Listing } from '../types/listing';
import type {
  ApiResponse,
  BackendListing,
  WishlistItemResponse,
  WishlistListingResponse,
  WishlistPayload,
} from './types';

const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api';
const ASSET_BASE = API_BASE.replace(/\/api\/?$/, '');

export interface WishlistEntry {
  wishlistId: number;
  listing: Listing;
}

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

function mapWishlistRow(row: WishlistListingResponse): WishlistEntry {
  return {
    wishlistId: row.wishlist_id,
    listing: mapBackendListing(row.listing),
  };
}

export const wishlistApi = {
  getAll: async (): Promise<{ data: ApiResponse<WishlistEntry[]> }> => {
    const res = await api.get<ApiResponse<WishlistListingResponse[]>>('/wishlist');
    const body: ApiResponse<WishlistEntry[]> = {
      success: res.data.success,
      error: res.data.error,
      data: res.data.data?.map(mapWishlistRow),
    };
    return { data: body };
  },

  create: (payload: WishlistPayload) =>
    api.post<ApiResponse<WishlistItemResponse>>('/wishlist', payload),

  remove: (wishlistRowId: number | string) =>
    api.delete<ApiResponse<{ message?: string }>>(`/wishlist/${wishlistRowId}`),
};
