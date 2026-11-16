import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';

const BASE_URL = 'http://localhost:8080/api/v1';

export const api = axios.create({
  baseURL: BASE_URL,
  withCredentials: true,
  headers: { 'Content-Type': 'application/json' },
});

const rawApi = axios.create({
  baseURL: BASE_URL,
  withCredentials: true,
  headers: { 'Content-Type': 'application/json' },
});

let refreshPromise: Promise<string> | null = null;

async function refreshAccessToken(): Promise<string> {
  const refresh = localStorage.getItem('refresh_token');
  if (!refresh) throw new Error('no refresh token');

  const { data } = await rawApi.post<{ refresh_token: string }>(
    '/auth/token-refresh',
    { refresh_token: refresh }
  );
  localStorage.setItem('refresh_token', data.refresh_token);
  return data.refresh_token;
}

function redirectToLogin() {
  if (typeof window === 'undefined') return;
  if (window.location.pathname === '/login') return;
  localStorage.removeItem('refresh_token');
  const next = encodeURIComponent(window.location.pathname + window.location.search);
  window.location.href = `/login?next=${next}`;
}

type RetryConfig = InternalAxiosRequestConfig & { _retry?: boolean };

api.interceptors.response.use(
  (res) => res,
  async (error: AxiosError) => {
    const original = error.config as RetryConfig | undefined;

    const url = original?.url ?? '';
    const isAuthEndpoint =
      url.includes('/auth/sign-in') ||
      url.includes('/auth/sign-up') ||
      url.includes('/auth/token-refresh');

    if (error.response?.status === 401 && original && !original._retry && !isAuthEndpoint) {
      original._retry = true;
      try {
        if (!refreshPromise) {
          refreshPromise = refreshAccessToken().finally(() => {
            refreshPromise = null;
          });
        }
        await refreshPromise;
        return api(original);
      } catch {
        redirectToLogin();
        return Promise.reject(error);
      }
    }

    return Promise.reject(error);
  }
);