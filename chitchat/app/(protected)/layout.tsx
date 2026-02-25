"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useUnit } from "effector-react";
import { $sessionValid, $sessionValidating, $isBlocked } from "@/lib/store/authStore";
import { validateSessionFx } from "@/lib/store/authSession";
import { checkBlockStatusFx } from "@/lib/store/authBlocked";
import { SpinnerContainer } from "../../utils/spinner/Spinner";

export default function ProtectedLayout({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const sessionValid = useUnit($sessionValid);
  const validating = useUnit($sessionValidating);
  const isBlocked = useUnit($isBlocked);

  // 1. При монтировании проверяем сессию и блокировку
  useEffect(() => {
    validateSessionFx();
    checkBlockStatusFx();
  }, []);

  // 2. Если пользователь заблокирован — редирект на /auth/blocked
  useEffect(() => {
    if (isBlocked === true) {
      router.replace('/auth/blocked');
    }
  }, [isBlocked, router]);

  // 3. Если сессия невалидна — редирект
  useEffect(() => {
    if (sessionValid === false) {
      router.replace('/');
    }
  }, [sessionValid, router]);

  // 4. Периодическая проверка сессии каждые 15 минут
  useEffect(() => {
    const sessionInterval = setInterval(() => {
      validateSessionFx();
    }, 15 * 60 * 1000); // 15 минут

    return () => clearInterval(sessionInterval);
  }, []);

  // 5. Периодическая проверка блокировки каждые час
  useEffect(() => {
    const blockInterval = setInterval(() => {
      checkBlockStatusFx();
    }, 60 * 60 * 1000); // 60 минут

    return () => clearInterval(blockInterval);
  }, []);

  // 6. Проверка при видимости таба (сессия + блокировка)
  useEffect(() => {
    const handleVisibilityChange = () => {
      if (!document.hidden) {
        validateSessionFx();
        checkBlockStatusFx();
      }
    };

    document.addEventListener('visibilitychange', handleVisibilityChange);
    return () => document.removeEventListener('visibilitychange', handleVisibilityChange);
  }, []);

  // Показываем спиннер пока проверяем
  if (validating && sessionValid === null) {
    return <SpinnerContainer />;
  }

  // Если не авторизован — не показываем
  if (sessionValid === false) {
    return null;
  }
  return <>{children}</>;
}