export type ListingStatus = 'active' | 'sold' | 'hidden';

export interface Listing {
  id: number;
  title: string;
  description: string;
  price: number;
  category: string;
  imageUrl: string;
  sellerName?: string;
  status?: ListingStatus;
}

