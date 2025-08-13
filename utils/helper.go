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

func InterfaceToString(params interface{}) string {
	if params == nil {
		return ""
	}

	return params.(string)
}

func InterfaceToInt(params interface{}) int {
	if params == nil {
		return 0
	}

	return params.(int)
}

func InterfaceToFloat(params interface{}) float64 {
	if params == nil {
		return 0.0
	}

	return params.(float64)
}
