package bookbeat_test

import (
	"testing"

	"github.com/ahobsonsayers/abs-tract/bookbeat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_DefaultFormats(t *testing.T) {
	client, err := bookbeat.NewClient("uk", "all", "all")

	require.NoError(t, err)
	require.NotNil(t, client)

	bookTypes := []string{"audiobook", "ebook"}
	for _, bookType := range bookTypes {
		assert.Contains(t, client.Formats(), bookType)
	}
}

func TestNewClient_DefaultLanguages(t *testing.T) {
	client, err := bookbeat.NewClient("uk", "all", "all")

	require.NoError(t, err)
	require.NotNil(t, client)

	languages := []string{
		"english", "german", "arabic", "basque", "catalan",
		"czech", "danish", "dutch", "estonian", "finnish",
		"french", "hungarian", "italian", "norwegian",
		"norwegiannynorsk", "polish", "portuguese", "russian",
		"spanish", "swedish", "turkish"}
	for _, language := range languages {
		assert.Contains(t, client.Languages(), language)
	}
}

func TestNewClient_AudiobookFormat(t *testing.T) {
	client, err := bookbeat.NewClient("uk", "audiobook", "all")

	require.NoError(t, err)
	require.NotNil(t, client)

	assert.Contains(t, client.Formats(), "audiobook")
}

func TestNewClient_EbookFormat(t *testing.T) {
	client, err := bookbeat.NewClient("uk", "ebook", "all")

	require.NoError(t, err)
	require.NotNil(t, client)

	assert.Contains(t, client.Formats(), "ebook")
}

func TestNewClient_SingleLanguage(t *testing.T) {
	var languageMap = map[string]string{
		"en": "english",
		"de": "german",
		"ar": "arabic",
		"eu": "basque",
		"ca": "catalan",
		"cs": "czech",
		"da": "danish",
		"nl": "dutch",
		"et": "estonian",
		"fi": "finnish",
		"fr": "french",
		"hu": "hungarian",
		"it": "italian",
		"nb": "norwegian",
		"nn": "norwegiannynorsk",
		"pl": "polish",
		"pt": "portuguese",
		"ru": "russian",
		"es": "spanish",
		"sv": "swedish",
		"tr": "turkish",
	}

	for languageCode, languageName := range languageMap {
		t.Run(languageName, func(t *testing.T) {
			client, err := bookbeat.NewClient("uk", "all", languageCode)

			require.NoError(t, err)
			require.NotNil(t, client)

			assert.Len(t, client.Languages(), 1)
			assert.Contains(t, client.Languages(), languageName)
		})
	}
}

func TestNewClient_MultipleLanguages(t *testing.T) {
	client, err := bookbeat.NewClient("uk", "all", "en,fr,hu,ar")

	require.NoError(t, err)
	require.NotNil(t, client)

	languages := []string{"english", "french", "hungarian", "arabic"}
	assert.Len(t, client.Languages(), len(languages))
	for _, language := range languages {
		assert.Contains(t, client.Languages(), language)
	}
}

func TestNewClient_Market(t *testing.T) {
	countryCodes := []string{"gr", "nl", "be", "fr", "es", "hu",
		"it", "ro", "ch", "at", "uk", "dk", "se", "no", "pl",
		"de", "pt", "lu", "ie", "mt", "cy", "fi", "bg", "lt",
		"lv", "ee", "hr", "si", "cz", "sk"}

	for _, countryCode := range countryCodes {
		t.Run(countryCode, func(t *testing.T) {
			client, err := bookbeat.NewClient(countryCode, "all", "en")

			require.NoError(t, err)
			require.NotNil(t, client)

			assert.Equal(t, client.Market(), countryCode)
		})
	}
}
