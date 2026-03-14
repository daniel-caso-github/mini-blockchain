interface Props {
  healthy: boolean;
}

export function HealthIndicator({ healthy }: Props) {
  return (
    <div className="flex items-center gap-2" title={healthy ? 'API: Online' : 'API: Offline'}>
      <div
        className={`w-3 h-3 rounded-full ${
          healthy ? 'bg-green-500 shadow-[0_0_6px_rgba(34,197,94,0.6)]' : 'bg-red-500 shadow-[0_0_6px_rgba(239,68,68,0.6)]'
        }`}
      />
      <span className="text-sm text-gray-400">{healthy ? 'Online' : 'Offline'}</span>
    </div>
  );
}
