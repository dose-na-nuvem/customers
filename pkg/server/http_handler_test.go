package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

		assert.Equal(t, "John Doe", name)

		return nil, nil
	}}

	mux := http.NewServeMux()
	mux.Handle("/", NewCustomerHandler(zap.NewNop(), st))
	ts := httptest.NewTLSServer(mux)

	defer ts.Close()

	client := ts.Client()

	// test
	ctx := context.Background()
	form := make(url.Values)
	form.Add("name", "John Doe")
	formReader := strings.NewReader(form.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ts.URL, formReader)
	require.NoError(t, err)
	req.Form = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// verify
	// verification is on the mock
	assert.True(t, called, "mock was expected to have been called")
}

func TestCustomerHandlerError(t *testing.T) {
	// prepare
	core, _ := observer.New(zap.InfoLevel)
	logger := zap.New(core)
	h := NewCustomerHandler(logger, nil)
	writer := httptest.NewRecorder()

	// test
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)
	require.NoError(t, err)
	h.ServeHTTP(writer, req)

	// verify
	response := writer.Result()
	defer response.Body.Close()
	assert.Equal(t, response.StatusCode, http.StatusNotImplemented)
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
		_, _ = m.createCustomerFunc(name)
	}

	return nil, nil
}
