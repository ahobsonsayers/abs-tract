package bookbeat

import (
	"time"
)

type RatingDistribution struct {
	Num1 int `json:"1"`
	Num2 int `json:"2"`
	Num3 int `json:"3"`
	Num4 int `json:"4"`
	Num5 int `json:"5"`
}

type Rating struct {
	RatingValue        float64            `json:"ratingValue"`
	NumberOfRatings    int                `json:"numberOfRatings"`
	RatingDistribution RatingDistribution `json:"ratingDistribution"`
}

type Badge struct {
	ID             string `json:"id"`
	TranslationKey string `json:"translationKey"`
	Type           string `json:"type"`
	Icon           string `json:"icon"`
}

type Genre struct {
	Genreid int    `json:"genreid"`
	Name    string `json:"name"`
}

type Series struct {
	Count      int    `json:"count"`
	Partindex  int    `json:"partindex"`
	Prevbookid int    `json:"prevbookid"`
	Nextbookid int    `json:"nextbookid"`
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Partnumber int    `json:"partnumber"`
	URL        string `json:"url"`
}

type CopyrightOwner struct {
	Year int    `json:"year"`
	Name string `json:"name"`
}

type Contributor struct {
	ID          int    `json:"id"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Displayname string `json:"displayname"`
	Role        string `json:"role"`
	Description string `json:"description"`
	Booksurl    string `json:"booksurl"`
}

type Edition struct {
	ID                    int              `json:"id"`
	Isbn                  string           `json:"isbn"`
	Format                string           `json:"format"`
	Language              string           `json:"language"`
	Published             time.Time        `json:"published"`
	BookBeatPublishDate   time.Time        `json:"bookBeatPublishDate"`
	BookBeatUnpublishDate time.Time        `json:"bookBeatUnpublishDate"`
	Availablefrom         time.Time        `json:"availablefrom"`
	Publisher             string           `json:"publisher"`
	CopyrightOwners       []CopyrightOwner `json:"copyrightOwners"`
	Contributors          []Contributor    `json:"contributors"`
	Previewenabled        bool             `json:"previewenabled"`
}

type BookResponse struct {
	ID                  int       `json:"id"`
	Title               string    `json:"title"`
	Subtitle            string    `json:"subtitle"`
	Originaltitle       string    `json:"originaltitle"`
	Author              string    `json:"author"`
	Shareurl            string    `json:"shareurl"`
	Summary             string    `json:"summary"`
	Grade               float64   `json:"grade"`
	Rating              Rating    `json:"rating"`
	NarratingRating     Rating    `json:"narratingRating"`
	Cover               string    `json:"cover"`
	Narrator            string    `json:"narrator"`
	Translator          string    `json:"translator"`
	Language            string    `json:"language"`
	Published           time.Time `json:"published"`
	Originalpublishyear int       `json:"originalpublishyear"`
	Ebooklength         float64   `json:"ebooklength"`
	Ebookduration       float64   `json:"ebookduration"`
	Audiobooklength     int       `json:"audiobooklength"` // in seconds
	Genres              []Genre   `json:"genres"`
	Editions            []Edition `json:"editions"`
	Upcomingeditions    []Edition `json:"upcomingeditions"`
	Markets             []string  `json:"markets"`
	Series              *Series   `json:"series"`
	Contenttypetags     []string  `json:"contenttypetags"`
	Relatedreadingsurl  string    `json:"relatedreadingsurl"`
	Nextcontenturl      string    `json:"nextcontenturl"`
	Type                int       `json:"Type"`
	Bookappviewurl      string    `json:"bookappviewurl"`
	Badges              []Badge   `json:"badges"`
}
