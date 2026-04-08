package dspm_config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigGood(t *testing.T) {
	cfg, err := GetConfig("config-good", "../../test/data/")
	assert.NoError(t, err)

	expectedCfg := Config{
		DataStore: DataStoreConfig{
			LocalDB: dataStoreLocalDBConfig{
				Enabled: false,
				Path:    "dspm.db",
			},
			Memory: dataStoreMemoryConfig{
				Enabled: true,
			},
		},
		Scraper: ScraperConfig{
			Aws: scraperAWSConfig{
				Enabled:    false,
				BucketName: "pgcrooks-dspm",
			},
			Local: scraperLocalConfig{
				Enabled: true,
				Path:    "test/data/",
			},
		},
	}

	assert.Equal(t, expectedCfg, cfg)
}

func TestConfigNoDS(t *testing.T) {
	_, err := GetConfig("config-no-ds", "../../test/data/")
	assert.EqualError(t, err, "no ds enabled")
}
