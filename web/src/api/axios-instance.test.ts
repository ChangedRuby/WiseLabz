import { afterAll, afterEach, beforeAll, describe, expect, it, vi } from 'vitest';
import { setupServer } from 'msw/node';
import { http, HttpResponse } from 'msw';
import {
  AXIOS_INSTANCE,
  setAccessToken,
  getAccessToken,
  setRefreshHandler,
} from './axios-instance';

const server = setupServer();

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }));
afterEach(() => {
  server.resetHandlers();
  // Reset token and refresh handler to clean state for next test
  setAccessToken(null);
  setRefreshHandler(async () => false);
});
afterAll(() => server.close());

describe('axios-instance', () => {
  describe('request interceptor', () => {
    it('attaches Authorization header when token is set', async () => {
      const capturedHeaders: Record<string, string> = {};
      server.use(
        http.get('/api/test', ({ request }) => {
          const auth = request.headers.get('Authorization');
          if (auth) capturedHeaders['Authorization'] = auth;
          return HttpResponse.json({ success: true });
        })
      );

      setAccessToken('abc123');
      await AXIOS_INSTANCE.get('/test');

      expect(capturedHeaders['Authorization']).toBe('Bearer abc123');
    });

    it('does not attach Authorization header when token is null', async () => {
      const capturedHeaders: Record<string, string | null> = {};
      server.use(
        http.get('/api/test', ({ request }) => {
          capturedHeaders['Authorization'] = request.headers.get('Authorization');
          return HttpResponse.json({ success: true });
        })
      );

      setAccessToken(null);
      await AXIOS_INSTANCE.get('/test');

      expect(capturedHeaders['Authorization']).toBeNull();
    });

    it('does not attach Authorization header when no token has been set', async () => {
      const capturedHeaders: Record<string, string | null> = {};
      server.use(
        http.get('/api/test', ({ request }) => {
          capturedHeaders['Authorization'] = request.headers.get('Authorization');
          return HttpResponse.json({ success: true });
        })
      );

      // Never call setAccessToken, so token remains null
      await AXIOS_INSTANCE.get('/test');

      expect(capturedHeaders['Authorization']).toBeNull();
    });

    it.each(['/healthz', '/readyz'])('uses the server root for %s', async (path) => {
      server.use(http.get(path, () => HttpResponse.json({ status: 'ok' })));

      await AXIOS_INSTANCE.get(path);
    });
  });

  describe('response interceptor - happy path', () => {
    it('passes through successful responses unchanged', async () => {
      const testData = { id: 1, name: 'test' };
      server.use(http.get('/api/test', () => HttpResponse.json(testData)));

      const response = await AXIOS_INSTANCE.get('/test');

      expect(response.data).toEqual(testData);
      expect(response.status).toBe(200);
    });

    it('does not trigger refresh handler on non-401 errors', async () => {
      const refreshHandler = vi.fn(async () => true);
      setRefreshHandler(refreshHandler);
      server.use(
        http.get('/api/test', () => HttpResponse.json({ error: 'Server error' }, { status: 500 }))
      );

      await expect(AXIOS_INSTANCE.get('/test')).rejects.toThrow();
      expect(refreshHandler).not.toHaveBeenCalled();
    });
  });

  describe('response interceptor - 401 refresh success', () => {
    it('calls refresh handler on 401 and retries original request when refresh succeeds', async () => {
      const refreshHandler = vi.fn(async () => true);
      setRefreshHandler(refreshHandler);
      setAccessToken('old_token');

      const testData = { id: 1, name: 'test' };
      let callCount = 0;

      server.use(
        http.get('/api/test', ({ request }) => {
          callCount++;
          const auth = request.headers.get('Authorization');
          // First call with old token gets 401
          if (callCount === 1 && auth === 'Bearer old_token') {
            return HttpResponse.json({ error: 'Unauthorized' }, { status: 401 });
          }
          // Second call (retry) returns success
          return HttpResponse.json(testData);
        })
      );

      const response = await AXIOS_INSTANCE.get('/test');

      expect(refreshHandler).toHaveBeenCalledTimes(1);
      expect(response.data).toEqual(testData);
      expect(callCount).toBe(2);
    });

    it('retries with the new token after refresh succeeds', async () => {
      const refreshHandler = vi.fn(async () => {
        setAccessToken('new_token');
        return true;
      });
      setRefreshHandler(refreshHandler);
      setAccessToken('old_token');

      const tokens: string[] = [];
      server.use(
        http.get('/api/test', ({ request }) => {
          const auth = request.headers.get('Authorization');
          if (auth) {
            tokens.push(auth);
          }
          if (tokens.length === 1) {
            return HttpResponse.json({ error: 'Unauthorized' }, { status: 401 });
          }
          return HttpResponse.json({ success: true });
        })
      );

      await AXIOS_INSTANCE.get('/test');

      expect(tokens).toEqual(['Bearer old_token', 'Bearer new_token']);
    });
  });

  describe('response interceptor - 401 refresh failure', () => {
    it('rejects original request when refresh handler returns false', async () => {
      const refreshHandler = vi.fn(async () => false);
      setRefreshHandler(refreshHandler);
      setAccessToken('expired_token');

      server.use(
        http.get('/api/test', () => HttpResponse.json({ error: 'Unauthorized' }, { status: 401 }))
      );

      await expect(AXIOS_INSTANCE.get('/test')).rejects.toThrow();
      expect(refreshHandler).toHaveBeenCalledTimes(1);
    });

    it('does not retry when refresh handler rejects', async () => {
      const refreshHandler = vi.fn(async () => {
        throw new Error('Refresh failed');
      });
      setRefreshHandler(refreshHandler);
      setAccessToken('expired_token');

      server.use(
        http.get('/api/test', () => HttpResponse.json({ error: 'Unauthorized' }, { status: 401 }))
      );

      await expect(AXIOS_INSTANCE.get('/test')).rejects.toThrow();
      expect(refreshHandler).toHaveBeenCalledTimes(1);
    });
  });

  describe('response interceptor - refresh endpoint guard', () => {
    it('does not call refresh handler when 401 occurs on /auth/refresh endpoint', async () => {
      const refreshHandler = vi.fn(async () => true);
      setRefreshHandler(refreshHandler);
      setAccessToken('token');

      server.use(
        http.post('/api/auth/refresh', () =>
          HttpResponse.json({ error: 'Unauthorized' }, { status: 401 })
        )
      );

      await expect(AXIOS_INSTANCE.post('/auth/refresh')).rejects.toThrow();
      expect(refreshHandler).not.toHaveBeenCalled();
    });

    it('does not retry refresh endpoint even if refresh handler is registered', async () => {
      const refreshHandler = vi.fn(async () => true);
      setRefreshHandler(refreshHandler);

      let callCount = 0;
      server.use(
        http.post('/api/auth/refresh', () => {
          callCount++;
          return HttpResponse.json({ error: 'Unauthorized' }, { status: 401 });
        })
      );

      await expect(AXIOS_INSTANCE.post('/auth/refresh')).rejects.toThrow();
      expect(callCount).toBe(1); // Only called once, not retried
      expect(refreshHandler).not.toHaveBeenCalled();
    });
  });

  describe('response interceptor - single retry guard', () => {
    it('does not call refresh handler again when request already has _retried flag', async () => {
      const refreshHandler = vi.fn(async () => true);
      setRefreshHandler(refreshHandler);

      let callCount = 0;
      server.use(
        http.get('/api/test', () => {
          callCount++;
          return HttpResponse.json({ error: 'Unauthorized' }, { status: 401 });
        })
      );

      // Make a request with _retried already set
      const config = {
        method: 'get' as const,
        url: '/test',
        _retried: true,
      };

      await expect(AXIOS_INSTANCE(config)).rejects.toThrow();
      expect(refreshHandler).not.toHaveBeenCalled();
      expect(callCount).toBe(1); // Only one call, no retry
    });

    it('prevents infinite retry loops on subsequent 401s', async () => {
      const refreshHandler = vi.fn(async () => true);
      setRefreshHandler(refreshHandler);
      setAccessToken('token');

      // Simulate a scenario where even after refresh, the endpoint still returns 401
      let callCount = 0;
      server.use(
        http.get('/api/test', () => {
          callCount++;
          return HttpResponse.json({ error: 'Unauthorized' }, { status: 401 });
        })
      );

      await expect(AXIOS_INSTANCE.get('/test')).rejects.toThrow();
      // Should be called twice: first request + one retry after refresh
      expect(callCount).toBe(2);
      // Refresh handler should only be called once
      expect(refreshHandler).toHaveBeenCalledTimes(1);
    });
  });

  describe('getAccessToken', () => {
    it('returns the current access token', () => {
      setAccessToken('test_token_123');
      expect(getAccessToken()).toBe('test_token_123');
    });

    it('returns null when no token is set', () => {
      setAccessToken(null);
      expect(getAccessToken()).toBeNull();
    });
  });

  describe('integration - full auth flow', () => {
    it('handles a complete auth flow: request -> 401 -> refresh -> retry -> success', async () => {
      const refreshHandler = vi.fn(async () => {
        setAccessToken('refreshed_token');
        return true;
      });
      setRefreshHandler(refreshHandler);
      setAccessToken('initial_token');

      const requestLog: Array<{ token: string; attempt: number }> = [];
      let attempt = 0;

      server.use(
        http.get('/api/protected-resource', ({ request }) => {
          attempt++;
          const token = request.headers.get('Authorization')?.replace('Bearer ', '') || 'none';
          requestLog.push({ token, attempt });

          if (attempt === 1) {
            return HttpResponse.json({ error: 'Token expired' }, { status: 401 });
          }
          return HttpResponse.json({ data: 'sensitive info' });
        })
      );

      const response = await AXIOS_INSTANCE.get('/protected-resource');

      expect(response.data).toEqual({ data: 'sensitive info' });
      expect(refreshHandler).toHaveBeenCalledTimes(1);
      expect(requestLog).toEqual([
        { token: 'initial_token', attempt: 1 },
        { token: 'refreshed_token', attempt: 2 },
      ]);
    });
  });
});
