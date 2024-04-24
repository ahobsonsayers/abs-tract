package utils

import "net/url"

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
