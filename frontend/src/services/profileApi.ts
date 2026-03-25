import { api } from './http';
import type { Listing } from '../types/listing';
import type { ApiResponse } from './types';

export const profileApi = {
  getMe: () => api.get<ApiResponse<unknown>>('/auth/me'),
  getMyListings: (userId: number | string) =>
    api.get<ApiResponse<Listing[]>>(`/users/${userId}/listings`),
};
