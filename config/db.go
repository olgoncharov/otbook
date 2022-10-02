package config

import (
	"errors"
	"os"
	"strconv"
)

const (
	masterRole  string = "master"
	replicaRole string = "replica"
)

type DBInstanceConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	Role     string `yaml:"role"`
}

func (c *Config) MasterDBConfig() (DBInstanceConfig, error) {
	configs := c.dbInstancesByRole(masterRole)
	if len(configs) != 1 {
		return DBInstanceConfig{}, errors.New("config must contain 1 master db")
	}

	return configs[0], nil
}

func (c *Config) ReplicaConfigs() []DBInstanceConfig {
	return c.dbInstancesByRole(replicaRole)
}

func (c *Config) dbInstancesByRole(role string) []DBInstanceConfig {
	var result []DBInstanceConfig
	for _, db := range c.DB {
		if db.Role == role {
			result = append(result, db)
		}
	}

	return result
}

func readDBInstancesConfigFromEnv() []DBInstanceConfig {
	var (
		result  []DBInstanceConfig
		counter int64 = 1

		host, port, user, password, dbName, role string
	)

	for {
		suffix := strconv.FormatInt(counter, 10)

		host = os.Getenv("DB_HOST_" + suffix)
		if host == "" {
			break
		}

		port = os.Getenv("DB_PORT_" + suffix)
		if port == "" {
			break
		}

		user = os.Getenv("DB_USER_" + suffix)
		if user == "" {
			break
		}

		password = os.Getenv("DB_PASSWORD_" + suffix)
		if password == "" {
			break
		}

		dbName = os.Getenv("DB_NAME_" + suffix)
		if dbName == "" {
			break
		}

		role = os.Getenv("DB_ROLE_" + suffix)
		if role == "" {
			break
		}

		result = append(result, DBInstanceConfig{
			Host:     host,
			Port:     port,
			User:     user,
			Password: password,
			DBName:   dbName,
			Role:     role,
		})

		counter++
	}

	return result
}
