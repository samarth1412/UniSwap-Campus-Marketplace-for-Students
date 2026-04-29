import { Link, useNavigate } from 'react-router-dom';
import { isAuthenticated, clearToken } from '../hooks/useAuth';
import { UniSwapLogo } from './UniSwapLogo';
import './Layout.css';

export function Layout({ children }: { children: React.ReactNode }) {
  const navigate = useNavigate();
  const authenticated = isAuthenticated();

  const handleLogout = () => {
    clearToken();
    navigate('/login');
  };

  return (
    <div className="layout-shell">
      <header className="layout-header">
        <nav className="layout-nav" aria-label="Main">
          <Link to={authenticated ? '/' : '/login'} className="layout-brand">
            <UniSwapLogo className="layout-brand__logo" />
            <span>UniSwap</span>
          </Link>
          {authenticated && (
            <>
              <Link to="/" className="layout-nav-link">
                Listings
              </Link>
              <Link to="/wishlist" className="layout-nav-link">
                Wishlist
              </Link>
              <Link to="/create" className="layout-nav-link">
                Create Listing
              </Link>
              <Link to="/my-listings" className="layout-nav-link">
                My Listings
              </Link>
            </>
          )}
          <span className="layout-nav-spacer" aria-hidden />
          <div className="layout-nav-actions">
            {authenticated ? (
              <button type="button" className="layout-btn layout-btn--ghost" onClick={handleLogout}>
                Logout
              </button>
            ) : (
              <>
                <Link to="/login" className="layout-nav-link">
                  Login
                </Link>
                <Link to="/register" className="layout-nav-link">
                  Register
                </Link>
              </>
            )}
          </div>
        </nav>
      </header>
      <main style={{ flex: 1, padding: '1.5rem' }}>{children}</main>
    </div>
  );
}
