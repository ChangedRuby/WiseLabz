import Axios, { type AxiosError, type AxiosRequestConfig } from 'axios';

/**
 * Shared axios instance used by every orval-generated call (configured as the
 * `mutator` in orval.config.ts).
 *
 * Implements the auth model:
 *   - access token kept IN MEMORY only (never localStorage),
 *   - refresh token rides an HttpOnly cookie (sent automatically via withCredentials),
 *   - a single silent refresh on 401, then one retry; only if THAT fails do we
 *     surface re-auth (handled by authStore / system.notice).
 *
 * `authStore` is not wired yet — the token accessors below are the seam it will
 * plug into. Until then they no-op and requests go out unauthenticated (fine
 * against MSW mocks).
 */

export const AXIOS_INSTANCE = Axios.create({
  baseURL: '/api',
  withCredentials: true, // send/receive the HttpOnly refresh cookie
});

// ─── access-token seam (authStore replaces these) ───────────────────────────
let accessToken: string | null = null;
export const setAccessToken = (token: string | null): void => {
  accessToken = token;
};
export const getAccessToken = (): string | null => accessToken;

/**
 * Called once on a 401 to mint a fresh access token from the refresh cookie.
 * authStore wires the real implementation; default returns false (no session).
 */
type RefreshFn = () => Promise<boolean>;
let refreshHandler: RefreshFn = async () => false;
export const setRefreshHandler = (fn: RefreshFn): void => {
  refreshHandler = fn;
};

// Attach the in-memory access token to every request.
AXIOS_INSTANCE.interceptors.request.use((config) => {
  if (config.url === '/healthz' || config.url === '/readyz') {
    config.baseURL = '';
  }
  const token = getAccessToken();
  if (token) {
    config.headers = config.headers ?? {};
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Single silent-refresh-and-retry on 401.
AXIOS_INSTANCE.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const original = error.config as (AxiosRequestConfig & { _retried?: boolean }) | undefined;
    if (
      error.response?.status === 401 &&
      original &&
      !original._retried &&
      !original.url?.includes('/auth/refresh')
    ) {
      original._retried = true;
      const refreshed = await refreshHandler();
      if (refreshed) {
        return AXIOS_INSTANCE(original);
      }
    }
    return Promise.reject(error);
  }
);

/**
 * orval mutator entrypoint. Generated hooks call this with a request config and
 * expect the unwrapped response body. Supports cancellation via React Query.
 */
export const customInstance = <T>(
  config: AxiosRequestConfig,
  options?: AxiosRequestConfig
): Promise<T> => {
  const source = Axios.CancelToken.source();
  const promise = AXIOS_INSTANCE({
    ...config,
    ...options,
    cancelToken: source.token,
  }).then(({ data }) => data as T);

  // React Query calls .cancel() to abort in-flight requests.
  (promise as Promise<T> & { cancel?: () => void }).cancel = () => {
    source.cancel('Query was cancelled');
  };

  return promise;
};

// orval references these type aliases in generated code.
export type ErrorType<Error> = AxiosError<Error>;
export type BodyType<BodyData> = BodyData;
