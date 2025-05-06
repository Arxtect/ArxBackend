package initializers

import (
	"github.com/Arxtect/ArxBackend/golangp/config"

	"github.com/meilisearch/meilisearch-go"
	"github.com/toheart/functrace"
)

var MeiliClient *meilisearch.Client

func InitMeiliClient(config *config.Config) {
	defer functrace.Trace([]interface {
	}{config})()
	MeiliClient = meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   config.MeiliHost,
		APIKey: config.MeiliKey,
	})
}
