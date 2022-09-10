package config

import "github.com/ilyakaznacheev/cleanenv"

type (
	Config struct {
		DB       DBConfig       `yaml:"database"`
		JWT      JWTConfig      `yaml:"jwt"`
		Password PasswordConfig `yaml:"password"`
		HTTP     HTTPConfig     `yaml:"http"`
	}

	DBConfig struct {
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     string `yaml:"port" env:"DB_PORT"`
		User     string `yaml:"user" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
		DBName   string `yaml:"db_name" env:"DB_NAME"`
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

	return &cfg, nil
}

func NewDefaultConfig() (*Config, error) {
	return NewConfigFromFile("./config/default_conf.yaml")
}

func (c *Config) DBHost() string {
	return c.DB.Host
}

func (c *Config) DBPort() string {
	return c.DB.Port
}

func (c *Config) DBUser() string {
	return c.DB.User
}

func (c *Config) DBPassword() string {
	return c.DB.Password
}

func (c *Config) DBName() string {
	return c.DB.DBName
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
