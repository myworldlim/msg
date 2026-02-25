'use client';

import Link from 'next/link';

export default function SecurityPolicy() {
  return (
    <footer style={{marginTop: 24, fontSize: 12, color: '#666'}}>
      <p>
        © ChitChat — Нажимая на кнопку "Продолжить", вы соглашаетесь с нашей
        <Link href="/auth/security-policy" style={{color: '#666', textDecoration: 'underline', marginLeft: 4, marginRight: 4}}>
          Политикой безопасности
        </Link>
        и
        <Link href="/auth/terms-of-use" style={{color: '#666', textDecoration: 'underline', marginLeft: 4}}>
          Условиями использования
        </Link>
        .
      </p>
    </footer>
  );
}
