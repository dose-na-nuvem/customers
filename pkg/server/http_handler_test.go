package server //nolint

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestCustomerHandler(t *testing.T) {
	// prepare
	mux := http.NewServeMux()
	mux.Handle("/", NewCustomerHandler(zap.NewNop()))
	ts := httptest.NewTLSServer(mux)

	defer ts.Close()

	client := ts.Client()

	// test
	res, err := client.Get(ts.URL) //nolint
	require.NoError(t, err)

	msg, err := io.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)

	// verify
	assert.Equal(t, "Recebemos um chamado!\n", string(msg))
}

func TestCustomerHandlerError(t *testing.T) {
	// prepare
	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)
	h := NewCustomerHandler(logger)
	writer := &ResponseWriterMock{}

	// test
	h.ServeHTTP(writer, nil)

	// verify
	assert.Len(t, logs.All(), 1)
	assert.Contains(t, "erro ao escrever mensagem ao cliente", logs.All()[0].Message)
}

type ResponseWriterMock struct {
	http.ResponseWriter
}

func (r *ResponseWriterMock) Write([]byte) (int, error) {
	return 0, errors.New("boo")
}
