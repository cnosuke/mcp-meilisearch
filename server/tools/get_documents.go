package tools

import (
	"encoding/json"

	"github.com/cockroachdb/errors"
	"github.com/meilisearch/meilisearch-go"
	mcp "github.com/metoro-io/mcp-golang"
	"go.uber.org/zap"
)

// GetDocumentsArgs - Arguments for get_documents tool
type GetDocumentsArgs struct {
	IndexUID string   `json:"index_uid" jsonschema:"description=The UID of the index to get documents from"`
	Limit    int64    `json:"limit,omitempty" jsonschema:"description=The maximum number of documents to return"`
	Offset   int64    `json:"offset,omitempty" jsonschema:"description=The number of documents to skip"`
	Fields   []string `json:"fields,omitempty" jsonschema:"description=The list of fields to retrieve"`
}

// RegisterGetDocumentsTool - Register the get_documents tool
func RegisterGetDocumentsTool(server *mcp.Server, client meilisearch.ServiceManager) error {
	zap.S().Debug("registering get_documents tool")
	err := server.RegisterTool("get_documents", "Get documents from a Meilisearch index",
		func(args GetDocumentsArgs) (*mcp.ToolResponse, error) {
			zap.S().Debug("executing get_documents",
				zap.String("index_uid", args.IndexUID),
				zap.Int64("limit", args.Limit),
				zap.Int64("offset", args.Offset))

			// Get the index
			index := client.Index(args.IndexUID)

			// Create request parameters
			request := &meilisearch.DocumentsQuery{
				Limit:  args.Limit,
				Offset: args.Offset,
			}

			// Add fields if provided
			if len(args.Fields) > 0 {
				request.Fields = args.Fields
			}

			// Get documents
			var documents meilisearch.DocumentsResult
			err := index.GetDocuments(request, &documents)
			if err != nil {
				zap.S().Error("failed to get documents",
					zap.String("index_uid", args.IndexUID),
					zap.Error(err))
				return nil, errors.Wrap(err, "failed to get documents")
			}

			// Convert documents to JSON
			jsonResult, err := json.Marshal(documents)
			if err != nil {
				zap.S().Error("failed to convert documents to JSON", zap.Error(err))
				return nil, errors.Wrap(err, "failed to convert documents to JSON")
			}

			return mcp.NewToolResponse(mcp.NewTextContent(string(jsonResult))), nil
		})

	if err != nil {
		zap.S().Error("failed to register get_documents tool", zap.Error(err))
		return errors.Wrap(err, "failed to register get_documents tool")
	}

	return nil
}
