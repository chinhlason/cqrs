package query

import (
	"cqrs-postgres-elastic-search-debezium/utils"
)

type IQueryService interface {
	Search(query string, index string) ([]Document, error)
}

type QService struct {
	es *ESClient
}

func NewQueryService(es *ESClient) IQueryService {
	return &QService{es: es}
}

func (s *QService) Search(query string, index string) ([]Document, error) {
	res, err := s.es.Search(utils.USER_INDEX, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}
