import { api } from './http';
import type { ApiResponse, AuthResponse, RegisterPayload } from './types';

export const authApi = {
  login: (email: string, password: string) =>
    api.post<ApiResponse<AuthResponse>>('/auth/login', { email, password }),
  register: (payload: RegisterPayload) =>
    api.post<ApiResponse<AuthResponse>>('/auth/register', payload),
  me: () => api.get<ApiResponse<AuthResponse['user']>>('/auth/me'),
};
