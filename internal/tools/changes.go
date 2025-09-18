package tools

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/avtion/weblate-mcp/internal/mcp"
	"github.com/avtion/weblate-mcp/internal/weblate"
)

// ChangesTools contains tools for change tracking
type ChangesTools struct {
	client *weblate.Client
}

// NewChangesTools creates a new ChangesTools instance
func NewChangesTools(client *weblate.Client) *ChangesTools {
	return &ChangesTools{client: client}
}

// ListRecentChangesArgs represents arguments for listRecentChanges tool
type ListRecentChangesArgs struct {
	Limit int `json:"limit,omitempty"`
}

// GetChangesByUserArgs represents arguments for getChangesByUser tool
type GetChangesByUserArgs struct {
	User  string `json:"user"`
	Limit int    `json:"limit,omitempty"`
}

// GetProjectChangesArgs represents arguments for getProjectChanges tool
type GetProjectChangesArgs struct {
	ProjectSlug string `json:"projectSlug"`
	Limit       int    `json:"limit,omitempty"`
}

// GetComponentChangesArgs represents arguments for getComponentChanges tool
type GetComponentChangesArgs struct {
	ProjectSlug   string `json:"projectSlug"`
	ComponentSlug string `json:"componentSlug"`
	Limit         int    `json:"limit,omitempty"`
}

// ListRecentChanges handles the listRecentChanges MCP tool
func (t *ChangesTools) ListRecentChanges(arguments json.RawMessage) (interface{}, error) {
	var args ListRecentChangesArgs
	if err := json.Unmarshal(arguments, &args); err != nil {
		return map[string]interface{}{
			"error": "Invalid arguments: " + err.Error(),
		}, nil
	}

	limit := args.Limit
	if limit <= 0 {
		limit = 50
	}

	params := map[string]string{
		"page_size": strconv.Itoa(limit),
	}

	changes, err := t.client.ListChanges(params)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Error listing changes: %v", err),
		}, nil
	}

	if len(changes) == 0 {
		return map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": "No recent changes found",
				},
			},
		}, nil
	}

	var changeStrings []string
	for i, change := range changes {
		if i >= 10 { // Limit display to first 10 changes
			changeStrings = append(changeStrings, fmt.Sprintf("... and %d more changes", len(changes)-10))
			break
		}

		author := change.Author
		if author == "" {
			author = change.User
		}
		if author == "" {
			author = "Unknown"
		}

		changeStrings = append(changeStrings, fmt.Sprintf("- **%s** by %s\n  Time: %s\n  Target: %s",
			change.ActionName,
			author,
			change.Timestamp.Format("2006-01-02 15:04:05"),
			change.Target))
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Found %d recent changes:\n\n%s",
					len(changes), joinStrings(changeStrings, "\n\n")),
			},
		},
	}, nil
}

// GetChangesByUser handles the getChangesByUser MCP tool
func (t *ChangesTools) GetChangesByUser(arguments json.RawMessage) (interface{}, error) {
	var args GetChangesByUserArgs
	if err := json.Unmarshal(arguments, &args); err != nil {
		return map[string]interface{}{
			"error": "Invalid arguments: " + err.Error(),
		}, nil
	}

	if args.User == "" {
		return map[string]interface{}{
			"error": "user is required",
		}, nil
	}

	limit := args.Limit
	if limit <= 0 {
		limit = 50
	}

	params := map[string]string{
		"user":      args.User,
		"page_size": strconv.Itoa(limit),
	}

	changes, err := t.client.ListChanges(params)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Error listing changes by user: %v", err),
		}, nil
	}

	if len(changes) == 0 {
		return map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("No changes found for user %s", args.User),
				},
			},
		}, nil
	}

	var changeStrings []string
	for i, change := range changes {
		if i >= 10 { // Limit display to first 10 changes
			changeStrings = append(changeStrings, fmt.Sprintf("... and %d more changes", len(changes)-10))
			break
		}

		changeStrings = append(changeStrings, fmt.Sprintf("- **%s**\n  Time: %s\n  Target: %s",
			change.ActionName,
			change.Timestamp.Format("2006-01-02 15:04:05"),
			change.Target))
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Found %d changes by user %s:\n\n%s",
					len(changes), args.User, joinStrings(changeStrings, "\n\n")),
			},
		},
	}, nil
}

// GetProjectChanges handles the getProjectChanges MCP tool
func (t *ChangesTools) GetProjectChanges(arguments json.RawMessage) (interface{}, error) {
	var args GetProjectChangesArgs
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

	limit := args.Limit
	if limit <= 0 {
		limit = 50
	}

	params := map[string]string{
		"project":   args.ProjectSlug,
		"page_size": strconv.Itoa(limit),
	}

	changes, err := t.client.ListChanges(params)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Error listing project changes: %v", err),
		}, nil
	}

	if len(changes) == 0 {
		return map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("No changes found for project %s", args.ProjectSlug),
				},
			},
		}, nil
	}

	var changeStrings []string
	for i, change := range changes {
		if i >= 10 { // Limit display to first 10 changes
			changeStrings = append(changeStrings, fmt.Sprintf("... and %d more changes", len(changes)-10))
			break
		}

		author := change.Author
		if author == "" {
			author = change.User
		}
		if author == "" {
			author = "Unknown"
		}

		changeStrings = append(changeStrings, fmt.Sprintf("- **%s** by %s\n  Time: %s\n  Target: %s",
			change.ActionName,
			author,
			change.Timestamp.Format("2006-01-02 15:04:05"),
			change.Target))
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Found %d changes in project %s:\n\n%s",
					len(changes), args.ProjectSlug, joinStrings(changeStrings, "\n\n")),
			},
		},
	}, nil
}

// GetComponentChanges handles the getComponentChanges MCP tool
func (t *ChangesTools) GetComponentChanges(arguments json.RawMessage) (interface{}, error) {
	var args GetComponentChangesArgs
	if err := json.Unmarshal(arguments, &args); err != nil {
		return map[string]interface{}{
			"error": "Invalid arguments: " + err.Error(),
		}, nil
	}

	if args.ProjectSlug == "" || args.ComponentSlug == "" {
		return map[string]interface{}{
			"error": "projectSlug and componentSlug are required",
		}, nil
	}

	limit := args.Limit
	if limit <= 0 {
		limit = 50
	}

	params := map[string]string{
		"component": args.ComponentSlug,
		"project":   args.ProjectSlug,
		"page_size": strconv.Itoa(limit),
	}

	changes, err := t.client.ListChanges(params)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Error listing component changes: %v", err),
		}, nil
	}

	if len(changes) == 0 {
		return map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("No changes found for component %s/%s", args.ProjectSlug, args.ComponentSlug),
				},
			},
		}, nil
	}

	var changeStrings []string
	for i, change := range changes {
		if i >= 10 { // Limit display to first 10 changes
			changeStrings = append(changeStrings, fmt.Sprintf("... and %d more changes", len(changes)-10))
			break
		}

		author := change.Author
		if author == "" {
			author = change.User
		}
		if author == "" {
			author = "Unknown"
		}

		changeStrings = append(changeStrings, fmt.Sprintf("- **%s** by %s\n  Time: %s\n  Target: %s",
			change.ActionName,
			author,
			change.Timestamp.Format("2006-01-02 15:04:05"),
			change.Target))
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Found %d changes in component %s/%s:\n\n%s",
					len(changes), args.ProjectSlug, args.ComponentSlug, joinStrings(changeStrings, "\n\n")),
			},
		},
	}, nil
}

// RegisterChangesTools registers all change-related MCP tools
func RegisterChangesTools(server *mcp.Server, client *weblate.Client) error {
	tools := NewChangesTools(client)

	// Register listRecentChanges tool
	err := server.RegisterTool("listRecentChanges", mcp.Tool{
		Description: "List recent changes across all projects",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Maximum number of changes to return (default: 50)",
					"default":     50,
				},
			},
		},
		Handler: tools.ListRecentChanges,
	})
	if err != nil {
		return fmt.Errorf("failed to register listRecentChanges tool: %w", err)
	}

	// Register getChangesByUser tool
	err = server.RegisterTool("getChangesByUser", mcp.Tool{
		Description: "Get recent changes by a specific user",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"user": map[string]interface{}{
					"type":        "string",
					"description": "The username to filter changes by",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Maximum number of changes to return (default: 50)",
					"default":     50,
				},
			},
			"required": []string{"user"},
		},
		Handler: tools.GetChangesByUser,
	})
	if err != nil {
		return fmt.Errorf("failed to register getChangesByUser tool: %w", err)
	}

	// Register getProjectChanges tool
	err = server.RegisterTool("getProjectChanges", mcp.Tool{
		Description: "Get recent changes for a specific project",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"projectSlug": map[string]interface{}{
					"type":        "string",
					"description": "The slug of the project",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Maximum number of changes to return (default: 50)",
					"default":     50,
				},
			},
			"required": []string{"projectSlug"},
		},
		Handler: tools.GetProjectChanges,
	})
	if err != nil {
		return fmt.Errorf("failed to register getProjectChanges tool: %w", err)
	}

	// Register getComponentChanges tool
	err = server.RegisterTool("getComponentChanges", mcp.Tool{
		Description: "Get recent changes for a specific component",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"projectSlug": map[string]interface{}{
					"type":        "string",
					"description": "The slug of the project",
				},
				"componentSlug": map[string]interface{}{
					"type":        "string",
					"description": "The slug of the component",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Maximum number of changes to return (default: 50)",
					"default":     50,
				},
			},
			"required": []string{"projectSlug", "componentSlug"},
		},
		Handler: tools.GetComponentChanges,
	})
	if err != nil {
		return fmt.Errorf("failed to register getComponentChanges tool: %w", err)
	}

	return nil
}