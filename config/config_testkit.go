package config

import (
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/templates"
)

var configYamlTemplate = `
{{ if .Author }}{{ .Author }}{{ end }}
{{ if .License }}{{ .License }}{{ end }}
config:
 currentContext: {{ if .CurrentContext }}{{ .CurrentContext }}{{ else }} 1st-anchorfiles {{ end }}
 contexts:
  - name: {{ if .FirstContextName }}{{ .FirstContextName }}{{ else }} 1st-anchorfiles {{ end }}
    context:
     repository:
      remote:
       url: {{ if .FirstContextRemoteRepoUrl }}{{ .FirstContextRemoteRepoUrl }}{{ else }} https://github.com/ZachiNachshon/dummy-repo-1.git {{ end }}      
       revision: {{ if .FirstContextRemoteRepoRevision }}{{ .FirstContextRemoteRepoRevision }}{{ else }} a123456789 {{ end }} 
       branch: {{ if .FirstContextRemoteRepoBranch }}{{ .FirstContextRemoteRepoBranch }}{{ else }} some-branch-1 {{ end }}
       clonePath: {{ if .FirstContextClonePath }}{{ .FirstContextClonePath }}{{ end }}
      local:
       path: {{ if .FirstContextLocalRepoPath }}{{ .FirstContextLocalRepoPath }}{{ else }} /path/to/ctx1/local/folder {{ end }}
  - name: {{ if .SecondContextName }}{{ .SecondContextName }}{{ else }} 2nd-anchorfiles {{ end }}
    context:
     repository:
      remote:
       url: {{ if .SecondContextRemoteRepoUrl }}{{ .SecondContextRemoteRepoUrl }}{{ else }} https://github.com/ZachiNachshon/dummy-repo-2.git {{ end }}      
       revision: {{ if .SecondContextRemoteRepoRevision }}{{ .SecondContextRemoteRepoRevision }}{{ else }} 987654321a {{ end }} 
       branch: {{ if .SecondContextRemoteRepoBranch }}{{ .SecondContextRemoteRepoBranch }}{{ else }} some-branch-2 {{ end }}
       clonePath: {{ if .SecondContextClonePath }}{{ .SecondContextClonePath }}{{ end }}
     local:
      path: {{ if .SecondContextLocalRepoPath }}{{ .SecondContextLocalRepoPath }}{{ else }} /path/to/ctx2/local/folder {{ end }}
`

type TemplateItems struct {
	Author                          string
	License                         string
	CurrentContext                  string
	FirstContextName                string
	FirstContextClonePath           string
	FirstContextRemoteRepoUrl       string
	FirstContextRemoteRepoRevision  string
	FirstContextRemoteRepoBranch    string
	FirstContextLocalRepoPath       string
	SecondContextName               string
	SecondContextClonePath          string
	SecondContextRemoteRepoUrl      string
	SecondContextRemoteRepoRevision string
	SecondContextRemoteRepoBranch   string
	SecondContextLocalRepoPath      string
}

var GetDefaultTestConfigText = func() string {
	var yamlConfig, err = templates.TemplateToText(configYamlTemplate, nil)
	if err != nil {
		// Stop testing process since tests environment has an issue
		logger.Fatalf("Failed to generate config template. error: %s", err)
	}
	return yamlConfig
}

var GetCustomTestConfigText = func(items TemplateItems) string {
	var yamlConfig, err = templates.TemplateToText(configYamlTemplate, &items)
	if err != nil {
		// Stop testing process since tests environment has an issue
		logger.Fatalf("Failed to generate config template with substitutions. error: %s", err)
	}
	return yamlConfig
}
