import type { Block } from '../api/types';

interface StatsPanelProps {
  blocks: Block[];
  difficulty: number;
  connected: boolean;
  total?: number;
}

export function StatsPanel({ blocks, difficulty, connected, total }: StatsPanelProps) {
  const lastBlock = blocks[blocks.length - 1];
  const avgNonce =
    blocks.length > 0
      ? Math.round(blocks.reduce((sum, b) => sum + b.nonce, 0) / blocks.length)
      : 0;

  const stats = [
    { label: 'Total Bloques', value: total ?? blocks.length },
    { label: 'Dificultad', value: difficulty },
    { label: 'Último Hash', value: lastBlock ? lastBlock.hash.slice(0, 12) + '...' : '-' },
    { label: 'Nonce Promedio', value: avgNonce.toLocaleString() },
  ];

  return (
    <div className="flex items-center gap-4 px-6 py-3 border-b border-gray-800 overflow-x-auto">
      {stats.map((s) => (
        <div
          key={s.label}
          className="flex flex-col bg-gray-900 rounded-lg px-4 py-2 min-w-[140px]"
        >
          <span className="text-xs text-gray-500">{s.label}</span>
          <span className="text-sm font-bold text-white font-mono">{s.value}</span>
        </div>
      ))}
      <div className="flex flex-col bg-gray-900 rounded-lg px-4 py-2 min-w-[140px]">
        <span className="text-xs text-gray-500">WebSocket</span>
        <span className={`text-sm font-bold ${connected ? 'text-green-400' : 'text-red-400'}`}>
          {connected ? 'Conectado' : 'Desconectado'}
        </span>
      </div>
    </div>
  );
}
