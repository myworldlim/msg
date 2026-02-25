// Минимальное обнаружение устройств без получения IP-адреса.
export const detectDevice = () => {
  // Простое определение устройства без вызовов внешних API.
  if (typeof window !== 'undefined') {
    console.debug('Устройство обнаружено');
  }
};