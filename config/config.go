package config

import "github.com/ilyakaznacheev/cleanenv"

type (
	Config struct {
		DB       []DBInstanceConfig `yaml:"database"`
		JWT      JWTConfig          `yaml:"jwt"`
		Password PasswordConfig     `yaml:"password"`
		HTTP     HTTPConfig         `yaml:"http"`
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
	return NewConfigFromFile("./config/default_conf.yaml")
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
