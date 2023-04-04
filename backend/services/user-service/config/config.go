package config

import (
	"userService/libs/utils"
)

var (
	userServiceHostEnvName = "USER_SERVICE_HOST"
	userServicePortEnvName = "USER_SERVICE_PORT"
	pgConnStringEnvName    = "PG_CONNECTION_STRING"
	jwtSecret              = "JWT_SECRET"
	refreshSecret          = "REFRESH_SECRET"
)

type Config struct {
	HostUser     string
	PortUser     string
	PgConnString string
	Jwt          string
	Refresh      string
}

func UserConfig() *Config {
	return &Config{
		HostUser:     utils.TrimEnv(userServiceHostEnvName),
		PortUser:     utils.TrimEnv(userServicePortEnvName),
		PgConnString: utils.TrimEnv(pgConnStringEnvName),
		Jwt:          utils.TrimEnv(jwtSecret),
		Refresh:      utils.TrimEnv(refreshSecret),
	}
}
