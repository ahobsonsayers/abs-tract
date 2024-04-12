package goodreads_test

import (
	"encoding/xml"
	"testing"

	"github.com/ahobsonsayers/abs-goodreads/goodreads"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalGenres(t *testing.T) {
	// XML of popular shelves for The Hobbit
	xmlString := `
	<popular_shelves>
		<shelf name="to-read" count="1159744"/>
		<shelf name="currently-reading" count="93681"/>
		<shelf name="fantasy" count="64880"/>
		<shelf name="classics" count="22453"/>
		<shelf name="fiction" count="18153"/>
		<shelf name="owned" count="9082"/>
		<shelf name="books-i-own" count="6317"/>
		<shelf name="classic" count="4091"/>
		<shelf name="adventure" count="3951"/>
	</popular_shelves>
	`

	var genres goodreads.Genres
	err := xml.Unmarshal([]byte(xmlString), &genres)
	require.NoError(t, err)

	expectedGenres := goodreads.Genres{"Fantasy", "Classics", "Fiction", "Classic", "Adventure"}
	require.Equal(t, expectedGenres, genres)
}
