package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/olgoncharov/otbook/config"
	"github.com/olgoncharov/otbook/internal/pkg/hash"
	"github.com/olgoncharov/otbook/internal/pkg/jwt"
	"github.com/olgoncharov/otbook/internal/repository/mysql"
	"github.com/rs/zerolog/log"
)

type (
	App struct {
		httpServer    *http.Server
		dbConnections []*sql.DB
	}

	configer interface {
		MasterDBConfig() (config.DBInstanceConfig, error)
		ReplicaConfigs() []config.DBInstanceConfig

		JWTAccessTokenTTL() uint64
		JWTRefreshTokenTTL() uint64
		JWTSigningKey() []byte

		PasswordHashGenerationCost() int

		HTTPServerAddr() string
	}
)

func NewApp(cfg configer) (*App, error) {
	masterDBConfig, err := cfg.MasterDBConfig()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	masterDB, err := initDB(masterDBConfig)
	if err != nil {
		return nil, fmt.Errorf("can't connect to db: %w", err)
	}
	writeRepo := mysql.NewRepository(masterDB)
	dbConnections := []*sql.DB{masterDB}

	var readRepo *mysql.Repository
	replicaConfigs := cfg.ReplicaConfigs()
	if len(replicaConfigs) > 0 {
		replicaDB, err := initDB(replicaConfigs[0])
		if err != nil {
			return nil, fmt.Errorf("can't connect to db: %w", err)
		}
		readRepo = mysql.NewRepository(replicaDB)
		dbConnections = append(dbConnections, replicaDB)
	} else {
		readRepo = writeRepo
	}

	passwordChecker := hash.NewHashChecker()
	hashGenerator := hash.NewHashGenerator(cfg.PasswordHashGenerationCost())
	tokenGenerator := jwt.NewTokenGenerator(cfg, time.Now)

	uc := initUsecases(cfg, writeRepo, readRepo, hashGenerator, passwordChecker, tokenGenerator, time.Now)
	httpServer := initHTTPServer(cfg, uc)

	return &App{
		httpServer:    httpServer,
		dbConnections: dbConnections,
	}, nil
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

	for _, db := range a.dbConnections {
		err = db.Close()
		if err != nil {
			log.Err(err)
		}
	}
}
