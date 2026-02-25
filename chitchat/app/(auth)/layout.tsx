import type { ReactNode } from 'react';
import Image from 'next/image';
import styles from './authLayout.module.css';
import SecurityPolicy from './auth/components/SecurityPolicy';

export default function AuthLayout({ children }: { children: ReactNode }) {
  return (
    <div className={styles.auth}>
      <div className={styles.authform}>
        <Image src="/logo.png" alt="Логотип" width={60} height={60} className={styles.logo} />
        <h2 className={styles.welcome}>ChitChat</h2>
        <div className={styles.authblock}>{children}</div>
        <SecurityPolicy />
      </div>
    </div>
  );
}
