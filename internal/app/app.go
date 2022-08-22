package app

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/olgoncharov/otbook/internal/pkg/hash"
	"github.com/olgoncharov/otbook/internal/pkg/jwt"
	"github.com/olgoncharov/otbook/internal/repository/mysql"
	"github.com/rs/zerolog/log"
)

type (
	App struct {
		httpServer *http.Server
		db         *sql.DB
	}

	configer interface {
		DBHost() string
		DBPort() string
		DBUser() string
		DBPassword() string
		DBName() string

		JWTAccessTokenTTL() uint64
		JWTRefreshTokenTTL() uint64
		JWTSigningKey() []byte

		PasswordHashGenerationCost() int

		HTTPServerAddr() string
	}
)

func NewApp(cfg configer) *App {
	db := InitDB(cfg)
	repo := mysql.NewRepository(db)

	nowFn := func() time.Time {
		return time.Now()
	}

	passwordChecker := hash.NewHashChecker()
	hashGenerator := hash.NewHashGenerator(cfg.PasswordHashGenerationCost())
	tokenGenerator := jwt.NewTokenGenerator(cfg, nowFn)

	uc := initUsecases(cfg, repo, hashGenerator, passwordChecker, tokenGenerator, nowFn)
	httpServer := initHTTPServer(cfg, uc)

	return &App{
		httpServer: httpServer,
		db:         db,
	}
}

func (a *App) Run() {
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Err(err)
		}
	}()
}

func (a *App) Shutdown(ctx context.Context) {
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		log.Err(err)
	}

	err = a.db.Close()
	if err != nil {
		log.Err(err)
	}
}
