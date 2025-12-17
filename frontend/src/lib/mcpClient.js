/**
 * MCP Client
 * 
 * This module provides the interface to the Go MCP server HTTP API.
 */

const API_BASE = 'http://localhost:8080/api';

/**
 * Upload an image file to the server
 * @param {File} file - The image file to upload
 * @returns {Promise<{success: boolean, file_path: string, file_name: string, size: number, mime_type: string}>}
 */
export async function uploadImage(file) {
  const formData = new FormData();
  formData.append('image', file);

  const response = await fetch(`${API_BASE}/upload`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || 'Upload failed');
  }

  return response.json();
}

/**
 * Run the full receipt analysis pipeline
 * @param {string} imagePath - Path to receipt image (server will auto-find or run Textract)
 * @returns {Promise<{textract: object, llmOutput: object}>}
 */
export async function runAnalysis(imagePath) {
  const response = await fetch(`${API_BASE}/analyze`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ image_path: imagePath }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || 'Analysis failed');
  }

  const result = await response.json();
  
  return {
    textract: result.textract,
    llmOutput: result.llm_output,
  };
}

/**
 * Check API health
 * @returns {Promise<{status: string, service: string, version: string}>}
 */
export async function checkHealth() {
  const response = await fetch(`${API_BASE}/health`);
  return response.json();
}

export default {
  uploadImage,
  runAnalysis,
  checkHealth,
};
