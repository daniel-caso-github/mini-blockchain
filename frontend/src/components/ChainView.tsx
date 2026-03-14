import type { Block } from '../api/types';
import { BlockCard } from './BlockCard';
import { ChainGraph } from './ChainGraph';

interface Props {
  blocks: Block[];
  onSelectBlock: (block: Block) => void;
  viewMode?: 'grid' | 'graph';
}

export function ChainView({ blocks, onSelectBlock, viewMode = 'graph' }: Props) {
  if (blocks.length === 0) {
    return (
      <div className="flex items-center justify-center h-48 text-gray-500">
        No blocks in the chain
      </div>
    );
  }

  if (viewMode === 'graph') {
    return <ChainGraph blocks={blocks} onSelectBlock={onSelectBlock} />;
  }

  return (
    <div className="w-full max-w-5xl mx-auto grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 px-4 py-4">
      {blocks.map((block) => (
        <BlockCard key={block.index} block={block} onClick={onSelectBlock} />
      ))}
    </div>
  );
}
