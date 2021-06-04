package kits

import (
	"github.com/ZachiNachshon/anchor/logger"
)

var configYamlTemplate = `
{{ if .Author }}{{ .Author }}{{ end }}
{{ if .License }}{{ .License }}{{ end }}
config:
  repositoryFiles:
    remote:
      url: {{ if .RemoteRepoUrl }}{{ .RemoteRepoUrl }}{{ else }} https://github.com/ZachiNachshon/dummy-repo.git {{ end }}      
      revision: {{ if .RemoteRepoRevision }}{{ .RemoteRepoRevision }}{{ else }} a123456789 {{ end }} 
      branch: {{ if .RemoteRepoBranch }}{{ .RemoteRepoBranch }}{{ else }} some-branch {{ end }}
      localPath: {{ if .RemoteRepoLocalPath }}{{ .RemoteRepoLocalPath }}{{ else }} /path/to/cloned/repo {{ end }}
    local:
      path: {{ if .LocalRepoPath }}{{ .LocalRepoPath }}{{ else }} /path/to/local/folder {{ end }}
`

type TemplateItems struct {
	Author              string
	License             string
	RemoteRepoUrl       string
	RemoteRepoRevision  string
	RemoteRepoBranch    string
	RemoteRepoLocalPath string
	LocalRepoPath       string
}

var GetDefaultTestConfigText = func() string {
	var yamlConfig, err = TemplateToText(configYamlTemplate, nil)
	if err != nil {
		// Stop testing process since tests environment has an issue
		logger.Fatalf("Failed to generate config template. error: %s", err)
	}
	return yamlConfig
}

var GetCustomTestConfigText = func(items TemplateItems) string {
	var yamlConfig, err = TemplateToText(configYamlTemplate, &items)
	if err != nil {
		// Stop testing process since tests environment has an issue
		logger.Fatalf("Failed to generate config template with substitutions. error: %s", err)
	}
	return yamlConfig
}
