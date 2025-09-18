package tools

import (
	"encoding/json"
	"fmt"

	"github.com/avtion/weblate-mcp/internal/mcp"
	"github.com/avtion/weblate-mcp/internal/weblate"
)

// ComponentsTools contains tools for component management
type ComponentsTools struct {
	client *weblate.Client
}

// NewComponentsTools creates a new ComponentsTools instance
func NewComponentsTools(client *weblate.Client) *ComponentsTools {
	return &ComponentsTools{client: client}
}

// ListComponentsArgs represents arguments for listComponents tool
type ListComponentsArgs struct {
	ProjectSlug string `json:"projectSlug"`
}

// ListComponents handles the listComponents MCP tool
func (t *ComponentsTools) ListComponents(arguments json.RawMessage) (interface{}, error) {
	var args ListComponentsArgs
	if err := json.Unmarshal(arguments, &args); err != nil {
		return map[string]interface{}{
			"error": "Invalid arguments: " + err.Error(),
		}, nil
	}

	if args.ProjectSlug == "" {
		return map[string]interface{}{
			"error": "projectSlug is required",
		}, nil
	}

	components, err := t.client.ListComponents(args.ProjectSlug)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Error listing components: %v", err),
		}, nil
	}

	var componentStrings []string
	for _, c := range components {
		componentStrings = append(componentStrings, fmt.Sprintf("- **%s** (%s)\n  Source Language: %s\n  URL: %s", 
			c.Name, c.Slug, c.SourceLanguage.Name, c.WebURL))
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Found %d components in project %s:\n\n%s", 
					len(components), args.ProjectSlug, joinStrings(componentStrings, "\n\n")),
			},
		},
	}, nil
}

// RegisterComponentsTools registers all component-related MCP tools
func RegisterComponentsTools(server *mcp.Server, client *weblate.Client) error {
	tools := NewComponentsTools(client)

	// Register listComponents tool
	err := server.RegisterTool("listComponents", mcp.Tool{
		Description: "List components in a specific project",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"projectSlug": map[string]interface{}{
					"type":        "string",
					"description": "The slug of the project",
				},
			},
			"required": []string{"projectSlug"},
		},
		Handler: tools.ListComponents,
	})
	if err != nil {
		return fmt.Errorf("failed to register listComponents tool: %w", err)
	}

	return nil
}