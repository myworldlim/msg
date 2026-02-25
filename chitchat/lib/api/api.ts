import { createEffect } from 'effector'
import { refreshFx, validateSessionFx } from '../store/authSession'

// Simple api effect that automatically attempts refresh on 401
export const apiFx = createEffect(async ({ url, init, skipRefresh }: { url: string; init?: RequestInit; skipRefresh?: boolean }) => {
  const doFetch = async () => {
    const res = await fetch(url, { credentials: 'include', ...init })
    console.log(`[apiFx] fetch(${url}) -> status ${res.status}`);
    return res
  }

  let res = await doFetch()
  
  // Не пытаемся refresh для auth эндпоинтов или если skipRefresh=true
  if (res.status !== 401 || skipRefresh || url.includes('/auth/')) return res

  console.log(`[apiFx] Got 401, attempting refresh...`);
  // try refresh
  try {
    await refreshFx()
    console.log(`[apiFx] Refresh successful, retrying fetch...`);
  } catch (e) {
    console.error(`[apiFx] Refresh failed:`, e);
    // refresh failed, validate session to clear state if needed
    try {
      await validateSessionFx()
    } catch (_) {}
    throw new Error('Unauthorized')
  }

  // retry once
  res = await doFetch()
  if (res.status === 401) throw new Error('Unauthorized')
  return res
})

export const apiFetch = async (url: string, init?: RequestInit) => {
  try {
    const res = await apiFx({ url, init })
    console.log(`apiFetch(${url}) status:`, res.status);
    return res
  } catch (e) {
    console.error(`apiFetch(${url}) error:`, e);
    throw e;
  }
}
