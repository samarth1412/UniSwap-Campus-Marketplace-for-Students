import axios from 'axios';
import {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from 'react';
import { useLocation } from 'react-router-dom';
import type { Listing } from '../types/listing';
import { getToken } from '../hooks/useAuth';
import { wishlistApi } from '../services/wishlistApi';
import type { ApiResponse } from '../services/types';
import { WishlistContext, type WishlistContextValue } from './WishlistContextCore';

function getErrorMessage(err: unknown): string {
  if (axios.isAxiosError(err)) {
    const data = err.response?.data as ApiResponse<unknown> | undefined;
    if (data?.error && typeof data.error === 'string') return data.error;
    return err.message;
  }
  if (err instanceof Error) return err.message;
  return 'Something went wrong';
}

export function WishlistProvider({ children }: { children: React.ReactNode }) {
  const location = useLocation();
  const [wishlistIdByListingId, setWishlistIdByListingId] = useState<Map<number, number>>(
    () => new Map()
  );
  const [wishlistListings, setWishlistListings] = useState<Listing[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [togglingListingId, setTogglingListingId] = useState<number | null>(null);

  const wishlistIdByListingIdRef = useRef(wishlistIdByListingId);
  wishlistIdByListingIdRef.current = wishlistIdByListingId;

  const refreshWishlist = useCallback(async () => {
    if (!getToken()) {
      setWishlistIdByListingId(new Map());
      setWishlistListings([]);
      setError(null);
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const res = await wishlistApi.getAll();
      if (!res.data.success || !res.data.data) {
        setError(res.data.error ?? 'Failed to load wishlist');
        setWishlistIdByListingId(new Map());
        setWishlistListings([]);
        return;
      }
      const map = new Map<number, number>();
      const listings: Listing[] = [];
      for (const entry of res.data.data) {
        map.set(entry.listing.id, entry.wishlistId);
        listings.push(entry.listing);
      }
      setWishlistIdByListingId(map);
      setWishlistListings(listings);
    } catch (err) {
      setError(getErrorMessage(err));
      setWishlistIdByListingId(new Map());
      setWishlistListings([]);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void refreshWishlist();
  }, [location.pathname, refreshWishlist]);

  const savedSet = useMemo(() => new Set(wishlistIdByListingId.keys()), [wishlistIdByListingId]);

  const savedIds = useMemo(() => Array.from(savedSet), [savedSet]);

  const isWishlisted = useCallback(
    (listingId: number) => savedSet.has(listingId),
    [savedSet]
  );

  const toggleWishlist = useCallback(
    async (listingId: number) => {
      if (!getToken()) {
        window.location.href = '/login';
        return;
      }

      setTogglingListingId(listingId);
      setError(null);
      try {
        const wishlistRowId = wishlistIdByListingIdRef.current.get(listingId);
        if (wishlistRowId !== undefined) {
          await wishlistApi.remove(wishlistRowId);
        } else {
          try {
            await wishlistApi.create({ listing_id: listingId });
          } catch (err) {
            if (axios.isAxiosError(err) && err.response?.status === 409) {
              // Already saved — sync with server
            } else {
              throw err;
            }
          }
        }
        await refreshWishlist();
      } catch (err) {
        setError(getErrorMessage(err));
        await refreshWishlist();
      } finally {
        setTogglingListingId(null);
      }
    },
    [refreshWishlist]
  );

  const value = useMemo<WishlistContextValue>(
    () => ({
      savedIds,
      wishlistListings,
      loading,
      error,
      togglingListingId,
      isWishlisted,
      toggleWishlist,
      refreshWishlist,
    }),
    [
      savedIds,
      wishlistListings,
      loading,
      error,
      togglingListingId,
      isWishlisted,
      toggleWishlist,
      refreshWishlist,
    ]
  );

  return <WishlistContext.Provider value={value}>{children}</WishlistContext.Provider>;
}
