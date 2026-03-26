import { api } from './http';
import type { Listing } from '../types/listing';
import type {
  ApiResponse,
  BackendListing,
  CreateListingPayload,
  UpdateListingPayload,
} from './types';

export interface ListingQueryParams {
  search?: string;
  category?: string;
  min_price?: number;
  max_price?: number;
  page?: number;
  limit?: number;
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
  getAll: (params?: ListingQueryParams) =>
    api.get<ApiResponse<Listing[]>>('/listings', { params }),
  getById: (id: number | string) =>
    mapSingleListingResponse(api.get<ApiResponse<BackendListing>>(`/listings/${id}`)),
  create: (payload: CreateListingPayload) =>
    api.post<ApiResponse<Listing>>('/listings', payload),
  update: (id: number | string, payload: UpdateListingPayload) =>
    mapSingleListingResponse(api.put<ApiResponse<BackendListing>>(`/listings/${id}`, payload)),
  remove: (id: number | string) =>
    api.delete<ApiResponse<{ message?: string }>>(`/listings/${id}`),
  report: (id: number | string, reason: string) =>
    api.post<ApiResponse<unknown>>(`/listings/${id}/report`, { reason }),
};
