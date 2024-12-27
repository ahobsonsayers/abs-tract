package kindle

import (
	"strings"
	"time"

	"github.com/ahobsonsayers/abs-tract/utils"
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

// parseBookInfoNodeValue parses the value of the book info node.
// Returns author and date published if found.
// See test for expected value format.
func parseBookInfoNodeValue(bookInfoNodeValue string) (string, *time.Time) {
	// Book info parts are separated by a |
	// Split the parts and sanitise them.
	bookInfoParts := strings.Split(bookInfoNodeValue, "|")
	for idx, part := range bookInfoParts {
		bookInfoParts[idx] = strings.TrimSpace(part)
	}

	// Get author and publish date from info parts.
	var author string
	var publishDate *time.Time
	for _, part := range bookInfoParts {
		partFields := strings.Fields(part)

		// Attempt to get author from part
		// Author is in a part beginning with "by".
		if len(partFields) > 1 && strings.EqualFold(partFields[0], "by") {
			author = strings.Join(partFields[1:], " ")
			continue
		}

		// Attempt to get publish date from part
		// Publish date is in a part that can be parsed as a date
		parsedPublishDate, err := time.Parse(publishDateLayout, part)
		if err == nil {
			publishDate = &parsedPublishDate
			continue
		}
	}

	return author, publishDate
}
