package config

import (
	"encoding/json"
	"os"

	"github.com/NikitaTsaralov/utils/connectors/logger"
	"github.com/go-playground/validator/v10"
	"sigs.k8s.io/yaml"
)

func LoadJSONConfig(path string, cfg any) {
	jsonContents, err := os.ReadFile(path)
	if err != nil {
		logger.Fatalf("can't read config file %q: %s", path, err)
	}

	parseConfig(jsonContents, cfg)
}

func LoadYAMLConfig(path string, cfg any) {
	yamlContents, err := os.ReadFile(path)
	if err != nil {
		logger.Fatalf("can't read config file %q: %s", path, err)
	}

	jsonContents, err := yaml.YAMLToJSON(yamlContents)
	if err != nil {
		logger.Infof("config content:\n%s", logger.Numerate(string(yamlContents)))
		logger.Fatalf("can't parse config file yaml %q: %s", path, err.Error())
	}

	parseConfig(jsonContents, cfg)
}

func parseConfig(data []byte, config any) {
	err := json.Unmarshal(data, &config)
	if err != nil {
		logger.Fatalf("can't decode config: %s", err.Error())
	}

	err = validator.New().Struct(config)
	if err != nil {
		logger.Fatalf("can't validate config: %s", err.Error())
	}
}
