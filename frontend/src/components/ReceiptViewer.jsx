/**
 * ReceiptViewer Component
 * Displays the uploaded receipt image with metadata
 */
export default function ReceiptViewer({ image, isLoading }) {
  if (!image && !isLoading) {
    return (
      <div className="h-full min-h-[300px] flex items-center justify-center bg-gray-900/50 rounded-xl border border-gray-800">
        <div className="text-center text-gray-500">
          <svg className="w-12 h-12 mx-auto mb-3 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          <p className="text-sm">No image uploaded</p>
        </div>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="h-full min-h-[300px] flex items-center justify-center bg-gray-900/50 rounded-xl border border-gray-800">
        <div className="text-center">
          <div className="w-10 h-10 mx-auto mb-3 border-2 border-cyan-500 border-t-transparent rounded-full animate-spin" />
          <p className="text-sm text-gray-400">Processing image...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full">
      {/* Image Container */}
      <div className="relative flex-1 bg-gray-900/50 rounded-xl border border-gray-800 overflow-hidden">
        <img
          src={image.dataUrl}
          alt="Receipt"
          className="w-full h-full object-contain p-2"
        />
        
        {/* Overlay gradient at bottom */}
        <div className="absolute bottom-0 left-0 right-0 h-16 bg-gradient-to-t from-gray-900/90 to-transparent" />
      </div>

      {/* Metadata */}
      <div className="mt-3 p-3 bg-gray-900/50 rounded-lg border border-gray-800">
        <div className="flex items-center justify-between text-xs">
          <div className="flex items-center gap-2">
            <span className="w-2 h-2 bg-emerald-500 rounded-full" />
            <span className="text-gray-400 truncate max-w-[150px]" title={image.name}>
              {image.name}
            </span>
          </div>
          <span className="text-gray-500">
            {(image.size / 1024).toFixed(1)} KB
          </span>
        </div>
        <div className="mt-1 text-xs text-gray-600">
          {image.type}
        </div>
      </div>
    </div>
  );
}

