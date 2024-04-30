package utils_test

import (
	"testing"

	"github.com/ahobsonsayers/abs-tract/utils"
	"github.com/stretchr/testify/require"
)

func TestCleanImageURL(t *testing.T) {
	imageUrl := "https://m.media-amazon.com/images/I/61Ng-W9EhBL.jpg"
	dirtyImageUrl := "https://m.media-amazon.com/images/I/61Ng-W9EhBL._AC_UY218_.jpg"
	cleanImageUrl := utils.SanitiseImageURL(dirtyImageUrl)
	require.Equal(t, imageUrl, cleanImageUrl)
}
