import { useState, useEffect, useCallback, useRef } from 'react';
import type { Block } from '../api/types';
import { fetchChain } from '../api/client';

const DEFAULT_PAGE_SIZE = 20;

export function useChain() {
  const [blocks, setBlocks] = useState<Block[]>([]);
  const [difficulty, setDifficulty] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const initialLoad = useRef(true);

  const loadPage = useCallback(async (p: number, size: number) => {
    try {
      setLoading(true);
      setError(null);
      const data = await fetchChain(p, size);
      setBlocks(data.chain);
      setDifficulty(data.difficulty);
      setPage(data.page);
      setTotalPages(data.total_pages);
      setTotal(data.total);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  }, []);

  const refresh = useCallback(() => loadPage(page, pageSize), [loadPage, page, pageSize]);

  // Initial load: go to last page
  useEffect(() => {
    (async () => {
      const data = await fetchChain(1, pageSize);
      if (data.total_pages > 1) {
        loadPage(data.total_pages, pageSize);
      } else {
        setBlocks(data.chain);
        setDifficulty(data.difficulty);
        setPage(data.page);
        setTotalPages(data.total_pages);
        setTotal(data.total);
        setLoading(false);
      }
      initialLoad.current = false;
    })();
  }, [loadPage, pageSize]);

  // When pageSize changes (after initial load), reload page 1
  useEffect(() => {
    if (initialLoad.current) return;
    loadPage(1, pageSize);
  }, [pageSize, loadPage]);

  const goToPage = useCallback((p: number) => {
    if (p < 1) return;
    loadPage(p, pageSize);
  }, [loadPage, pageSize]);

  const addBlock = useCallback((block: Block) => {
    setTotal((prev) => prev + 1);
    setTotalPages(() => {
      const newTotal = total + 1;
      return Math.ceil(newTotal / pageSize);
    });
    setBlocks((prev) => {
      const isLastPage = prev.length === 0 || prev[prev.length - 1].index === block.index - 1;
      if (!isLastPage) return prev;
      if (prev.some((b) => b.index === block.index)) return prev;
      const updated = [...prev, block];
      if (updated.length > pageSize) {
        return updated.slice(1);
      }
      return updated;
    });
  }, [total, pageSize]);

  return { blocks, difficulty, setDifficulty, loading, error, refresh, addBlock, page, totalPages, total, goToPage, pageSize, setPageSize };
}
