package elasticsearch

import (
	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/piqba/common/pkg/config"
	"strings"
	"time"
)

func NewElasticDb(elasticUrls ...string) (*elasticsearch.Client, error) {
	address := elasticUrls
	var err error
	if elasticUrls == nil {
		address = strings.Split(
			config.LoadEnvOrFallback(
				"APP_ELASTIC_URL",
				"http://localhost:9092"),
			",",
		)
	}
	retryBackoff := backoff.NewExponentialBackOff()

	cfg := elasticsearch.Config{
		Addresses:    address,
		DisableRetry: false,
		// Retry on 429 TooManyRequests statuses
		RetryOnStatus: []int{502, 503, 504, 429},
		// Configure the backoff function
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
		// Retry up to 10 attempts
		MaxRetries:          10,
		CompressRequestBody: true,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return es, nil
}
