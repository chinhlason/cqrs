package command

import (
	"encoding/json"
	"net/http"
)

type IHandler interface {
	InsertUser(res http.ResponseWriter, req *http.Request)
	UpdateUser(res http.ResponseWriter, req *http.Request)

	InsertOrder(res http.ResponseWriter, req *http.Request)
	UpdateOrder(res http.ResponseWriter, req *http.Request)
}

type Handler struct {
	service IService
}

func NewHandler(service IService) IHandler {
	return &Handler{service}
}

func (h *Handler) InsertUser(res http.ResponseWriter, req *http.Request) {
	var user User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.service.InsertUser(user.Name, user.Email)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(`{"message": "success"}`))
}

func (h *Handler) UpdateUser(res http.ResponseWriter, req *http.Request) {
	var user User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.service.UpdateUser(user.Id, user.Name, user.Email)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`{"message": "success"}`))
}

func (h *Handler) InsertOrder(res http.ResponseWriter, req *http.Request) {
	var order Order
	if err := json.NewDecoder(req.Body).Decode(&order); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.service.InsertOrder(order.UserId, order.Product)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(`{"message": "success"}`))
}

func (h *Handler) UpdateOrder(res http.ResponseWriter, req *http.Request) {
	var order Order
	if err := json.NewDecoder(req.Body).Decode(&order); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.service.UpdateOrder(order.Id, order.Product)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`{"message": "success"}`))
}
