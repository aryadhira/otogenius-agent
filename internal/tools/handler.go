package tools

import (
	"fmt"
	"reflect"

	"github.com/aryadhira/otogenius-agent/internal/models"
)

func GetToolDispatcher() map[string]interface{} {
	return map[string]interface{}{
		"read_wiki": ReadWiki,
	}
}

func RegisterTools() []models.Tool {
	return []models.Tool{
		{
			Type:     "function",
			Function: ReadWikiToolDescription(),
		},
	}
}

func ToolCalling(message models.Message) ([]models.Message, error) {
	toolCallHistory := []models.Message{}

	for _, toolCall := range message.ToolCalls {
		fmt.Printf("LLM requested call function: %s\n", toolCall.Function.Name)
		toolCallHistory = append(toolCallHistory, models.Message{
			Role:      "assistant",
			ToolCalls: []models.FunctionCall{toolCall},
		})
		toolDispatcher := GetToolDispatcher()

		fn := toolDispatcher[toolCall.Function.Name]
		if fn == nil {
			return nil, fmt.Errorf("tools %s not available", toolCall.Function.Name)
		}

		fnVal := reflect.ValueOf(fn)
		if fnVal.Kind() != reflect.Func {
			return nil, fmt.Errorf("tools %s not a function", toolCall.Function.Name)
		}

		fnCallRes := make([]reflect.Value, 0)
		if fnVal.Type().NumIn() > 0 {
			in, err := functionInputParser(toolCall.Function.Name, fnVal.Type(), toolCall)
			if err != nil {
				return nil, fmt.Errorf("can't parse parameter %s", toolCall.Function.Name)
			}

			fnCallRes = fnVal.Call(in)
		} else {
			emptyParam := make([]reflect.Value, 0)
			fnCallRes = fnVal.Call(emptyParam)
		}

		fnRes, fnErr := functionOutputParser(fnCallRes, fnVal.Type())

		if fnErr != nil {
			fmt.Printf("Error executing function %s: %v\n", toolCall.Function.Name, fnErr)
			toolCallHistory = append(toolCallHistory, models.Message{
				Role:    "user",
				Content: fmt.Sprintf("Error calling %s: %v", toolCall.Function.Name, fnErr),
			})
		} else {
			toolCallHistory = append(toolCallHistory, models.Message{
				Role:      "tool",
				Content:   fnRes,
				ToolCalls: []models.FunctionCall{},
			})
		}

	}

	return toolCallHistory, nil
}

func functionInputParser(funcName string, fnType reflect.Type, toolCall models.FunctionCall) ([]reflect.Value, error) {
	var inputParam []reflect.Value
	var err error

	switch funcName {
	default:
		err = fmt.Errorf("function name is empty")

	}

	return inputParam, err
}

func functionOutputParser(fnRes []reflect.Value, fnType reflect.Type) (string, error) {
	result := ""
	var err error

	numOut := fnType.NumOut()
	for i := 0; i < numOut; i++ {
		returnType := fnType.Out(i)
		if returnType.Kind() == fnRes[i].Kind() && !fnRes[i].Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			result = fnRes[i].String()
		} else if fnRes[i].Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) && !fnRes[i].IsNil() {
			err = fnRes[i].Interface().(error)
		}
	}

	return result, err
}
