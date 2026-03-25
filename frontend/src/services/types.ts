export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

export interface AuthResponse {
  token: string;
  user?: { id: string; email: string; full_name?: string };
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
