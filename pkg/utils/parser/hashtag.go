package parser

import (
	"strings"
)

type hashtagParser struct {
	mrkr  string
	sep   string
	pfx   string
	sfx   string
	lines []string
}

func NewHashtagParser() Parser {
	return &hashtagParser{
		mrkr: "EOS",
		sep:  "\n",
		pfx:  "# ",
		sfx:  "\\",
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
			line = x.mrkr
		}

		x.lines = append(x.lines, line)
	}

	return nil
}

func (x *hashtagParser) Find(text string) string {
	var result string
	var found bool

	for _, line := range x.lines {

		if strings.HasPrefix(line, text) {
			result += line + x.sep
			found = true
			continue
		}

		// Append additional input strings
		if found && line != x.mrkr {
			result += line + x.sep
			continue
		}

		// If reached  end of input
		if found && line == "EOS" {
			break
		}
	}

	return result
}
