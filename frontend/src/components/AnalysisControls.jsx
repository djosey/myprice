/**
 * AnalysisControls Component
 * Provides controls for running the receipt analysis pipeline
 */
export default function AnalysisControls({ 
  onRunAnalysis, 
  onReset,
  isAnalyzing,
  hasImage,
  analysisComplete,
  currentStep
}) {
  const steps = [
    { id: 1, label: 'Upload Image', icon: '📷' },
    { id: 2, label: 'Load OCR', icon: '🔍' },
    { id: 3, label: 'LLM Analysis', icon: '🤖' },
    { id: 4, label: 'Complete', icon: '✅' },
  ];

  return (
    <div className="bg-gradient-to-r from-gray-900/80 via-gray-900/60 to-gray-900/80 backdrop-blur-sm rounded-2xl border border-gray-800 p-6">
      {/* Progress Steps */}
      <div className="flex items-center justify-between mb-6">
        {steps.map((step, index) => (
          <div key={step.id} className="flex items-center">
            {/* Step indicator */}
            <div className={`
              relative flex items-center justify-center w-10 h-10 rounded-full
              transition-all duration-300
              ${currentStep >= step.id 
                ? 'bg-cyan-500/20 border-2 border-cyan-500' 
                : 'bg-gray-800 border-2 border-gray-700'
              }
              ${currentStep === step.id && isAnalyzing ? 'animate-pulse-glow' : ''}
            `}>
              {currentStep > step.id ? (
                <svg className="w-5 h-5 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
              ) : (
                <span className="text-sm">{step.icon}</span>
              )}
              
              {/* Pulse ring for active step */}
              {currentStep === step.id && isAnalyzing && (
                <span className="absolute inset-0 rounded-full border-2 border-cyan-400 animate-ping opacity-30" />
              )}
            </div>

            {/* Step label */}
            <span className={`
              hidden sm:block ml-2 text-xs font-medium
              ${currentStep >= step.id ? 'text-gray-300' : 'text-gray-600'}
            `}>
              {step.label}
            </span>

            {/* Connector line */}
            {index < steps.length - 1 && (
              <div className={`
                hidden sm:block w-8 lg:w-16 h-0.5 mx-2
                transition-colors duration-300
                ${currentStep > step.id ? 'bg-cyan-500' : 'bg-gray-700'}
              `} />
            )}
          </div>
        ))}
      </div>

      {/* Action Buttons */}
      <div className="flex flex-col sm:flex-row gap-3">
        {/* Run Analysis Button */}
        <button
          onClick={onRunAnalysis}
          disabled={!hasImage || isAnalyzing}
          className={`
            flex-1 px-6 py-3 rounded-xl font-semibold
            flex items-center justify-center gap-2
            transition-all duration-300
            ${!hasImage || isAnalyzing
              ? 'bg-gray-800 text-gray-500 cursor-not-allowed'
              : 'bg-gradient-to-r from-cyan-600 to-cyan-500 text-white hover:from-cyan-500 hover:to-cyan-400 shadow-lg shadow-cyan-500/20 hover:shadow-cyan-500/40'
            }
          `}
        >
          {isAnalyzing ? (
            <>
              <svg className="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
              </svg>
              <span>Analyzing...</span>
            </>
          ) : (
            <>
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
              <span>Run Analysis</span>
            </>
          )}
        </button>

        {/* Reset Button */}
        <button
          onClick={onReset}
          disabled={isAnalyzing}
          className={`
            px-6 py-3 rounded-xl font-medium
            flex items-center justify-center gap-2
            transition-all duration-300
            ${isAnalyzing
              ? 'bg-gray-800 text-gray-600 cursor-not-allowed'
              : 'bg-gray-800 text-gray-300 hover:bg-gray-700 hover:text-white border border-gray-700'
            }
          `}
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          <span>Reset</span>
        </button>
      </div>

      {/* Status message */}
      {isAnalyzing && (
        <div className="mt-4 text-center">
          <p className="text-sm text-cyan-400 animate-pulse">
            {currentStep === 2 && '🔍 Loading OCR data...'}
            {currentStep === 3 && '🤖 Running LLM analysis...'}
          </p>
        </div>
      )}

      {analysisComplete && !isAnalyzing && (
        <div className="mt-4 text-center">
          <p className="text-sm text-emerald-400">
            ✅ Analysis complete! Review the results below.
          </p>
        </div>
      )}
    </div>
  );
}

