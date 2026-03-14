import { useState } from 'react';
import { mineBlock } from '../api/client';

interface Props {
  onMined: () => void;
}

export function MineForm({ onMined }: Props) {
  const [data, setData] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!data.trim() || loading) return;

    try {
      setLoading(true);
      setError(null);
      await mineBlock(data.trim());
      setData('');
      onMined();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Mining failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex items-center gap-2">
      <input
        type="text"
        value={data}
        onChange={(e) => setData(e.target.value)}
        placeholder="Block data..."
        className="px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-blue-500 w-56"
        disabled={loading}
      />
      <button
        type="submit"
        disabled={!data.trim() || loading}
        className="px-4 py-2 bg-blue-600 hover:bg-blue-500 disabled:bg-gray-700 disabled:text-gray-500 text-white text-sm font-medium rounded-lg transition-colors cursor-pointer disabled:cursor-not-allowed flex items-center gap-2"
      >
        {loading && (
          <svg className="animate-spin h-4 w-4" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
        )}
        {loading ? 'Mining...' : 'Mine Block'}
      </button>
      {error && <span className="text-xs text-red-400">{error}</span>}
    </form>
  );
}
