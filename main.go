package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/avtion/weblate-mcp/internal/config"
	"github.com/avtion/weblate-mcp/internal/mcp"
	"github.com/avtion/weblate-mcp/internal/server"
)

func main() {
	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		// Don't fail if .env doesn't exist, just continue
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create MCP server
	mcpServer := mcp.NewServer("weblate-mcp-server", "1.3.0")

	// Set server description
	mcpServer.SetDescription(`This is a Weblate MCP server that provides tools for managing translations.

Available tools:
Translation Management:
- listProjects: List all available Weblate projects
- listComponents: List components in a specific project
- listLanguages: List languages available in a specific project
- searchStringInProject: Search for translations containing specific text
- getTranslationForKey: Get translation value for a specific key
- writeTranslation: Write or update a translation value
- searchTranslationsByKey: Search for translations by key pattern
- findTranslationsForKey: Find all translations for a specific key
- listTranslationKeys: List all translation keys in a project
- searchTranslationKeys: Search for translation keys by pattern

Change Tracking & History:
- listRecentChanges: List recent changes across all projects
- getProjectChanges: Get recent changes for a specific project
- getComponentChanges: Get recent changes for a specific component
- getChangesByUser: Get recent changes by a specific user

Translation Statistics Dashboard:
- getProjectStatistics: Get comprehensive project statistics with completion rates
- getComponentStatistics: Get detailed statistics for a specific component
- getProjectDashboard: Get full dashboard overview with all component statistics
- getTranslationStatistics: Get statistics for specific translation (project/component/language)
- getComponentLanguageProgress: Get translation progress for all languages in a component
- getLanguageStatistics: Get statistics for a language across all projects
- getUserStatistics: Get contribution statistics for a specific user`)

	// Create and register Weblate tools
	if err := server.RegisterTools(mcpServer, cfg); err != nil {
		log.Fatalf("Failed to register tools: %v", err)
	}

	// Run the server with STDIO transport
	if err := mcpServer.ServeStdio(); err != nil {
		log.Fatalf("Server error: %v", err)
		os.Exit(1)
	}
}