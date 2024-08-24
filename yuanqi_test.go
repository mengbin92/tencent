package tencent

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/mengbin92/tencent/yuanqi"
)

func TestYuanqi(t *testing.T) {
	os.Setenv("TENCENT_YUANQI_TOKEN", "xxxxxxx")
	request := &yuanqi.Request{
		AssistantID: "ZHBpsfvxlrKE",
		UserID:      uuid.NewString(),
		Stream:      false,
		Messages: []yuanqi.RequestMessage{
			{
				Role: "user",
				Content: []yuanqi.MessageContent{
					{
						Type: "text",
						Text: "你是谁",
					},
				},
			},
		},
	}

	client := NewClient()
	yuanqiReq, err := client.requestFactory.Build(context.Background(), http.MethodPost, yuanqiURL, request)
	if err != nil {
		t.Error(err)
	}

	yuanqiResp := &yuanqi.Response{}
	err = client.SendRequest(yuanqiReq, yuanqiResp)
	if err != nil {
		t.Error(err)
	}

	t.Log(yuanqiResp)
}
