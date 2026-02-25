"use client";

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { getAuthStatus } from '@/lib/auth/getauth';

import Spinner3 from '@/utils/spinner/Spinner3';

export default function Page() {
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [serverResp, setServerResp] = useState<any>(null);

  useEffect(() => {
    let mounted = true;

    const run = async () => {
      try {
        const result = await getAuthStatus();

        if (!mounted) return;

        // Покажем ответ сервера в UI для отладки
        setServerResp(result);

        if (result && result.hasSession) {
          router.replace('/protected');
        } else {
          router.replace('/preview');
        }
      } catch (e) {
        // В случае ошибки — отправляем на preview
        router.replace('/preview');
      } finally {
        if (mounted) setLoading(false);
      }
    };

    run();

    return () => {
      mounted = false;
    };
  }, [router]);

  return (
    <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh' }}>
      <div>
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh' }}>
            <Spinner3 />
        </div>
        {serverResp && (
          <pre style={{marginTop: 12, textAlign: 'left', maxWidth: 800}}>
            {JSON.stringify(serverResp, null, 2)}
          </pre>
        )}
      </div>
    </div>
  );
}