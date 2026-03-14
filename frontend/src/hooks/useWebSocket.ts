import { useEffect, useRef, useState, useCallback } from 'react';
import type { Block } from '../api/types';

interface WSBlockMessage {
  type: 'new_block';
  block: Block;
}

interface WSDifficultyMessage {
  type: 'difficulty_adjusted';
  difficulty: number;
}

type WSMessage = WSBlockMessage | WSDifficultyMessage;

interface UseWebSocketOptions {
  onBlock: (block: Block) => void;
  onDifficulty?: (difficulty: number) => void;
}

export function useWebSocket({ onBlock, onDifficulty }: UseWebSocketOptions) {
  const [connected, setConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);
  const retriesRef = useRef(0);
  const onBlockRef = useRef(onBlock);
  const onDifficultyRef = useRef(onDifficulty);
  onBlockRef.current = onBlock;
  onDifficultyRef.current = onDifficulty;

  const connect = useCallback(() => {
    const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:';
    const ws = new WebSocket(`${protocol}//${location.host}/ws`);
    wsRef.current = ws;

    ws.onopen = () => {
      setConnected(true);
      retriesRef.current = 0;
    };

    ws.onmessage = (event) => {
      try {
        const msg: WSMessage = JSON.parse(event.data);
        if (msg.type === 'new_block') {
          onBlockRef.current(msg.block);
        } else if (msg.type === 'difficulty_adjusted') {
          onDifficultyRef.current?.(msg.difficulty);
        }
      } catch {
        // ignore malformed messages
      }
    };

    ws.onclose = () => {
      setConnected(false);
      wsRef.current = null;
      const delay = Math.min(1000 * Math.pow(2, retriesRef.current), 30000);
      retriesRef.current++;
      setTimeout(connect, delay);
    };

    ws.onerror = () => {
      ws.close();
    };
  }, []);

  useEffect(() => {
    connect();
    return () => {
      wsRef.current?.close();
    };
  }, [connect]);

  return { connected };
}
