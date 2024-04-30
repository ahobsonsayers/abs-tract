package kindle_test

import (
	"net/http"
	"testing"

	"github.com/ahobsonsayers/abs-tract/kindle"
	"github.com/stretchr/testify/require"
)

func TestNewNoParameters(t *testing.T) {
	client, err := kindle.NewClient(nil, nil)
	require.NoError(t, err)
	require.Equal(t, kindle.DefaultAmazonURL, client.URL().String())
}

func TestNewWithParameters(t *testing.T) {
	httpClient := &http.Client{}
	countryCode := "es"
	client, err := kindle.NewClient(httpClient, &countryCode)
	require.NoError(t, err)
	require.Equal(t, "https://www.amazon.es", client.URL().String())
}
