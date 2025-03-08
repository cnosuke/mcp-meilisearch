package tools

import (
	"encoding/json"

	"github.com/cockroachdb/errors"
	"github.com/meilisearch/meilisearch-go"
	mcp "github.com/metoro-io/mcp-golang"
	"go.uber.org/zap"
)

// AddDocumentsArgs - Arguments for add_documents tool
type AddDocumentsArgs struct {
	IndexUID   string          `json:"index_uid" jsonschema:"description=The UID of the index to add documents to"`
	Documents  json.RawMessage `json:"documents" jsonschema:"description=The documents to add to the index"`
	PrimaryKey string          `json:"primary_key,omitempty" jsonschema:"description=The primary key of the documents"`
}

// RegisterAddDocumentsTool - Register the add_documents tool
func RegisterAddDocumentsTool(server *mcp.Server, client meilisearch.ServiceManager) error {
	zap.S().Debug("registering add_documents tool")
	err := server.RegisterTool("add_documents", "Add documents to a Meilisearch index",
		func(args AddDocumentsArgs) (*mcp.ToolResponse, error) {
			zap.S().Debug("executing add_documents",
				zap.String("index_uid", args.IndexUID),
				zap.String("primary_key", args.PrimaryKey))

			// Get the index
			index := client.Index(args.IndexUID)

			// Parse the documents
			var documents []map[string]interface{}
			if err := json.Unmarshal(args.Documents, &documents); err != nil {
				zap.S().Error("failed to parse documents",
					zap.Error(err))
				return nil, errors.Wrap(err, "failed to parse documents")
			}

			// Add documents
			var task *meilisearch.TaskInfo
			var err error
			if args.PrimaryKey != "" {
				task, err = index.AddDocuments(documents, args.PrimaryKey)
			} else {
				task, err = index.AddDocuments(documents)
			}

			if err != nil {
				zap.S().Error("failed to add documents",
					zap.String("index_uid", args.IndexUID),
					zap.Error(err))
				return nil, errors.Wrap(err, "failed to add documents")
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
		zap.S().Error("failed to register add_documents tool", zap.Error(err))
		return errors.Wrap(err, "failed to register add_documents tool")
	}

	return nil
}
