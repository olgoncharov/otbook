package app

import (
	"time"

	"github.com/olgoncharov/otbook/internal/pkg/hash"
	"github.com/olgoncharov/otbook/internal/pkg/jwt"
	"github.com/olgoncharov/otbook/internal/repository/mysql"
	"github.com/olgoncharov/otbook/internal/usecase/access/command/login"
	refreshToken "github.com/olgoncharov/otbook/internal/usecase/access/command/refresh_token"
	addFriend "github.com/olgoncharov/otbook/internal/usecase/friends/command/add"
	deleteFriend "github.com/olgoncharov/otbook/internal/usecase/friends/command/delete"
	friendsList "github.com/olgoncharov/otbook/internal/usecase/friends/query/list"
	createPost "github.com/olgoncharov/otbook/internal/usecase/post/command/create"
	postsList "github.com/olgoncharov/otbook/internal/usecase/post/query/full_list"
	getPost "github.com/olgoncharov/otbook/internal/usecase/post/query/single_post"
	"github.com/olgoncharov/otbook/internal/usecase/profile/command/create"
	updateUserProfile "github.com/olgoncharov/otbook/internal/usecase/profile/command/update"
	profilesList "github.com/olgoncharov/otbook/internal/usecase/profile/query/full_list"
	profilesSearch "github.com/olgoncharov/otbook/internal/usecase/profile/query/search"
	getUserProfile "github.com/olgoncharov/otbook/internal/usecase/profile/query/user_profile"
)

type useCases struct {
	signup            *create.Handler
	login             *login.Handler
	refreshToken      *refreshToken.Handler
	getUserProfile    *getUserProfile.Handler
	updateUserProfile *updateUserProfile.Handler
	profilesList      *profilesList.Handler
	profilesSearch    *profilesSearch.Handler
	friendsList       *friendsList.Handler
	addFriend         *addFriend.Handler
	deleteFriend      *deleteFriend.Handler
	postsList         *postsList.Handler
	getPost           *getPost.Handler
	createPost        *createPost.Handler
}

func initUsecases(
	cfg configer,
	writeRepo *mysql.Repository,
	readRepo *mysql.Repository,
	hashGenerator *hash.HashGenerator,
	passwordChecker *hash.HashChecker,
	tokenGenerator *jwt.TokenGenerator,
	nowFn func() time.Time,
) useCases {
	return useCases{
		signup:            create.NewHandler(writeRepo, hashGenerator),
		login:             login.NewHandler(writeRepo, passwordChecker, tokenGenerator, cfg, nowFn),
		refreshToken:      refreshToken.NewHandler(writeRepo, tokenGenerator, cfg, nowFn),
		getUserProfile:    getUserProfile.NewHandler(readRepo),
		updateUserProfile: updateUserProfile.NewHandler(writeRepo),
		profilesList:      profilesList.NewHandler(readRepo),
		profilesSearch:    profilesSearch.NewHandler(readRepo),
		friendsList:       friendsList.NewHandler(readRepo),
		addFriend:         addFriend.NewHandler(writeRepo, readRepo),
		deleteFriend:      deleteFriend.NewHandler(writeRepo, readRepo),
		postsList:         postsList.NewHandler(readRepo),
		getPost:           getPost.NewHandler(readRepo),
		createPost:        createPost.NewHandler(writeRepo, nowFn),
	}
}
