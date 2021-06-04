package parser

type PromptItems struct {
	Items []PromptItem `yaml:"promptItems"`
}

type PromptItem struct {
	Id    string `yaml:"id"`
	Title string `yaml:"title"`
	File  string `yaml:"file"`
}

type Parser interface {
	Parse(yamlText string) (*PromptItems, error)
	//Find(text string) string
}
