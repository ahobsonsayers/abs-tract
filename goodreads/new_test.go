package goodreads_test

import (
	"net/http"
	"testing"

	"github.com/ahobsonsayers/abs-tract/goodreads"
	"github.com/stretchr/testify/require"
)

func TestNewNoParameters(t *testing.T) {
	client, err := goodreads.NewClient(nil, nil, nil)
	require.NoError(t, err)
	require.Equal(t, goodreads.DefaultGoodreadsUrl, client.URL().String())
}

func TestNewWithParameters(t *testing.T) {
	httpClient := &http.Client{}
	goodreadUrl := "http://example.com"
	apiKey := "test"
	client, err := goodreads.NewClient(httpClient, &goodreadUrl, &apiKey)
	require.NoError(t, err)
	require.Equal(t, goodreadUrl, client.URL().String())
}
