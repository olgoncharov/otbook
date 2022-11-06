package utils

import (
	"fmt"
)

type LinkBuilder struct {
	profileLinkPattern string
	postLinkPattern    string
}

func NewLinkBuilder(profileLinkPattern, postLinkPattern string) *LinkBuilder {
	return &LinkBuilder{
		profileLinkPattern: profileLinkPattern,
		postLinkPattern:    postLinkPattern,
	}
}

func (b *LinkBuilder) BuildProfileLink(username string) string {
	return fmt.Sprintf(b.profileLinkPattern, username)
}

func (b *LinkBuilder) BuildPostLink(postID uint64) string {
	return fmt.Sprintf(b.postLinkPattern, postID)
}
