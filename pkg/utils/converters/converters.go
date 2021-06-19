package converters

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"gopkg.in/yaml.v3"
)

func UnmarshalYamlToObj(yamlText string, out interface{}) error {
	if err := yaml.Unmarshal([]byte(yamlText), out); err != nil {
		return fmt.Errorf("failed to unmarshal YAML text to config object. error: %s", err.Error())
	}
	return nil
}

func UnmarshalObjToYaml(obj interface{}) (string, error) {
	if out, err := yaml.Marshal(obj); err != nil {
		return "", fmt.Errorf("failed to unmarshal YAML text to config object. error: %s", err.Error())
	} else {
		return string(out), nil
	}
}

func YamlToConfigObj(yamlText string) config.AnchorConfig {
	cfg := config.AnchorConfig{}
	if err := UnmarshalYamlToObj(yamlText, &cfg); err != nil {
		logger.Fatalf("Failed to generate config template. error: %s", err.Error())
	}
	return cfg
}

func ConfigObjToYaml(cfg config.AnchorConfig) (string, error) {
	if out, err := UnmarshalObjToYaml(cfg); err != nil {
		logger.Errorf("Failed to generate config template. error: %s", err.Error())
		return "", nil
	} else {
		return out, nil
	}
}
