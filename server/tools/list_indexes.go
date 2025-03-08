package tools

import (
	"encoding/json"

	"github.com/cockroachdb/errors"
	"github.com/meilisearch/meilisearch-go"
	mcp "github.com/metoro-io/mcp-golang"
	"go.uber.org/zap"
)

// ListIndexesArgs - Arguments for list_indexes tool
type ListIndexesArgs struct {
	// No arguments needed
}

// RegisterListIndexesTool - Register the list_indexes tool
func RegisterListIndexesTool(server *mcp.Server, client meilisearch.ServiceManager) error {
	zap.S().Debug("registering list_indexes tool")
	err := server.RegisterTool("list_indexes", "List all Meilisearch indexes",
		func(args ListIndexesArgs) (*mcp.ToolResponse, error) {
			zap.S().Debug("executing list_indexes")

			// Get all indexes
			indexes, err := client.ListIndexes(&meilisearch.IndexesQuery{
				Limit: 100, // Set a reasonable limit
			})
			if err != nil {
				zap.S().Error("failed to get indexes", zap.Error(err))
				return nil, errors.Wrap(err, "failed to get indexes")
			}

			// Convert indexes to JSON
			jsonResult, err := json.Marshal(indexes.Results)
			if err != nil {
				zap.S().Error("failed to convert indexes to JSON", zap.Error(err))
				return nil, errors.Wrap(err, "failed to convert indexes to JSON")
			}

			return mcp.NewToolResponse(mcp.NewTextContent(string(jsonResult))), nil
		})

	if err != nil {
		zap.S().Error("failed to register list_indexes tool", zap.Error(err))
		return errors.Wrap(err, "failed to register list_indexes tool")
	}

	return nil
}
