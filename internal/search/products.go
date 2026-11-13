package search

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	opensearchapi "github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

func SearchProducts(
	ctx context.Context,
	client *opensearchapi.Client,
	query string,
) (*opensearchapi.SearchResp, error) {
	body := fmt.Sprintf(`{
		"query": {
			"multi_match": {
				"query": %q,
				"fields": ["name", "description"],
				"type": "bool_prefix"
			}
		},
		"size": 20
	}`, query)

	searchReq := opensearchapi.SearchReq{
		Indices: []string{"products"},
		Body:    strings.NewReader(body),
	}

	return client.Search(ctx, &searchReq)
}

func IndexDoc(
	ctx context.Context,
	client *opensearchapi.Client,
	index string,
	doc io.Reader,
) (*opensearchapi.IndexResp, error) {
	return client.Index(ctx, opensearchapi.IndexReq{
		Index: index,
		Body:  doc,
	})
}

func BulkIndex(
	ctx context.Context,
	client *opensearchapi.Client,
	index string,
	docs []string,
) error {
	var buf strings.Builder
	for _, doc := range docs {
		buf.WriteString(fmt.Sprintf(`{"index":{"_index":%q}}`+"\n", index))
		buf.WriteString(doc + "\n")
	}

	resp, err := client.Bulk(ctx, opensearchapi.BulkReq{
		Body: strings.NewReader(buf.String()),
	})
	if err != nil {
		return fmt.Errorf("bulk request failed: %w", err)
	}

	// CRITICAL: Bulk always returns HTTP 200, even if every item failed.
	if resp.Errors {
		failed := 0
		for _, item := range resp.Items {
			for _, details := range item {
				if details.Status >= 300 {
					failed++
					log.Printf("bulk item failed: status=%d type=%s reason=%s",
						details.Status, details.Error.Type, details.Error.Reason)
				}
			}
		}
		if failed > len(resp.Items)/2 {
			return fmt.Errorf("bulk mostly failed: %d/%d items", failed, len(resp.Items))
		}
	}
	return nil
}
