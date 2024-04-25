package utils

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

// CloneURL clones a url. Copied directly from net/http internals
// See: https://github.com/golang/go/blob/go1.19/src/net/http/clone.go#L22
func CloneURL(u *url.URL) *url.URL {
	if u == nil {
		return nil
	}
	u2 := new(url.URL)
	*u2 = *u
	if u.User != nil {
		u2.User = new(url.Userinfo)
		*u2.User = *u.User
	}
	return u2
}

// SanitiseImageURL removes all extensions but the final extension from images urls.
// Some image urls contain additional '._modifier_.' extensions before the
// final extension. This function will strip these modifier extensions leaving
// the url of the original image.
func SanitiseImageURL(imageURL string) string {
	// Attempt to parse the image URL and return an empty string if it fails
	parsedUrl, err := url.Parse(imageURL)
	if err != nil {
		return ""
	}

	// Extract the file name from the URL path
	dirPath, fileName := path.Split(parsedUrl.Path)
	fileParts := strings.Split(fileName, ".")

	// Return the original URL if the file name has not extension
	if len(fileParts) == 1 {
		return imageURL
	}

	// Construct a new file name with only the last extension
	newFileName := fmt.Sprintf("%s.%s", fileParts[0], fileParts[len(fileParts)-1])
	parsedUrl.Path = path.Join(dirPath, newFileName)

	return parsedUrl.String()
}
