import { createContext } from 'react';
import type { Listing } from '../types/listing';

export type WishlistContextValue = {
  savedIds: number[];
  wishlistListings: Listing[];
  loading: boolean;
  error: string | null;
  togglingListingId: number | null;
  isWishlisted: (listingId: number) => boolean;
  toggleWishlist: (listingId: number) => Promise<void>;
  refreshWishlist: () => Promise<void>;
};

export const WishlistContext = createContext<WishlistContextValue | undefined>(undefined);
