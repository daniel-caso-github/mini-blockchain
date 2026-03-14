import { useState } from 'react';
import type { Block } from './api/types';
import { useChain } from './hooks/useChain';
import { useHealth } from './hooks/useHealth';
import { useWebSocket } from './hooks/useWebSocket';
import { useBlockSearch } from './hooks/useBlockSearch';
import { HealthIndicator } from './components/HealthIndicator';
import { ChainView } from './components/ChainView';
import { BlockDetail } from './components/BlockDetail';
import { MineForm } from './components/MineForm';
import { ValidateButton } from './components/ValidateButton';
import { StatsPanel } from './components/StatsPanel';
import { SearchBar } from './components/SearchBar';
import { GridPagination } from './components/Pagination';

function App() {
  const { blocks, difficulty, setDifficulty, loading, error, refresh, addBlock, page, totalPages, total, goToPage, pageSize, setPageSize } = useChain();
  const { healthy } = useHealth();
  const { connected } = useWebSocket({ onBlock: addBlock, onDifficulty: setDifficulty });
  const { query, setQuery, filtered } = useBlockSearch(blocks);
  const [selectedBlock, setSelectedBlock] = useState<Block | null>(null);
  const [viewMode, setViewMode] = useState<'grid' | 'graph'>('graph');

  return (
    <div className="min-h-screen bg-gray-950 text-gray-100 flex flex-col">
      {/* Header */}
      <header className="flex items-center justify-between px-6 py-4 border-b border-gray-800">
        <div>
          <h1 className="text-xl font-bold text-white">Mini Blockchain Explorer</h1>
          <p className="text-xs text-gray-500">Difficulty: {difficulty}</p>
        </div>
        <div className="flex items-center gap-3">
          <button
            onClick={() => setViewMode(viewMode === 'graph' ? 'grid' : 'graph')}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-gray-700 bg-gray-800/50 text-xs text-gray-300 hover:border-gray-500 hover:text-white transition-colors cursor-pointer"
            title={viewMode === 'graph' ? 'Switch to grid view' : 'Switch to graph view'}
          >
            {viewMode === 'graph' ? (
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <rect x="3" y="3" width="7" height="7" /><rect x="14" y="3" width="7" height="7" /><rect x="3" y="14" width="7" height="7" /><rect x="14" y="14" width="7" height="7" />
              </svg>
            ) : (
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <circle cx="5" cy="12" r="3" /><circle cx="19" cy="12" r="3" /><line x1="8" y1="12" x2="16" y2="12" />
              </svg>
            )}
            {viewMode === 'graph' ? 'Grid' : 'Graph'}
          </button>
          <HealthIndicator healthy={healthy} />
        </div>
      </header>

      {/* Stats */}
      <StatsPanel blocks={blocks} difficulty={difficulty} connected={connected} total={total} />

      {/* Search */}
      <SearchBar query={query} setQuery={setQuery} total={blocks.length} filtered={filtered.length} />

      {/* Chain */}
      <main className="flex-1 flex items-center py-6 w-full">
        {loading ? (
          <div className="w-full text-center text-gray-500">Loading chain...</div>
        ) : error ? (
          <div className="w-full text-center text-red-400">{error}</div>
        ) : (
          <ChainView blocks={filtered} onSelectBlock={setSelectedBlock} viewMode={viewMode} />
        )}
      </main>

      {/* Pagination bottom */}
      <GridPagination page={page} totalPages={totalPages} total={total} pageSize={pageSize} onPageChange={goToPage} onPageSizeChange={setPageSize} />

      {/* Controls */}
      <footer className="flex items-center justify-between px-6 py-4 border-t border-gray-800">
        <MineForm onMined={refresh} />
        <ValidateButton />
      </footer>

      {/* Block detail modal */}
      {selectedBlock && (
        <BlockDetail block={selectedBlock} onClose={() => setSelectedBlock(null)} />
      )}
    </div>
  );
}

export default App;
