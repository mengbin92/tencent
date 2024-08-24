package yuanqi

import (
	"bytes"
	"context"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
)

type File struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type MessageContent struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
	FileURL File   `json:"file_url"`
}

type RequestMessage struct {
	Role    string           `json:"role"`
	Content []MessageContent `json:"content"`
}

type Request struct {
	AssistantID string           `json:"assistant_id"`
	UserID      string           `json:"user_id"`
	Stream      bool             `json:"stream"`
	Messages    []RequestMessage `json:"messages"`
}

type Yuanqi struct {
	token string
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Step struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCallID string     `json:"tool_call_id"`
	ToolCalls  []ToolCall `json:"tool_calls"`
	Usage      Usage      `json:"usage"`
	TimeCost   int        `json:"time_cost"`
}

type ToolCall struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Type      string `json:"type"`
	Arguments string `json:"arguments"`
}

type ResponseMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Steps   []Step `json:"steps"`
}

type Delta struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCallID string     `json:"tool_call_id"`
	ToolCalls  []ToolCall `json:"tool_calls"`
	TimeCost   int        `json:"time_cost"`
}

type Choice struct {
	Index           int             `json:"index"`
	FinishReason    string          `json:"finish_reason"`
	ModerationLevel string          `json:"moderation_level"`
	Message         ResponseMessage `json:"message"`
	Delta           Delta           `json:"delta"`
}

type Response struct {
	ID          string   `json:"id"`
	Created     int      `json:"created"`
	Choices     []Choice `json:"choices"`
	AssistantID string   `json:"assistant_id"`
	Usage       Usage    `json:"usage"`
}

func NewYuanqi(token string) *Yuanqi {
	return &Yuanqi{token: token}
}

// The Build method of the httpRequestFactory struct builds and returns an http.Request object given the specified parameters.
// It takes in parameters:
// - ctx, of type context.Context which is the execution context of the function.
// - method, of type string which represents the HTTP method to be used for the request.
// - url, of type string which is the URL that the request should be sent to.
// - request, of type any, which specifies the data to be sent in the request body. It uses Sonic to marshal the request data into bytes.
// It returns:
// - A pointer to the created *http.Request object which may consist of the specified context, method, url and request data.
// - err, an error object which will contain any errors that occurred during assembly of http.Request object.
func (y *Yuanqi) Build(ctx context.Context, method, url string, request any) (req *http.Request, err error) {
	// Check if request data is nil.
	if request == nil {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	} else {
		//Marshal the request data using Sonic marshaler.
		requestBytes, err := sonic.Marshal(request)
		if err != nil {
			return nil, errors.Wrap(err, "marshal request error")
		}
		//Use the marshaled bytes to create a new request with the specified context, method and URL.
		req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(requestBytes))
		if err != nil {
			return nil, errors.Wrap(err, "create http request error")
		}
	}

	// Add the required headers to the request.
	req.Header.Set("Authorization", "Bearer "+y.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Source", "openapi")
	return
}
