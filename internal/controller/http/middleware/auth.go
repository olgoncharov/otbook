package middleware

import (
	"net/http"
	"strings"

	"github.com/olgoncharov/otbook/internal/controller/http/utils"
	"github.com/rs/zerolog"
)

type (
	tokenValidator interface {
		ValidateAccessToken(tokenString string) (string, error)
	}

	JWTAuthMiddleware struct {
		validator tokenValidator
		logger    zerolog.Logger
	}
)

func NewJWTMiddleware(validator tokenValidator, logger zerolog.Logger) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
		validator: validator,
		logger:    logger,
	}
}

func (m *JWTAuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authFields := strings.Fields(r.Header.Get("Authorization"))
		if len(authFields) < 2 {
			m.logger.Debug().Str("path", r.URL.Path).Msg("credentials are not provided")
			utils.WriteJSONError(w, "credentials are not provided", http.StatusUnauthorized)

			return
		}

		if authFields[0] != "Bearer" {
			m.logger.Debug().Str("path", r.URL.Path).Msg("invalid credentials")
			utils.WriteJSONError(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		username, err := m.validator.ValidateAccessToken(authFields[1])

		if err != nil || username == "" {
			m.logger.Debug().Str("path", r.URL.Path).Err(err).Msg("invalid token")
			utils.WriteJSONError(w, "invalid token", http.StatusUnauthorized)
			return
		}

		newCtx := utils.AddUsernameToContext(r.Context(), username)
		next.ServeHTTP(w, r.Clone(newCtx))
	})
}
