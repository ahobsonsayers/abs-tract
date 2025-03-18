package kindle_test

import (
	"testing"

	"github.com/ahobsonsayers/abs-tract/kindle"
	"github.com/stretchr/testify/require"
)

func TestNewNoParameters(t *testing.T) {
	client, err := kindle.NewClient(nil)
	require.NoError(t, err)
	require.Equal(t, kindle.DefaultAmazonURL, client.URL().String())
}

func TestNewWithParameters(t *testing.T) {
	countryCode := "es"
	client, err := kindle.NewClient(&countryCode)
	require.NoError(t, err)
	require.Equal(t, "https://www.amazon.es", client.URL().String())
}
