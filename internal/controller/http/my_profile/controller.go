package myprofile

import (
	"context"
	"errors"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/olgoncharov/otbook/internal/controller/http/utils"
	"github.com/olgoncharov/otbook/internal/entity"
	"github.com/olgoncharov/otbook/internal/pkg/types"
	"github.com/olgoncharov/otbook/internal/usecase/profile/command/update"
	"github.com/olgoncharov/otbook/internal/usecase/profile/query/user_profile"
	"github.com/rs/zerolog"
)

type (
	getUseCase interface {
		Handle(ctx context.Context, query user_profile.Query) (*entity.Profile, error)
	}

	updateUseCase interface {
		Handle(ctx context.Context, command update.Command) error
	}

	Controller struct {
		getUseCase    getUseCase
		updateUseCase updateUseCase
		logger        zerolog.Logger
	}
)

func NewController(getUCase getUseCase, updateUCase updateUseCase, logger zerolog.Logger) *Controller {
	return &Controller{
		getUseCase:    getUCase,
		updateUseCase: updateUCase,
		logger:        logger,
	}
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.serveGET(w, r)
	case http.MethodPost:
		c.servePOST(w, r)
	}
}

func (c *Controller) serveGET(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := utils.GetUsernameFromContext(ctx)

	profile, err := c.getUseCase.Handle(ctx, user_profile.Query{
		Username: username,
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

func (c *Controller) servePOST(w http.ResponseWriter, r *http.Request) {
	bodyDecoder := jsoniter.NewDecoder(r.Body)
	body := &requestBody{}
	err := bodyDecoder.Decode(body)

	if err != nil {
		c.logger.Error().Err(err).Msg("")
		utils.WriteJSONError(w, "invalid request body", http.StatusBadRequest)

		return
	}

	if err = body.validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)

		return
	}

	err = c.updateUseCase.Handle(r.Context(), update.Command{
		Username:  utils.GetUsernameFromContext(r.Context()),
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Birthdate: body.Birthdate.Time,
		City:      body.City,
		Sex:       body.Sex,
		Hobby:     body.Hobby,
	})

	if errors.Is(err, update.ErrUserNotFound) {
		w.WriteHeader(http.StatusNotFound)
		utils.WriteJSONError(w, err.Error(), http.StatusNotFound)

		return
	}

	if err != nil {
		c.logger.Error().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJSONError(w, utils.ErrInternal.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}
