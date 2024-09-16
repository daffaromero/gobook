package config

import (
	"fmt"
	"log"

	"github.com/daffaromero/gobook/services/common/utils"
)

var EndpointPrefix = utils.GetEnv("ENDPOINT_PREFIX")

type ServerConfig struct {
	HTTP       string
	HTTPAddr   string
	HTTPPort   string
	ConsulAddr string
	Name       string
}

func NewServerConfig() ServerConfig {
	httpAddr := utils.GetEnv("HTTP_ADDR")
	if httpAddr == "" {
		log.Fatal("HTTP_ADDR environment variable is not set")
	}
	port := utils.GetEnv("HTTP_PORT")
	if port == "" {
		log.Fatal("HTTP_PORT environment variable is not set")
	}
	consulAddr := utils.GetEnv("CONSUL_ADDR")
	if consulAddr == "" {
		log.Fatal("CONSUL_ADDR environment variable is not set")
	}
	name := utils.GetEnv("SERVICE_NAME")
	if name == "" {
		log.Fatal("SERVICE_NAME environment variable is not set")
	}
	return ServerConfig{
		HTTP:       fmt.Sprintf("%s:%s", httpAddr, port),
		HTTPAddr:   httpAddr,
		HTTPPort:   port,
		ConsulAddr: consulAddr,
		Name:       name,
	}
}
