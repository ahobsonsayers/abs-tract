package kindle

import (
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
)

var (
	bookCoverSetExpr  = xpath.MustCompile(`.//img/@srcset`)
	bookFormatExpr    = xpath.MustCompile(`.//a[contains(text(), "Kindle Edition")]//text()`)
	bookInfoExpr      = xpath.MustCompile(`.//div[contains(@class, "a-color-secondary")]`)
	bookTitleExpr     = xpath.MustCompile(`.//h2`)
	searchResultsExpr = xpath.MustCompile(`//div[contains(@class, "s-result-list")]//div[@data-index and @data-asin]`)
)

type Book struct {
	ASIN   string
	Title  string
	Author string
	Cover  string
}

func BookFromSearchResultHTML(resultNode *html.Node) *Book {
	if !isKindleBook(resultNode) {
		return nil
	}

	asin := bookAsin(resultNode)
	if asin == "" {
		return nil
	}

	title := bookTitle(resultNode)
	if title == "" {
		return nil
	}

	cover := bookCover(resultNode)
	author := bookAuthor(resultNode)

	return &Book{
		ASIN:   asin,
		Title:  title,
		Author: author,
		Cover:  cover,
	}
}

func BooksFromHTML(searchNode *html.Node) ([]Book, error) {
	resultNodes := htmlquery.QuerySelectorAll(searchNode, searchResultsExpr)

	books := make([]Book, 0, len(resultNodes))
	for _, resultNode := range resultNodes {
		book := BookFromSearchResultHTML(resultNode)
		if book != nil {
			books = append(books, *book)
		}
	}

	return books, nil
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

func bookAsin(bookNode *html.Node) string {
	return htmlquery.SelectAttr(bookNode, "data-asin")
}

func bookTitle(bookNode *html.Node) string {
	titleNode := htmlquery.QuerySelector(bookNode, bookTitleExpr)
	titleNodeValue := htmlquery.InnerText(titleNode)
	return strings.TrimSpace(titleNodeValue)
}

func bookCover(bookNode *html.Node) string {
	coverSetAttr := htmlquery.QuerySelector(bookNode, bookCoverSetExpr)
	coverSetAttrValue := htmlquery.InnerText(coverSetAttr)
	return bookCoverFromCoverSetAttrValue(coverSetAttrValue)
}

func bookAuthor(bookNode *html.Node) string {
	infoNode := htmlquery.QuerySelector(bookNode, bookInfoExpr)
	bookInfoNodeValue := htmlquery.InnerText(infoNode)
	return bookAuthorFromInfoNodeValue(bookInfoNodeValue)
}

func bookCoverFromCoverSetAttrValue(coverSetAttrValue string) string {
	// Covers are separated by , and contain a zoom suffix e.g. 2x
	coverUrlsWithZoom := strings.Split(coverSetAttrValue, ",")
	if len(coverUrlsWithZoom) == 0 {
		return ""
	}

	// Get cover urls without zoom
	coverUrls := make([]string, 0, len(coverUrlsWithZoom))
	for _, coverUrlWithZoom := range coverUrlsWithZoom {
		coverUrl := strings.Fields(coverUrlWithZoom)[0]
		coverUrls = append(coverUrls, coverUrl)
	}

	// Get largest cover (the last in the cover set)
	largestCover := coverUrls[len(coverUrls)-1]

	return largestCover
}

func bookAuthorFromInfoNodeValue(bookInfoNodeValue string) string {
	// Book info parts are separated by | the of which is the author
	bookInfoParts := strings.Split(bookInfoNodeValue, "|")
	bookAuthorPart := bookInfoParts[0]

	// Strip out the "by" from the author part
	bookAuthorFields := strings.Fields(bookAuthorPart)
	if len(bookAuthorFields) > 1 && strings.EqualFold(bookAuthorFields[0], "by") {
		bookAuthorFields = bookAuthorFields[1:]
	}

	// Rejoin author fields
	bookAuthor := strings.Join(bookAuthorFields, " ")

	return bookAuthor
}
