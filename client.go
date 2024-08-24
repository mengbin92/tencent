package tencent

import (
	"net/http"
	"os"

	"github.com/bytedance/sonic"
	"github.com/mengbin92/tencent/yuanqi"
	"github.com/pkg/errors"
)

type Client struct {
	httpClient     *http.Client
	requestFactory RequestFactory
}

func NewClient() *Client {
	token := os.Getenv("TENCENT_YUANQI_TOKEN")
	return &Client{
		httpClient:     http.DefaultClient,
		requestFactory: yuanqi.NewYuanqi(token),
	}
}

func (c *Client) SendRequest(req *http.Request, v any) error {
	// Send the HTTP request and handle errors
	res, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "send HTTP request error")
	}
	defer res.Body.Close()

	// body,err := io.ReadAll(res.Body)
	// if err != nil{
	// 	return errors.Wrap(err, "read response body error")
	// }
	// fmt.Println(string(body))

	if res.StatusCode != http.StatusOK {
		var errResp ErrResponse
		err = sonic.ConfigDefault.NewDecoder(res.Body).Decode(&errResp)
		if err != nil {
			return errors.Wrap(err, "unmarshal Tencent Response error")
		}
		return errors.Wrap(&errResp, "Tencent API error")
	}

	// Handle the response from openAI
	if v != nil {
		if err := sonic.ConfigDefault.NewDecoder(res.Body).Decode(v); err != nil {
			return errors.Wrap(err, "unmarshal Tencent Response error")
		}
	}
	return nil
}
