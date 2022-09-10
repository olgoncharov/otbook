package signup

import (
	"context"
	"errors"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/olgoncharov/otbook/internal/controller/http/utils"
	"github.com/olgoncharov/otbook/internal/usecase/profile/command/create"
	"github.com/rs/zerolog"
)

type (
	useCase interface {
		Handle(ctx context.Context, command create.Command) error
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
	bodyDecoder := jsoniter.NewDecoder(r.Body)
	body := &requestBody{}
	err := bodyDecoder.Decode(body)

	if err != nil {
		c.logger.Error().Err(err).Msg("")
		utils.WriteJSONError(w, "invalid request body", http.StatusBadRequest)

		return
	}

	code, err := c.serve(r.Context(), body)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), code)

		return
	}

	w.WriteHeader(code)
	w.Write([]byte("{}"))
}

func (c *Controller) serve(ctx context.Context, reqBody *requestBody) (int, error) {
	if err := reqBody.validate(); err != nil {
		return http.StatusBadRequest, err
	}

	err := c.useCase.Handle(ctx, create.Command{
		Username:  reqBody.Username,
		Password:  reqBody.Password,
		FirstName: reqBody.FirstName,
		LastName:  reqBody.LastName,
		Birthdate: reqBody.Birthdate.Time,
		City:      reqBody.City,
		Sex:       reqBody.Sex,
		Hobby:     reqBody.Hobby,
	})

	if errors.Is(err, create.ErrUserAlreadyExists) {
		return http.StatusBadRequest, err
	}

	if err != nil {
		c.logger.Error().Err(err).Msg("")

		return http.StatusInternalServerError, utils.ErrInternal
	}

	return http.StatusCreated, nil
}
