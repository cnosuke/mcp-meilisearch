package tools

import (
	"encoding/json"

	"github.com/cockroachdb/errors"
	"github.com/meilisearch/meilisearch-go"
	mcp "github.com/metoro-io/mcp-golang"
	"go.uber.org/zap"
)

// SearchArgs - Arguments for search tool
type SearchArgs struct {
	IndexUID string   `json:"index_uid" jsonschema:"description=The UID of the index to search"`
	Query    string   `json:"query" jsonschema:"description=The search query string"`
	Limit    int64    `json:"limit,omitempty" jsonschema:"description=The maximum number of search results to return"`
	Offset   int64    `json:"offset,omitempty" jsonschema:"description=The number of search results to skip"`
	Filter   string   `json:"filter,omitempty" jsonschema:"description=The filter query string"`
	Sort     []string `json:"sort,omitempty" jsonschema:"description=The list of attributes to sort by"`
}

// RegisterSearchTool - Register the search tool
func RegisterSearchTool(server *mcp.Server, client meilisearch.ServiceManager) error {
	zap.S().Debug("registering search tool")
	err := server.RegisterTool("search", "Search for documents in a Meilisearch index",
		func(args SearchArgs) (*mcp.ToolResponse, error) {
			zap.S().Debug("executing search",
				zap.String("index_uid", args.IndexUID),
				zap.String("query", args.Query))

			// Get the index
			index := client.Index(args.IndexUID)

			// Create search request
			searchRequest := &meilisearch.SearchRequest{
				Query:  args.Query,
				Limit:  args.Limit,
				Offset: args.Offset,
			}

			// Add filter if provided
			if args.Filter != "" {
				searchRequest.Filter = args.Filter
			}

			// Add sort if provided
			if len(args.Sort) > 0 {
				searchRequest.Sort = args.Sort
			}

			// Perform search
			searchResult, err := index.Search(args.Query, searchRequest)
			if err != nil {
				zap.S().Error("failed to search index",
					zap.String("index_uid", args.IndexUID),
					zap.String("query", args.Query),
					zap.Error(err))
				return nil, errors.Wrap(err, "failed to search index")
			}

			// Convert search results to JSON
			jsonResult, err := json.Marshal(searchResult)
			if err != nil {
				zap.S().Error("failed to convert search results to JSON", zap.Error(err))
				return nil, errors.Wrap(err, "failed to convert search results to JSON")
			}

			return mcp.NewToolResponse(mcp.NewTextContent(string(jsonResult))), nil
		})

	if err != nil {
		zap.S().Error("failed to register search tool", zap.Error(err))
		return errors.Wrap(err, "failed to register search tool")
	}

	return nil
}
