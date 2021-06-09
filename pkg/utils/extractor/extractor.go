package extractor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"github.com/ZachiNachshon/anchor/pkg/utils/parser"
	"io/ioutil"
)

type extractor struct{}

func New() Extractor {
	return &extractor{}
}

func (e *extractor) ExtractPromptItems(instructionsPath string, p parser.Parser) (*parser.Instructions, error) {
	if !ioutils.IsValidPath(instructionsPath) {
		return nil, fmt.Errorf("invalid instructions path. path: %s", instructionsPath)
	}

	if contentByte, err := ioutil.ReadFile(instructionsPath); err != nil {
		return nil, err
	} else {
		var yamlText = string(contentByte)

		if items, err := p.Parse(yamlText); err != nil {
			return nil, err
		} else {
			return items, nil
		}
	}
}

//
//func replaceDockerCommandPlaceholders(content string, path string) string {
//	// In case the Dockerfile is referenced by a custom path
//	if strings.Contains(content, "/Dockerfile") {
//		return content
//	} else {
//		content = strings.ReplaceAll(content, "Dockerfile", path)
//		return content
//	}
//}
//
//func missingDockerCmdMsg(dirname string) string {
//	//return fmt.Sprintf("Missing '%v' on %v Dockerfile instructions\n", dirname)
//	return ""
//}
//
//func (e *extractor) ManifestCmd(identifier string) (string, error) {
//	if manifestFilePath, err := locator.DirLocator.Manifest(identifier); err != nil {
//		return "", err
//	} else {
//		var result = ""
//		if contentByte, err := ioutil.ReadFile(manifestFilePath); err != nil {
//			return "", err
//		} else {
//			// Load .env file
//			//config.LoadEnvVars(identifier)
//
//			var manifestContent = string(contentByte)
//
//			p := parser.NewHashtagParser()
//			if err := p.Parse(manifestContent); err != nil {
//				return "", errors.Errorf("Failed to parse: %v, err: %v", manifestFilePath, err.Error())
//			}
//
//			//if result = p.Find(string(manifestCommand)); result != "" {
//			//	result = strings.TrimSuffix(result, "\n")
//			//}
//			//
//			//if len(result) == 0 {
//			//	return "", errors.Errorf(missingManifestCmdMsg(manifestCommand, identifier))
//			//}
//		}
//		return result, nil
//	}
//}
//
//func missingManifestCmdMsg(dirname string) string {
//	//return fmt.Sprintf("Missing '%v' on %v K8s manifest\n", manifestCommand, dirname)
//	return ""
//}
//
//func (e *extractor) ManifestContent(identifier string) (bool, error) {
//	if _, err := locator.DirLocator.Manifest(identifier); err != nil {
//		return false, err
//	} else {
//		return false, nil
//		//if contentByte, err := ioutil.ReadFile(manifestFilePath); err != nil {
//		//	return false, err
//		//} else {
//		//	// Load .env file
//		//	config.LoadEnvVars(identifier)
//		//
//		//	var manifestContent = string(contentByte)
//		//	stateful := strings.Contains(manifestContent, string(manifestCommand))
//		//	return stateful, nil
//		//}
//	}
//}
