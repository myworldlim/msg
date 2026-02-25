"use client";

import { useRouter } from 'next/navigation';

export default function ProtectedVisiblePage() {
  const router = useRouter();

  return (
    <main style={{display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh'}}>
      <div style={{textAlign: 'center'}}>
        <h1>Пользователь авторизован</h1>
        <p>Это видимая страница для защищённого сегмента /protected.</p>
        <button 
          onClick={() => router.push('/auth/logout')}
          style={{
            marginTop: '20px',
            padding: '10px 20px',
            backgroundColor: '#dc3545',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer'
          }}
        >
          Выйти
        </button>
      </div>
    </main>
  );
}
