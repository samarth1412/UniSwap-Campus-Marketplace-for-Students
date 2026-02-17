/**
 * API service layer - FE-7
 * Axios instance with base URL and JWT attached to requests.
 * Handles 401 by clearing token and redirecting to login.
 */
import axios, { type AxiosInstance } from 'axios';

const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api';

export const api: AxiosInstance = axios.create({
  baseURL: API_BASE,
  headers: { 'Content-Type': 'application/json' },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Backend returns { success, data: { token, user } }
export interface AuthResponse {
  token: string;
  user?: { id: string; email: string; full_name?: string };
}

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

export interface RegisterPayload {
  full_name: string;
  email: string;
  university?: string;
  password: string;
}

export const authApi = {
  login: (email: string, password: string) =>
    api.post<ApiResponse<AuthResponse>>('/auth/login', { email, password }),
  register: (payload: RegisterPayload) =>
    api.post<ApiResponse<AuthResponse>>('/auth/register', payload),
};
