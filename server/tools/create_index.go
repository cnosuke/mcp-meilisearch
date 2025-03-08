package tools

import (
	"encoding/json"

	"github.com/cockroachdb/errors"
	"github.com/meilisearch/meilisearch-go"
	mcp "github.com/metoro-io/mcp-golang"
	"go.uber.org/zap"
)

// CreateIndexArgs - Arguments for create_index tool
type CreateIndexArgs struct {
	UID        string `json:"uid" jsonschema:"description=The unique identifier for the index"`
	PrimaryKey string `json:"primary_key,omitempty" jsonschema:"description=The primary key of the documents in the index"`
}

// RegisterCreateIndexTool - Register the create_index tool
func RegisterCreateIndexTool(server *mcp.Server, client meilisearch.ServiceManager) error {
	zap.S().Debug("registering create_index tool")
	err := server.RegisterTool("create_index", "Create a new Meilisearch index",
		func(args CreateIndexArgs) (*mcp.ToolResponse, error) {
			zap.S().Debug("executing create_index",
				zap.String("uid", args.UID),
				zap.String("primary_key", args.PrimaryKey))

			// Create index
			indexConfig := meilisearch.IndexConfig{
				Uid:        args.UID,
				PrimaryKey: args.PrimaryKey,
			}
			task, err := client.CreateIndex(&indexConfig)
			if err != nil {
				zap.S().Error("failed to create index",
					zap.String("uid", args.UID),
					zap.Error(err))
				return nil, errors.Wrap(err, "failed to create index")
			}

			// Convert task to JSON
			jsonResult, err := json.Marshal(task)
			if err != nil {
				zap.S().Error("failed to convert task to JSON", zap.Error(err))
				return nil, errors.Wrap(err, "failed to convert task to JSON")
			}

			return mcp.NewToolResponse(mcp.NewTextContent(string(jsonResult))), nil
		})

	if err != nil {
		zap.S().Error("failed to register create_index tool", zap.Error(err))
		return errors.Wrap(err, "failed to register create_index tool")
	}

	return nil
}
