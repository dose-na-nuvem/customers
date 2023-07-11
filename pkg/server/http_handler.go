package server

import (
	"encoding/json"
	"net/http"

	"github.com/dose-na-nuvem/customers/pkg/model"
	"github.com/dose-na-nuvem/customers/pkg/telemetry"
	"go.uber.org/zap"
)

var _ http.Handler = (*CustomerHandler)(nil)

type CustomerStore interface {
	CreateCustomer(string) (*model.Customer, error)
	ListCustomers() ([]model.Customer, error)
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
	case http.MethodGet:
		h.listCustomers(w, r)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (h *CustomerHandler) createCustomer(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.logger.Warn("erro ao varrer dados do post")
		w.WriteHeader(http.StatusInternalServerError)
	}
	name := r.PostForm.Get("name")

	_, span := telemetry.GetTracer().Start(r.Context(), "create-customer")
	defer span.End()
	_, err = h.store.CreateCustomer(name)
	if err != nil {
		h.logger.Warn("falha ao criar um customer", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *CustomerHandler) listCustomers(w http.ResponseWriter, r *http.Request) {
	_, span := telemetry.GetTracer().Start(r.Context(), "list-customers")
	defer span.End()

	c, err := h.store.ListCustomers()
	if err != nil {
		h.logger.Warn("Falha ao consultar customers", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	b, err := json.Marshal(c)
	if err != nil {
		h.logger.Warn("Falha ao serializar customers", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		h.logger.Error("Falha ao escrever resposta da requisição")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
