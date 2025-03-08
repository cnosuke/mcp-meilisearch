package tools

import (
	"github.com/meilisearch/meilisearch-go"
	mcp "github.com/metoro-io/mcp-golang"
)

// RegisterAllTools - Register all tools with the server
func RegisterAllTools(server *mcp.Server, client meilisearch.ServiceManager) error {
	// Register health check tool
	if err := RegisterHealthCheckTool(server, client); err != nil {
		return err
	}

	// Register index management tools
	if err := RegisterListIndexesTool(server, client); err != nil {
		return err
	}

	if err := RegisterCreateIndexTool(server, client); err != nil {
		return err
	}

	// Register search tools
	if err := RegisterSearchTool(server, client); err != nil {
		return err
	}

	// Register document management tools
	if err := RegisterGetDocumentsTool(server, client); err != nil {
		return err
	}

	if err := RegisterAddDocumentsTool(server, client); err != nil {
		return err
	}

	return nil
}
