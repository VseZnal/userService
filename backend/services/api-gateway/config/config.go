package config

import (
	"userService/libs/utils"
)

var (
	userServiceHostEnvName = "USER_SERVICE_HOST"
	userServicePortEnvName = "USER_SERVICE_PORT"
	gatewayHostEnvName     = "API_HOST"
	gatewayPortEnvName     = "API_PORT"
	jwtSecret              = "JWT_SECRET"
	refreshSecret          = "REFRESH_SECRET"
	cors                   = "cors"
)

type Config struct {
	HostUser    string
	PortUser    string
	HostGateway string
	PortGateway string
	Jwt         string
	Refresh     string
	Cors        string
}

func PublicConfig() *Config {
	return &Config{
		HostUser:    utils.TrimEnv(userServiceHostEnvName),
		PortUser:    utils.TrimEnv(userServicePortEnvName),
		HostGateway: utils.TrimEnv(gatewayHostEnvName),
		PortGateway: utils.TrimEnv(gatewayPortEnvName),
		Jwt:         utils.TrimEnv(jwtSecret),
		Refresh:     utils.TrimEnv(refreshSecret),
		Cors:        utils.TrimEnv(cors),
	}
}
