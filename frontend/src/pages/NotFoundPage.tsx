import { Link } from 'react-router-dom';

export function NotFoundPage() {
  return (
    <div style={{ textAlign: 'center', padding: '2rem' }}>
      <h1>404</h1>
      <p>Page not found.</p>
      <Link to="/">Go to home</Link>
    </div>
  );
}
