package tools

import (
	"encoding/json"
	"fmt"

	"github.com/avtion/weblate-mcp/internal/mcp"
	"github.com/avtion/weblate-mcp/internal/weblate"
)

// ProjectsTools contains tools for project management
type ProjectsTools struct {
	client *weblate.Client
}

// NewProjectsTools creates a new ProjectsTools instance
func NewProjectsTools(client *weblate.Client) *ProjectsTools {
	return &ProjectsTools{client: client}
}

// ListProjects handles the listProjects MCP tool
func (t *ProjectsTools) ListProjects(arguments json.RawMessage) (interface{}, error) {
	projects, err := t.client.ListProjects()
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Error listing projects: %v", err),
		}, nil
	}

	var projectStrings []string
	for _, p := range projects {
		webURL := p.WebURL
		if webURL == "" && p.Web != "" {
			webURL = p.Web
		}
		projectStrings = append(projectStrings, fmt.Sprintf("- **%s** (%s)\n  URL: %s", p.Name, p.Slug, webURL))
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Found %d projects:\n\n%s", len(projects), 
					joinStrings(projectStrings, "\n\n")),
			},
		},
	}, nil
}

// RegisterProjectsTools registers all project-related MCP tools
func RegisterProjectsTools(server *mcp.Server, client *weblate.Client) error {
	tools := NewProjectsTools(client)

	// Register listProjects tool
	err := server.RegisterTool("listProjects", mcp.Tool{
		Description: "List all available Weblate projects",
		Parameters: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
		Handler: tools.ListProjects,
	})
	if err != nil {
		return fmt.Errorf("failed to register listProjects tool: %w", err)
	}

	return nil
}

// joinStrings is a simple utility to join strings
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}