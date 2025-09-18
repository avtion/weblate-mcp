package tools

import (
	"encoding/json"
	"fmt"

	"github.com/avtion/weblate-mcp/internal/mcp"
	"github.com/avtion/weblate-mcp/internal/weblate"
)

// StatisticsTools contains tools for statistics
type StatisticsTools struct {
	client *weblate.Client
}

// NewStatisticsTools creates a new StatisticsTools instance
func NewStatisticsTools(client *weblate.Client) *StatisticsTools {
	return &StatisticsTools{client: client}
}

// GetProjectStatisticsArgs represents arguments for getProjectStatistics tool
type GetProjectStatisticsArgs struct {
	ProjectSlug string `json:"projectSlug"`
}

// GetComponentStatisticsArgs represents arguments for getComponentStatistics tool
type GetComponentStatisticsArgs struct {
	ProjectSlug   string `json:"projectSlug"`
	ComponentSlug string `json:"componentSlug"`
}

// GetTranslationStatisticsArgs represents arguments for getTranslationStatistics tool
type GetTranslationStatisticsArgs struct {
	ProjectSlug   string `json:"projectSlug"`
	ComponentSlug string `json:"componentSlug"`
	LanguageCode  string `json:"languageCode"`
}

// GetProjectStatistics handles the getProjectStatistics MCP tool
func (t *StatisticsTools) GetProjectStatistics(arguments json.RawMessage) (interface{}, error) {
	var args GetProjectStatisticsArgs
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

	stats, err := t.client.GetProjectStatistics(args.ProjectSlug)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Error getting project statistics: %v", err),
		}, nil
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Project Statistics for **%s**:\n\n"+
					"📊 **Translation Progress**:\n"+
					"- Translated: %.1f%% (%d strings)\n"+
					"- Approved: %.1f%% (%d strings)\n"+
					"- Not translated: %.1f%% (%d strings)\n"+
					"- Read-only: %.1f%% (%d strings)\n"+
					"- Total strings: %d\n\n"+
					"🔗 **Links**:\n"+
					"- Web URL: %s\n"+
					"- Repository: %s",
					stats.Name,
					stats.TranslatedPercent, stats.Translated,
					stats.ApprovedPercent, stats.Approved,
					stats.NottranslatedPercent, stats.Nottranslated,
					stats.ReadonlyPercent, stats.Readonly,
					stats.Total,
					stats.WebURL,
					stats.RepositoryURL),
			},
		},
	}, nil
}

// GetComponentStatistics handles the getComponentStatistics MCP tool
func (t *StatisticsTools) GetComponentStatistics(arguments json.RawMessage) (interface{}, error) {
	var args GetComponentStatisticsArgs
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

	stats, err := t.client.GetComponentStatistics(args.ProjectSlug, args.ComponentSlug)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Error getting component statistics: %v", err),
		}, nil
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Component Statistics for **%s**:\n\n"+
					"📊 **Translation Progress**:\n"+
					"- Translated: %.1f%% (%d strings)\n"+
					"- Approved: %.1f%% (%d strings)\n"+
					"- Not translated: %.1f%% (%d strings)\n"+
					"- Read-only: %.1f%% (%d strings)\n"+
					"- Total strings: %d\n\n"+
					"🔗 **Links**:\n"+
					"- Web URL: %s\n"+
					"- Repository: %s",
					stats.Name,
					stats.TranslatedPercent, stats.Translated,
					stats.ApprovedPercent, stats.Approved,
					stats.NottranslatedPercent, stats.Nottranslated,
					stats.ReadonlyPercent, stats.Readonly,
					stats.Total,
					stats.WebURL,
					stats.RepositoryURL),
			},
		},
	}, nil
}

// GetTranslationStatistics handles the getTranslationStatistics MCP tool
func (t *StatisticsTools) GetTranslationStatistics(arguments json.RawMessage) (interface{}, error) {
	var args GetTranslationStatisticsArgs
	if err := json.Unmarshal(arguments, &args); err != nil {
		return map[string]interface{}{
			"error": "Invalid arguments: " + err.Error(),
		}, nil
	}

	if args.ProjectSlug == "" || args.ComponentSlug == "" || args.LanguageCode == "" {
		return map[string]interface{}{
			"error": "projectSlug, componentSlug, and languageCode are required",
		}, nil
	}

	stats, err := t.client.GetTranslationStatistics(args.ProjectSlug, args.ComponentSlug, args.LanguageCode)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Error getting translation statistics: %v", err),
		}, nil
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Translation Statistics for **%s** (%s):\n\n"+
					"📊 **Translation Progress**:\n"+
					"- Translated: %.1f%% (%d strings)\n"+
					"- Approved: %.1f%% (%d strings)\n"+
					"- Not translated: %.1f%% (%d strings)\n"+
					"- Read-only: %.1f%% (%d strings)\n"+
					"- Total strings: %d\n\n"+
					"🔗 **Links**:\n"+
					"- Web URL: %s\n"+
					"- Repository: %s",
					stats.Name, args.LanguageCode,
					stats.TranslatedPercent, stats.Translated,
					stats.ApprovedPercent, stats.Approved,
					stats.NottranslatedPercent, stats.Nottranslated,
					stats.ReadonlyPercent, stats.Readonly,
					stats.Total,
					stats.WebURL,
					stats.RepositoryURL),
			},
		},
	}, nil
}

// GetProjectDashboard provides a comprehensive dashboard view
func (t *StatisticsTools) GetProjectDashboard(arguments json.RawMessage) (interface{}, error) {
	var args GetProjectStatisticsArgs
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

	// Get project info and statistics
	project, err := t.client.GetProject(args.ProjectSlug)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Error getting project: %v", err),
		}, nil
	}

	projectStats, err := t.client.GetProjectStatistics(args.ProjectSlug)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Error getting project statistics: %v", err),
		}, nil
	}

	// Get components and their statistics
	components, err := t.client.ListComponents(args.ProjectSlug)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Error getting components: %v", err),
		}, nil
	}

	var componentStrings []string
	for _, component := range components {
		stats, err := t.client.GetComponentStatistics(args.ProjectSlug, component.Slug)
		if err != nil {
			componentStrings = append(componentStrings, fmt.Sprintf("- **%s**: Error getting statistics", component.Name))
			continue
		}

		progressBar := t.createProgressBar(stats.TranslatedPercent)
		componentStrings = append(componentStrings,
			fmt.Sprintf("- **%s** %s\n  %.1f%% translated (%d/%d strings)",
				component.Name, progressBar, stats.TranslatedPercent, stats.Translated, stats.Total))
	}

	// Create overall progress bar
	overallProgressBar := t.createProgressBar(projectStats.TranslatedPercent)

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("# 📊 Project Dashboard: %s\n\n"+
					"## 🎯 Overall Progress\n"+
					"%s **%.1f%%** translated (%d/%d strings)\n\n"+
					"📈 **Detailed Statistics:**\n"+
					"- ✅ Translated: %.1f%% (%d strings)\n"+
					"- ✔️ Approved: %.1f%% (%d strings)\n"+
					"- ❌ Not translated: %.1f%% (%d strings)\n"+
					"- 🔒 Read-only: %.1f%% (%d strings)\n\n"+
					"## 🔧 Components (%d total):\n\n%s\n\n"+
					"🔗 **Project Links:**\n"+
					"- Web: %s\n"+
					"- Repository: %s",
					project.Name,
					overallProgressBar, projectStats.TranslatedPercent, projectStats.Translated, projectStats.Total,
					projectStats.TranslatedPercent, projectStats.Translated,
					projectStats.ApprovedPercent, projectStats.Approved,
					projectStats.NottranslatedPercent, projectStats.Nottranslated,
					projectStats.ReadonlyPercent, projectStats.Readonly,
					len(components), joinStrings(componentStrings, "\n"),
					project.WebURL, projectStats.RepositoryURL),
			},
		},
	}, nil
}

// createProgressBar creates a visual progress bar
func (t *StatisticsTools) createProgressBar(percent float64) string {
	const barLength = 20
	filled := int(percent / 100.0 * barLength)
	if filled > barLength {
		filled = barLength
	}

	bar := ""
	for i := 0; i < barLength; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	return fmt.Sprintf("[%s]", bar)
}

// RegisterStatisticsTools registers all statistics-related MCP tools
func RegisterStatisticsTools(server *mcp.Server, client *weblate.Client) error {
	tools := NewStatisticsTools(client)

	// Register getProjectStatistics tool
	err := server.RegisterTool("getProjectStatistics", mcp.Tool{
		Description: "Get comprehensive project statistics with completion rates",
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
		Handler: tools.GetProjectStatistics,
	})
	if err != nil {
		return fmt.Errorf("failed to register getProjectStatistics tool: %w", err)
	}

	// Register getComponentStatistics tool
	err = server.RegisterTool("getComponentStatistics", mcp.Tool{
		Description: "Get detailed statistics for a specific component",
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
			},
			"required": []string{"projectSlug", "componentSlug"},
		},
		Handler: tools.GetComponentStatistics,
	})
	if err != nil {
		return fmt.Errorf("failed to register getComponentStatistics tool: %w", err)
	}

	// Register getTranslationStatistics tool
	err = server.RegisterTool("getTranslationStatistics", mcp.Tool{
		Description: "Get statistics for specific translation (project/component/language)",
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
				"languageCode": map[string]interface{}{
					"type":        "string",
					"description": "The language code",
				},
			},
			"required": []string{"projectSlug", "componentSlug", "languageCode"},
		},
		Handler: tools.GetTranslationStatistics,
	})
	if err != nil {
		return fmt.Errorf("failed to register getTranslationStatistics tool: %w", err)
	}

	// Register getProjectDashboard tool
	err = server.RegisterTool("getProjectDashboard", mcp.Tool{
		Description: "Get full dashboard overview with all component statistics",
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
		Handler: tools.GetProjectDashboard,
	})
	if err != nil {
		return fmt.Errorf("failed to register getProjectDashboard tool: %w", err)
	}

	return nil
}