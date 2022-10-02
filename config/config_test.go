package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	configPath := t.TempDir() + "/test_config.yaml"
	err := os.WriteFile(configPath, []byte(`
database:
  - host: "localhost"
    port: "3306"
    user: "admin"
    password: "admin"
    db_name: "otbook"
    role: master

  - host: "localhost"
    port: "3307"
    user: "admin"
    password: "admin"
    db_name: "otbook"
    role: replica

  - host: "localhost"
    port: "3308"
    user: "admin"
    password: "admin"
    db_name: "otbook"
    role: replica

jwt:
  access_token_ttl: 3600
  refresh_token_ttl: 108000
  signing_key: nnUdOVjyvgtsTUPguUrm

password:
  hash_generation_cost: 5

http:
  server_addr: ":8000"
`), 0666)

	if err != nil {
		t.Fatalf("can't create file with test config: %s", err.Error())
	}

	t.Run("only yaml config", func(t *testing.T) {
		cfg, err := NewConfigFromFile(configPath)
		require.NoError(t, err)

		assert.EqualValues(t, 3600, cfg.JWTAccessTokenTTL())
		assert.EqualValues(t, 108000, cfg.JWTRefreshTokenTTL())
		assert.EqualValues(t, "nnUdOVjyvgtsTUPguUrm", cfg.JWTSigningKey())
		assert.EqualValues(t, 5, cfg.PasswordHashGenerationCost())
		assert.EqualValues(t, ":8000", cfg.HTTPServerAddr())

		masterDBConfig, err := cfg.MasterDBConfig()
		require.NoError(t, err)
		assert.Equal(t, DBInstanceConfig{
			Host:     "localhost",
			Port:     "3306",
			User:     "admin",
			Password: "admin",
			DBName:   "otbook",
			Role:     "master",
		}, masterDBConfig)

		replicaConfigs := cfg.ReplicaConfigs()
		require.Len(t, replicaConfigs, 2)
		assert.Equal(t, DBInstanceConfig{
			Host:     "localhost",
			Port:     "3307",
			User:     "admin",
			Password: "admin",
			DBName:   "otbook",
			Role:     "replica",
		}, replicaConfigs[0])
		assert.Equal(t, DBInstanceConfig{
			Host:     "localhost",
			Port:     "3308",
			User:     "admin",
			Password: "admin",
			DBName:   "otbook",
			Role:     "replica",
		}, replicaConfigs[1])
	})

	t.Run("yaml + env", func(t *testing.T) {
		t.Setenv("JWT_ACESS_TOKEN_TTL", "5000")
		t.Setenv("JWT_REFRESH_TOKEN_TTL", "10000")
		t.Setenv("JWT_SIGNING_KEY", "signing_key_from_env")
		t.Setenv("PASSWORD_HASH_GENERATION_COST", "7")
		t.Setenv("HTTP_SERVER_ADDR", ":5000")
		t.Setenv("DB_HOST_1", "db_alpha")
		t.Setenv("DB_PORT_1", "5001")
		t.Setenv("DB_USER_1", "user1")
		t.Setenv("DB_PASSWORD_1", "password1")
		t.Setenv("DB_NAME_1", "db1")
		t.Setenv("DB_ROLE_1", "master")
		t.Setenv("DB_HOST_2", "db_bravo")
		t.Setenv("DB_PORT_2", "5002")
		t.Setenv("DB_USER_2", "user2")
		t.Setenv("DB_PASSWORD_2", "password2")
		t.Setenv("DB_NAME_2", "db2")
		t.Setenv("DB_ROLE_2", "replica")
		t.Setenv("DB_HOST_3", "db_charlie")
		t.Setenv("DB_PORT_3", "5003")
		t.Setenv("DB_USER_3", "user3")
		t.Setenv("DB_PASSWORD_3", "password3")
		t.Setenv("DB_NAME_3", "db3")
		t.Setenv("DB_ROLE_3", "replica")

		cfg, err := NewConfigFromFile(configPath)
		require.NoError(t, err)

		assert.EqualValues(t, 5000, cfg.JWTAccessTokenTTL())
		assert.EqualValues(t, 10000, cfg.JWTRefreshTokenTTL())
		assert.EqualValues(t, "signing_key_from_env", cfg.JWTSigningKey())
		assert.EqualValues(t, 7, cfg.PasswordHashGenerationCost())
		assert.EqualValues(t, ":5000", cfg.HTTPServerAddr())

		masterDBConfig, err := cfg.MasterDBConfig()
		require.NoError(t, err)
		assert.Equal(t, DBInstanceConfig{
			Host:     "db_alpha",
			Port:     "5001",
			User:     "user1",
			Password: "password1",
			DBName:   "db1",
			Role:     "master",
		}, masterDBConfig)

		replicaConfigs := cfg.ReplicaConfigs()
		require.Len(t, replicaConfigs, 2)
		assert.Equal(t, DBInstanceConfig{
			Host:     "db_bravo",
			Port:     "5002",
			User:     "user2",
			Password: "password2",
			DBName:   "db2",
			Role:     "replica",
		}, replicaConfigs[0])

		assert.Equal(t, DBInstanceConfig{
			Host:     "db_charlie",
			Port:     "5003",
			User:     "user3",
			Password: "password3",
			DBName:   "db3",
			Role:     "replica",
		}, replicaConfigs[1])
	})
}
