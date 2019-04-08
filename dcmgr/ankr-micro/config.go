package ankrmicro

import (
	"fmt"
	"os"
)

type Config struct {
	Rabbitmq     string
	DatabaseHost string
	DatabaseName string
	Listen       string
}

var config Config

func init() {
	config.DatabaseHost = "localhost:27018"
	config.Rabbitmq = "localhost:5672"
	config.Listen = ":50051"
	config = LoadConfigFromEnv()

}

func GetConfig() Config {
	return LoadConfigFromEnv()
}

func LoadConfigFromEnv() Config {
	value := os.Getenv("MICRO_BROKER_ADDRESS")
	if len(value) > 0 {
		config.Rabbitmq = value
	}

	value = os.Getenv("DB_HOST")
	if len(value) > 0 {
		config.DatabaseHost = value
	}

	value = os.Getenv("DB_NAME")

	if len(value) > 0 {
		config.DatabaseName = value
	}

	value = os.Getenv("MICRO_SERVER_ADDRESS")

	if len(value) > 0 {
		config.Listen = value
	}

	return config
}

func (config *Config) Show() {
	fmt.Printf("RabbitMQ : %s \n", config.Rabbitmq)
	fmt.Printf("DB_HOST  : %s  \n", config.DatabaseHost)
	fmt.Printf("DB_Name  : %s  \n", config.DatabaseName)
	fmt.Printf("Listen   : %s \n", config.Listen)

}
