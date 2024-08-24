package yuanqi

import (
	"bytes"
	"context"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
)

// File 结构体表示文件信息，包含类型和URL
type File struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// MessageContent 结构体表示消息内容，包含类型、文本和文件URL
type MessageContent struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
	FileURL File   `json:"file_url"`
}

// RequestMessage 结构体表示请求消息，包含角色和内容列表
type RequestMessage struct {
	Role    string           `json:"role"`
	Content []MessageContent `json:"content"`
}

// Request 结构体表示请求，包含助手ID、用户ID、是否流式传输和消息列表
type Request struct {
	AssistantID string           `json:"assistant_id"`
	UserID      string           `json:"user_id"`
	Stream      bool             `json:"stream"`
	Messages    []RequestMessage `json:"messages"`
}

// Yuanqi 结构体表示Yuanqi对象，包含token
type Yuanqi struct {
	token string
}

// Usage 结构体表示使用情况，包含提示tokens、完成tokens和总tokens
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Step 结构体表示步骤，包含角色、内容、工具调用ID、工具调用列表、使用情况和时间成本
type Step struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCallID string     `json:"tool_call_id"`
	ToolCalls  []ToolCall `json:"tool_calls"`
	Usage      Usage      `json:"usage"`
	TimeCost   int        `json:"time_cost"`
}

// ToolCall 结构体表示工具调用，包含ID、类型和函数
type ToolCall struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// Function 结构体表示函数，包含名称、描述、类型和参数
type Function struct {
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Type      string `json:"type"`
	Arguments string `json:"arguments"`
}

// ResponseMessage 结构体表示响应消息，包含角色、内容和步骤列表
type ResponseMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Steps   []Step `json:"steps"`
}

// Delta 结构体表示增量，包含角色、内容、工具调用ID、工具调用列表和时间成本
type Delta struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCallID string     `json:"tool_call_id"`
	ToolCalls  []ToolCall `json:"tool_calls"`
	TimeCost   int        `json:"time_cost"`
}

// Choice 结构体表示选择，包含索引、完成原因、审核级别、消息和增量
type Choice struct {
	Index           int             `json:"index"`
	FinishReason    string          `json:"finish_reason"`
	ModerationLevel string          `json:"moderation_level"`
	Message         ResponseMessage `json:"message"`
	Delta           Delta           `json:"delta"`
}

// Response 结构体表示响应，包含ID、创建时间、选择列表、助手ID和使用情况
type Response struct {
	ID          string   `json:"id"`
	Created     int      `json:"created"`
	Choices     []Choice `json:"choices"`
	AssistantID string   `json:"assistant_id"`
	Usage       Usage    `json:"usage"`
}

// NewYuanqi 函数创建并返回一个新的Yuanqi对象
func NewYuanqi(token string) *Yuanqi {
	return &Yuanqi{token: token}
}

// Build 方法构建并返回一个http.Request对象，给定指定的参数
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
	req.Header.Set("X-source", "openapi")
	return
}
