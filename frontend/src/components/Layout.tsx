import { Link, useNavigate } from 'react-router-dom';
import { isAuthenticated, clearToken } from '../hooks/useAuth';

export function Layout({ children }: { children: React.ReactNode }) {
  const navigate = useNavigate();
  const authenticated = isAuthenticated();

  const handleLogout = () => {
    clearToken();
    navigate('/login');
  };

  return (
    <div style={{ minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
      <header style={{ width: '100%', borderBottom: '1px solid #ddd', boxSizing: 'border-box' }}>
        <nav style={{
          padding: '1rem 1.5rem',
          display: 'flex',
          gap: '1.5rem',
          alignItems: 'center',
        }}>
        <Link to="/" style={{ fontWeight: 700, textDecoration: 'none', color: '#7c3aed', fontSize: '1.125rem' }}>
          UniSwap
        </Link>
        <Link to="/">Listings</Link>
        {authenticated ? (
          <>
            <Link to="/create">Create Listing</Link>
            <button type="button" onClick={handleLogout} style={{ marginLeft: 'auto' }}>
              Logout
            </button>
          </>
        ) : (
          <>
            <Link to="/login">Login</Link>
            <Link to="/register">Register</Link>
          </>
        )}
        </nav>
      </header>
      <main style={{ flex: 1, padding: '1.5rem' }}>{children}</main>
    </div>
  );
}
