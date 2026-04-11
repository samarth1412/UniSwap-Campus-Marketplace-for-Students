import { api } from './http';
import type { Listing } from '../types/listing';
import type {
  ApiResponse,
  BackendListing,
  CreateListingPayload,
  UpdateListingPayload,
} from './types';

const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api';
const ASSET_BASE = API_BASE.replace(/\/api\/?$/, '');

export interface ListingQueryParams {
  search?: string;
  category?: string;
  min_price?: number;
  max_price?: number;
  page?: number;
  limit?: number;
}

export interface PaginatedListings {
  items: Listing[];
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

interface BackendPaginatedListings {
  items: BackendListing[];
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

function isBackendPaginatedListings(value: unknown): value is BackendPaginatedListings {
  return (
    typeof value === 'object' &&
    value !== null &&
    'items' in value &&
    Array.isArray((value as BackendPaginatedListings).items)
  );
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
    const response = await api.get<ApiResponse<BackendPaginatedListings | BackendListing[]>>('/listings', {
      params,
    });

    if (response.data.success && response.data.data) {
      const backend = response.data.data;
      const mapped: PaginatedListings = Array.isArray(backend)
        ? {
            items: backend.map(mapBackendListing),
            page: params?.page ?? 1,
            limit: params?.limit ?? (backend.length || 1),
            total: backend.length,
            total_pages: 1,
          }
        : isBackendPaginatedListings(backend)
          ? {
              items: backend.items.map(mapBackendListing),
              page: backend.page,
              limit: backend.limit,
              total: backend.total,
              total_pages: backend.total_pages,
            }
          : {
              items: [],
              page: 1,
              limit: params?.limit ?? 10,
              total: 0,
              total_pages: 1,
            };
      // Keep response shape consistent with our PaginatedListings type.
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
};
