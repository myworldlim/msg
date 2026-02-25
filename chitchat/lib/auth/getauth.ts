import { detectDevice } from '../store/deviceStore';
import { checkSessionFx } from '../store/authSession'

export const getAuthStatus = async () => {
  try {
    detectDevice();
    const result = await checkSessionFx()
    console.debug('getAuthStatus -> result:', result)
    return result
  } catch (error) {
    console.debug('Auth check failed, no session:', error)
    return { hasSession: false }
  }
}