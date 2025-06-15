package bookbeat

type SearchLinkSelf struct {
	Href string `json:"href"`
}

type SearchLink struct {
	Self SearchLinkSelf `json:"self"`
}

type SearchBook struct {
	ID    int        `json:"id"`
	Links SearchLink `json:"_links"`
}

type SearchEmbedded struct {
	Books []SearchBook `json:"books"`
}

type SearchResponse struct {
	QueryUrl string         `json:"queryurl"`
	Count    int            `json:"count"`
	Embedded SearchEmbedded `json:"_embedded"`
	IsCapped bool           `json:"iscapped"`
}
