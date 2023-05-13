package server

import (
	"io"
	"net/http"

	"go.uber.org/zap"
)

var _ http.Handler = (*CustomerHandler)(nil)

type CustomerHandler struct {
	logger *zap.Logger
}

func NewCustomerHandler(logger *zap.Logger) *CustomerHandler {
	return &CustomerHandler{
		logger: logger,
	}
}

func (h *CustomerHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_, err := io.WriteString(w, "Recebemos um chamado!\n")
	if err != nil {
		h.logger.Warn("erro ao escrever mensagem ao cliente", zap.Error(err))
	}
}
