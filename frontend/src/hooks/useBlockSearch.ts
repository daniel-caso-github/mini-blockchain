import { useState, useMemo } from 'react';
import type { Block } from '../api/types';

export function useBlockSearch(blocks: Block[]) {
  const [query, setQuery] = useState('');

  const filtered = useMemo(() => {
    if (!query.trim()) return blocks;
    const q = query.toLowerCase();
    return blocks.filter(
      (b) =>
        b.index.toString() === q ||
        b.hash.toLowerCase().includes(q) ||
        b.data.toLowerCase().includes(q)
    );
  }, [blocks, query]);

  return { query, setQuery, filtered };
}
