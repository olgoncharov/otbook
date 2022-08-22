package myfriends

import (
	"context"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/olgoncharov/otbook/internal/controller/http/utils"
	"github.com/olgoncharov/otbook/internal/entity"
	"github.com/olgoncharov/otbook/internal/pkg/types"
	"github.com/olgoncharov/otbook/internal/usecase/friends/query/list"
	"github.com/rs/zerolog"
)

type (
	useCase interface {
		Handle(ctx context.Context, query list.Query) (*list.Result, error)
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
	username := utils.GetUsernameFromContext(r.Context())
	limit, offset := utils.GetLimitOffsetFromURL(r.URL)

	result, err := c.useCase.Handle(r.Context(), list.Query{
		Username: username,
		Limit:    limit,
		Offset:   offset,
	})

	if err != nil {
		c.logger.Error().Err(err).Msg("")
		utils.WriteJSONError(w, utils.ErrInternal.Error(), http.StatusInternalServerError)

		return
	}

	resp := response{
		List:       make([]profileInfo, len(result.Friends)),
		TotalCount: result.TotalCount,
	}

	for i, p := range result.Friends {
		resp.List[i] = convertDomainProfileToResponse(p)
	}

	responseEncoder := jsoniter.NewEncoder(w)
	responseEncoder.Encode(resp)
}

func convertDomainProfileToResponse(domainProfile entity.Profile) profileInfo {
	return profileInfo{
		Username:  domainProfile.Username,
		FirstName: domainProfile.FirstName,
		LastName:  domainProfile.LastName,
		Birthdate: types.Date{Time: domainProfile.Birthdate},
		City:      domainProfile.City,
		Sex:       domainProfile.Sex,
		Hobby:     domainProfile.Hobby,
	}
}
