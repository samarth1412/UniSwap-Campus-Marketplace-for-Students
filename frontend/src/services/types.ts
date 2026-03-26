export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

export interface AuthResponse {
  token: string;
  user?: BackendUser;
}

export interface RegisterPayload {
  full_name: string;
  email: string;
  university?: string;
  password: string;
}

export interface CreateListingPayload {
  title: string;
  description: string;
  price: number;
  category: string;
  imageUrl?: string;
}

export interface UpdateListingPayload {
  title: string;
  description: string;
  price: number;
  category: string;
}

export interface WishlistPayload {
  listing_id: number;
}

export interface BackendListing {
  id: number;
  user_id: number;
  title: string;
  description: string;
  price: number;
  category: string;
  primary_image_url?: string;
  created_at: string;
}

export interface BackendUser {
  id: number;
  full_name: string;
  email: string;
  university?: string;
  created_at?: string;
  updated_at?: string;
}
