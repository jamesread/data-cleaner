package config

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"path"
)

type Config struct {
	Csv             *CsvConfig

	Connections     map[string]map[string]string

	Extract         *ExtractConfig
	Transform       *TransformConfig
	Load            *LoadConfig

	Network         *NetworkConfig
}

type ExtractConfig struct {
	ImportDirectory string
}

type LoadConfig struct {
    Destination string
	ColumnMap map[int]string
}

type TransformConfig struct {
	Replacements    *ReplacementsConfig
}

type CsvConfig struct {
	Header bool
}

type ReplacementsConfig struct {
	Exact map[string]string
	Regex map[string]string
}

type NetworkConfig struct {
	BindGrpc  string
	BindRest  string
	BindProxy string
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		config = ReloadConfig()
	}

	return config
}

func newDefaultConfig() *Config {
	return &Config{
		Extract: &ExtractConfig {
			ImportDirectory: "/opt/import/",
		},
		Network: &NetworkConfig{
			BindGrpc:  "127.0.0.1:50051",
			BindRest:  "127.0.0.1:8081",
			BindProxy: "0.0.0.0:8080",
		},
		Csv: &CsvConfig{
			Header: true,
		},
		Transform: &TransformConfig{
			Replacements: &ReplacementsConfig{
				Exact: map[string]string{},
				Regex: map[string]string{},
			},
		},
	}
}

func ReloadConfig() *Config {
	config = newDefaultConfig()

	envconf := os.Getenv("DATA_CLEANER_CONFIG")
	filename := ""

	if envconf != "" {
		filename = envconf
	} else {
		filename = path.Join(os.Getenv("HOME"), ".data-cleaner-config.yaml")
	}

	log.Infof("Loading config from %s", filename)

	file, err := os.ReadFile(filename)

	if err != nil {
		log.Warnf("Could not load config file: %v", err)
		return config
	}

	err = yaml.UnmarshalStrict(file, &config)

	if err != nil {
		log.Warnf("Could not load config file: %v", err)
	}

	return config
}
