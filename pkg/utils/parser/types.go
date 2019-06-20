package parser

type Parser interface {
	Parse(text string) error
	Find(text string) string
}
