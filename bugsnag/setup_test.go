package bugsnag

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type testVars struct {
	ctx    context.Context
	mux    *http.ServeMux
	server *httptest.Server
	client Client
}

func setup(t *testing.T) *testVars {
	t.Helper()

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client, err := NewClient(WithBaseURL(server.URL))

	require.NoError(t, err)

	t.Cleanup(func() {
		server.Close()
	})

	return &testVars{
		ctx:    context.Background(),
		mux:    mux,
		server: server,
		client: *client,
	}
}
