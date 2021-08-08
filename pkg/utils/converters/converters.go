package converters

import (
	"fmt"
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
