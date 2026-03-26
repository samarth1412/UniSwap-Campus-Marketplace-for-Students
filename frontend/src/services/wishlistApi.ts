import { api } from './http';
import type { Listing } from '../types/listing';
import type {
  ApiResponse,
  BackendListing,
  WishlistItemResponse,
  WishlistListingResponse,
  WishlistPayload,
} from './types';

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
    imageUrl:
      listing.primary_image_url?.trim() ||
      'https://via.placeholder.com/640x420?text=Listing+Image',
  };
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
