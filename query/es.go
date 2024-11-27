package query

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"strings"
)

func GetESClient(host string) (*elasticsearch.Client, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{host},
		Transport: nil,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

type ESClient struct {
	client *elasticsearch.Client
}

func NewESClient(client *elasticsearch.Client) *ESClient {
	return &ESClient{
		client: client,
	}
}

func (es *ESClient) CreateIndex(index string) error {
	res, err := es.client.Indices.Get([]string{index})
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		res, err := es.client.Indices.Create(index)
		if err != nil {
			return err
		}
		defer res.Body.Close()
	}
	return nil
}

func (es *ESClient) InsertDocument(index, id string, data interface{}) error {
	err := es.CreateIndex(index)
	if err != nil {
		return err
	}
	dataJson, err := json.Marshal(data)
	js := string(dataJson)

	ind, err := es.client.Index(
		index,
		strings.NewReader(js),
		es.client.Index.WithDocumentID(id),
		es.client.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	fmt.Println(ind)
	return nil
}

func (es *ESClient) UpdateDocument(id, index string, data interface{}) error {
	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}
	updateData := fmt.Sprintf(`{"doc": %s, "doc_as_upsert": true}`, dataJson)
	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader([]byte(updateData)),
	}

	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return nil
	}

	defer res.Body.Close()

	if res.IsError() {
		return errors.New("error updating document")
	}

	return nil
}

func (es *ESClient) DeleteDocument(id, index string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
	}

	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		return errors.New("error deleting document")
	}

	return nil
}

func (es *ESClient) Search(index, query string) ([]Document, error) {
	q := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":     query,
				"fields":    []string{"name", "email", "orders.product"},
				"fuzziness": "AUTO",
			},
		},
	}

	qJson, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}

	req := esapi.SearchRequest{
		Index: []string{index},
		Body:  strings.NewReader(string(qJson)),
	}

	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	var documents []Document

	if hits, ok := r["hits"].(map[string]interface{}); ok {
		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				if hitData, ok := hit.(map[string]interface{}); ok {
					doc := Document{
						Id:     hitData["_id"].(string),
						Type:   hitData["_type"].(string),
						Index:  hitData["_index"].(string),
						Source: hitData["_source"],
					}
					documents = append(documents, doc)
				}
			}
		}
	}

	return documents, nil
}

func (es *ESClient) GetDocument(id, index string) (interface{}, error) {
	res, err := es.client.Get(index, id)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	return r, nil
}
