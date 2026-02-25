"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { apiFetch } from '@/lib/api/api';
import styles from "../../auth.module.css";
import ErrorPassword from "../error/ErrorPassword";
import RecoverPassword from "../recover/RecoverPassword";

export default function PasswordLogin() {
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    // Валидация
    if (!password.trim()) {
      setError("Пароль не может быть пустым");
      setLoading(false);
      return;
    }

    // Получаем auth-info из localStorage
    const authInfoStr = localStorage.getItem("auth-info");
    if (!authInfoStr) {
      setError("Ошибка: информация о пользователе не найдена. Начните заново.");
      setLoading(false);
      return;
    }

    let authInfo;
    try {
      authInfo = JSON.parse(authInfoStr);
    } catch (err) {
      setError("Ошибка: некорректные данные пользователя");
      setLoading(false);
      return;
    }

    const userUid = authInfo.isUserUid;
    if (!userUid) {
      setError("Ошибка: user UID не найден. Начните заново.");
      setLoading(false);
      return;
    }

    try {
      const response = await apiFetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/password/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          userUid: userUid,
          password: password.trim(),
        }),
      });

      const data = await response.json();

      if (response.ok) {
        // Обновляем localStorage - пользователь авторизован
        authInfo.isAuthenticated = 'authenticated';
        authInfo.isPassword = true; // Пароль введен правильно
        authInfo.isProtection = data.protection; // сохраняем статус protection с сервера
        if (!data.protection) {
          authInfo.isSecret = false; // если без защиты - секретное слово не нужно
        }
        localStorage.setItem("auth-info", JSON.stringify(authInfo));

        // Проверяем нужна ли дополнительная защита
        if (data.protection) {
          // Если включена защита - переходим к вводу секретного слова
          router.push("/auth/secret/login");
        } else {
          // Если защита отключена - сразу в защищенную зону
          router.push("/protected");
        }
      } else {
        setError(data.error || "Неверный пароль");
        setPassword(""); // Очищаем поле пароля при ошибке
        
        // Обновляем localStorage - пароль неверный
        authInfo.isPassword = false;
        localStorage.setItem("auth-info", JSON.stringify(authInfo));
      }
    } catch (err) {
      setError("Ошибка сети: " + (err as Error).message);
      setPassword(""); // Очищаем поле пароля при ошибке
      
      // Обновляем localStorage - ошибка сети
      if (authInfo) {
        authInfo.isPassword = false;
        localStorage.setItem("auth-info", JSON.stringify(authInfo));
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.sendaccess}>
      <form onSubmit={handleSubmit} className={styles.formaccess}>
        <div>
          <h2>Введите пароль</h2>
          <p style={{ fontSize: "0.9em", color: "#666", marginBottom: "20px" }}>
            Введите пароль для входа в аккаунт
          </p>
        </div>

        <ErrorPassword 
          onPasswordChange={setPassword}
          disabled={loading}
        />

        {error && <div className={styles.error}>{error}</div>}

        <button type="submit" className={styles.button} disabled={loading || !password.trim()}>
          {loading ? "Проверка..." : "Продолжить"}
        </button>
        
        <RecoverPassword />
      </form>
    </div>
  );
}