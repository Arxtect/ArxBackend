package initializers

import (
	"github.com/Arxtect/ArxBackend/golangp/config"

	"github.com/meilisearch/meilisearch-go"
)

var MeiliClient *meilisearch.Client

func InitMeiliClient(config *config.Config) {
	MeiliClient = meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   config.MeiliHost,
		APIKey: config.MeiliKey,
	})
}
