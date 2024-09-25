package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"strings"
)

type Place struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type Elasticsearch struct {
	Es *elasticsearch.Client
}

func NewElasticsearch() (*Elasticsearch, error) {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	return &Elasticsearch{es}, nil
}

func (e *Elasticsearch) GetPlaces(limit int, offset int) ([]Place, int, error) {
	query := map[string]interface{}{
		"size": limit,
		"from": offset,
	}
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, 0, err
	}

	req := esapi.SearchRequest{
		Index:          []string{"places"},
		Body:           strings.NewReader(string(queryJSON)),
		TrackTotalHits: true,
	}
	res, err := req.Do(context.Background(), e.Es)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = res.Body.Close() }()
	if res.IsError() {
		return nil, 0, fmt.Errorf("Elasticsearch search request failed: %s", res.String())
	}
	var resBody map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, 0, err
	}
	hits := resBody["hits"].(map[string]interface{})["hits"].([]interface{})
	places := make([]Place, 0, len(hits))

	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		placeBytes, err := json.Marshal(source)
		if err != nil {
			continue
		}

		var place Place
		if err := json.Unmarshal(placeBytes, &place); err != nil {
			continue
		}
		places = append(places, place)
	}

	totalHits := int(resBody["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	return places, totalHits, nil
}
