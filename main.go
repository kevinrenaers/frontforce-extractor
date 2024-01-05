package main

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"renaers.be/frontforce/internal"
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
	frontforce, err := internal.NewFrontforce()
	if err != nil {
		panic(fmt.Errorf("fatal setting up frontforce: %w", err))
	}
	frontforce.StartUpdater()
}
