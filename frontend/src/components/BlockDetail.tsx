import type { Block } from '../api/types';

interface Props {
  block: Block;
  onClose: () => void;
}

export function BlockDetail({ block, onClose }: Props) {
  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
      onClick={onClose}
    >
      <div
        className="bg-gray-900 border border-gray-700 rounded-2xl p-6 max-w-lg w-full mx-4 shadow-2xl"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-bold text-white">
            Block #{block.index}
            {block.index === 0 && (
              <span className="ml-2 px-2 py-0.5 text-xs font-semibold uppercase rounded bg-amber-500/20 text-amber-400">
                Genesis
              </span>
            )}
          </h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-white text-2xl leading-none cursor-pointer"
          >
            &times;
          </button>
        </div>

        <div className="space-y-3">
          <Field label="Timestamp" value={block.timestamp} />
          <Field label="Data" value={block.data || '(empty)'} />
          <Field label="Hash" value={block.hash} mono />
          <Field label="Previous Hash" value={block.prev_hash} mono />
          <Field label="Nonce" value={block.nonce.toLocaleString()} />
        </div>
      </div>
    </div>
  );
}

function Field({ label, value, mono }: { label: string; value: string; mono?: boolean }) {
  return (
    <div>
      <div className="text-xs text-gray-500 uppercase tracking-wider mb-0.5">{label}</div>
      <div
        className={`text-sm text-gray-200 break-all ${mono ? 'font-mono bg-gray-800 rounded px-2 py-1' : ''}`}
      >
        {value}
      </div>
    </div>
  );
}
