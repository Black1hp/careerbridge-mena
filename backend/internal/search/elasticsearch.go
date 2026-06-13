package search

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/olivere/elastic/v7"
)

var Client *elastic.Client

type ESOpportunity struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Country     string `json:"country"`
	Deadline    string `json:"deadline,omitempty"`
	URL         string `json:"url"`
	Source      string `json:"source"`
	Description string `json:"description"`
	Eligibility string `json:"eligibility"`
	Funding     string `json:"funding"`
}

const indexName = "opportunities"

func Connect() error {
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://localhost:9200"
	}

	var err error
	Client, err = elastic.NewClient(
		elastic.SetURL(esURL),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(true),
	)
	if err != nil {
		return fmt.Errorf("elasticsearch connect: %w", err)
	}

	ctx := context.Background()
	exists, err := Client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return fmt.Errorf("index check: %w", err)
	}

	if !exists {
		mapping := `{
			"settings": {
				"number_of_shards": 1,
				"number_of_replicas": 0,
				"analysis": {
					"analyzer": {
						"arabic_english": {
							"type": "custom",
							"tokenizer": "standard",
							"filter": ["lowercase", "asciifolding"]
						}
					}
				}
			},
			"mappings": {
				"properties": {
					"title": { "type": "text", "analyzer": "arabic_english" },
					"type": { "type": "keyword" },
					"country": { "type": "keyword" },
					"deadline": { "type": "date", "format": "epoch_millis||yyyy-MM-dd'T'HH:mm:ssZ||strict_date_optional_time" },
					"url": { "type": "keyword", "index": false },
					"source": { "type": "keyword" },
					"description": { "type": "text", "analyzer": "arabic_english" },
					"eligibility": { "type": "text", "analyzer": "arabic_english" },
					"funding": { "type": "text", "analyzer": "arabic_english" }
				}
			}
		}`

		_, err = Client.CreateIndex(indexName).Body(mapping).Do(ctx)
		if err != nil {
			return fmt.Errorf("create index: %w", err)
		}
	}

	return nil
}

func IndexOpportunity(ctx context.Context, opp ESOpportunity) error {
	_, err := Client.Index().
		Index(indexName).
		Id(fmt.Sprintf("%d", opp.ID)).
		BodyJson(opp).
		Do(ctx)
	return err
}

func BulkIndex(ctx context.Context, opps []ESOpportunity) error {
	bulk := Client.Bulk().Index(indexName)
	for _, opp := range opps {
		bulk.Add(elastic.NewBulkIndexRequest().
			Id(fmt.Sprintf("%d", opp.ID)).
			Doc(opp))
	}
	_, err := bulk.Do(ctx)
	return err
}

type SearchRequest struct {
	Query   string
	Type    string
	Country string
	Page    int
	Limit   int
}

type SearchResult struct {
	Total      int
	Opportunities []ESOpportunity
}

func Search(ctx context.Context, req SearchRequest) (*SearchResult, error) {
	q := elastic.NewBoolQuery()

	if req.Query != "" {
		q = q.Must(
			elastic.NewMultiMatchQuery(req.Query, "title", "description", "eligibility").
				Fuzziness("AUTO").
				Operator("and"),
		)
	} else {
		q = q.Must(elastic.NewMatchAllQuery())
	}

	if req.Type != "" {
		q = q.Filter(elastic.NewTermQuery("type", req.Type))
	}
	if req.Country != "" {
		q = q.Filter(elastic.NewTermQuery("country", req.Country))
	}

	from := (req.Page - 1) * req.Limit
	if from < 0 {
		from = 0
	}

	searchResult, err := Client.Search().
		Index(indexName).
		Query(q).
		Sort("deadline", true).
		Sort("_score", false).
		From(from).
		Size(req.Limit).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	var opps []ESOpportunity
	for _, hit := range searchResult.Hits.Hits {
		var opp ESOpportunity
		if err := json.Unmarshal(hit.Source, &opp); err != nil {
			continue
		}
		opps = append(opps, opp)
	}

	return &SearchResult{
		Total:         int(searchResult.Hits.TotalHits.Value),
		Opportunities: opps,
	}, nil
}
