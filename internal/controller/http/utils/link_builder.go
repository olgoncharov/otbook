package utils

import (
	"fmt"
	"net/url"
)

type LinkBuilder struct {
	scheme             string
	profileLinkPattern string
}

func NewLinkBuilder(scheme, profileLinkPattern string) *LinkBuilder {
	return &LinkBuilder{
		scheme:             scheme,
		profileLinkPattern: profileLinkPattern,
	}
}

func (b *LinkBuilder) BuildProfileLink(host, username string) string {
	u := url.URL{
		Scheme: b.scheme,
		Host:   host,
		Path:   fmt.Sprintf(b.profileLinkPattern, username),
	}

	return u.String()
}
