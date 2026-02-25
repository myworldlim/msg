"use client";

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { logoutFx, resetAuth } from '@/lib/store/authStore';

export default function LogoutPage() {
  const router = useRouter();

  useEffect(() => {
    const performLogout = async () => {
      try {
        await logoutFx();
        resetAuth();
        localStorage.removeItem('auth-info');
        router.replace('/');
      } catch (e) {
        // Даже если logout failed, очищаем локальное состояние
        resetAuth();
        localStorage.removeItem('auth-info');
        router.replace('/');
      }
    };
    performLogout();
  }, [router]);

  return (
    <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh' }}>
      <div>Выход из системы...</div>
    </div>
  );
}