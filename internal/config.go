package scanner_int

import (
	"fmt"

	"github.com/spf13/viper"
)

func GetConfig() (Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("unable to read config file, %v", err)
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("unable to decode into struct, %v", err)
	}
	return cfg, nil
}
