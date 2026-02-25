"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { apiFetch } from '@/lib/api/api';
import styles from "../../auth.module.css";

export default function SecretLogin() {
  const [secretWord, setSecretWord] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    // Валидация
    if (!secretWord.trim()) {
      setError("Кодовое слово не может быть пустым");
      setLoading(false);
      return;
    }

    // Получаем auth-info из localStorage
    const authInfoStr = localStorage.getItem("auth-info");
    if (!authInfoStr) {
      setError("Ошибка: информация о пользователе не найдена");
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
      setError("Ошибка: user UID не найден");
      setLoading(false);
      return;
    }

    try {
      const response = await apiFetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/secret/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          userUid: userUid,
          secretWord: secretWord.trim(),
        }),
      });

      const data = await response.json();

      if (response.ok) {
        // Обновляем localStorage - пользователь полностью авторизован
        authInfo.isSecret = true;
        localStorage.setItem("auth-info", JSON.stringify(authInfo));
        
        // Редирект в защищенную зону
        router.push("/protected");
      } else {
        setError(data.error || "Неверное кодовое слово");
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
          <h2>Введите кодовое слово</h2>
          <p style={{ fontSize: "0.9em", color: "#666", marginBottom: "20px" }}>
            Введите кодовое слово для завершения входа
          </p>
        </div>

        <div>
          <label htmlFor="secretWord" style={{ display: "block", marginBottom: "8px" }}>
            Кодовое слово
          </label>
          <input
            id="secretWord"
            type="text"
            value={secretWord}
            onChange={(e) => setSecretWord(e.target.value)}
            placeholder="Введите кодовое слово"
            className={styles.input}
            disabled={loading}
            autoFocus
          />
        </div>

        {error && <div className={styles.error}>{error}</div>}

        <button type="submit" className={styles.button} disabled={loading}>
          {loading ? "Проверка..." : "Войти"}
        </button>
      </form>
    </div>
  );
}