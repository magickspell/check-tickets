package config

import (
	logg "backend/logger"
	"fmt"
	"os"
	"reflect"
)

type Config struct {
	Host  string
	DbURL string
}

func GetConfig(logger *logg.Logger) *Config {
	cfg := &Config{
		Host:  os.Getenv("GO_HOST"),
		DbURL: os.Getenv("GO_DB_URL"),
	}

	logger.OuteputLog(logg.LogPayload{Info: "start parsing config"})
	keys := reflect.TypeOf(*cfg)
	values := reflect.ValueOf(*cfg)
	for i := 0; i < keys.NumField(); i++ {
		key := keys.Field(i)
		value := values.Field(i)
		fmt.Printf("config item: [%v] = %v (length = %v)\n", key.Name, value, value.Len())
		if value.Len() == 0 {
			logger.OuteputLog(logg.LogPayload{Error: fmt.Errorf("env has empty value")})
		}
	}

	return cfg
}
