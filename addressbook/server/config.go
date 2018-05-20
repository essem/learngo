package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initConfig() {
	viper.AddConfigPath("./config")
	viper.SetConfigName("default")

	err := viper.ReadInConfig()
	if err != nil {
		log.WithError(err).Fatal("Failed to read default config")
	}

	viper.SetConfigName("local")
	err = viper.MergeInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.WithError(err).Fatal("Failed to read local config")
		}
	}
}
