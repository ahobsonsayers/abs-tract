package kindle

import (
	"strings"
	"time"

	"github.com/ahobsonsayers/abs-goodreads/utils"
	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
)

const publishDateLayout = "Jan 2, 2006"

var (
	bookCoverSetExpr  = xpath.MustCompile(`.//img/@srcset`)
	bookFormatExpr    = xpath.MustCompile(`.//a[contains(text(), "Kindle Edition")]//text()`)
	bookInfoExpr      = xpath.MustCompile(`.//div[contains(@class, "a-color-secondary")]`)
	bookTitleExpr     = xpath.MustCompile(`.//h2`)
	searchResultsExpr = xpath.MustCompile(`//div[contains(@class, "s-result-list")]//div[@data-index and @data-asin]`)
)

type Book struct {
	ASIN        string
	Title       string
	Author      string
	Cover       string
	PublishDate *time.Time
}

// BooksFromHTML parses and returns the books from the html of a search results page
func BooksFromHTML(searchNode *html.Node) ([]Book, error) {
	resultNodes := htmlquery.QuerySelectorAll(searchNode, searchResultsExpr)

	books := make([]Book, 0, len(resultNodes))
	for _, resultNode := range resultNodes {
		if !isKindleBook(resultNode) {
			continue
		}

		book := BookFromHTML(resultNode)
		if book != nil {
			books = append(books, *book)
		}
	}

	return books, nil
}

// BookFromHTML parses and returns a book from the html
// of a book result on the search results page
func BookFromHTML(bookNode *html.Node) *Book {
	asin := bookAsin(bookNode)
	if asin == "" {
		return nil
	}

	title := bookTitle(bookNode)
	if title == "" {
		return nil
	}

	cover := bookCover(bookNode)
	author, publishDate := bookInfo(bookNode)

	return &Book{
		ASIN:        asin,
		Title:       title,
		Author:      author,
		Cover:       cover,
		PublishDate: publishDate,
	}
}

func isKindleBook(bookNode *html.Node) bool {
	bookFormatNode := htmlquery.QuerySelector(bookNode, bookFormatExpr)
	if bookFormatNode == nil {
		return false
	}
	bookFormatNodeValue := htmlquery.InnerText(bookFormatNode)

	bookFormat := strings.ToLower(bookFormatNodeValue)
	return strings.Contains(bookFormat, "kindle")
}

// bookAsin gets the book asim.
func bookAsin(bookNode *html.Node) string {
	return htmlquery.SelectAttr(bookNode, "data-asin")
}

// bookTitle gets the book title.
func bookTitle(bookNode *html.Node) string {
	titleNode := htmlquery.QuerySelector(bookNode, bookTitleExpr)
	titleNodeValue := htmlquery.InnerText(titleNode)
	return strings.TrimSpace(titleNodeValue)
}

// bookCover gets the book cover.
func bookCover(bookNode *html.Node) string {
	coverSetAttr := htmlquery.QuerySelector(bookNode, bookCoverSetExpr)
	coverSetAttrValue := htmlquery.InnerText(coverSetAttr)
	return parseBookCoversAttrValue(coverSetAttrValue)
}

// bookInfo gets additional book info.
// Return author and publish date (if found)
func bookInfo(bookNode *html.Node) (string, *time.Time) {
	infoNode := htmlquery.QuerySelector(bookNode, bookInfoExpr)
	bookInfoNodeValue := htmlquery.InnerText(infoNode)
	return parseBookInfoNodeValue(bookInfoNodeValue)
}

// parseBookCoversAttrValue parses the value of the book cover set attribute.
// Returns the url of the original/full-size book cover.
// See test for expected value format.
func parseBookCoversAttrValue(coverSetAttrValue string) string {
	// Get first cover url from the cover set.
	// This will not be original full size cover
	modifiedCoverUrl := strings.Fields(coverSetAttrValue)[0]
	originalCoverUrl := utils.SanitiseImageURL(modifiedCoverUrl)
	return originalCoverUrl
}

// parseBookInfoNodeValue parses the value of the book info node
// Returns author, publisher and date published.
// See test for expected value format.
func parseBookInfoNodeValue(bookInfoNodeValue string) (string, *time.Time) {
	// Book info parts are separated by | the of which is the author
	bookInfoParts := strings.Split(bookInfoNodeValue, "|")
	bookAuthorPart := strings.TrimSpace(bookInfoParts[0])
	publishDatePart := strings.TrimSpace(bookInfoParts[2])

	// Strip out the "by" from the author part
	bookAuthorFields := strings.Fields(bookAuthorPart)
	if len(bookAuthorFields) > 1 && strings.EqualFold(bookAuthorFields[0], "by") {
		bookAuthorFields = bookAuthorFields[1:]
	}

	// Rejoin author fields
	bookAuthor := strings.Join(bookAuthorFields, " ")

	// Parse the date string according to the defined layout
	var publishDate *time.Time
	parsedPublishDate, err := time.Parse(publishDateLayout, publishDatePart)
	if err == nil {
		publishDate = &parsedPublishDate
	}

	return bookAuthor, publishDate
}
