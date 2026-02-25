"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { apiFetch } from '@/lib/api/api';
import styles from "../../auth.module.css";

export default function CreateSecretPage() {
  const [secretWord, setSecretWord] = useState("");
  const [confirmSecret, setConfirmSecret] = useState("");
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

    if (secretWord.length < 3 || secretWord.length > 100) {
      setError("Кодовое слово должно быть от 3 до 100 символов");
      setLoading(false);
      return;
    }

    if (secretWord !== confirmSecret) {
      setError("Кодовые слова не совпадают");
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
      const response = await apiFetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/secret/create`, {
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
        authInfo.isAuthenticated = 'authenticated';
        authInfo.isSecret = true; // секретное слово создано
        localStorage.setItem("auth-info", JSON.stringify(authInfo));
        
        // Редирект в защищенную зону
        router.push("/protected");
      } else {
        setError(data.error || "Ошибка при создании кодового слова");
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
          <h2>Создать кодовое слово</h2>
          <p style={{ fontSize: "0.9em", color: "#666", marginBottom: "20px" }}>
            Придумайте кодовое слово или фразу для дополнительной защиты аккаунта
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
          />
        </div>

        <div>
          <label htmlFor="confirmSecret" style={{ display: "block", marginBottom: "8px" }}>
            Подтвердите кодовое слово
          </label>
          <input
            id="confirmSecret"
            type="text"
            value={confirmSecret}
            onChange={(e) => setConfirmSecret(e.target.value)}
            placeholder="Повторите кодовое слово"
            className={styles.input}
            disabled={loading}
          />
        </div>

        {error && <div className={styles.error}>{error}</div>}

        <button type="submit" className={styles.button} disabled={loading}>
          {loading ? "Создание..." : "Сохранить кодовое слово"}
        </button>
      </form>
    </div>
  );
}