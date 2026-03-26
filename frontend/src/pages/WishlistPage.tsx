import { useMemo } from 'react';
import { Link } from 'react-router-dom';
import { MOCK_LISTINGS } from '../data/mockListings';
import { ListingCard } from '../components/ListingCard';
import { useWishlist } from '../context/WishlistContext';
import './WishlistPage.css';

export function WishlistPage() {
  const { savedIds } = useWishlist();

  // Frontend-only placeholder: match saved ids against mock listings.
  // When backend wishlist APIs are added, replace MOCK_LISTINGS with fetched saved listings.
  const savedIdSet = useMemo(() => new Set(savedIds), [savedIds]);
  const savedListings = useMemo(() => MOCK_LISTINGS.filter((l) => savedIdSet.has(l.id)), [savedIdSet]);

  return (
    <div className="wishlist-page">
      <header className="wishlist-page__header">
        <div className="wishlist-page__header-text">
          <h1 className="wishlist-page__title">Wishlist</h1>
          <p className="wishlist-page__subtitle">
            {savedListings.length > 0
              ? `You saved ${savedListings.length} listing${savedListings.length === 1 ? '' : 's'}.`
              : 'Save listings by tapping the heart on a card.'}
          </p>
        </div>
      </header>

      {savedListings.length === 0 ? (
        <div className="wishlist-page__empty">
          <p className="wishlist-page__empty-title">No wishlist items yet</p>
          <p className="wishlist-page__empty-text">Browse listings and click the heart to add them here.</p>
          <Link to="/" className="wishlist-page__empty-link">
            Browse listings
          </Link>
        </div>
      ) : (
        <section className="wishlist-page__grid" aria-label="Saved listings">
          {savedListings.map((listing) => (
            <ListingCard key={listing.id} listing={listing} />
          ))}
        </section>
      )}
    </div>
  );
}

