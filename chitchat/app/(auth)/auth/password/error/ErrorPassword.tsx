"use client";

import { useState, useEffect } from "react";
import Spinner from "../../../../../utils/spinner/Spinner";
import styles from "../../auth.module.css";

interface ErrorPasswordProps {
  onPasswordChange: (password: string) => void;
  disabled?: boolean;
}

export default function ErrorPassword({ onPasswordChange, disabled }: ErrorPasswordProps) {
  const [status, setStatus] = useState<'loading' | 'input' | 'blocked'>('loading');
  const [password, setPassword] = useState("");
  const [timeRemaining, setTimeRemaining] = useState(0);
  const [failedAttempts, setFailedAttempts] = useState(0);

  useEffect(() => {
    checkErrorStatus();
  }, []);

  useEffect(() => {
    let interval: NodeJS.Timeout;
    if (status === 'blocked' && timeRemaining > 0) {
      interval = setInterval(() => {
        setTimeRemaining(prev => {
          if (prev <= 1) {
            setStatus('input');
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
    }
    return () => clearInterval(interval);
  }, [status, timeRemaining]);

  const checkErrorStatus = async () => {
    try {
      const authInfo = localStorage.getItem('auth-info');
      if (!authInfo) return;
      
      const parsed = JSON.parse(authInfo);
      const userUid = parsed.isUserUid;
      
      if (!userUid) return;

      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/password/error`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ userUid }),
      });

      const data = await res.json();
      
      if (data.success) {
        setFailedAttempts(data.failed_attempts || 0);
        
        if (data.error_active && data.time_remaining > 0) {
          setStatus('blocked');
          setTimeRemaining(data.time_remaining);
        } else {
          setStatus('input');
        }
      } else {
        setStatus('input');
      }
    } catch (error) {
      setStatus('input');
    }
  };

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setPassword(value);
    onPasswordChange(value);
  };

  const formatTime = (seconds: number) => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;
    
    if (hours > 0) {
      return `${hours}ч ${minutes}м ${secs}с`;
    } else if (minutes > 0) {
      return `${minutes}м ${secs}с`;
    } else {
      return `${secs}с`;
    }
  };

  if (status === 'loading') {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', padding: '20px' }}>
        <Spinner />
      </div>
    );
  }

  if (status === 'blocked') {
    return (
      <div>
        <div style={{ textAlign: 'center', padding: '20px', color: '#ff4136' }}>
          <h3>Аккаунт временно заблокирован</h3>
          <p>Слишком много неудачных попыток входа ({failedAttempts})</p>
          <p>Повторите попытку через: <strong>{formatTime(timeRemaining)}</strong></p>
        </div>
      </div>
    );
  }

  return (
    <div>
      <label htmlFor="password" style={{ display: "block", marginBottom: "8px" }}>
        Пароль
      </label>
      <input
        id="password"
        type="password"
        value={password}
        onChange={handlePasswordChange}
        placeholder="Введите пароль"
        className={styles.input}
        disabled={disabled}
        autoFocus
      />
      {failedAttempts > 0 && (
        <p style={{ fontSize: '0.9em', color: '#ff4136', marginTop: '5px' }}>
          Неудачных попыток: {failedAttempts}/5
        </p>
      )}
    </div>
  );
}