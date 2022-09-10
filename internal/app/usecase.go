package app

import (
	"time"

	"github.com/olgoncharov/otbook/internal/pkg/hash"
	"github.com/olgoncharov/otbook/internal/pkg/jwt"
	"github.com/olgoncharov/otbook/internal/repository/mysql"
	"github.com/olgoncharov/otbook/internal/usecase/access/command/login"
	refreshToken "github.com/olgoncharov/otbook/internal/usecase/access/command/refresh_token"
	becomeFriends "github.com/olgoncharov/otbook/internal/usecase/friends/command/become_friends"
	friendsList "github.com/olgoncharov/otbook/internal/usecase/friends/query/list"
	"github.com/olgoncharov/otbook/internal/usecase/profile/command/create"
	updateUserProfile "github.com/olgoncharov/otbook/internal/usecase/profile/command/update"
	profilesList "github.com/olgoncharov/otbook/internal/usecase/profile/query/list"
	getUserProfile "github.com/olgoncharov/otbook/internal/usecase/profile/query/user_profile"
)

type useCases struct {
	signup            *create.Handler
	login             *login.Handler
	refreshToken      *refreshToken.Handler
	getUserProfile    *getUserProfile.Handler
	updateUserProfile *updateUserProfile.Handler
	profilesList      *profilesList.Handler
	friendsList       *friendsList.Handler
	becomeFriends     *becomeFriends.Handler
}

func initUsecases(
	cfg configer,
	repo *mysql.Repository,
	hashGenerator *hash.HashGenerator,
	passwordChecker *hash.HashChecker,
	tokenGenerator *jwt.TokenGenerator,
	nowFn func() time.Time,
) useCases {
	return useCases{
		signup:            create.NewHandler(repo, hashGenerator),
		login:             login.NewHandler(repo, passwordChecker, tokenGenerator, cfg, nowFn),
		refreshToken:      refreshToken.NewHandler(repo, tokenGenerator, cfg, nowFn),
		getUserProfile:    getUserProfile.NewHandler(repo),
		updateUserProfile: updateUserProfile.NewHandler(repo),
		profilesList:      profilesList.NewHandler(repo),
		friendsList:       friendsList.NewHandler(repo),
		becomeFriends:     becomeFriends.NewHandler(repo, repo),
	}
}
