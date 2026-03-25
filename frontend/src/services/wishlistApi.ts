import { api } from './http';
import type { Listing } from '../types/listing';
import type { ApiResponse, WishlistPayload } from './types';

export const wishlistApi = {
  getAll: () => api.get<ApiResponse<Listing[]>>('/wishlist'),
  create: (payload: WishlistPayload) =>
    api.post<ApiResponse<unknown>>('/wishlist', payload),
  remove: (id: number | string) =>
    api.delete<ApiResponse<unknown>>(`/wishlist/${id}`),
};
