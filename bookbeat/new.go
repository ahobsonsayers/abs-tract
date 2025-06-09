package bookbeat

import (
	"strings"

	"github.com/imroc/req/v3"
)

// NewClient creates a new BookBeat client with specified market, format and language filters
// market is the country code for the BookBeat market (e.g., "uk", "de", "ee")
// formatStr can be "audiobook", "ebook", or "all" for both formats
// languageCodes is a comma-separated list of ISO language codes (e.g., "en,de,fr")
func NewClient(market, formatStr, languagesCodes string) (*Bookbeat, error) {
	return NewClientWithURLs(market, formatStr, languagesCodes, SearchBaseURL)
}

// NewClientWithURLs creates a new BookBeat client with custom base URLs (useful for testing)
// market is the country code for the BookBeat market (e.g., "uk", "de", "ee")
// formatStr can be "audiobook", "ebook", or "all" for both formats
// languageCodes is a comma-separated list of ISO language codes (e.g., "en,de,fr")
// searchBaseURL is the custom base URL for the search API
func NewClientWithURLs(market, formatStr, languagesCodes, searchBaseURL string) (*Bookbeat, error) {
	// Set format filter - default to both audiobook and ebook
	formats := []string{"audiobook", "ebook"}
	if formatStr != "all" {
		formats = []string{formatStr}
	}

	// Set default language codes
	if languagesCodes == "all" {
		languagesCodes = "ar,ca,cs,da,de,en,es,et,eu,fi,fr,hu,it,nb,nl,nn,pl,pt,ru,sv,tr"
	}
	// Parse language codes from comma-separated string
	codes := strings.Split(languagesCodes, ",")
	// Convert language codes to BookBeat language identifiers
	languages := make([]string, 0, len(codes))
	for _, code := range codes {
		if code = strings.TrimSpace(code); code != "" {
			// Only add if the language code is supported
			if language, exists := languageMap[code]; exists {
				languages = append(languages, language)
			}
		}
	}

	// Create and configure the BookBeat client
	return &Bookbeat{
		httpClient:    req.C().ImpersonateFirefox(),
		searchBaseURL: searchBaseURL,
		market:        market,
		formats:       formats,
		languages:     languages,
	}, nil
}
