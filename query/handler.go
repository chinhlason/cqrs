package query

import (
	"encoding/json"
	"net/http"
)

type IQueryHandler interface {
	Search(res http.ResponseWriter, req *http.Request)
}

type QHandler struct {
	svc IQueryService
}

func NewQueryHandler(svc IQueryService) IQueryHandler {
	return &QHandler{svc}
}

func (h *QHandler) Search(res http.ResponseWriter, req *http.Request) {
	query := req.URL.Query().Get("query")
	index := req.URL.Query().Get("index")
	result, err := h.svc.Search(query, index)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resultJson)
}
