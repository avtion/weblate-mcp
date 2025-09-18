package tools

import (
	"encoding/json"
	"fmt"

	"github.com/avtion/weblate-mcp/internal/mcp"
	"github.com/avtion/weblate-mcp/internal/weblate"
)

// LanguagesTools contains tools for language management
type LanguagesTools struct {
	client *weblate.Client
}

// NewLanguagesTools creates a new LanguagesTools instance
func NewLanguagesTools(client *weblate.Client) *LanguagesTools {
	return &LanguagesTools{client: client}
}

// ListLanguagesArgs represents arguments for listLanguages tool
type ListLanguagesArgs struct {
	ProjectSlug string `json:"projectSlug"`
}

// ListLanguages handles the listLanguages MCP tool
func (t *LanguagesTools) ListLanguages(arguments json.RawMessage) (interface{}, error) {
	var args ListLanguagesArgs
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

	languages, err := t.client.ListLanguages(args.ProjectSlug)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Error listing languages: %v", err),
		}, nil
	}

	var languageStrings []string
	for _, l := range languages {
		languageStrings = append(languageStrings, fmt.Sprintf("- **%s** (%s)\n  Direction: %s\n  URL: %s", 
			l.Name, l.Code, l.Direction, l.WebURL))
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Found %d languages in project %s:\n\n%s", 
					len(languages), args.ProjectSlug, joinStrings(languageStrings, "\n\n")),
			},
		},
	}, nil
}

// RegisterLanguagesTools registers all language-related MCP tools
func RegisterLanguagesTools(server *mcp.Server, client *weblate.Client) error {
	tools := NewLanguagesTools(client)

	// Register listLanguages tool
	err := server.RegisterTool("listLanguages", mcp.Tool{
		Description: "List languages available in a specific project",
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
		Handler: tools.ListLanguages,
	})
	if err != nil {
		return fmt.Errorf("failed to register listLanguages tool: %w", err)
	}

	return nil
}