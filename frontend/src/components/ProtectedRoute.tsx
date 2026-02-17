/**
 * FE-6: Protected Route Component
 * Redirects to /login if no token; otherwise renders children.
 */
import { Navigate, useLocation } from 'react-router-dom';
import { isAuthenticated } from '../hooks/useAuth';

interface ProtectedRouteProps {
  children: React.ReactNode;
}

export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const location = useLocation();
  if (!isAuthenticated()) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }
  return <>{children}</>;
}
