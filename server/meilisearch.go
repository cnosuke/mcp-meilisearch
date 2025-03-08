package server

import (
	"github.com/cnosuke/mcp-meilisearch/config"
	"github.com/cockroachdb/errors"
	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
)

// MeilisearchServer - Meilisearch server structure
type MeilisearchServer struct {
	Client meilisearch.ServiceManager
	cfg    *config.Config
}

// NewMeilisearchServer - Create a new Meilisearch server
func NewMeilisearchServer(cfg *config.Config) (*MeilisearchServer, error) {
	zap.S().Info("creating new Meilisearch server",
		zap.String("host", cfg.Meilisearch.Host))

	// Create a new Meilisearch client
	var client meilisearch.ServiceManager
	if cfg.Meilisearch.APIKey != "" {
		client = meilisearch.New(cfg.Meilisearch.Host, meilisearch.WithAPIKey(cfg.Meilisearch.APIKey))
		zap.S().Info("Meilisearch client created with API key")
	} else {
		client = meilisearch.New(cfg.Meilisearch.Host)
		zap.S().Info("Meilisearch client created without API key")
	}

	// Test connection
	zap.S().Debug("testing Meilisearch connection")
	_, err := client.Health()
	if err != nil {
		zap.S().Error("failed to connect to Meilisearch",
			zap.String("host", cfg.Meilisearch.Host),
			zap.Error(err))
		return nil, errors.Wrap(err, "failed to connect to Meilisearch")
	}
	zap.S().Info("successfully connected to Meilisearch")

	return &MeilisearchServer{
		Client: client,
		cfg:    cfg,
	}, nil
}
