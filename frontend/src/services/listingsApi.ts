import { api } from './http';
import type { Listing } from '../types/listing';
import type {
  ApiResponse,
  BackendListing,
  BackendPaginatedListings,
  CreateListingPayload,
  ListingQueryParams,
  UpdateListingPayload,
} from './types';

const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api';
const ASSET_BASE = API_BASE.replace(/\/api\/?$/, '');

export interface PaginatedListings {
  items: Listing[];
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

function mapBackendListing(listing: BackendListing): Listing {
  return {
    id: listing.id,
    userId: listing.user_id,
    title: listing.title,
    description: listing.description,
    price: listing.price,
    category: listing.category,
    imageUrl: listingImageUrl(listing.primary_image_url),
    createdAt: listing.created_at,
  };
}

function listingImageUrl(imageUrl?: string): string {
  const trimmed = imageUrl?.trim();
  if (!trimmed) return 'https://placehold.co/640x420?text=Listing+Image';
  if (/^https?:\/\//i.test(trimmed)) return trimmed;
  return `${ASSET_BASE}${trimmed.startsWith('/') ? '' : '/'}${trimmed}`;
}

async function mapSingleListingResponse(
  request: Promise<{ data: ApiResponse<BackendListing> }>
): Promise<{ data: ApiResponse<Listing> }> {
  const response = await request;
  if (response.data.success && response.data.data) {
    response.data.data = mapBackendListing(response.data.data) as never;
  }
  return response as unknown as { data: ApiResponse<Listing> };
}

export const listingsApi = {
  getAll: async (
    params?: ListingQueryParams
  ): Promise<{ data: ApiResponse<PaginatedListings> }> => {
    const response = await api.get<ApiResponse<BackendPaginatedListings>>('/listings', {
      params,
    });

    if (response.data.success && response.data.data) {
      const backend = response.data.data;
      const mapped: PaginatedListings = {
        items: backend.items.map(mapBackendListing),
        page: backend.page,
        limit: backend.limit,
        total: backend.total,
        total_pages: backend.total_pages,
      };
      response.data.data = mapped as never;
    }

    return response as unknown as { data: ApiResponse<PaginatedListings> };
  },
  getById: (id: number | string) =>
    mapSingleListingResponse(api.get<ApiResponse<BackendListing>>(`/listings/${id}`)),
  create: (payload: CreateListingPayload) =>
    mapSingleListingResponse(api.post<ApiResponse<BackendListing>>('/listings', payload)),
  update: (id: number | string, payload: UpdateListingPayload) =>
    mapSingleListingResponse(api.put<ApiResponse<BackendListing>>(`/listings/${id}`, payload)),
  uploadImages: (id: number | string, files: File[]) => {
    const formData = new FormData();
    for (const file of files) {
      formData.append('files', file);
    }
    return api.post<ApiResponse<unknown>>(`/listings/${id}/images`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
  },
  remove: (id: number | string) =>
    api.delete<ApiResponse<{ message?: string }>>(`/listings/${id}`),
  report: (id: number | string, reason: string) =>
    api.post<ApiResponse<unknown>>(`/listings/${id}/report`, { reason }),
  createContactRequest: (id: number | string, message: string) =>
    api.post<ApiResponse<unknown>>(`/listings/${id}/contact`, { message }),
};
