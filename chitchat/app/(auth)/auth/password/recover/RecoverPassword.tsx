"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Spinner from "../../../../../utils/spinner/Spinner";

export default function RecoverPassword() {
  const [status, setStatus] = useState<'loading' | 'hidden' | 'show'>('loading');
  const [recoveryMethod, setRecoveryMethod] = useState<string>('');
  const [recoveryContact, setRecoveryContact] = useState<string>('');
  const router = useRouter();

  useEffect(() => {
    checkRecoveryStatus();
  }, []);

  const checkRecoveryStatus = async () => {
    try {
      const authInfo = localStorage.getItem('auth-info');
      if (!authInfo) {
        setStatus('hidden');
        return;
      }
      
      const parsed = JSON.parse(authInfo);
      const userUid = parsed.isUserUid;
      
      if (!userUid) {
        setStatus('hidden');
        return;
      }

      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/password/recover`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ userUid }),
      });

      const data = await res.json();
      
      if (data.success && data.recovery_available) {
        setStatus('show');
        setRecoveryMethod(data.recovery_method || '');
        setRecoveryContact(data.recovery_contact || '');
      } else {
        setStatus('hidden');
      }
    } catch (error) {
      setStatus('hidden');
    }
  };

  const handleRecoverClick = () => {
    router.push('/auth/password/recover');
  };

  if (status === 'loading') {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', padding: '10px' }}>
        <Spinner size="small" />
      </div>
    );
  }

  if (status === 'hidden') {
    return null;
  }

  return (
    <div style={{ marginTop: '20px', textAlign: 'center' }}>
      <button
        type="button"
        onClick={handleRecoverClick}
        style={{
          background: 'none',
          border: 'none',
          color: '#0070f3',
          textDecoration: 'underline',
          cursor: 'pointer',
          fontSize: '0.9em'
        }}
      >
        Забыли пароль?
      </button>
      {recoveryContact && (
        <p style={{ fontSize: '0.8em', color: '#666', marginTop: '5px' }}>
          Восстановление через {recoveryMethod}: {recoveryContact}
        </p>
      )}
    </div>
  );
}