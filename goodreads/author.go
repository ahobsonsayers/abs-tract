package goodreads

type authorCommon struct {
	Id   string `xml:"id"`
	Name string `xml:"name"`
}

type AuthorSummary struct {
	authorCommon
}

type AuthorDetails struct {
	authorCommon
	Role             string `xml:"role"`
	ImageURL         string `xml:"image_url"`
	SmallImageURL    string `xml:"small_image_url"`
	Link             string `xml:"link"`
	AverageRating    string `xml:"average_rating"`
	RatingsCount     string `xml:"ratings_count"`
	TextReviewsCount string `xml:"text_reviews_count"`
}
