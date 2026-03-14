const PAGE_SIZE_OPTIONS = [5, 10, 20, 50];

interface Props {
  page: number;
  totalPages: number;
  total: number;
  pageSize: number;
  onPageChange: (page: number) => void;
  onPageSizeChange: (size: number) => void;
}

export function GridPagination({ page, totalPages, total, pageSize, onPageChange, onPageSizeChange }: Props) {
  const pages: (number | '...')[] = [];
  for (let i = 1; i <= totalPages; i++) {
    if (i === 1 || i === totalPages || (i >= page - 1 && i <= page + 1)) {
      pages.push(i);
    } else if (pages[pages.length - 1] !== '...') {
      pages.push('...');
    }
  }

  return (
    <div className="flex items-center justify-center gap-2 py-3">
      <button
        onClick={() => onPageChange(page - 1)}
        disabled={page <= 1}
        className="px-2 py-1 text-xs rounded border border-gray-700 bg-gray-800/50 text-gray-300 hover:border-gray-500 hover:text-white disabled:opacity-30 disabled:cursor-not-allowed transition-colors cursor-pointer"
      >
        Prev
      </button>
      {pages.map((p, i) =>
        p === '...' ? (
          <span key={`dots-${i}`} className="text-gray-500 text-xs px-1">...</span>
        ) : (
          <button
            key={p}
            onClick={() => onPageChange(p)}
            className={`px-2.5 py-1 text-xs rounded border transition-colors cursor-pointer ${
              p === page
                ? 'border-blue-500 bg-blue-500/20 text-blue-300'
                : 'border-gray-700 bg-gray-800/50 text-gray-300 hover:border-gray-500 hover:text-white'
            }`}
          >
            {p}
          </button>
        ),
      )}
      <button
        onClick={() => onPageChange(page + 1)}
        disabled={page >= totalPages}
        className="px-2 py-1 text-xs rounded border border-gray-700 bg-gray-800/50 text-gray-300 hover:border-gray-500 hover:text-white disabled:opacity-30 disabled:cursor-not-allowed transition-colors cursor-pointer"
      >
        Next
      </button>
      <span className="text-xs text-gray-500 ml-2">{total} blocks</span>
      <select
        value={pageSize}
        onChange={(e) => onPageSizeChange(Number(e.target.value))}
        className="ml-2 px-1.5 py-1 text-xs rounded border border-gray-700 bg-gray-800 text-gray-300 cursor-pointer focus:outline-none focus:border-gray-500"
      >
        {PAGE_SIZE_OPTIONS.map((s) => (
          <option key={s} value={s}>{s} / page</option>
        ))}
      </select>
    </div>
  );
}
