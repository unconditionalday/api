package parser

import (
	"regexp"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(content string) string {
	// remove bloated text from the content
	content = removeBloat(content)

	// remove html tags from the content
	content = removeHTML(content)

	// remove new lines from the content
	content = removeNewLines(content)

	return content
}

// function to remove bloated text from the content using regex
func removeBloat(content string) string {
	// regex to match the text to remove
	re := regexp.MustCompile(`\[(.*?)\]`)

	// common bloated text to remove
	bloat := []string{"Continue reading...", " Read more...", " Read more", " Read more at", " Read more at:"}

	// transform the bloated text to regex
	for _, b := range bloat {
		re = regexp.MustCompile(re.String() + "|" + regexp.QuoteMeta(b))
	}

	return re.ReplaceAllString(content, "")
}

// function to remove html tags from the content using regex
func removeHTML(content string) string {
	// regex to remove html tags
	re := regexp.MustCompile(`<.*?>`)

	return re.ReplaceAllString(content, "")
}

// function to remove new lines from the content using regex
func removeNewLines(content string) string {
	// regex to match the text to remove
	re := regexp.MustCompile(`\r?\n`)

	return re.ReplaceAllString(content, "")
}
