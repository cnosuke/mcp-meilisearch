package tools

import (
	"encoding/json"

	"github.com/cockroachdb/errors"
	"github.com/meilisearch/meilisearch-go"
	mcp "github.com/metoro-io/mcp-golang"
	"go.uber.org/zap"
)

// HealthCheckArgs - Arguments for health_check tool
type HealthCheckArgs struct {
	// No arguments needed
}

// RegisterHealthCheckTool - Register the health_check tool
func RegisterHealthCheckTool(server *mcp.Server, client meilisearch.ServiceManager) error {
	zap.S().Debug("registering health_check tool")
	err := server.RegisterTool("health_check", "Check Meilisearch server health status",
		func(args HealthCheckArgs) (*mcp.ToolResponse, error) {
			zap.S().Debug("executing health_check")

			// Get health information
			health, err := client.Health()
			if err != nil {
				zap.S().Error("failed to get health status", zap.Error(err))
				return nil, errors.Wrap(err, "failed to get health status")
			}

			// Convert health to JSON
			jsonResult, err := json.Marshal(health)
			if err != nil {
				zap.S().Error("failed to convert health status to JSON", zap.Error(err))
				return nil, errors.Wrap(err, "failed to convert health status to JSON")
			}

			return mcp.NewToolResponse(mcp.NewTextContent(string(jsonResult))), nil
		})

	if err != nil {
		zap.S().Error("failed to register health_check tool", zap.Error(err))
		return errors.Wrap(err, "failed to register health_check tool")
	}

	return nil
}
