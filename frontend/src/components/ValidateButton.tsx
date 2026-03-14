import { useState } from 'react';
import { validateChain } from '../api/client';

export function ValidateButton() {
  const [result, setResult] = useState<boolean | null>(null);
  const [loading, setLoading] = useState(false);

  const handleValidate = async () => {
    try {
      setLoading(true);
      const res = await validateChain();
      setResult(res.valid);
      setTimeout(() => setResult(null), 5000);
    } catch {
      setResult(false);
      setTimeout(() => setResult(null), 5000);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex items-center gap-2">
      <button
        onClick={handleValidate}
        disabled={loading}
        className="px-4 py-2 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 text-white text-sm font-medium rounded-lg transition-colors cursor-pointer disabled:cursor-not-allowed"
      >
        {loading ? 'Validating...' : 'Validate Chain'}
      </button>
      {result !== null && (
        <span className={`text-lg ${result ? 'text-green-400' : 'text-red-400'}`}>
          {result ? '✓ Valid' : '✗ Invalid'}
        </span>
      )}
    </div>
  );
}
