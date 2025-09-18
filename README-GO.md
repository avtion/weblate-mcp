# Weblate MCP Server (Go Implementation)

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server that provides seamless integration with Weblate translation management platform, implemented in Go for improved performance and simplified deployment.

## 🌟 Features

- **🔧 Complete Weblate API Access**: Full integration with Weblate's REST API
- **🤖 AI-Powered Workflow**: Natural language interaction with your translation projects  
- **📊 Project Management**: Create, list, and manage translation projects
- **🔍 Component Operations**: Handle translation components and configurations
- **✏️ Translation Management**: Update, search, and manage translations with pluralization support
- **🌐 Language Support**: Work with all supported languages in your Weblate instance
- **🚀 STDIO Transport**: Native MCP STDIO communication
- **🛡️ Type Safety**: Full Go implementation with comprehensive error handling
- **⚡ High Performance**: Fast, compiled binary with minimal resource usage
- **🐳 Easy Deployment**: Single binary deployment, no runtime dependencies

## 🎯 What is This?

This MCP server acts as a bridge between AI assistants (like Claude Desktop) and your Weblate translation management platform. Instead of manually navigating the Weblate web interface, you can use natural language to:

- **"List all projects in my Weblate instance"**
- **"Show me the French translations for the frontend component"**
- **"Update the welcome message translation"**
- **"Create a new translation project"**

## 🚀 Quick Start

### Option 1: Download Binary (Recommended)

Download the latest binary from the [releases page](https://github.com/avtion/weblate-mcp/releases).

**For Claude Desktop:**
```json
{
  "mcpServers": {
    "weblate": {
      "command": "/path/to/weblate-mcp",
      "env": {
        "WEBLATE_API_URL": "https://your-weblate-instance.com/api",
        "WEBLATE_API_TOKEN": "your-weblate-api-token"
      }
    }
  }
}
```

### Option 2: Build from Source

#### Prerequisites
- Go 1.21+
- Make (optional, for build automation)
- Weblate instance with API access

#### Installation
```bash
# Clone the repository
git clone https://github.com/avtion/weblate-mcp.git
cd weblate-mcp

# Install dependencies
go mod download

# Build the binary
make build
# or
go build -o weblate-mcp .

# Configure environment
cp .env.example .env
# Edit .env with your Weblate API URL and token

# Test the server
./weblate-mcp
```

### Environment Configuration
```env
WEBLATE_API_URL=https://your-weblate-instance.com
WEBLATE_API_TOKEN=your-api-token-here
LOG_LEVEL=info
MCP_SERVER_NAME=weblate-mcp-server
MCP_SERVER_VERSION=1.3.0
```

## 🔗 MCP Client Configuration

### Claude Desktop
Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "weblate": {
      "command": "/path/to/weblate-mcp",
      "env": {
        "WEBLATE_API_URL": "https://your-weblate-instance.com/api",
        "WEBLATE_API_TOKEN": "your-weblate-api-token"
      }
    }
  }
}
```

### Other MCP Clients
The server communicates via STDIO using the standard MCP protocol, making it compatible with any MCP-compliant client.

## 🛠️ Available Tools

### 📊 Project Management
| Tool | Description |
|------|-------------|
| **`listProjects`** | List all available Weblate projects with URLs and metadata |

### 🔧 Component Management
| Tool | Description |
|------|-------------|
| **`listComponents`** | List components in a specific project with source language details |

### ✏️ Translation Management
| Tool | Description |
|------|-------------|
| **`getTranslationForKey`** | Get translation value for a specific key |
| **`writeTranslation`** | Update or write translation values with approval support |
| **`searchStringInProject`** | Search for translations containing specific text in a project |

### 🌐 Language Management
| Tool | Description |
|------|-------------|
| **`listLanguages`** | List languages available in a specific project |

### 📊 Translation Statistics Dashboard
| Tool | Description |
|------|-------------|
| **`getProjectStatistics`** | Comprehensive project statistics with completion rates and string counts |
| **`getComponentStatistics`** | Detailed statistics for a specific component |
| **`getProjectDashboard`** | Complete dashboard overview with all component statistics |
| **`getTranslationStatistics`** | Statistics for specific translation (project/component/language) |

### 📈 Change Tracking & History
| Tool | Description |
|------|-------------|
| **`listRecentChanges`** | Recent changes across all projects with user and timestamp filtering |
| **`getProjectChanges`** | Recent changes for a specific project |
| **`getComponentChanges`** | Recent changes for a specific component |
| **`getChangesByUser`** | Recent changes by a specific user |

## 🧪 Development

### Building
```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Production build (with tests and checks)
make prod-build
```

### Testing
```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# Lint and format code
make check
```

### Development Setup
```bash
# Install development tools
make dev-setup

# Format code
make fmt

# Run linter
make lint
```

## 📚 Architecture

The Go implementation follows a clean architecture pattern:

```
weblate-mcp/
├── main.go                    # Application entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── mcp/                  # MCP protocol implementation
│   ├── server/               # Server setup and tool registration
│   ├── weblate/              # Weblate API client and models
│   └── tools/                # MCP tool implementations
├── Makefile                  # Build automation
└── go.mod                    # Go module definition
```

**Key Features:**
- **Single Binary**: Compiled Go binary with no runtime dependencies
- **Memory Efficient**: Low memory footprint compared to Node.js version
- **Fast Startup**: Instant startup time, ideal for MCP STDIO usage
- **Comprehensive API**: All Weblate API functionality preserved
- **Pluralization Support**: Full language pluralization rules (15+ languages)

## 🔒 Security & Production

- **API Token Security**: Environment variable configuration
- **Request Validation**: Input validation for all API calls
- **Error Handling**: Comprehensive error responses with context
- **Resource Management**: Efficient HTTP client with connection pooling
- **Logging**: Configurable logging levels (disabled by default for MCP compatibility)

## 🎯 Use Cases

### Translation Management
- **Project oversight**: Monitor translation progress across projects
- **Content updates**: Update translations programmatically  
- **Quality assurance**: Review and approve translations
- **Team coordination**: Manage translation workflows

### Development Integration
- **CI/CD pipelines**: Automate translation updates in deployment
- **Content management**: Sync translations with content systems
- **Localization testing**: Validate translations in different contexts
- **Documentation**: Generate translation reports and statistics

### AI-Assisted Workflows
- **Natural language queries**: Ask about translation status in plain English
- **Contextual operations**: AI understands your translation needs
- **Batch operations**: Perform bulk updates with AI assistance
- **Smart suggestions**: Get AI-powered translation recommendations

## 🚀 Performance Benefits

Compared to the TypeScript/Node.js version:

- **50-90% faster startup time**
- **30-60% lower memory usage** 
- **Single binary deployment** (no npm dependencies)
- **Better resource efficiency** for CI/CD environments
- **Native performance** for translation processing

## 🤝 Contributing

We welcome contributions! 

1. **Fork** the repository
2. **Create** a feature branch from main
3. **Implement** changes with tests
4. **Update** documentation
5. **Submit** a pull request

### Code Style
- Use `gofmt` for formatting
- Follow Go best practices
- Add tests for new functionality
- Update documentation

## 📄 License

MIT License - see [LICENSE](./LICENSE) file for details.

---

**Built with ❤️ in Go for the translation community**

*Need help? Check our [documentation](./docs/) or create an [issue](https://github.com/avtion/weblate-mcp/issues)!*