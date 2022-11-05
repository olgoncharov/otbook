package config

import (
	"path"
	"runtime"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		DB       []DBInstanceConfig `yaml:"database"`
		Redis    RedisConfig        `yaml:"redis"`
		JWT      JWTConfig          `yaml:"jwt"`
		Password PasswordConfig     `yaml:"password"`
		HTTP     HTTPConfig         `yaml:"http"`
		Feed     FeedConfig         `yaml:"feed"`
	}

	JWTConfig struct {
		AccessTokenTTL  uint64 `yaml:"access_token_ttl" env:"JWT_ACESS_TOKEN_TTL"`
		RefreshTokenTTL uint64 `yaml:"refresh_token_ttl" env:"JWT_REFRESH_TOKEN_TTL"`
		SigningKey      string `yaml:"signing_key" env:"JWT_SIGNING_KEY"`
	}

	PasswordConfig struct {
		HashGenerationCost int `yaml:"hash_generation_cost" env:"PASSWORD_HASH_GENERATION_COST"`
	}

	HTTPConfig struct {
		ServerAddr string `yaml:"server_addr" env:"HTTP_SERVER_ADDR"`
	}

	RedisConfig struct {
		Addr     string `yaml:"addr" env:"REDIS_ADDR"`
		Password string `yaml:"password" env:"REDIS_PASSWORD"`
		DB       uint64 `ysml:"db" env:"REDIS_DB"`
	}

	FeedConfig struct {
		Limit           int  `yaml:"limit" env:"POST_FEED_LIMIT"`
		IsCacheDisabled bool `yaml:"is_cache_disabled" env:"POST_FEED_CACHE_DISABLED"`
	}
)

func NewConfigFromFile(filePath string) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(filePath, &cfg)
	if err != nil {
		return nil, err
	}

	if dbInstancesConfig := readDBInstancesConfigFromEnv(); dbInstancesConfig != nil {
		cfg.DB = dbInstancesConfig
	}

	return &cfg, nil
}

func NewDefaultConfig() (*Config, error) {
	_, currentFileName, _, _ := runtime.Caller(0)
	currentDir := path.Dir(currentFileName)
	return NewConfigFromFile(path.Join(currentDir, "default_conf.yaml"))
}

func (c *Config) JWTAccessTokenTTL() uint64 {
	return c.JWT.AccessTokenTTL
}

func (c *Config) JWTRefreshTokenTTL() uint64 {
	return c.JWT.RefreshTokenTTL
}

func (c *Config) JWTSigningKey() []byte {
	return []byte(c.JWT.SigningKey)
}

func (c *Config) PasswordHashGenerationCost() int {
	return c.Password.HashGenerationCost
}

func (c *Config) HTTPServerAddr() string {
	return c.HTTP.ServerAddr
}

func (c *Config) RedisAddr() string {
	return c.Redis.Addr
}

func (c *Config) RedisPassword() string {
	return c.Redis.Password
}

func (c *Config) RedisDB() uint64 {
	return c.Redis.DB
}

func (c *Config) PostFeedLimit() int {
	return c.Feed.Limit
}

func (c *Config) IsFeedCacheDisabled() bool {
	return c.Feed.IsCacheDisabled
}
