package app

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	friends "github.com/olgoncharov/otbook/internal/controller/http/friends"
	"github.com/olgoncharov/otbook/internal/controller/http/login"
	"github.com/olgoncharov/otbook/internal/controller/http/middleware"
	myfriends "github.com/olgoncharov/otbook/internal/controller/http/my_friends"
	myprofile "github.com/olgoncharov/otbook/internal/controller/http/my_profile"
	"github.com/olgoncharov/otbook/internal/controller/http/posts"
	"github.com/olgoncharov/otbook/internal/controller/http/profile"
	profileslist "github.com/olgoncharov/otbook/internal/controller/http/profiles_list"
	profilesSearch "github.com/olgoncharov/otbook/internal/controller/http/profiles_search"
	refreshtoken "github.com/olgoncharov/otbook/internal/controller/http/refresh_token"
	"github.com/olgoncharov/otbook/internal/controller/http/signup"
	"github.com/olgoncharov/otbook/internal/controller/http/utils"
	"github.com/olgoncharov/otbook/internal/pkg/jwt"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
)

func initHTTPServer(cfg configer, uc useCases) *http.Server {
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("component", "http").
		Logger()

	linkBuilder := utils.NewLinkBuilder(
		"/api/v1/profiles/%s",
		"/api/v1/posts/%d",
	)

	tokenValidator := jwt.NewTokenValidator(cfg)
	jwtMiddleware := middleware.NewJWTMiddleware(tokenValidator, logger)

	router := mux.NewRouter()
	router.Use(middleware.SetJSONContentType)

	subRouterNoAuth := router.PathPrefix("/api/v1").Subrouter()
	subRouterAuth := router.PathPrefix("/api/v1").Subrouter()
	subRouterAuth.Use(jwtMiddleware.Middleware)

	subRouterNoAuth.Handle(
		"/signup",
		signup.NewController(uc.signup, logger.With().Str("path", "/signup").Logger()),
	).Methods(http.MethodPost)

	subRouterNoAuth.Handle(
		"/auth/login",
		login.NewController(uc.login, logger.With().Str("path", "/auth/login").Logger()),
	).Methods(http.MethodPost)

	subRouterNoAuth.Handle(
		"/profiles",
		profileslist.NewController(uc.profilesList, logger.With().Str("path", "/profiles").Logger()),
	).Methods(http.MethodGet)

	subRouterNoAuth.Handle(
		"/profiles/search",
		profilesSearch.NewController(uc.profilesSearch, linkBuilder, logger.With().Str("path", "/profiles/search").Logger()),
	).Methods(http.MethodPost)

	subRouterNoAuth.Handle(
		"/profiles/{username}",
		profile.NewController(uc.getUserProfile, logger.With().Str("path", "/profiles/{username}").Logger()),
	).Methods(http.MethodGet)

	subRouterAuth.Handle(
		"/auth/refresh",
		refreshtoken.NewController(uc.refreshToken, logger.With().Str("path", "/auth/refresh").Logger()),
	)

	subRouterAuth.Handle(
		"/me/profile",
		myprofile.NewController(uc.getUserProfile, uc.updateUserProfile, logger.With().Str("path", "/me/profile").Logger()),
	).Methods(http.MethodGet, http.MethodPost)

	subRouterAuth.Handle(
		"/me/friends",
		myfriends.NewController(uc.friendsList, logger.With().Str("path", "/me/friends").Logger()),
	).Methods(http.MethodGet)

	subRouterAuth.Handle(
		"/profiles/{username}/friends",
		friends.NewController(uc.friendsList, uc.addFriend, uc.deleteFriend, logger.With().Str("path", "/profiles/{username}/friends").Logger()),
	).Methods(http.MethodGet, http.MethodPost, http.MethodDelete)

	subRouterAuth.Handle(
		"/posts",
		posts.NewListController(uc.postsList, uc.createPost, linkBuilder, logger.With().Str("path", "/posts").Logger()),
	).Methods(http.MethodGet, http.MethodPost)

	subRouterAuth.Handle(
		"/posts/{id}",
		posts.NewObjectController(uc.getPost, linkBuilder, logger.With().Str("path", "/posts/{id}").Logger()),
	).Methods(http.MethodGet)

	return &http.Server{
		Addr:    cfg.HTTPServerAddr(),
		Handler: cors.AllowAll().Handler(router),
	}
}
