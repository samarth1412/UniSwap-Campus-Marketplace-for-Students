import { Link } from 'react-router-dom';
import { ListingCard } from '../components/ListingCard';
import { isAuthenticated } from '../hooks/useAuth';
import { useWishlist } from '../context/WishlistContext';
import './WishlistPage.css';

export function WishlistPage() {
  const { wishlistListings, loading, error } = useWishlist();
  const authed = isAuthenticated();

  if (!authed) {
    return (
      <div className="wishlist-page">
        <header className="wishlist-page__header">
          <div className="wishlist-page__header-text">
            <h1 className="wishlist-page__title">Wishlist</h1>
            <p className="wishlist-page__subtitle">Sign in to view and manage saved listings.</p>
          </div>
        </header>
        <div className="wishlist-page__empty">
          <p className="wishlist-page__empty-title">Log in to use your wishlist</p>
          <p className="wishlist-page__empty-text">Your saved items are stored on your account.</p>
          <Link to="/login" className="wishlist-page__empty-link">
            Log in
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="wishlist-page">
      <header className="wishlist-page__header">
        <div className="wishlist-page__header-text">
          <h1 className="wishlist-page__title">Wishlist</h1>
          <p className="wishlist-page__subtitle">
            {loading && wishlistListings.length === 0
              ? 'Loading your saved listings…'
              : wishlistListings.length > 0
                ? `You saved ${wishlistListings.length} listing${wishlistListings.length === 1 ? '' : 's'}.`
                : 'Save listings by tapping the heart on a card.'}
          </p>
        </div>
      </header>

      {error && (
        <div className="wishlist-page__error" role="alert">
          {error}
        </div>
      )}

      {loading && wishlistListings.length === 0 && !error ? (
        <p className="wishlist-page__loading" aria-live="polite">
          Loading wishlist…
        </p>
      ) : wishlistListings.length === 0 ? (
        <div className="wishlist-page__empty">
          <p className="wishlist-page__empty-title">No wishlist items yet</p>
          <p className="wishlist-page__empty-text">Browse listings and click the heart to add them here.</p>
          <Link to="/" className="wishlist-page__empty-link">
            Browse listings
          </Link>
        </div>
      ) : (
        <section className="wishlist-page__grid" aria-label="Saved listings">
          {wishlistListings.map((listing) => (
            <ListingCard key={listing.id} listing={listing} />
          ))}
        </section>
      )}
    </div>
  );
}
