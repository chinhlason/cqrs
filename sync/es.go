package sync

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
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

func NewESClient(client *elasticsearch.Client) ESClient {
	return ESClient{
		client: client,
	}
}

func (es *ESClient) CreateIndex(index string) error {
	res, err := esapi.IndicesExistsRequest{
		Index: []string{index},
	}.Do(context.Background(), es.client)
	if err != nil {
		es.client.Indices.Create(index)
	}
	fmt.Println("res", res)
	return nil
}

func (es *ESClient) InsertDocument(index string, data interface{}) error {
	err := es.CreateIndex(index)
	if err != nil {
		return err
	}
	return nil
}
