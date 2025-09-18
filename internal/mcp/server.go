package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Server represents an MCP server
type Server struct {
	name        string
	version     string
	description string
	tools       map[string]Tool
	reader      *bufio.Scanner
	writer      io.Writer
}

// Tool represents an MCP tool
type Tool struct {
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Handler     func(json.RawMessage) (interface{}, error)
}

// Request represents an MCP request
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

// Response represents an MCP response
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// Error represents an MCP error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ListToolsResult represents the result of listing tools
type ListToolsResult struct {
	Tools []ToolInfo `json:"tools"`
}

// ToolInfo represents tool metadata
type ToolInfo struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// CallToolParams represents parameters for calling a tool
type CallToolParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// CallToolResult represents the result of calling a tool
type CallToolResult struct {
	Content []Content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

// Content represents MCP content
type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// NewServer creates a new MCP server
func NewServer(name, version string) *Server {
	return &Server{
		name:    name,
		version: version,
		tools:   make(map[string]Tool),
		reader:  bufio.NewScanner(os.Stdin),
		writer:  os.Stdout,
	}
}

// SetDescription sets the server description
func (s *Server) SetDescription(description string) {
	s.description = description
}

// RegisterTool registers a tool with the server
func (s *Server) RegisterTool(name string, tool Tool) error {
	s.tools[name] = tool
	return nil
}

// ServeStdio runs the server using STDIO transport
func (s *Server) ServeStdio() error {
	for s.reader.Scan() {
		line := s.reader.Text()
		if line == "" {
			continue
		}

		var req Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			s.sendError(nil, -32700, "Parse error")
			continue
		}

		s.handleRequest(req)
	}

	return s.reader.Err()
}

// handleRequest handles an incoming MCP request
func (s *Server) handleRequest(req Request) {
	switch req.Method {
	case "initialize":
		s.handleInitialize(req)
	case "tools/list":
		s.handleListTools(req)
	case "tools/call":
		s.handleCallTool(req)
	default:
		s.sendError(req.ID, -32601, "Method not found")
	}
}

// handleInitialize handles the initialize request
func (s *Server) handleInitialize(req Request) {
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    s.name,
			"version": s.version,
		},
	}

	s.sendResponse(req.ID, result)
}

// handleListTools handles the tools/list request
func (s *Server) handleListTools(req Request) {
	var tools []ToolInfo
	for name, tool := range s.tools {
		tools = append(tools, ToolInfo{
			Name:        name,
			Description: tool.Description,
			InputSchema: tool.Parameters,
		})
	}

	result := ListToolsResult{Tools: tools}
	s.sendResponse(req.ID, result)
}

// handleCallTool handles the tools/call request
func (s *Server) handleCallTool(req Request) {
	var params CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.sendError(req.ID, -32602, "Invalid params")
		return
	}

	tool, exists := s.tools[params.Name]
	if !exists {
		s.sendError(req.ID, -32601, "Tool not found")
		return
	}

	result, err := tool.Handler(params.Arguments)
	if err != nil {
		s.sendError(req.ID, -32000, err.Error())
		return
	}

	// Convert result to CallToolResult format
	var callResult CallToolResult

	if resultMap, ok := result.(map[string]interface{}); ok {
		if errorMsg, hasError := resultMap["error"]; hasError {
			callResult = CallToolResult{
				Content: []Content{{Type: "text", Text: fmt.Sprintf("Error: %v", errorMsg)}},
				IsError: true,
			}
		} else if content, hasContent := resultMap["content"]; hasContent {
			if contentArray, ok := content.([]map[string]interface{}); ok {
				for _, c := range contentArray {
					if text, ok := c["text"].(string); ok {
						callResult.Content = append(callResult.Content, Content{Type: "text", Text: text})
					}
				}
			}
		}
	}

	// If no content was extracted, create a default response
	if len(callResult.Content) == 0 {
		callResult.Content = []Content{{Type: "text", Text: fmt.Sprintf("%v", result)}}
	}

	s.sendResponse(req.ID, callResult)
}

// sendResponse sends a successful response
func (s *Server) sendResponse(id interface{}, result interface{}) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}

	data, _ := json.Marshal(resp)
	fmt.Fprintln(s.writer, string(data))
}

// sendError sends an error response
func (s *Server) sendError(id interface{}, code int, message string) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Error: &Error{
			Code:    code,
			Message: message,
		},
	}

	data, _ := json.Marshal(resp)
	fmt.Fprintln(s.writer, string(data))
}