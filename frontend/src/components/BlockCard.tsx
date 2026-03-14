import type { Block } from '../api/types';

interface Props {
  block: Block;
  onClick: (block: Block) => void;
}

function truncateHash(hash: string): string {
  return hash.length > 12 ? `${hash.slice(0, 6)}...${hash.slice(-6)}` : hash;
}

export function BlockCard({ block, onClick }: Props) {
  const isGenesis = block.index === 0;

  return (
    <button
      onClick={() => onClick(block)}
      className={`w-full p-4 rounded-xl border text-left transition-all hover:scale-[1.02] hover:shadow-lg cursor-pointer ${
        isGenesis
          ? 'border-amber-500/50 bg-amber-500/5 hover:border-amber-400'
          : 'border-gray-700 bg-gray-800/50 hover:border-gray-500'
      }`}
    >
      <div className="flex items-center justify-between mb-2">
        <span className="text-xs font-bold text-gray-400">BLOCK</span>
        <span className={`text-lg font-bold ${isGenesis ? 'text-amber-400' : 'text-white'}`}>
          #{block.index}
        </span>
      </div>

      {isGenesis && (
        <span className="inline-block px-2 py-0.5 mb-2 text-[10px] font-semibold uppercase tracking-wider rounded bg-amber-500/20 text-amber-400">
          Genesis
        </span>
      )}

      <div className="space-y-1.5 text-xs">
        <div>
          <span className="text-gray-500">Hash: </span>
          <span className="font-mono text-green-400">{truncateHash(block.hash)}</span>
        </div>
        <div>
          <span className="text-gray-500">Prev: </span>
          <span className={`font-mono ${isGenesis ? 'text-gray-500' : 'text-green-400'}`}>
            {isGenesis ? '—' : truncateHash(block.prev_hash)}
          </span>
        </div>
        <div>
          <span className="text-gray-500">Data: </span>
          <span className="text-gray-300 truncate block">{block.data || '—'}</span>
        </div>
        <div>
          <span className="text-gray-500">Nonce: </span>
          <span className="text-gray-300">{block.nonce.toLocaleString()}</span>
        </div>
      </div>
    </button>
  );
}
