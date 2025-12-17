import { useState, useCallback, useEffect } from 'react';
import ReceiptUpload from '../components/ReceiptUpload';
import ReceiptViewer from '../components/ReceiptViewer';
import JsonViewer from '../components/JsonViewer';
import AnalysisControls from '../components/AnalysisControls';
import { uploadImage, runAnalysis, checkHealth } from '../lib/mcpClient';

/**
 * DemoPage Component
 * Main container for the receipt analysis demo
 * Manages all state and coordinates the analysis workflow
 */
export default function DemoPage() {
  // State management
  const [image, setImage] = useState(null);
  const [uploadedPath, setUploadedPath] = useState(null);
  const [textractData, setTextractData] = useState(null);
  const [llmOutput, setLlmOutput] = useState(null);
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [currentStep, setCurrentStep] = useState(1);
  const [error, setError] = useState(null);
  const [apiConnected, setApiConnected] = useState(false);

  // Check API health on mount
  useEffect(() => {
    checkHealth()
      .then(() => setApiConnected(true))
      .catch(() => setApiConnected(false));
  }, []);

  // Handle image upload
  const handleImageUpload = useCallback(async (imageData) => {
    console.log('[DEBUG] handleImageUpload called with:', imageData.name);
    setImage(imageData);
    setCurrentStep(1);
    setTextractData(null);
    setLlmOutput(null);
    setError(null);

    // Upload the image to the server
    if (apiConnected) {
      try {
        console.log('[DEBUG] Uploading image to server...');
        const result = await uploadImage(imageData.file);
        console.log('[DEBUG] Upload result:', result);
        setUploadedPath(result.file_path);
        console.log('[DEBUG] Set uploadedPath to:', result.file_path);
      } catch (err) {
        console.error('[DEBUG] Image upload failed:', err);
        setUploadedPath(null);
      }
    } else {
      console.warn('[DEBUG] API not connected, skipping upload');
    }
  }, [apiConnected]);

  // Run the analysis pipeline
  const handleRunAnalysis = useCallback(async () => {
    if (!image) return;

    console.log('[DEBUG] handleRunAnalysis called');
    console.log('[DEBUG] uploadedPath:', uploadedPath);
    console.log('[DEBUG] image.name:', image.name);
    
    const imagePath = uploadedPath || image.name;
    console.log('[DEBUG] Using imagePath for analysis:', imagePath);

    setIsAnalyzing(true);
    setError(null);

    try {
      // Step 2: Load OCR
      setCurrentStep(2);
      await new Promise(resolve => setTimeout(resolve, 300));

      // Step 3: Run LLM analysis
      setCurrentStep(3);
      
      // Call the API - it will auto-detect or run Textract
      console.log('[DEBUG] Calling runAnalysis with:', imagePath);
      const result = await runAnalysis(imagePath);
      console.log('[DEBUG] Analysis result:', result);

      // Update state with results
      setTextractData(result.textract);
      setLlmOutput(result.llmOutput);
      
      // Step 4: Complete
      setCurrentStep(4);
    } catch (err) {
      console.error('[DEBUG] Analysis failed:', err);
      setError(err.message || 'Analysis failed. Is the API server running?');
      setCurrentStep(1);
    } finally {
      setIsAnalyzing(false);
    }
  }, [image, uploadedPath]);

  // Reset all state
  const handleReset = useCallback(() => {
    setImage(null);
    setUploadedPath(null);
    setTextractData(null);
    setLlmOutput(null);
    setIsAnalyzing(false);
    setCurrentStep(1);
    setError(null);
  }, []);

  return (
    <div className="min-h-screen bg-[#0a0a0f]">
      {/* Header */}
      <header className="border-b border-gray-800 bg-gray-900/50 backdrop-blur-sm sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-cyan-500 to-emerald-500 flex items-center justify-center">
                <span className="text-xl">🧾</span>
              </div>
              <div>
                <h1 className="text-xl font-bold text-white">MyPrice</h1>
                <p className="text-xs text-gray-500">MCP Receipt Analyzer</p>
              </div>
            </div>
            
            <div className="flex items-center gap-2 text-xs">
              <span className={`px-2 py-1 rounded-full border ${
                apiConnected 
                  ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20'
                  : 'bg-red-500/10 text-red-400 border-red-500/20'
              }`}>
                {apiConnected ? 'API Connected' : 'API Offline'}
              </span>
              <span className="px-2 py-1 bg-gray-800 text-gray-400 rounded-full">
                v0.1.0
              </span>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* API Warning */}
        {!apiConnected && (
          <div className="mb-6 p-4 bg-amber-500/10 border border-amber-500/30 rounded-xl text-amber-400 text-sm">
            <div className="flex items-center gap-2">
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <span>
                API server not running. Start it with: <code className="bg-gray-800 px-2 py-0.5 rounded">go run ./cmd/api</code>
              </span>
            </div>
          </div>
        )}

        {/* Controls Section */}
        <section className="mb-8">
          <AnalysisControls
            onRunAnalysis={handleRunAnalysis}
            onReset={handleReset}
            isAnalyzing={isAnalyzing}
            hasImage={!!image}
            analysisComplete={currentStep === 4}
            currentStep={currentStep}
          />
        </section>

        {/* Error Display */}
        {error && (
          <div className="mb-6 p-4 bg-red-500/10 border border-red-500/30 rounded-xl text-red-400 text-sm">
            <div className="flex items-center gap-2">
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span>{error}</span>
            </div>
          </div>
        )}

        {/* Three Column Grid */}
        <section className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Column 1: Receipt Image */}
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold text-white flex items-center gap-2">
                <span className="w-6 h-6 rounded-lg bg-cyan-500/20 text-cyan-400 flex items-center justify-center text-sm">1</span>
                Receipt Image
              </h2>
            </div>
            
            {!image ? (
              <ReceiptUpload 
                onImageUpload={handleImageUpload}
                isDisabled={isAnalyzing}
              />
            ) : (
              <ReceiptViewer 
                image={image}
                isLoading={isAnalyzing && currentStep < 2}
              />
            )}
          </div>

          {/* Column 2: Textract OCR */}
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold text-white flex items-center gap-2">
                <span className="w-6 h-6 rounded-lg bg-emerald-500/20 text-emerald-400 flex items-center justify-center text-sm">2</span>
                Textract OCR
              </h2>
            </div>
            
            <JsonViewer
              data={textractData}
              title="OCR Results"
              icon="🔍"
              accentColor="emerald"
              isLoading={isAnalyzing && currentStep === 2}
              emptyMessage="Run analysis to see OCR results"
            />
          </div>

          {/* Column 3: LLM Output */}
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold text-white flex items-center gap-2">
                <span className="w-6 h-6 rounded-lg bg-amber-500/20 text-amber-400 flex items-center justify-center text-sm">3</span>
                Structured Output
              </h2>
            </div>
            
            <JsonViewer
              data={llmOutput}
              title="Parsed Receipt"
              icon="🤖"
              accentColor="amber"
              isLoading={isAnalyzing && currentStep === 3}
              emptyMessage="Run analysis to see structured output"
            />
          </div>
        </section>

        {/* Summary Card (shown when analysis is complete) */}
        {llmOutput && currentStep === 4 && (
          <section className="mt-8">
            <div className="bg-gradient-to-r from-gray-900 via-gray-800 to-gray-900 rounded-2xl border border-gray-700 p-6">
              <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                <span>📊</span> Receipt Summary
              </h3>
              
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div className="bg-gray-900/50 rounded-xl p-4 border border-gray-800">
                  <p className="text-xs text-gray-500 mb-1">Vendor</p>
                  <p className="text-lg font-semibold text-white">{llmOutput.vendor || 'Unknown'}</p>
                </div>
                <div className="bg-gray-900/50 rounded-xl p-4 border border-gray-800">
                  <p className="text-xs text-gray-500 mb-1">Date</p>
                  <p className="text-lg font-semibold text-white">{llmOutput.date || 'N/A'}</p>
                </div>
                <div className="bg-gray-900/50 rounded-xl p-4 border border-gray-800">
                  <p className="text-xs text-gray-500 mb-1">Items</p>
                  <p className="text-lg font-semibold text-white">{llmOutput.items?.length || 0}</p>
                </div>
                <div className="bg-gray-900/50 rounded-xl p-4 border border-gray-800">
                  <p className="text-xs text-gray-500 mb-1">Total</p>
                  <p className="text-lg font-semibold text-emerald-400">
                    ${typeof llmOutput.total === 'number' ? llmOutput.total.toFixed(2) : '0.00'}
                  </p>
                </div>
              </div>

              {/* Items list */}
              {llmOutput.items && llmOutput.items.length > 0 && (
                <div className="mt-4 p-4 bg-gray-900/50 rounded-xl border border-gray-800">
                  <p className="text-xs text-gray-500 mb-3">Line Items</p>
                  <div className="space-y-2">
                    {llmOutput.items.map((item, index) => (
                      <div key={index} className="flex items-center justify-between text-sm">
                        <span className="text-gray-300">
                          {item.qty > 1 && <span className="text-gray-500">{item.qty}× </span>}
                          {item.name}
                        </span>
                        <span className="text-gray-400">${item.price?.toFixed(2) || '0.00'}</span>
                      </div>
                    ))}
                    {llmOutput.subtotal > 0 && (
                      <div className="border-t border-gray-700 pt-2 mt-2 flex items-center justify-between text-sm font-medium">
                        <span className="text-gray-300">Subtotal</span>
                        <span className="text-gray-300">${llmOutput.subtotal?.toFixed(2)}</span>
                      </div>
                    )}
                    {llmOutput.tax > 0 && (
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-gray-400">Tax</span>
                        <span className="text-gray-400">${llmOutput.tax?.toFixed(2)}</span>
                      </div>
                    )}
                    {llmOutput.total > 0 && (
                      <div className="flex items-center justify-between text-base font-bold">
                        <span className="text-white">Total</span>
                        <span className="text-emerald-400">${llmOutput.total?.toFixed(2)}</span>
                      </div>
                    )}
                  </div>
                </div>
              )}
            </div>
          </section>
        )}
      </main>

      {/* Footer */}
      <footer className="border-t border-gray-800 mt-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex items-center justify-between text-xs text-gray-500">
            <span>Built with Go MCP + React + Tailwind</span>
            <span>🧾 MyPrice Demo</span>
          </div>
        </div>
      </footer>
    </div>
  );
}
