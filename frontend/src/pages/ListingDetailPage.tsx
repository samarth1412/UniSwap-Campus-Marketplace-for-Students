import { useParams, Link } from 'react-router-dom';

/**
 * Listing detail - full page for /listing/:id.
 * Will be wired to API by Bhumi (FE-9).
 */
export function ListingDetailPage() {
  const { id } = useParams<{ id: string }>();
  return (
    <div>
      <p><Link to="/">← Back to listings</Link></p>
      <h1>Listing #{id}</h1>
      <p>Detail view will show title, description, price, image, seller. (Connected in FE-9.)</p>
    </div>
  );
}
