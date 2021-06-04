package converters

import (
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"gopkg.in/yaml.v3"
)

func UnmarshalYamlToObj(yamlText string, out interface{}) error {
	if err := yaml.Unmarshal([]byte(yamlText), out); err != nil {
		logger.Errorf("Failed to unmarshal YAML text to config object. error: %s", err.Error())
		return err
	}
	return nil
}

func YamlToConfigObj(yamlText string) config.AnchorConfig {
	cfg := config.AnchorConfig{}
	if err := UnmarshalYamlToObj(yamlText, &cfg); err != nil {
		logger.Fatalf("Failed to generate config template. error: %s", err.Error())
	}
	return cfg
}
