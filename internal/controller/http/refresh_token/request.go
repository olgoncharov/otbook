package refreshtoken

import "errors"

type requestBody struct {
	RefreshToken string `json:"refreshToken"`
}

func (r *requestBody) validate() error {
	if r.RefreshToken == "" {
		return errors.New("empty refresh token")
	}

	return nil
}
