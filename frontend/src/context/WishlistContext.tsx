import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';

type WishlistContextValue = {
  savedIds: number[];
  isWishlisted: (listingId: number) => boolean;
  toggleWishlist: (listingId: number) => void;
  clearWishlist: () => void;
};

const STORAGE_KEY = 'wishlist_ids';

function loadSavedIds(): number[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return [];
    const parsed = JSON.parse(raw) as unknown;
    if (!Array.isArray(parsed)) return [];
    const ids = parsed
      .map((v) => (typeof v === 'number' ? v : Number(v)))
      .filter((n) => Number.isFinite(n));
    // Ensure uniqueness while keeping order.
    return Array.from(new Set(ids));
  } catch {
    return [];
  }
}

const WishlistContext = createContext<WishlistContextValue | undefined>(undefined);

export function WishlistProvider({ children }: { children: React.ReactNode }) {
  const [savedIds, setSavedIds] = useState<number[]>(loadSavedIds);

  const savedSet = useMemo(() => new Set(savedIds), [savedIds]);

  useEffect(() => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(savedIds));
  }, [savedIds]);

  const isWishlisted = useCallback(
    (listingId: number) => {
      return savedSet.has(listingId);
    },
    [savedSet]
  );

  const toggleWishlist = useCallback((listingId: number) => {
    setSavedIds((prev) => {
      const next = prev.includes(listingId) ? prev.filter((id) => id !== listingId) : [...prev, listingId];
      return Array.from(new Set(next));
    });
  }, []);

  const clearWishlist = useCallback(() => {
    setSavedIds([]);
  }, []);

  const value = useMemo<WishlistContextValue>(
    () => ({ savedIds, isWishlisted, toggleWishlist, clearWishlist }),
    [savedIds, isWishlisted, toggleWishlist, clearWishlist]
  );

  return <WishlistContext.Provider value={value}>{children}</WishlistContext.Provider>;
}

export function useWishlist() {
  const ctx = useContext(WishlistContext);
  if (!ctx) throw new Error('useWishlist must be used within WishlistProvider');
  return ctx;
}

