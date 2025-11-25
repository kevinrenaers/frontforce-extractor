package main

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"renaers.be/frontforce/internal/configuration"
	"renaers.be/frontforce/internal/frontforce"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var config configuration.Configuration

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}

	frontforce, err := frontforce.NewFrontforce(config)
	if err != nil {
		panic(fmt.Errorf("fatal setting up frontforce: %w", err))
	}
	frontforce.Start()
}
