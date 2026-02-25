import { createEffect } from 'effector';

// New: check session effect (uses /auth/session/check)
export const checkSessionFx = createEffect(async () => {
  // Используем обычный fetch чтобы избежать зацикливания с apiFetch -> refreshFx
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/session/check`, {
    method: 'GET',
    credentials: 'include',
  });
  const json = await res.json().catch(() => ({ hasSession: false }));
  if (!res.ok) {
    return { hasSession: false };
  }
  return json;
});

// New: refresh effect (uses /auth/session/refresh)
export const refreshFx = createEffect(async () => {
  // Keep raw fetch here to avoid recursion with apiFetch -> refreshFx
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/session/refresh`, {
    method: 'POST',
    credentials: 'include',
  });
  if (!res.ok) throw new Error('refresh failed');
  return await res.json().catch(() => ({}));
});

// NEW: Validate session and clear if invalid (centralized)
export const validateSessionFx = createEffect(async () => {
  try {
    const result = await checkSessionFx();
    
    // Если сервер говорит нет сессии
    if (!result.hasSession) {
      // Очищаем клиент
      const { resetAuth } = await import('./authStore');
      resetAuth();
      if (typeof window !== 'undefined') {
        localStorage.removeItem('auth-info');
      }
      return { valid: false };
    }
    
    // Сессия валидна
    return { valid: true, ...result };
  } catch (e) {
    // При ошибке сети или сервера тоже очищаем
    const { resetAuth } = await import('./authStore');
    resetAuth();
    if (typeof window !== 'undefined') {
      localStorage.removeItem('auth-info');
    }
    return { valid: false };
  }
});