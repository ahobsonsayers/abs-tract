package goodreads_test

import (
	"encoding/xml"
	"testing"

	"github.com/ahobsonsayers/abs-tract/goodreads"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalGenres(t *testing.T) {
	// XML of popular shelves for The Hobbit
	xmlString := `
	<popular_shelves>
		<shelf name="to-read" count="1161504"/>
		<shelf name="currently-reading" count="93764"/>
		<shelf name="fantasy" count="64931"/>
		<shelf name="classics" count="22453"/>
		<shelf name="fiction" count="18170"/>
		<shelf name="owned" count="9096"/>
		<shelf name="books-i-own" count="6320"/>
		<shelf name="classic" count="4094"/>
		<shelf name="adventure" count="3954"/>
		<shelf name="young-adult" count="3184"/>
		<shelf name="favourites" count="3050"/>
		<shelf name="tolkien" count="2839"/>
		<shelf name="physical-tbr" count="2417"/>
		<shelf name="high-fantasy" count="2078"/>
		<shelf name="owned-books" count="2007"/>
	</popular_shelves>
	`

	var genres goodreads.Genres
	err := xml.Unmarshal([]byte(xmlString), &genres)
	require.NoError(t, err)

	expectedGenres := goodreads.Genres{"Fantasy", "Classic", "Fiction"}
	require.Equal(t, expectedGenres, genres)
}

func TestBookUnmarshalBrTagReplacement(t *testing.T) {
	testXML := `
	<GoodreadsResponse>
		<book>
			<id>123</id>
			<title>Test Book</title>
			<description><![CDATA[Test description<br />2. line<br/>3. line<br>4. line]]></description>
			<work>
				<original_title>Test Book</original_title>
				<ratings_sum>100</ratings_sum>
				<ratings_count>10</ratings_count>
			</work>
			<popular_shelves>
				<shelf name="fiction" count="100"/>
			</popular_shelves>
		</book>
	</GoodreadsResponse>
	`

	var response struct {
		Book goodreads.Book `xml:"book"`
	}

	err := xml.Unmarshal([]byte(testXML), &response)
	require.NoError(t, err)

	description := response.Book.BestEdition.Description
	t.Logf("Description after processing: %q", description)

	// Verify that <br> tags have been correctly converted to newlines
	require.Contains(t, description, "Test description\n2. line\n3. line\n4. line")

	// Verify that no HTML br tags remain
	require.NotContains(t, description, "<br")
	require.NotContains(t, description, "br>")
}
