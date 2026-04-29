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
}

export interface UpdateListingPayload {
  title: string;
  description: string;
  price: number;
  category: string;
}

export interface ListingQueryParams {
  search?: string;
  category?: string;
  min_price?: number;
  max_price?: number;
  page?: number;
  limit?: number;
}

export interface WishlistPayload {
  listing_id: number;
}

/** Response body for POST /wishlist (wishlist row id is used for DELETE /wishlist/:id). */
export interface WishlistItemResponse {
  id: number;
  user_id: number;
  listing_id: number;
  created_at: string;
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

export interface BackendPaginatedListings {
  items: BackendListing[];
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

/** Single row from GET /wishlist. */
export interface WishlistListingResponse {
  wishlist_id: number;
  saved_at: string;
  listing: BackendListing;
}

export interface BackendUser {
  id: number;
  full_name: string;
  email: string;
  university?: string;
  created_at?: string;
  updated_at?: string;
}

export interface BackendReceivedContactRequest {
  id: number;
  listing_id: number;
  listing_title: string;
  buyer_id: number;
  buyer_name: string;
  buyer_email: string;
  message: string;
  status: string;
  created_at: string;
}
