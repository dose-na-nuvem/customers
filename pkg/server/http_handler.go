package server

import (
	"net/http"

	"github.com/dose-na-nuvem/customers/pkg/model"
	"go.uber.org/zap"
)

var _ http.Handler = (*CustomerHandler)(nil)

type CustomerStore interface {
	CreateCustomer(string) (*model.Customer, error)
}

type CustomerHandler struct {
	logger *zap.Logger
	store  CustomerStore
}

func NewCustomerHandler(logger *zap.Logger, store CustomerStore) *CustomerHandler {
	return &CustomerHandler{
		logger: logger,
		store:  store,
	}
}

func (h *CustomerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createCustomer(w, r)
	default:
		// TODO: isso esta dando nil pointer em um teste
		// w.WriteHeader(http.StatusNotImplemented)
	}
}

func (h *CustomerHandler) createCustomer(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.logger.Warn("erro ao varrer dados do post")
		w.WriteHeader(http.StatusInternalServerError)
	}
	name := r.PostForm.Get("name")
	_, err = h.store.CreateCustomer(name)
	if err != nil {
		h.logger.Warn("falha ao criar um customer", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
