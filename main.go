// Package main implements an MCP server for multimodal receipt processing.
//
// This server exposes tools for loading images, parsing Textract OCR output,
// and writing structured receipt data to disk. It is designed to be used
// with an LLM that orchestrates the receipt extraction workflow.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"myprice/tools"
)

const (
	serverName    = "myprice-mcp"
	serverVersion = "0.1.0"
)

func main() {
	// Create the MCP server
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    serverName,
			Version: serverVersion,
		},
		&mcp.ServerOptions{
			HasTools: true,
		},
	)

	// Register tools using the typed AddTool function
	mcp.AddTool(server, tools.LoadImageTool(), tools.HandleLoadImage)
	mcp.AddTool(server, tools.LoadTextractTool(), tools.HandleLoadTextract)
	mcp.AddTool(server, tools.WriteOutputTool(), tools.HandleWriteOutput)

	log.Printf("Registered tools: load_image, load_textract, write_output")

	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down MCP server...")
		cancel()
	}()

	// Run the server over stdio
	log.Printf("Starting %s v%s MCP server over stdio...\n", serverName, serverVersion)

	transport := &mcp.StdioTransport{}
	if err := server.Run(ctx, transport); err != nil {
		if ctx.Err() != nil {
			// Context was cancelled, graceful shutdown
			log.Println("Server shutdown complete")
			return
		}
		log.Fatalf("Server error: %v", err)
	}
}
