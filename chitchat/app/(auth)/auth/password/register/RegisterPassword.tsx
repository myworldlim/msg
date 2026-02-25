"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { apiFetch } from '@/lib/api/api';
import styles from "../../auth.module.css";

export default function RegisterPassword() {
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [protection, setProtection] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  const validatePassword = (pwd: string): { isValid: boolean; error?: string } => {
    if (!pwd) {
      return { isValid: false, error: "Пароль не может быть пустым" };
    }
    if (pwd.length < 8) {
      return { isValid: false, error: "Пароль должен быть минимум 8 символов" };
    }
    if (pwd.length > 128) {
      return { isValid: false, error: "Пароль не должен превышать 128 символов" };
    }
    return { isValid: true };
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    // Валидация пароля
    const validation = validatePassword(password);
    if (!validation.isValid) {
      setError(validation.error || "Некорректный пароль");
      setLoading(false);
      return;
    }

    // Проверка совпадения паролей
    if (password !== confirmPassword) {
      setError("Пароли не совпадают");
      setLoading(false);
      return;
    }

    // Получаем auth-info из localStorage
    const authInfoStr = localStorage.getItem("auth-info");
    if (!authInfoStr) {
      setError("Ошибка: информация о пользователе не найдена. Пожалуйста, начните заново.");
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
      setError("Ошибка: user UID не найден. Пожалуйста, начните заново.");
      setLoading(false);
      return;
    }

    try {
      // Отправляем пароль на регистрацию
      const response = await apiFetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/password/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          userUid: userUid,
          password: password,
          protection: protection,
        }),
      });

      const data = await response.json();

      if (response.ok) {
        // Обновляем localStorage — пароль успешно установлен
        authInfo.isPassword = true;
        authInfo.isProtection = protection; // сохраняем статус protection
        authInfo.isSecret = false; // пока секретное слово не создано
        if (!protection) {
          authInfo.isAuthenticated = 'authenticated'; // если без защиты - сразу авторизован
        }
        localStorage.setItem("auth-info", JSON.stringify(authInfo));

        // Редирект в зависимости от protection
        if (data.protection) {
          router.push("/auth/secret/register"); // Создание secret word
        } else {
          router.push("/protected"); // Прямо в приложение
        }
      } else {
        setError(data.error || "Ошибка при регистрации пароля");
      }
    } catch (err) {
      setError("Ошибка сети: " + (err as Error).message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.sendaccess}>
      <form onSubmit={handleSubmit} className={styles.formaccess}>
        <div>
          <h2>Создание пароля</h2>
          <p style={{ fontSize: "0.9em", color: "#666" }}>
            Минимум 8 символов, максимум 128 символов
          </p>
        </div>

        <div>
          <label htmlFor="password" style={{ display: "block", marginBottom: "8px" }}>
            Пароль
          </label>
          <input
            id="password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Введите пароль"
            className={styles.input}
            disabled={loading}
          />
        </div>

        <div>
          <label htmlFor="confirmPassword" style={{ display: "block", marginBottom: "8px" }}>
            Подтвердите пароль
          </label>
          <input
            id="confirmPassword"
            type="password"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            placeholder="Подтвердите пароль"
            className={styles.input}
            disabled={loading}
          />
        </div>

        <div style={{ marginTop: "16px", display: "flex", alignItems: "center", gap: "8px" }}>
          <input
            id="protection"
            type="checkbox"
            checked={protection}
            onChange={(e) => setProtection(e.target.checked)}
            disabled={loading}
            style={{ width: "20px", height: "20px", cursor: "pointer" }}
          />
          <label htmlFor="protection" style={{ cursor: "pointer", margin: 0 }}>
            Дополнительная защита
          </label>
        </div>

        {error && <div className={styles.error}>{error}</div>}

        <button type="submit" className={styles.button} disabled={loading}>
          {loading ? "Регистрация..." : "Сохранить пароль"}
        </button>
      </form>
    </div>
  );
}
