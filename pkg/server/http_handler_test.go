package server //nolint

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCustomerHandler(t *testing.T) {
	// prepare
	mux := http.NewServeMux()
	mux.Handle("/", NewCustomerHandler(zap.NewNop()))
	ts := httptest.NewServer(mux)

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
