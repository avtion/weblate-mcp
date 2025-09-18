package weblate

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/avtion/weblate-mcp/internal/config"
)

// Client provides access to the Weblate API
type Client struct {
	config *config.Config
	resty  *resty.Client
}

// NewClient creates a new Weblate API client
func NewClient(cfg *config.Config) *Client {
	client := resty.New().
		SetBaseURL(cfg.WeblateAPIURL).
		SetHeader("Authorization", fmt.Sprintf("Token %s", cfg.WeblateAPIToken)).
		SetHeader("Content-Type", "application/json").
		SetTimeout(10 * time.Second)

	// Disable logging for MCP STDIO compatibility
	if cfg.LogLevel != "debug" {
		client.SetDisableWarn(true)
		// Create a no-op logger by setting a custom writer
		client.SetLogger(nil)
	}

	return &Client{
		config: cfg,
		resty:  client,
	}
}

// GetRestyClient returns the underlying resty client for custom requests
func (c *Client) GetRestyClient() *resty.Client {
	return c.resty
}

// ListProjects lists all projects
func (c *Client) ListProjects() ([]Project, error) {
	var response PaginatedResponse[Project]
	_, err := c.resty.R().
		SetResult(&response).
		Get("/projects/")
	
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	
	return response.Results, nil
}

// GetProject gets a specific project by slug
func (c *Client) GetProject(slug string) (*Project, error) {
	var project Project
	resp, err := c.resty.R().
		SetResult(&project).
		Get(fmt.Sprintf("/projects/%s/", slug))
	
	if err != nil {
		return nil, fmt.Errorf("failed to get project %s: %w", slug, err)
	}
	
	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("project %s not found", slug)
	}
	
	return &project, nil
}

// ListComponents lists components for a project
func (c *Client) ListComponents(projectSlug string) ([]Component, error) {
	var response PaginatedResponse[Component]
	_, err := c.resty.R().
		SetResult(&response).
		Get(fmt.Sprintf("/projects/%s/components/", projectSlug))
	
	if err != nil {
		return nil, fmt.Errorf("failed to list components for project %s: %w", projectSlug, err)
	}
	
	return response.Results, nil
}

// GetComponent gets a specific component
func (c *Client) GetComponent(projectSlug, componentSlug string) (*Component, error) {
	var component Component
	resp, err := c.resty.R().
		SetResult(&component).
		Get(fmt.Sprintf("/projects/%s/components/%s/", projectSlug, componentSlug))
	
	if err != nil {
		return nil, fmt.Errorf("failed to get component %s/%s: %w", projectSlug, componentSlug, err)
	}
	
	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("component %s/%s not found", projectSlug, componentSlug)
	}
	
	return &component, nil
}

// ListLanguages lists languages for a project
func (c *Client) ListLanguages(projectSlug string) ([]Language, error) {
	var response PaginatedResponse[Language]
	_, err := c.resty.R().
		SetResult(&response).
		Get(fmt.Sprintf("/projects/%s/languages/", projectSlug))
	
	if err != nil {
		return nil, fmt.Errorf("failed to list languages for project %s: %w", projectSlug, err)
	}
	
	return response.Results, nil
}

// ListTranslations lists translations for a component
func (c *Client) ListTranslations(projectSlug, componentSlug string) ([]Translation, error) {
	var response PaginatedResponse[Translation]
	_, err := c.resty.R().
		SetResult(&response).
		Get(fmt.Sprintf("/projects/%s/components/%s/translations/", projectSlug, componentSlug))
	
	if err != nil {
		return nil, fmt.Errorf("failed to list translations for component %s/%s: %w", projectSlug, componentSlug, err)
	}
	
	return response.Results, nil
}

// GetTranslation gets a specific translation
func (c *Client) GetTranslation(projectSlug, componentSlug, languageCode string) (*Translation, error) {
	var translation Translation
	resp, err := c.resty.R().
		SetResult(&translation).
		Get(fmt.Sprintf("/projects/%s/components/%s/translations/%s/", projectSlug, componentSlug, languageCode))
	
	if err != nil {
		return nil, fmt.Errorf("failed to get translation %s/%s/%s: %w", projectSlug, componentSlug, languageCode, err)
	}
	
	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("translation %s/%s/%s not found", projectSlug, componentSlug, languageCode)
	}
	
	return &translation, nil
}

// ListUnits lists translation units with optional filters
func (c *Client) ListUnits(projectSlug, componentSlug, languageCode string, params map[string]string) ([]Unit, error) {
	var response PaginatedResponse[Unit]
	req := c.resty.R().SetResult(&response)
	
	if params != nil {
		req.SetQueryParams(params)
	}
	
	_, err := req.Get(fmt.Sprintf("/projects/%s/components/%s/translations/%s/units/", projectSlug, componentSlug, languageCode))
	
	if err != nil {
		return nil, fmt.Errorf("failed to list units for translation %s/%s/%s: %w", projectSlug, componentSlug, languageCode, err)
	}
	
	return response.Results, nil
}

// UpdateUnit updates a translation unit
func (c *Client) UpdateUnit(unitID int, target []string, state *int) (*Unit, error) {
	body := map[string]interface{}{
		"target": target,
	}
	
	if state != nil {
		body["state"] = *state
	}
	
	var unit Unit
	_, err := c.resty.R().
		SetBody(body).
		SetResult(&unit).
		Patch(fmt.Sprintf("/units/%d/", unitID))
	
	if err != nil {
		return nil, fmt.Errorf("failed to update unit %d: %w", unitID, err)
	}
	
	return &unit, nil
}

// ListChanges lists recent changes
func (c *Client) ListChanges(params map[string]string) ([]Change, error) {
	var response PaginatedResponse[Change]
	req := c.resty.R().SetResult(&response)
	
	if params != nil {
		req.SetQueryParams(params)
	}
	
	_, err := req.Get("/changes/")
	
	if err != nil {
		return nil, fmt.Errorf("failed to list changes: %w", err)
	}
	
	return response.Results, nil
}

// GetProjectStatistics gets statistics for a project
func (c *Client) GetProjectStatistics(projectSlug string) (*Statistics, error) {
	var stats Statistics
	_, err := c.resty.R().
		SetResult(&stats).
		Get(fmt.Sprintf("/projects/%s/statistics/", projectSlug))
	
	if err != nil {
		return nil, fmt.Errorf("failed to get project statistics for %s: %w", projectSlug, err)
	}
	
	return &stats, nil
}

// GetComponentStatistics gets statistics for a component
func (c *Client) GetComponentStatistics(projectSlug, componentSlug string) (*Statistics, error) {
	var stats Statistics
	_, err := c.resty.R().
		SetResult(&stats).
		Get(fmt.Sprintf("/projects/%s/components/%s/statistics/", projectSlug, componentSlug))
	
	if err != nil {
		return nil, fmt.Errorf("failed to get component statistics for %s/%s: %w", projectSlug, componentSlug, err)
	}
	
	return &stats, nil
}

// GetTranslationStatistics gets statistics for a translation
func (c *Client) GetTranslationStatistics(projectSlug, componentSlug, languageCode string) (*Statistics, error) {
	var stats Statistics
	_, err := c.resty.R().
		SetResult(&stats).
		Get(fmt.Sprintf("/projects/%s/components/%s/translations/%s/statistics/", projectSlug, componentSlug, languageCode))
	
	if err != nil {
		return nil, fmt.Errorf("failed to get translation statistics for %s/%s/%s: %w", projectSlug, componentSlug, languageCode, err)
	}
	
	return &stats, nil
}