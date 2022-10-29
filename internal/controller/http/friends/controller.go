package friends

import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"github.com/olgoncharov/otbook/internal/controller/http/utils"
	"github.com/olgoncharov/otbook/internal/entity"
	"github.com/olgoncharov/otbook/internal/pkg/types"
	addFriend "github.com/olgoncharov/otbook/internal/usecase/friends/command/add"
	"github.com/olgoncharov/otbook/internal/usecase/friends/query/list"
	"github.com/rs/zerolog"
)

type (
	getFriendsUseCase interface {
		Handle(ctx context.Context, query list.Query) (*list.Result, error)
	}
	addFriendUseCase interface {
		Handle(ctx context.Context, command addFriend.Command) error
	}

	Controller struct {
		getFriendsUseCase getFriendsUseCase
		addFriendUseCase  addFriendUseCase
		logger            zerolog.Logger
	}
)

func NewController(getUCase getFriendsUseCase, addFriendUCase addFriendUseCase, logger zerolog.Logger) *Controller {
	return &Controller{
		getFriendsUseCase: getUCase,
		addFriendUseCase:  addFriendUCase,
		logger:            logger,
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

func (c *Controller) servePOST(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqVars := mux.Vars(r)

	err := c.addFriendUseCase.Handle(ctx, addFriend.Command{
		User:      utils.GetUsernameFromContext(ctx),
		NewFriend: reqVars["username"],
	})

	if errors.Is(err, addFriend.ErrUserNotFound) {
		utils.WriteJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	if errors.Is(err, addFriend.ErrSelfFriendship) || errors.Is(err, addFriend.ErrAlreadyFriends) {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		c.logger.Error().Err(err).Msg("")
		utils.WriteJSONError(w, utils.ErrInternal.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (c *Controller) serveGET(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)
	limit, offset := utils.GetLimitOffsetFromURL(r.URL)

	result, err := c.getFriendsUseCase.Handle(r.Context(), list.Query{
		Username: reqVars["username"],
		Limit:    limit,
		Offset:   offset,
	})

	if err != nil {
		c.logger.Error().Err(err).Msg("")
		utils.WriteJSONError(w, utils.ErrInternal.Error(), http.StatusInternalServerError)

		return
	}

	resp := getResponse{
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
