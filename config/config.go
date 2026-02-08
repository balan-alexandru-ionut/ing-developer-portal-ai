package config

import (
	"ai-test/logger"
	"ai-test/util"
	"ai-test/util/level"

	"github.com/spf13/viper"
)

var log = logger.NewLogger()

var C = &Config{}

func ReadConfigFile() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	util.HandleError("Could not read config file: %v", err, level.FATAL)

	err = viper.Unmarshal(C)
	util.HandleError("Could not read config file: %v", err, level.FATAL)
}
