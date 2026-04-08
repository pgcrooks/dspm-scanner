package dspm_config

import (
	"fmt"

	"github.com/spf13/viper"
)

type dataStoreLocalDBConfig struct {
	Enabled bool
	Path    string `mapstructure:"path"`
}

type dataStoreMemoryConfig struct {
	Enabled bool
}

type DataStoreConfig struct {
	LocalDB dataStoreLocalDBConfig `mapstructure:"local_db"`
	Memory  dataStoreMemoryConfig
}

type finderAWSConfig struct {
	Enabled    bool
	BucketName string `mapstructure:"bucket_name"`
}

type finderLocalConfig struct {
	Enabled bool
	Path    string `mapstructure:"path"`
}

type FinderConfig struct {
	Aws   finderAWSConfig
	Local finderLocalConfig
}

type Config struct {
	DataStore DataStoreConfig
	Finder    FinderConfig
}

func GetConfig(configFile string, configPath string) (Config, error) {
	viper.SetConfigName(configFile)
	viper.AddConfigPath(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("unable to read config file, %v", err)
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("unable to decode into struct, %v", err)
	}

	// Validators
	if !cfg.DataStore.Memory.Enabled && !cfg.DataStore.LocalDB.Enabled {
		return Config{}, fmt.Errorf("no ds enabled")
	}
	if cfg.DataStore.Memory.Enabled && cfg.DataStore.LocalDB.Enabled {
		return Config{}, fmt.Errorf("only one ds may be enabled")
	}

	return cfg, nil
}
