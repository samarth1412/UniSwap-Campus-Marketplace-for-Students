/**
 * FE-3: Login Page UI
 * Email + password form, validation, error display.
 */
import { useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { authApi } from '../services/api';
import { setToken } from '../hooks/useAuth';

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

export function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const from = (location.state as { from?: { pathname: string } })?.from?.pathname ?? '/';

  const valid = email.trim() !== '' && EMAIL_RE.test(email) && password.length >= 6;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    if (!valid) return;
    setLoading(true);
    try {
      const { data } = await authApi.login(email.trim(), password);
      const token = data.data?.token;
      if (!token) throw new Error('No token received');
      setToken(token);
      navigate(from, { replace: true });
    } catch (err: unknown) {
      const ax = err as { response?: { data?: { error?: string } } };
      const msg = ax.response?.data?.error ?? 'Login failed. Try again.';
      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: 400, margin: '0 auto' }}>
      <h1>Login</h1>
      <form onSubmit={handleSubmit}>
        <div style={{ marginBottom: '1rem' }}>
          <label htmlFor="login-email">Email</label>
          <input
            id="login-email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="you@university.edu"
            required
            style={{ display: 'block', width: '100%', padding: '0.5rem', marginTop: '0.25rem' }}
          />
        </div>
        <div style={{ marginBottom: '1rem' }}>
          <label htmlFor="login-password">Password</label>
          <input
            id="login-password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            minLength={6}
            required
            style={{ display: 'block', width: '100%', padding: '0.5rem', marginTop: '0.25rem' }}
          />
        </div>
        {error && <p style={{ color: 'crimson', marginBottom: '1rem' }}>{error}</p>}
        <button type="submit" disabled={!valid || loading}>
          {loading ? 'Signing in…' : 'Sign in'}
        </button>
      </form>
      <p style={{ marginTop: '1rem' }}>
        Don&apos;t have an account? <Link to="/register">Register</Link>
      </p>
    </div>
  );
}
