import type { Block, ChainResponse, ValidateResponse, HealthResponse } from './types';

export async function fetchChain(page = 1, limit = 20): Promise<ChainResponse> {
  const res = await fetch(`/chain?page=${page}&limit=${limit}`);
  if (!res.ok) throw new Error('Failed to fetch chain');
  return res.json();
}

export async function mineBlock(data: string): Promise<Block> {
  const res = await fetch('/mine', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ data }),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || 'Failed to mine block');
  }
  return res.json();
}

export async function validateChain(): Promise<ValidateResponse> {
  const res = await fetch('/validate');
  if (!res.ok) throw new Error('Failed to validate chain');
  return res.json();
}

export async function fetchBlock(id: number): Promise<Block> {
  const res = await fetch(`/block/${id}`);
  if (!res.ok) throw new Error('Block not found');
  return res.json();
}

export async function checkHealth(): Promise<HealthResponse> {
  const res = await fetch('/health');
  if (!res.ok) throw new Error('Health check failed');
  return res.json();
}
