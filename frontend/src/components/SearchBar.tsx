interface SearchBarProps {
  query: string;
  setQuery: (q: string) => void;
  total: number;
  filtered: number;
}

export function SearchBar({ query, setQuery, total, filtered }: SearchBarProps) {
  return (
    <div className="flex items-center gap-3 px-6 py-3">
      <div className="relative flex-1 max-w-md">
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Buscar por índice, hash o data..."
          className="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2 text-sm text-gray-100 placeholder-gray-500 focus:outline-none focus:border-blue-500"
        />
        {query && (
          <button
            onClick={() => setQuery('')}
            className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-300 text-sm px-1"
          >
            X
          </button>
        )}
      </div>
      {query && (
        <span className="text-xs text-gray-500">
          {filtered} de {total} bloques
        </span>
      )}
    </div>
  );
}
