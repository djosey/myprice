import { useState, useMemo } from 'react';

/**
 * JsonViewer Component
 * Displays JSON data with syntax highlighting and collapsible sections
 */
export default function JsonViewer({ 
  data, 
  title, 
  icon,
  accentColor = 'cyan',
  isLoading = false,
  emptyMessage = 'No data available'
}) {
  const [isExpanded, setIsExpanded] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');

  const colorClasses = {
    cyan: {
      border: 'border-cyan-500/30',
      bg: 'bg-cyan-500/5',
      text: 'text-cyan-400',
      accent: 'bg-cyan-500',
    },
    emerald: {
      border: 'border-emerald-500/30',
      bg: 'bg-emerald-500/5',
      text: 'text-emerald-400',
      accent: 'bg-emerald-500',
    },
    amber: {
      border: 'border-amber-500/30',
      bg: 'bg-amber-500/5',
      text: 'text-amber-400',
      accent: 'bg-amber-500',
    },
  };

  const colors = colorClasses[accentColor] || colorClasses.cyan;

  const formattedJson = useMemo(() => {
    if (!data) return null;
    return JSON.stringify(data, null, 2);
  }, [data]);

  const highlightedJson = useMemo(() => {
    if (!formattedJson) return null;
    
    // Apply syntax highlighting
    return formattedJson
      .replace(/"([^"]+)":/g, '<span class="text-purple-400">"$1"</span>:')
      .replace(/: "([^"]+)"/g, ': <span class="text-emerald-400">"$1"</span>')
      .replace(/: (\d+\.?\d*)/g, ': <span class="text-amber-400">$1</span>')
      .replace(/: (true|false)/g, ': <span class="text-cyan-400">$1</span>')
      .replace(/: (null)/g, ': <span class="text-gray-500">$1</span>');
  }, [formattedJson]);

  const copyToClipboard = async () => {
    if (formattedJson) {
      await navigator.clipboard.writeText(formattedJson);
    }
  };

  if (isLoading) {
    return (
      <div className={`rounded-xl border ${colors.border} ${colors.bg} overflow-hidden`}>
        {/* Header */}
        <div className="px-4 py-3 border-b border-gray-800 flex items-center justify-between">
          <div className="flex items-center gap-2">
            {icon && <span className="text-lg">{icon}</span>}
            <span className="font-medium text-gray-200">{title}</span>
          </div>
        </div>
        
        {/* Loading skeleton */}
        <div className="p-4 space-y-2">
          {[...Array(8)].map((_, i) => (
            <div 
              key={i} 
              className="h-4 bg-gray-800 rounded animate-pulse"
              style={{ width: `${Math.random() * 40 + 60}%` }}
            />
          ))}
        </div>
      </div>
    );
  }

  if (!data) {
    return (
      <div className={`rounded-xl border ${colors.border} ${colors.bg} overflow-hidden`}>
        {/* Header */}
        <div className="px-4 py-3 border-b border-gray-800 flex items-center justify-between">
          <div className="flex items-center gap-2">
            {icon && <span className="text-lg">{icon}</span>}
            <span className="font-medium text-gray-200">{title}</span>
          </div>
        </div>
        
        {/* Empty state */}
        <div className="p-8 flex flex-col items-center justify-center text-gray-500">
          <svg className="w-10 h-10 mb-3 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <p className="text-sm">{emptyMessage}</p>
        </div>
      </div>
    );
  }

  return (
    <div className={`rounded-xl border ${colors.border} ${colors.bg} overflow-hidden`}>
      {/* Header */}
      <div className="px-4 py-3 border-b border-gray-800 flex items-center justify-between">
        <div className="flex items-center gap-2">
          {icon && <span className="text-lg">{icon}</span>}
          <span className="font-medium text-gray-200">{title}</span>
          <span className={`w-2 h-2 ${colors.accent} rounded-full`} />
        </div>
        
        <div className="flex items-center gap-2">
          {/* Copy button */}
          <button
            onClick={copyToClipboard}
            className="p-1.5 text-gray-500 hover:text-gray-300 hover:bg-gray-800 rounded transition-colors"
            title="Copy JSON"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
          </button>
          
          {/* Expand/Collapse button */}
          <button
            onClick={() => setIsExpanded(!isExpanded)}
            className="p-1.5 text-gray-500 hover:text-gray-300 hover:bg-gray-800 rounded transition-colors"
          >
            <svg 
              className={`w-4 h-4 transition-transform ${isExpanded ? 'rotate-180' : ''}`} 
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M19 9l-7 7-7-7" />
            </svg>
          </button>
        </div>
      </div>

      {/* Content */}
      <div className={`transition-all duration-300 ${isExpanded ? 'max-h-[600px]' : 'max-h-0'} overflow-hidden`}>
        <pre 
          className="p-4 text-xs font-mono overflow-auto max-h-[550px] leading-relaxed"
          dangerouslySetInnerHTML={{ __html: highlightedJson }}
        />
      </div>

      {/* Stats footer */}
      <div className="px-4 py-2 border-t border-gray-800 flex items-center justify-between text-xs text-gray-500">
        <span>{Object.keys(data).length} keys</span>
        <span>{formattedJson.length} chars</span>
      </div>
    </div>
  );
}

