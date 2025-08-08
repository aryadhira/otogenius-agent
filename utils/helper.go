package utils

import (
	"fmt"
	"strings"
)

func SplitIntoChunks(text string, maxWords int) []string {
	words := strings.Fields(text)
	var chunks []string
	for i := 0; i < len(words); i += maxWords {
		end := i + maxWords
		if end > len(words) {
			end = len(words)
		}
		chunk := strings.Join(words[i:end], " ")
		chunks = append(chunks, chunk)
	}
	return chunks
}

func ConvertEmbeddingToString(embedding []float32) string {
	strEmb := make([]string, len(embedding))
	for i, v := range embedding {
		strEmb[i] = fmt.Sprintf("%f", v)
	}
	vectorStr := "[" + strings.Join(strEmb, ",") + "]"

	return vectorStr
}
