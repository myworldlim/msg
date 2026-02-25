"use client";

import { useRouter } from 'next/navigation';
import styles from '../preview.module.css';

export default function PreviewVisibleSegment() {
  const router = useRouter();

  return (
    <main className={styles.container}>
      <h1>Добро пожаловать</h1>
      <p>Для продолжения необходимо войти или зарегистрироваться</p>
      <button 
        className={styles.button}
        onClick={() => router.push('/auth/open')}
      >
        Вход / Регистрация
      </button>
    </main>
  );
}