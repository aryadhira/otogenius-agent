package tools

import (
	"bufio"
	"os"
	"strings"

	"github.com/aryadhira/otogenius-agent/internal/models"
)

func ReadWikiToolDescription() models.Function {
	return models.Function{
		Name:        "read_wiki",
		Description: "Provides detailed definitions and characteristics for car categories like Sedan, SUV, MPV, and Hatchback, including common features, typical uses, and examples, to aid in accurate categorization based on user descriptions",
		Parameters:  map[string]interface{}{},
	}
}

func ReadWiki() (string, error) {
	file, err := os.Open("docs/wiki.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var sb strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		sb.WriteString(line)
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return sb.String(), nil
}
