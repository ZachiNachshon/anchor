package config

func emptyContext() *Context {
	return &Context{
		Context: &ContextItem{
			Repository: &Repository{
				Remote: &Remote{},
				Local:  &Local{},
			},
		},
	}
}

type AnchorConfig struct {
	Config  *Config `yaml:"config"`
	Author  string  `yaml:"author"`
	License string  `yaml:"license"`
}

type Config struct {
	CurrentContext string     `yaml:"currentContext"`
	Contexts       []*Context `yaml:"contexts"`
	ActiveContext  *Context   `yaml:"-"` // being set programmatically
}

type Context struct {
	Name    string       `yaml:"name"`
	Context *ContextItem `yaml:"context"`
}

type ContextItem struct {
	Repository *Repository `yaml:"repository"`
}

type Repository struct {
	Remote *Remote `yaml:"remote"`
	Local  *Local  `yaml:"local"`
}

type Remote struct {
	Url        string `yaml:"url"`
	Revision   string `yaml:"revision"`
	Branch     string `yaml:"branch"`
	ClonePath  string `yaml:"clonePath"`
	AutoUpdate bool   `yaml:"autoUpdate"`
}

type Local struct {
	Path string `yaml:"path"`
}
