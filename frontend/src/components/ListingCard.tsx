import { useCallback } from 'react';
import { Link } from 'react-router-dom';
import type { Listing } from '../types/listing';
import { formatUsd } from '../utils/currency';
import { useWishlist } from '../context/useWishlist';
import './ListingCard.css';

type ListingCardProps = {
  listing: Listing;
};

const FALLBACK_LISTING_IMAGE = 'https://placehold.co/640x420?text=No+Image';

function HeartOutlineIcon() {
  return (
    <svg className="listing-card__heart-svg" viewBox="0 0 24 24" fill="none" aria-hidden>
      <path
        d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12z"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinejoin="round"
      />
    </svg>
  );
}

function HeartFilledIcon() {
  return (
    <svg className="listing-card__heart-svg" viewBox="0 0 24 24" aria-hidden>
      <path
        fill="currentColor"
        d="M11.645 20.91l-.007-.003-.022-.012a15.247 15.247 0 01-.383-.218 25.18 25.18 0 01-4.244-3.17C4.688 15.36 2.25 12.174 2.25 8.25 2.25 5.322 4.714 3 7.688 3A5.5 5.5 0 0112 5.052 5.5 5.5 0 0116.313 3c2.973 0 5.437 2.322 5.437 5.25 0 3.925-2.438 7.111-4.739 9.256a25.175 25.175 0 01-4.244 3.17l-.022.012-.007.003-.002.001h-.002z"
      />
    </svg>
  );
}

export function ListingCard({ listing }: ListingCardProps) {
  const { isWishlisted, toggleWishlist, togglingListingId } = useWishlist();
  const wishlisted = isWishlisted(listing.id);
  const toggling = togglingListingId === listing.id;
  const imageSrc = listing.imageUrl?.trim() ? listing.imageUrl : FALLBACK_LISTING_IMAGE;

  const handleWishlistClick = useCallback(
    (e: React.MouseEvent<HTMLButtonElement>) => {
      e.preventDefault();
      e.stopPropagation();
      if (toggling) return;
      void toggleWishlist(listing.id);
    },
    [listing.id, toggleWishlist, toggling]
  );

  return (
    <Link to={`/listing/${listing.id}`} className="listing-card-link">
      <article className="listing-card">
        <div className="listing-card__image-wrap">
          <img
            src={imageSrc}
            alt={listing.title}
            className="listing-card__image"
            loading="lazy"
            onError={(event) => {
              const image = event.currentTarget;
              if (image.src !== FALLBACK_LISTING_IMAGE) {
                image.src = FALLBACK_LISTING_IMAGE;
              }
            }}
          />
          <button
            type="button"
            className={`listing-card__wishlist${wishlisted ? ' listing-card__wishlist--active' : ''}${toggling ? ' listing-card__wishlist--busy' : ''}`}
            onClick={handleWishlistClick}
            disabled={toggling}
            aria-pressed={wishlisted}
            aria-busy={toggling}
            aria-label={wishlisted ? 'Remove from wishlist' : 'Add to wishlist'}
          >
            {wishlisted ? <HeartFilledIcon /> : <HeartOutlineIcon />}
          </button>
        </div>
        <div className="listing-card__body">
          <span className="listing-card__category">{listing.category}</span>
          <h2 className="listing-card__title">{listing.title}</h2>
          <span className="listing-card__price">{formatUsd(listing.price)}</span>
        </div>
      </article>
    </Link>
  );
}
