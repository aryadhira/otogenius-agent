package models

type Message struct {
	Role      string         `json:"role"`
	Content   interface{}    `json:"content"`
	ToolCalls []FunctionCall `json:"tool_calls"`
}

type LlmRequest struct {
	Messages       []Message              `json:"messages"`
	Temperature    float64                `json:"temperature"`
	MaxTokens      int                    `json:"max_tokens"`
	Stream         bool                   `json:"stream"`
	Tools          []Tool                 `json:"tools"`
	ToolChoice     string                 `json:"tool_choice"`
	ResponseFormat map[string]string      `json:"response_format"`
	JsonSchema     map[string]interface{} `json:"json_schema"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	LogProbs     *int    `json:"logprobs"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type LlmResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}