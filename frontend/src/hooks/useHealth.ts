import { useState, useEffect } from 'react';
import { checkHealth } from '../api/client';

export function useHealth() {
  const [healthy, setHealthy] = useState(false);

  useEffect(() => {
    const check = async () => {
      try {
        await checkHealth();
        setHealthy(true);
      } catch {
        setHealthy(false);
      }
    };

    check();
    const interval = setInterval(check, 10000);
    return () => clearInterval(interval);
  }, []);

  return { healthy };
}
