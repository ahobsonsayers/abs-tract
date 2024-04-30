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
