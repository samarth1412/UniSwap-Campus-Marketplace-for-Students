import { api } from './http';
import type { Listing } from '../types/listing';
import type {
  ApiResponse,
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

export const listingsApi = {
  getAll: (params?: ListingQueryParams) =>
    api.get<ApiResponse<Listing[]>>('/listings', { params }),
  getById: (id: number | string) =>
    api.get<ApiResponse<Listing>>(`/listings/${id}`),
  create: (payload: CreateListingPayload) =>
    api.post<ApiResponse<Listing>>('/listings', payload),
  update: (id: number | string, payload: UpdateListingPayload) =>
    api.put<ApiResponse<Listing>>(`/listings/${id}`, payload),
  remove: (id: number | string) =>
    api.delete<ApiResponse<{ message?: string }>>(`/listings/${id}`),
  report: (id: number | string, reason: string) =>
    api.post<ApiResponse<unknown>>(`/listings/${id}/report`, { reason }),
};
