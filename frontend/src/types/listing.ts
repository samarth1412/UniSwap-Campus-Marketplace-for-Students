export type ListingStatus = 'active' | 'sold' | 'hidden';

export interface Listing {
  id: number;
  userId?: number;
  title: string;
  description: string;
  price: number;
  category: string;
  imageUrl: string;
  createdAt?: string;
  sellerName?: string;
  status?: ListingStatus;
}

