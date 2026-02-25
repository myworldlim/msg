import { createEffect } from 'effector';

// NEW: Check if user is blocked via backend endpoint
export const checkBlockStatusFx = createEffect(async () => {
  try {
    // Получаем userUid из localStorage
    if (typeof window !== 'undefined') {
      const authInfo = localStorage.getItem('auth-info');
      if (!authInfo) {
        return { blocked: false };
      }
      
      const parsed = JSON.parse(authInfo);
      const userUid = parsed.isUserUid;
      
      if (!userUid) {
        return { blocked: false };
      }
      
      // Используем прямой запрос к Go backend
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/blocked`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ userUid }),
      });
      
      // Проверяем статус
      if (!res.ok) {
        return { blocked: false };
      }
      
      // Проверяем Content-Type перед парсингом
      const contentType = res.headers.get('content-type');
      if (!contentType?.includes('application/json')) {
        return { blocked: false };
      }
      
      // Парсим JSON
      let json;
      try {
        json = await res.json();
      } catch (e) {
        return { blocked: false };
      }
      
      // Проверяем успех запроса
      if (!json.success) {
        return { blocked: false };
      }
      
      // Если пользователь заблокирован
      if (json.blocked === true) {
        // Очищаем cookies через прямой запрос к Go backend
        try {
          await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/logout`, {
            method: 'POST',
            credentials: 'include',
          }).catch(() => {});
        } catch (_) {}
        
        // Очищаем localStorage, кроме userUid и identifier
        const cleaned = {
          isUserUid: parsed.isUserUid,
          isIdentifier: parsed.isIdentifier,
          isBlocked: true,
        };
        localStorage.setItem('auth-info', JSON.stringify(cleaned));
        
        const { resetAuth } = await import('./authStore');
        resetAuth();
        
        return { blocked: true, reason: json.reason };
      }
      
      // Если не заблокирован — обновляем статус в localStorage
      parsed.isBlocked = false;
      localStorage.setItem('auth-info', JSON.stringify(parsed));
      
      return { blocked: false };
    }
    
    return { blocked: false };
  } catch (e) {
    return { blocked: false };
  }
});