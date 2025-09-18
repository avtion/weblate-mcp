package tools

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/avtion/weblate-mcp/internal/mcp"
	"github.com/avtion/weblate-mcp/internal/weblate"
)

// PluralizationRule represents a language pluralization rule
type PluralizationRule struct {
	Forms int    `json:"forms"`
	Rule  string `json:"rule"`
}

// TranslationsTools contains tools for translation management
type TranslationsTools struct {
	client             *weblate.Client
	pluralizationRules map[string]PluralizationRule
}

// NewTranslationsTools creates a new TranslationsTools instance
func NewTranslationsTools(client *weblate.Client) *TranslationsTools {
	// Define pluralization rules (ported from TypeScript)
	rules := map[string]PluralizationRule{
		// 2 forms: singular (n=1), plural (n!=1)
		"en": {Forms: 2, Rule: "n != 1"}, // English
		"de": {Forms: 2, Rule: "n != 1"}, // German
		"es": {Forms: 2, Rule: "n != 1"}, // Spanish
		"fr": {Forms: 2, Rule: "n > 1"},  // French
		"it": {Forms: 2, Rule: "n != 1"}, // Italian
		"pt": {Forms: 2, Rule: "n != 1"}, // Portuguese
		"nl": {Forms: 2, Rule: "n != 1"}, // Dutch
		"da": {Forms: 2, Rule: "n != 1"}, // Danish
		"sv": {Forms: 2, Rule: "n != 1"}, // Swedish
		"no": {Forms: 2, Rule: "n != 1"}, // Norwegian

		// 3 forms: one, few, many/other
		"cs": {Forms: 3, Rule: "(n==1) ? 0 : (n>=2 && n<=4) ? 1 : 2"},                                                           // Czech
		"sk": {Forms: 3, Rule: "(n==1) ? 0 : (n>=2 && n<=4) ? 1 : 2"},                                                           // Slovak
		"pl": {Forms: 3, Rule: "(n==1) ? 0 : (n%10>=2 && n%10<=4 && (n%100<10 || n%100>=20)) ? 1 : 2"},                        // Polish
		"hr": {Forms: 3, Rule: "(n%10==1 && n%100!=11) ? 0 : (n%10>=2 && n%10<=4 && (n%100<10 || n%100>=20)) ? 1 : 2"},        // Croatian
		"sr": {Forms: 3, Rule: "(n%10==1 && n%100!=11) ? 0 : (n%10>=2 && n%10<=4 && (n%100<10 || n%100>=20)) ? 1 : 2"},        // Serbian

		// 4 forms: one, few, many, other
		"sl": {Forms: 4, Rule: "(n%100==1) ? 0 : (n%100==2) ? 1 : (n%100==3 || n%100==4) ? 2 : 3"}, // Slovenian

		// 6 forms: zero, one, two, few, many, other
		"ar": {Forms: 6, Rule: "(n==0) ? 0 : (n==1) ? 1 : (n==2) ? 2 : (n%100>=3 && n%100<=10) ? 3 : (n%100>=11) ? 4 : 5"}, // Arabic

		// Default fallback for unknown languages
		"default": {Forms: 2, Rule: "n != 1"},
	}

	return &TranslationsTools{
		client:             client,
		pluralizationRules: rules,
	}
}

// GetTranslationForKeyArgs represents arguments for getTranslationForKey tool
type GetTranslationForKeyArgs struct {
	ProjectSlug   string `json:"projectSlug"`
	ComponentSlug string `json:"componentSlug"`
	LanguageCode  string `json:"languageCode"`
	Key           string `json:"key"`
}

// WriteTranslationArgs represents arguments for writeTranslation tool
type WriteTranslationArgs struct {
	ProjectSlug    string `json:"projectSlug"`
	ComponentSlug  string `json:"componentSlug"`
	LanguageCode   string `json:"languageCode"`
	Key            string `json:"key"`
	Value          string `json:"value"`
	MarkAsApproved bool   `json:"markAsApproved,omitempty"`
}

// SearchStringInProjectArgs represents arguments for searchStringInProject tool
type SearchStringInProjectArgs struct {
	ProjectSlug string `json:"projectSlug"`
	SearchValue string `json:"searchValue"`
	SearchIn    string `json:"searchIn,omitempty"` // "source", "target", or "both"
}

// GetTranslationForKey handles the getTranslationForKey MCP tool
func (t *TranslationsTools) GetTranslationForKey(arguments json.RawMessage) (interface{}, error) {
	var args GetTranslationForKeyArgs
	if err := json.Unmarshal(arguments, &args); err != nil {
		return map[string]interface{}{
			"error": "Invalid arguments: " + err.Error(),
		}, nil
	}

	if args.ProjectSlug == "" || args.ComponentSlug == "" || args.LanguageCode == "" || args.Key == "" {
		return map[string]interface{}{
			"error": "projectSlug, componentSlug, languageCode, and key are required",
		}, nil
	}

	// Search for the translation unit by key
	params := map[string]string{
		"q": fmt.Sprintf("project:%s component:%s language:%s source:\"%s\"", 
			args.ProjectSlug, args.ComponentSlug, args.LanguageCode, args.Key),
		"page_size": "1",
	}

	units, err := t.client.ListUnits(args.ProjectSlug, args.ComponentSlug, args.LanguageCode, params)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Error searching for translation: %v", err),
		}, nil
	}

	if len(units) == 0 {
		return map[string]interface{}{
			"error": fmt.Sprintf("Translation not found for key \"%s\"", args.Key),
		}, nil
	}

	unit := units[0]
	targetValue := strings.Join(unit.Target, "")
	if len(unit.Target) > 1 {
		// For plural forms, show all forms
		targetValue = strings.Join(unit.Target, " | ")
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Translation for key \"%s\":\n\nSource: %s\nTarget: %s\nState: %s\nApproved: %t",
					args.Key,
					strings.Join(unit.Source, " | "),
					targetValue,
					t.getStateString(unit.State),
					unit.Approved),
			},
		},
	}, nil
}

// WriteTranslation handles the writeTranslation MCP tool
func (t *TranslationsTools) WriteTranslation(arguments json.RawMessage) (interface{}, error) {
	var args WriteTranslationArgs
	if err := json.Unmarshal(arguments, &args); err != nil {
		return map[string]interface{}{
			"error": "Invalid arguments: " + err.Error(),
		}, nil
	}

	if args.ProjectSlug == "" || args.ComponentSlug == "" || args.LanguageCode == "" || args.Key == "" || args.Value == "" {
		return map[string]interface{}{
			"error": "projectSlug, componentSlug, languageCode, key, and value are required",
		}, nil
	}

	// First, find the translation unit by key
	params := map[string]string{
		"q": fmt.Sprintf("project:%s component:%s language:%s source:\"%s\"", 
			args.ProjectSlug, args.ComponentSlug, args.LanguageCode, args.Key),
		"page_size": "1",
	}

	units, err := t.client.ListUnits(args.ProjectSlug, args.ComponentSlug, args.LanguageCode, params)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Error searching for translation: %v", err),
		}, nil
	}

	if len(units) == 0 {
		return map[string]interface{}{
			"error": fmt.Sprintf("Translation unit not found for key \"%s\"", args.Key),
		}, nil
	}

	unit := units[0]

	// Parse plural forms correctly for the target field using language-specific rules
	targetArray := t.parsePluralForms(args.Value, unit.Source, args.LanguageCode)

	// Update the translation using the units API
	state := 20 // 20 = translated
	if args.MarkAsApproved {
		state = 30 // 30 = approved
	}

	updatedUnit, err := t.client.UpdateUnit(unit.ID, targetArray, &state)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Error updating translation: %v", err),
		}, nil
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Successfully updated translation for key \"%s\":\n\nNew value: %s\nState: %s\nApproved: %t",
					args.Key,
					strings.Join(updatedUnit.Target, " | "),
					t.getStateString(updatedUnit.State),
					updatedUnit.Approved),
			},
		},
	}, nil
}

// SearchStringInProject handles the searchStringInProject MCP tool
func (t *TranslationsTools) SearchStringInProject(arguments json.RawMessage) (interface{}, error) {
	var args SearchStringInProjectArgs
	if err := json.Unmarshal(arguments, &args); err != nil {
		return map[string]interface{}{
			"error": "Invalid arguments: " + err.Error(),
		}, nil
	}

	if args.ProjectSlug == "" || args.SearchValue == "" {
		return map[string]interface{}{
			"error": "projectSlug and searchValue are required",
		}, nil
	}

	searchIn := args.SearchIn
	if searchIn == "" {
		searchIn = "both"
	}

	var results []weblate.Unit
	uniqueUnits := make(map[int]weblate.Unit)

	// Search in source if requested
	if searchIn == "source" || searchIn == "both" {
		params := map[string]string{
			"q": fmt.Sprintf("project:%s source:\"%s\"", args.ProjectSlug, args.SearchValue),
			"page_size": "100",
		}

		units, err := t.client.ListUnits(args.ProjectSlug, "", "", params)
		if err != nil {
			return map[string]interface{}{
				"error": fmt.Sprintf("Error searching in source: %v", err),
			}, nil
		}

		for _, unit := range units {
			uniqueUnits[unit.ID] = unit
		}
	}

	// Search in target if requested
	if searchIn == "target" || searchIn == "both" {
		params := map[string]string{
			"q": fmt.Sprintf("project:%s target:\"%s\"", args.ProjectSlug, args.SearchValue),
			"page_size": "100",
		}

		units, err := t.client.ListUnits(args.ProjectSlug, "", "", params)
		if err != nil {
			return map[string]interface{}{
				"error": fmt.Sprintf("Error searching in target: %v", err),
			}, nil
		}

		for _, unit := range units {
			uniqueUnits[unit.ID] = unit
		}
	}

	// Convert map to slice
	for _, unit := range uniqueUnits {
		results = append(results, unit)
	}

	if len(results) == 0 {
		return map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("No translations found containing \"%s\" in project %s", args.SearchValue, args.ProjectSlug),
				},
			},
		}, nil
	}

	var resultStrings []string
	for i, unit := range results {
		if i >= 10 { // Limit to first 10 results
			resultStrings = append(resultStrings, fmt.Sprintf("... and %d more results", len(results)-10))
			break
		}

		sourceText := strings.Join(unit.Source, " | ")
		targetText := strings.Join(unit.Target, " | ")
		
		resultStrings = append(resultStrings, fmt.Sprintf("- **Location**: %s\n  **Source**: %s\n  **Target**: %s\n  **State**: %s",
			unit.Location, sourceText, targetText, t.getStateString(unit.State)))
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Found %d translations containing \"%s\" in project %s:\n\n%s",
					len(results), args.SearchValue, args.ProjectSlug, strings.Join(resultStrings, "\n\n")),
			},
		},
	}, nil
}

// getExpectedPluralForms returns the expected number of plural forms for a language
func (t *TranslationsTools) getExpectedPluralForms(languageCode string) int {
	rule, exists := t.pluralizationRules[languageCode]
	if !exists {
		rule = t.pluralizationRules["default"]
	}
	return rule.Forms
}

// parsePluralForms parses plural forms from a concatenated string into an array
func (t *TranslationsTools) parsePluralForms(value string, sourceArray []string, languageCode string) []string {
	// If source is not an array or has only one element, treat as singular
	if len(sourceArray) <= 1 {
		return []string{value}
	}

	// Get expected plural forms count based on language rules
	expectedPluralCount := t.getExpectedPluralForms(languageCode)

	// Fallback to source array length if it's different (might be source language specific)
	targetPluralCount := expectedPluralCount
	if len(sourceArray) > expectedPluralCount {
		targetPluralCount = len(sourceArray)
	}

	// For plural forms, split on the pattern where %d starts a new plural form
	// This handles cases like "%d day%d days" -> ["%d day", "%d days"]
	percentDPattern := regexp.MustCompile(`(?=%d)`)
	percentDParts := percentDPattern.Split(value, -1)
	
	// Filter empty parts and add %d back to the beginning (except first part)
	var parts []string
	for i, part := range percentDParts {
		if len(part) > 0 {
			if i > 0 {
				part = "%d" + part
			}
			parts = append(parts, part)
		}
	}

	if len(parts) == targetPluralCount {
		return parts
	} else if len(parts) > 1 {
		// Method 2: Try to intelligently redistribute parts
		if len(parts) > targetPluralCount {
			// Too many parts - combine excess with last expected part
			result := parts[:targetPluralCount-1]
			lastPart := strings.Join(parts[targetPluralCount-1:], "")
			result = append(result, lastPart)
			return result
		}
		// Too few parts - use what we have
		return parts
	}

	// Fallback: if no %d patterns found or other issues, duplicate the value
	result := make([]string, targetPluralCount)
	for i := range result {
		result[i] = value
	}
	return result
}

// getStateString converts state integer to human-readable string
func (t *TranslationsTools) getStateString(state int) string {
	switch state {
	case 0:
		return "Empty"
	case 10:
		return "Needs editing"
	case 20:
		return "Translated"
	case 30:
		return "Approved"
	case 100:
		return "Read-only"
	default:
		return fmt.Sprintf("Unknown (%d)", state)
	}
}

// RegisterTranslationsTools registers all translation-related MCP tools
func RegisterTranslationsTools(server *mcp.Server, client *weblate.Client) error {
	tools := NewTranslationsTools(client)

	// Register getTranslationForKey tool
	err := server.RegisterTool("getTranslationForKey", mcp.Tool{
		Description: "Get translation value for a specific key",
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
				"key": map[string]interface{}{
					"type":        "string",
					"description": "The translation key to look up",
				},
			},
			"required": []string{"projectSlug", "componentSlug", "languageCode", "key"},
		},
		Handler: tools.GetTranslationForKey,
	})
	if err != nil {
		return fmt.Errorf("failed to register getTranslationForKey tool: %w", err)
	}

	// Register writeTranslation tool
	err = server.RegisterTool("writeTranslation", mcp.Tool{
		Description: "Write or update a translation value",
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
				"key": map[string]interface{}{
					"type":        "string",
					"description": "The translation key to update",
				},
				"value": map[string]interface{}{
					"type":        "string",
					"description": "The translation value",
				},
				"markAsApproved": map[string]interface{}{
					"type":        "boolean",
					"description": "Whether to mark the translation as approved",
					"default":     false,
				},
			},
			"required": []string{"projectSlug", "componentSlug", "languageCode", "key", "value"},
		},
		Handler: tools.WriteTranslation,
	})
	if err != nil {
		return fmt.Errorf("failed to register writeTranslation tool: %w", err)
	}

	// Register searchStringInProject tool
	err = server.RegisterTool("searchStringInProject", mcp.Tool{
		Description: "Search for translations containing specific text in a project",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"projectSlug": map[string]interface{}{
					"type":        "string",
					"description": "The slug of the project",
				},
				"searchValue": map[string]interface{}{
					"type":        "string",
					"description": "The text to search for",
				},
				"searchIn": map[string]interface{}{
					"type":        "string",
					"description": "Where to search: 'source', 'target', or 'both'",
					"default":     "both",
					"enum":        []string{"source", "target", "both"},
				},
			},
			"required": []string{"projectSlug", "searchValue"},
		},
		Handler: tools.SearchStringInProject,
	})
	if err != nil {
		return fmt.Errorf("failed to register searchStringInProject tool: %w", err)
	}

	return nil
}