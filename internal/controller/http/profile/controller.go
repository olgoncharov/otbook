package profile

import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"github.com/olgoncharov/otbook/internal/controller/http/utils"
	"github.com/olgoncharov/otbook/internal/entity"
	"github.com/olgoncharov/otbook/internal/pkg/types"
	"github.com/olgoncharov/otbook/internal/usecase/profile/query/user_profile"
	"github.com/rs/zerolog"
)

type (
	useCase interface {
		Handle(ctx context.Context, query user_profile.Query) (*entity.Profile, error)
	}

	Controller struct {
		useCase useCase
		logger  zerolog.Logger
	}
)

func NewController(uCase useCase, logger zerolog.Logger) *Controller {
	return &Controller{
		useCase: uCase,
		logger:  logger,
	}
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)

	profile, err := c.useCase.Handle(r.Context(), user_profile.Query{
		Username: reqVars["username"],
	})

	if errors.Is(err, user_profile.ErrUserNotFound) {
		utils.WriteJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		c.logger.Error().Err(err).Msg("")
		utils.WriteJSONError(w, utils.ErrInternal.Error(), http.StatusInternalServerError)

		return
	}

	responseEncoder := jsoniter.NewEncoder(w)
	responseEncoder.Encode(response{
		Username:  profile.Username,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Birthdate: types.Date{Time: profile.Birthdate},
		City:      profile.City,
		Sex:       profile.Sex,
		Hobby:     profile.Hobby,
	})
}
