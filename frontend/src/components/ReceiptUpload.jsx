import { useCallback, useRef } from 'react';

/**
 * ReceiptUpload Component
 * Handles drag-and-drop and click-to-upload for receipt images
 */
export default function ReceiptUpload({ onImageUpload, isDisabled }) {
  const fileInputRef = useRef(null);

  const handleDrop = useCallback((e) => {
    e.preventDefault();
    e.stopPropagation();
    
    if (isDisabled) return;
    
    const files = e.dataTransfer?.files;
    if (files && files.length > 0) {
      handleFile(files[0]);
    }
  }, [isDisabled]);

  const handleDragOver = useCallback((e) => {
    e.preventDefault();
    e.stopPropagation();
  }, []);

  const handleClick = () => {
    if (!isDisabled && fileInputRef.current) {
      fileInputRef.current.click();
    }
  };

  const handleFileChange = (e) => {
    const files = e.target.files;
    if (files && files.length > 0) {
      handleFile(files[0]);
    }
  };

  const handleFile = (file) => {
    if (!file.type.startsWith('image/')) {
      alert('Please upload an image file');
      return;
    }

    const reader = new FileReader();
    reader.onload = (e) => {
      onImageUpload({
        file,
        dataUrl: e.target.result,
        name: file.name,
        size: file.size,
        type: file.type
      });
    };
    reader.readAsDataURL(file);
  };

  return (
    <div
      onClick={handleClick}
      onDrop={handleDrop}
      onDragOver={handleDragOver}
      className={`
        relative group cursor-pointer
        border-2 border-dashed rounded-xl
        transition-all duration-300 ease-out
        min-h-[200px] flex flex-col items-center justify-center
        ${isDisabled 
          ? 'border-gray-700 bg-gray-900/50 cursor-not-allowed opacity-50' 
          : 'border-cyan-500/40 bg-gradient-to-br from-cyan-950/20 to-transparent hover:border-cyan-400 hover:bg-cyan-950/30'
        }
      `}
    >
      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        onChange={handleFileChange}
        className="hidden"
        disabled={isDisabled}
      />
      
      {/* Upload Icon */}
      <div className={`
        w-16 h-16 rounded-full mb-4
        flex items-center justify-center
        transition-all duration-300
        ${isDisabled 
          ? 'bg-gray-800' 
          : 'bg-cyan-500/10 group-hover:bg-cyan-500/20 group-hover:scale-110'
        }
      `}>
        <svg 
          className={`w-8 h-8 ${isDisabled ? 'text-gray-600' : 'text-cyan-400'}`}
          fill="none" 
          stroke="currentColor" 
          viewBox="0 0 24 24"
        >
          <path 
            strokeLinecap="round" 
            strokeLinejoin="round" 
            strokeWidth={1.5} 
            d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" 
          />
        </svg>
      </div>

      {/* Text */}
      <p className={`text-sm font-medium mb-1 ${isDisabled ? 'text-gray-500' : 'text-gray-300'}`}>
        {isDisabled ? 'Upload disabled' : 'Drop receipt image here'}
      </p>
      <p className={`text-xs ${isDisabled ? 'text-gray-600' : 'text-gray-500'}`}>
        or click to browse
      </p>

      {/* Decorative corner accents */}
      {!isDisabled && (
        <>
          <div className="absolute top-2 left-2 w-4 h-4 border-l-2 border-t-2 border-cyan-500/30 rounded-tl" />
          <div className="absolute top-2 right-2 w-4 h-4 border-r-2 border-t-2 border-cyan-500/30 rounded-tr" />
          <div className="absolute bottom-2 left-2 w-4 h-4 border-l-2 border-b-2 border-cyan-500/30 rounded-bl" />
          <div className="absolute bottom-2 right-2 w-4 h-4 border-r-2 border-b-2 border-cyan-500/30 rounded-br" />
        </>
      )}
    </div>
  );
}

