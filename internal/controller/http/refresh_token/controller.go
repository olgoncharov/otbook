package refreshtoken

import (
	"context"
	"errors"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/olgoncharov/otbook/internal/controller/http/utils"
	refreshtoken "github.com/olgoncharov/otbook/internal/usecase/access/command/refresh_token"
	"github.com/rs/zerolog"
)

type (
	useCase interface {
		Handle(ctx context.Context, command refreshtoken.Command) (*refreshtoken.Result, error)
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

	responseEncoder := jsoniter.NewEncoder(w)

	resp, code, err := c.serve(r.Context(), body)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), code)

		return
	}

	w.WriteHeader(code)
	responseEncoder.Encode(resp)
}

func (c *Controller) serve(ctx context.Context, reqBody *requestBody) (*response, int, error) {
	if err := reqBody.validate(); err != nil {
		return nil, http.StatusBadRequest, err
	}

	username := utils.GetUsernameFromContext(ctx)
	result, err := c.useCase.Handle(ctx, refreshtoken.Command{
		Username:     username,
		RefreshToken: reqBody.RefreshToken,
	})

	if errors.Is(err, refreshtoken.ErrInvalidTokenGiven) {
		return nil, http.StatusBadRequest, err
	}

	if errors.Is(err, refreshtoken.ErrExpiredTokenGiven) {
		return nil, http.StatusBadRequest, err
	}

	if err != nil {
		c.logger.Error().Err(err).Msg("")

		return nil, http.StatusInternalServerError, utils.ErrInternal
	}

	return &response{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, http.StatusOK, nil
}
