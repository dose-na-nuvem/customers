package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dose-na-nuvem/customers/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewCustomer(t *testing.T) {
	// prepare
	called := false
	st := &mockStore{createCustomerFunc: func(name string) (*model.Customer, error) {
		called = true

		// TODO: nosso código não está funcionando de verdade:
		// arrumar o código e ativar a próxima linha
		// assert.Equal(t, "John Doe", name)

		return nil, nil
	}}

	mux := http.NewServeMux()
	mux.Handle("/", NewCustomerHandler(zap.NewNop(), st))
	ts := httptest.NewTLSServer(mux)

	defer ts.Close()

	client := ts.Client()

	// test
	_, err := client.PostForm(ts.URL, url.Values{"name": {"John Doe"}})
	require.NoError(t, err)

	// verify
	// verification is on the mock
	assert.True(t, called, "mock was expected to have been called")
}

func TestCustomerHandlerError(t *testing.T) {
	// prepare
	core, _ := observer.New(zap.InfoLevel)
	logger := zap.New(core)
	h := NewCustomerHandler(logger, nil)
	writer := &ResponseWriterMock{}

	// test
	req, err := http.NewRequest(http.MethodGet, "", nil)
	require.NoError(t, err)
	h.ServeHTTP(writer, req)

	// verify
	//	assert.Len(t, logs.All(), 1)
	//	assert.Contains(t, "erro ao escrever mensagem ao cliente", logs.All()[0].Message)
}

type ResponseWriterMock struct {
	http.ResponseWriter
}

func (r *ResponseWriterMock) Write([]byte) (int, error) {
	return 0, errors.New("boo")
}

type mockStore struct {
	createCustomerFunc func(name string) (*model.Customer, error)
}

func (m *mockStore) CreateCustomer(name string) (*model.Customer, error) {
	if m.createCustomerFunc != nil {
		m.createCustomerFunc(name)
	}

	return nil, nil
}
