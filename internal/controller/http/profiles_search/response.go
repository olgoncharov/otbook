package profilessearch

type (
	profileLinks struct {
		Self string `json:"self"`
	}

	profileInfo struct {
		Username  string       `json:"username"`
		FirstName string       `json:"firstName"`
		LastName  string       `json:"lastName"`
		Links     profileLinks `json:"links"`
	}

	response struct {
		List []profileInfo
	}
)
