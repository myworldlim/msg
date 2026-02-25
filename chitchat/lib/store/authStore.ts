import { createStore, createEvent, createEffect } from 'effector';
import { apiFetch } from '../api/api';
import { checkSessionFx, refreshFx, validateSessionFx } from './authSession';
import { checkBlockStatusFx } from './authBlocked';

// Types
export type AuthStatus = 'not_defined' | 'authenticated' | 'not_authenticated';
export type IdentifierStatus = 'not_defined' | 'valid' | 'invalid';
export type BlockStatus = 'not_defined' | 'blocked' | 'not_blocked';
export type ExistsStatus = 'not_defined' | 'exists' | 'not_exists';
export type ProtectionStatus = 'not_defined' | 'protected' | 'not_protected';
export type SecretStatus = 'not_defined' | 'has_secret' | 'no_secret';
export type UserAgentStatus = 'not_defined' | 'valid' | 'invalid';

export interface AuthState {
  isAuthenticated: AuthStatus;
  isIdentifier: IdentifierStatus;
  isBlocked: BlockStatus;
  isExists: ExistsStatus;
  isProtection: ProtectionStatus;
  isSecret: SecretStatus;
  isUserAgent: UserAgentStatus;
}

// Effects
export const checkAuthStatusEffect = createEffect(async () => {
  // Явно запрашиваем backend (Go) на порту 8181.
  const BACKEND_URL = process.env.NEXT_PUBLIC_BACKEND_URL || 'http://localhost:8181/';
  let response: Response;
  try {
    response = await fetch(BACKEND_URL, {
      method: 'GET',
      credentials: 'include',
    });
  } catch (networkError) {
    console.debug('Network error while contacting backend:', networkError);
    return { hasSession: false, sessionToken: false, sessionRefresh: false };
  }

  // Логируем статус для отладки
  console.debug('checkAuthStatusEffect: response.status=', response.status);

  let data: any = { has_session: false, session_token: false, session_refresh: false };
  try {
    data = await response.json();
  } catch (e) {
    console.warn('Failed to parse JSON from backend auth check', e);
  }

  console.debug('checkAuthStatusEffect: backend body=', data);

  return {
    hasSession: Boolean(data.has_session),
    sessionToken: Boolean(data.session_token),
    sessionRefresh: Boolean(data.session_refresh),
    // также возвращаем raw для удобства отладки
    raw: data
  } as any;
});

// Logout effect
export const logoutFx = createEffect(async () => {
  try {
    const res = await apiFetch(`/auth/logout`, {
      method: 'POST',
    });
    if (!res.ok) throw new Error('logout failed');
    return await res.json().catch(() => ({}));
  } catch (e) {
    // propagate error to caller — UI can still call resetAuth()
    throw e;
  }
});

// Events
export const resetAuth = createEvent();
export const updateAuthStatus = createEvent<Partial<AuthState>>();

// NEW: Store to track session validity
export const $sessionValid = createStore<boolean | null>(null)
  .on(validateSessionFx.done, (_, { result }) => result.valid)
  .on(validateSessionFx.fail, () => false)
  .on(resetAuth, () => null);

// NEW: Store to track blocked status
export const $isBlocked = createStore<boolean>(false)
  .on(checkBlockStatusFx.done, (_, { result }) => result.blocked)
  .on(checkBlockStatusFx.fail, () => false)
  .on(resetAuth, () => false);

// NEW: Store for session validation in progress
export const $sessionValidating = validateSessionFx.pending;

// Store
export const $auth = createStore<AuthState>({
  isAuthenticated: 'not_defined',
  isIdentifier: 'not_defined',
  isBlocked: 'not_defined',
  isExists: 'not_defined',
  isProtection: 'not_defined',
  isSecret: 'not_defined',
  isUserAgent: 'not_defined'
})
  .on(updateAuthStatus, (state, payload) => ({ ...state, ...payload }))
  .on(resetAuth, () => ({
    isAuthenticated: 'not_defined',
    isIdentifier: 'not_defined',
    isBlocked: 'not_defined',
    isExists: 'not_defined',
    isProtection: 'not_defined',
    isSecret: 'not_defined',
    isUserAgent: 'not_defined'
  }))
  .on(checkAuthStatusEffect.done, (state, { result }) => ({
    ...state,
    isAuthenticated: result.hasSession ? 'authenticated' : 'not_authenticated'
  }))
  .on(logoutFx.done, () => ({
    isAuthenticated: 'not_authenticated',
    isIdentifier: 'not_defined',
    isBlocked: 'not_defined',
    isExists: 'not_defined',
    isProtection: 'not_defined',
    isSecret: 'not_defined',
    isUserAgent: 'not_defined'
  }))
  .on(validateSessionFx.done, (state, { result }) => {
    // Если сессия невалидна — очищаем состояние
    if (!result.valid) {
      return {
        isAuthenticated: 'not_authenticated',
        isIdentifier: 'not_defined',
        isBlocked: 'not_defined',
        isExists: 'not_defined',
        isProtection: 'not_defined',
        isSecret: 'not_defined',
        isUserAgent: 'not_defined'
      };
    }
    return state;
  })
  .on(validateSessionFx.fail, () => ({
    isAuthenticated: 'not_authenticated',
    isIdentifier: 'not_defined',
    isBlocked: 'not_defined',
    isExists: 'not_defined',
    isProtection: 'not_defined',
    isSecret: 'not_defined',
    isUserAgent: 'not_defined'
  }));

// Persist / hydrate auth store to localStorage
const AUTH_LOCAL_KEY = 'auth-info';
if (typeof window !== 'undefined') {
  try {
    const raw = localStorage.getItem(AUTH_LOCAL_KEY);
    if (raw) {
      const parsed = JSON.parse(raw);
      // Если parsed содержит поля статусов, применим их
      updateAuthStatus(parsed);
    }
  } catch (e) {
    console.warn('Failed to hydrate auth from localStorage', e);
  }

  $auth.watch((state) => {
    try {
      localStorage.setItem(AUTH_LOCAL_KEY, JSON.stringify(state));
    } catch (e) {
      // ignore
    }
  });
}