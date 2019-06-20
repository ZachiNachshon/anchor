package parser

import "strings"

type hashtagParser struct {
	sep   string
	pfx   string
	sfx   string
	lines []string
}

func NewHashtagParser() Parser {
	return &hashtagParser{
		sep: "\n",
		pfx: "# ",
		sfx: "\\",
	}
}

func (x *hashtagParser) Parse(text string) error {
	lines := strings.Split(text, x.sep)

	for _, line := range lines {
		if strings.HasPrefix(line, x.pfx) {
			line = strings.TrimPrefix(line, x.pfx)
			//line = strings.TrimSuffix(line, x.sfx)
			line = strings.TrimSpace(line)
		} else {
			line = "EOS"
		}

		if line != "" {
			x.lines = append(x.lines, line+"\n")
		}
	}

	return nil
}

func (x *hashtagParser) Find(text string) string {
	var result string
	var found bool

	for _, line := range x.lines {

		if strings.HasPrefix(line, text) {
			result += line + " "
			found = true
			continue
		}

		if found && line != "EOS" {
			result += line + " "
			continue
		}

		if found && line == "EOS" {
			break
		}
	}

	return result
}
