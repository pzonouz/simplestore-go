package search

import (
	"log"
	"os"

	opensearch "github.com/opensearch-project/opensearch-go/v4"
	opensearchapi "github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

func NewClient() *opensearchapi.Client {
	address := os.Getenv("OPENSEARCH_URL")
	if address == "" {
		// Host-facing: localhost because you published port 9200.
		// If the Go app later runs inside the Docker network, use "http://opensearch:9200".
		address = "http://localhost:9200"
	}

	client, err := opensearchapi.NewClient(opensearchapi.Config{
		Client: opensearch.Config{
			Addresses: []string{address},
		},
	})
	if err != nil {
		log.Fatalf("failed to create OpenSearch client: %v", err)
	}
	return client
}
