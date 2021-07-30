package templates

import (
	"bytes"
	"encoding/json"
	"github.com/ZachiNachshon/anchor/logger"
	"text/template"
)

func TemplateToText(templateText string, templateData interface{}) (string, error) {
	if out, err := parseConfigTemplate(templateText, templateData); err != nil {
		logger.Errorf("Failed to parse template to text. error: %s", err.Error())
		return "", err
	} else {
		return out, err
	}
}

func structToMap(input interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	inputBytes, _ := json.Marshal(input)
	err := json.Unmarshal(inputBytes, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func parseConfigTemplate(templateText string, data interface{}) (string, error) {
	if t, err := template.New("anchorTemplates").Funcs(FuncMap).Parse(templateText); err != nil {
		return "", err
	} else {
		templateItemsMap, convertErr := structToMap(data)
		if convertErr != nil {
			return "", convertErr
		}

		var buffer bytes.Buffer
		if err := t.Execute(&buffer, templateItemsMap); err != nil {
			return "", err
		} else {
			return buffer.String(), nil
		}
	}
}
