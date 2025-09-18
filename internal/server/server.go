package server

import (
	"github.com/avtion/weblate-mcp/internal/config"
	"github.com/avtion/weblate-mcp/internal/mcp"
	"github.com/avtion/weblate-mcp/internal/tools"
	"github.com/avtion/weblate-mcp/internal/weblate"
)

// RegisterTools registers all MCP tools with the server
func RegisterTools(server *mcp.Server, cfg *config.Config) error {
	// Create Weblate client
	client := weblate.NewClient(cfg)

	// Register all tool categories
	if err := tools.RegisterProjectsTools(server, client); err != nil {
		return err
	}

	if err := tools.RegisterComponentsTools(server, client); err != nil {
		return err
	}

	if err := tools.RegisterLanguagesTools(server, client); err != nil {
		return err
	}

	if err := tools.RegisterTranslationsTools(server, client); err != nil {
		return err
	}

	if err := tools.RegisterChangesTools(server, client); err != nil {
		return err
	}

	if err := tools.RegisterStatisticsTools(server, client); err != nil {
		return err
	}

	return nil
}