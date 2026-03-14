import { useMemo, useCallback } from 'react';
import type { Block } from '../api/types';

interface Props {
  blocks: Block[];
  onSelectBlock: (block: Block) => void;
}

function truncateHash(hash: string): string {
  return hash.length > 12 ? `${hash.slice(0, 6)}...${hash.slice(-6)}` : hash;
}

/** Derive a hue from the first 6 hex chars of a hash */
function hashToHue(hash: string): number {
  const hex = hash.replace(/^0x/, '').slice(0, 6);
  const num = parseInt(hex, 16);
  return isNaN(num) ? 200 : num % 360;
}

/** Position nodes in a spiral from center outward */
function getNodePosition(
  index: number,
  total: number,
  width: number,
  height: number,
): { x: number; y: number } {
  if (index === 0) return { x: width / 2, y: height / 2 };

  const angle = index * 0.8;
  const radius = 40 + index * 28;
  const x = width / 2 + radius * Math.cos(angle);
  const y = height / 2 + radius * Math.sin(angle);
  return { x, y };
}

const NODE_RADIUS = 28;

/** Mobile vertical fallback (simple list) */
function MobileVerticalNode({
  block,
  isLast,
  onClick,
}: {
  block: Block;
  isLast: boolean;
  onClick: (block: Block) => void;
}) {
  const isGenesis = block.index === 0;
  const hue = hashToHue(block.hash);

  return (
    <div className="flex flex-col items-center w-full">
      <button
        onClick={() => onClick(block)}
        className={`w-16 h-16 rounded-full border-2 flex items-center justify-center font-bold text-sm transition-all cursor-pointer ${
          isGenesis
            ? 'border-amber-400 text-amber-300 shadow-[0_0_14px_rgba(245,158,11,0.35)]'
            : 'border-gray-600 text-white hover:border-gray-400'
        }`}
        style={
          isGenesis
            ? undefined
            : {
                borderColor: `hsl(${hue}, 55%, 50%)`,
                boxShadow: `0 0 8px hsla(${hue}, 55%, 50%, 0.2)`,
              }
        }
      >
        #{block.index}
      </button>
      <span className="text-[9px] font-mono text-gray-500 mt-1">{truncateHash(block.hash)}</span>
      {!isLast && (
        <svg width="2" height="32" className="my-1">
          <line
            x1="1"
            y1="0"
            x2="1"
            y2="32"
            stroke="#4b5563"
            strokeWidth="2"
            strokeDasharray="4 3"
            className="dash-flow-v"
          />
        </svg>
      )}
    </div>
  );
}

export function ChainGraph({ blocks, onSelectBlock }: Props) {
  // Compute SVG dimensions based on block count
  const { svgWidth, svgHeight } = useMemo(() => {
    if (blocks.length <= 1) return { svgWidth: 400, svgHeight: 400 };
    const lastRadius = 40 + (blocks.length - 1) * 28;
    const dim = Math.max(400, (lastRadius + NODE_RADIUS + 40) * 2);
    return { svgWidth: dim, svgHeight: dim };
  }, [blocks.length]);

  // Precompute all node positions
  const positions = useMemo(
    () => blocks.map((_, i) => getNodePosition(i, blocks.length, svgWidth, svgHeight)),
    [blocks, svgWidth, svgHeight],
  );

  const handleClick = useCallback(
    (block: Block) => {
      onSelectBlock(block);
    },
    [onSelectBlock],
  );

  if (blocks.length === 0) {
    return (
      <div className="flex items-center justify-center h-48 text-gray-500">No blocks in the chain</div>
    );
  }

  return (
    <>
      <style>{`
        @keyframes dashFlow {
          to { stroke-dashoffset: -16; }
        }
        .dash-line {
          animation: dashFlow 0.8s linear infinite;
        }
        @keyframes dashFlowV {
          to { stroke-dashoffset: -14; }
        }
        .dash-flow-v {
          animation: dashFlowV 0.6s linear infinite;
        }
        .graph-node {
          cursor: pointer;
        }
        .graph-node:hover {
          filter: brightness(1.3);
        }
      `}</style>

      <div className="w-full px-4 py-4">
        {/* Desktop: SVG network graph */}
        <div className="hidden sm:flex justify-center overflow-auto max-h-[70vh]">
          <svg
            width={svgWidth}
            height={svgHeight}
            viewBox={`0 0 ${svgWidth} ${svgHeight}`}
            className="select-none"
          >
            {/* Layer 1: Connection lines */}
            {blocks.map((block, i) => {
              if (i === 0) return null;
              const from = positions[i - 1];
              const to = positions[i];
              return (
                <line
                  key={`line-${block.index}`}
                  x1={from.x}
                  y1={from.y}
                  x2={to.x}
                  y2={to.y}
                  stroke="#4b5563"
                  strokeWidth="1.5"
                  strokeDasharray="6 4"
                  className="dash-line"
                  opacity={0.6}
                />
              );
            })}

            {/* Layer 2: Nodes */}
            {blocks.map((block, i) => {
              const pos = positions[i];
              const isGenesis = block.index === 0;
              const hue = hashToHue(block.hash);
              const strokeColor = isGenesis ? '#f59e0b' : `hsl(${hue}, 55%, 50%)`;
              const glowColor = isGenesis
                ? 'rgba(245, 158, 11, 0.4)'
                : `hsla(${hue}, 55%, 50%, 0.3)`;

              return (
                <g
                  key={`node-${block.index}`}
                  className="graph-node"
                  transform={`translate(${pos.x}, ${pos.y})`}
                  onClick={() => handleClick(block)}
                >
                  {/* Glow circle — genesis only */}
                  {isGenesis && (
                    <circle
                      r={NODE_RADIUS + 6}
                      fill="none"
                      stroke={glowColor}
                      strokeWidth="3"
                      opacity={0.5}
                    />
                  )}

                  {/* Main circle */}
                  <circle
                    r={NODE_RADIUS}
                    fill="rgba(17, 24, 39, 0.85)"
                    stroke={strokeColor}
                    strokeWidth={isGenesis ? 3 : 2}
                  />

                  {/* Block number */}
                  <text
                    textAnchor="middle"
                    dominantBaseline="central"
                    fill={isGenesis ? '#fbbf24' : '#e5e7eb'}
                    fontSize="13"
                    fontWeight="bold"
                    fontFamily="monospace"
                    style={{ pointerEvents: 'none' }}
                  >
                    #{block.index}
                  </text>

                  {/* Genesis label */}
                  {isGenesis && (
                    <text
                      textAnchor="middle"
                      y={NODE_RADIUS + 16}
                      fill="#f59e0b"
                      fontSize="10"
                      fontWeight="600"
                      opacity={0.8}
                      style={{ pointerEvents: 'none' }}
                    >
                      GENESIS
                    </text>
                  )}
                </g>
              );
            })}
          </svg>
        </div>

        {/* Mobile: vertical layout */}
        <div className="flex flex-col items-center sm:hidden gap-0">
          {blocks.map((block, i) => (
            <MobileVerticalNode
              key={block.index}
              block={block}
              isLast={i === blocks.length - 1}
              onClick={onSelectBlock}
            />
          ))}
        </div>
      </div>
    </>
  );
}
