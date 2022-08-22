package utils

import (
	"net/url"
	"strconv"
)

const (
	limitQueryKey  = "limit"
	offsetQueryKey = "offset"
)

func GetLimitOffsetFromURL(u *url.URL) (uint, uint) {
	var limit, offset int

	limitRaw := u.Query().Get(limitQueryKey)
	if limitRaw != "" {
		limit, _ = strconv.Atoi(limitRaw)
	}

	offsetRaw := u.Query().Get(offsetQueryKey)
	if offsetRaw != "" {
		offset, _ = strconv.Atoi(offsetRaw)
	}

	return uint(limit), uint(offset)
}
