"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { apiFetch } from '@/lib/api/api';
import styles from "../auth.module.css";

export default function Open() {
  const [identifier, setIdentifier] = useState("");
  const router = useRouter();
  const [error, setError] = useState("");

  const validateIdentifier = (value: string) => {
    const trimmedValue = value.trim();
    if (!trimmedValue) {
      return { isValid: false, error: "Поле не может быть пустым" };
    }

    // If user typed an @ — treat as email candidate only
    const looksLikeEmail = /@/.test(trimmedValue);
    if (looksLikeEmail) {
      const isEmail = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(trimmedValue);
      if (!isEmail) {
        return { isValid: false, error: "Введите корректный email" };
      }
      return { isValid: true, type: "email", value: trimmedValue.toLowerCase() };
    }

    // Otherwise treat as phone: allow digits, spaces, hyphens and parentheses, optional leading +
    const phoneAllowed = /^\+?[0-9()\s-]{8,25}$/.test(trimmedValue);
    if (!phoneAllowed) {
      return { isValid: false, error: "Введите корректный номер телефона" };
    }
    // count digits only
    const digits = trimmedValue.replace(/[^0-9]/g, "");
    if (digits.length < 8 || digits.length > 15) {
      return { isValid: false, error: "Введите корректный номер телефона" };
    }
    const normalized = (trimmedValue.startsWith("+") ? "+" : "") + digits;
    return { isValid: true, type: "number", value: normalized };
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    const validation = validateIdentifier(identifier);

    if (!validation.isValid) {
      setError(validation.error || "Некорректный ввод");
      return;
    }

    try {
      const response = await apiFetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/open`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ identifier: validation.value, type: validation.type }),
      });

      const data = await response.json();

      if (response.ok) {
        // Сохраняем auth-info
        const authInfo = {
          isUserUid: data.isUserUid || "",
          isExists: typeof data.isExists === 'undefined' ? false : !!data.isExists,
          isPassword: typeof data.isPassword === 'undefined' ? false : !!data.isPassword, // есть ли пароль у пользователя
          isBlocked: typeof data.isBlocked === 'undefined' ? false : !!data.isBlocked,
          isIdentifier: validation.value,
          isUserAgent: typeof navigator !== 'undefined' ? navigator.userAgent : "not_defined",
          isAuthenticated: 'not_defined', // пока не авторизован
        };
        localStorage.setItem("auth-info", JSON.stringify(authInfo));

        // Если пользователь заблокирован — перенаправляем на страницу блокировки
        if (data.isBlocked) {
          router.push('/auth/blocked');
          return;
        }

        // Перенаправление в зависимости от статуса пароля
        // isExists=true → пароль установлен → вход (login)
        // isExists=false → пароль не установлен → регистрация (register)
        router.push(data.isExists ? "/auth/password/login" : "/auth/password/register");
      } else {
        setError(data.error || "Неизвестная ошибка сервера");
      }
    } catch (err) {
      setError("Ошибка сети: " + (err as Error).message);
    }
  };

  return (
    <div className={styles.sendaccess}>
      <form onSubmit={handleSubmit} className={styles.formaccess}>
        <div>
          <input
            type="text"
            value={identifier}
            onChange={(e) => setIdentifier(e.target.value)}
            placeholder="Email или телефон"
            className={styles.input}
          />
          {error && <div className={styles.error}>{error}</div>}
        </div>
        <button type="submit" className={styles.button}>
          Продолжить
        </button>
      </form>
    </div>
  );
}
