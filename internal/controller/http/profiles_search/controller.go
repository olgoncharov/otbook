package profilessearch

import (
	"context"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/olgoncharov/otbook/internal/controller/http/utils"
	"github.com/olgoncharov/otbook/internal/usecase/profile/query/search"
	"github.com/rs/zerolog"
)

type (
	useCase interface {
		Handle(ctx context.Context, query search.Query) ([]search.ProfileInfo, error)
	}

	linkBuilder interface {
		BuildProfileLink(host, username string) string
	}

	Controller struct {
		useCase     useCase
		linkBuilder linkBuilder
		logger      zerolog.Logger
	}
)

func NewController(uCase useCase, lBuilder linkBuilder, logger zerolog.Logger) *Controller {
	return &Controller{
		useCase:     uCase,
		linkBuilder: lBuilder,
		logger:      logger,
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

	err = body.validate()
	if err != nil {
		c.logger.Error().Err(err).Msg("")
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)

		return
	}

	foundProfiles, err := c.useCase.Handle(r.Context(), search.Query{
		FirstNamePrefix: body.FirstName,
		LastNamePrefix:  body.LastName,
	})

	if err != nil {
		c.logger.Error().Err(err).Msg("")
		utils.WriteJSONError(w, utils.ErrInternal.Error(), http.StatusInternalServerError)

		return
	}

	resp := response{
		List: make([]profileInfo, len(foundProfiles)),
	}

	for i, p := range foundProfiles {
		resp.List[i] = profileInfo{
			Username:  p.Username,
			FirstName: p.FirstName,
			LastName:  p.LastName,
			Links: profileLinks{
				Self: c.linkBuilder.BuildProfileLink(r.Host, p.Username),
			},
		}
	}

	responseEncoder := jsoniter.NewEncoder(w)
	responseEncoder.Encode(resp)
}
