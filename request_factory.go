package tencent

import (
	"context"
	"net/http"
)

type RequestFactory interface {
	Build(ctx context.Context, method, url string, request any) (*http.Request, error)
}
