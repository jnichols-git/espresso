package adapter

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/cloudretic/matcha/pkg/adapter"
)

func LambdaAdapter(s http.Handler) func(ctx context.Context, in events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return func(ctx context.Context, in events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		adp := &MatchaV2Adapter{}
		w, req, out, err := adp.Adapt(in)
		if err != nil {
			return *out, err
		}
		s.ServeHTTP(w, req)
		fmt.Println(out, err)
		return *out, nil
	}
}

type MatchaV2Adapter struct {
	adapter.Adapter[events.APIGatewayV2HTTPRequest, *events.APIGatewayV2HTTPResponse]
}

type ResponseWriter struct {
	target *events.APIGatewayV2HTTPResponse
}

func (rw *ResponseWriter) Header() http.Header {
	return http.Header(rw.target.MultiValueHeaders)
}

func (rw *ResponseWriter) Write(body []byte) (int, error) {
	rw.target.Body = string(body)
	if rw.target.StatusCode < 200 {
		rw.WriteHeader(http.StatusOK)
	}
	return len(body), nil
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.target.StatusCode = statusCode
}

func (mv2a *MatchaV2Adapter) Adapt(pr events.APIGatewayV2HTTPRequest) (http.ResponseWriter, *http.Request, *events.APIGatewayV2HTTPResponse, error) {
	resp := &events.APIGatewayV2HTTPResponse{}
	w := &ResponseWriter{
		target: resp,
	}
	req, err := http.NewRequest("", "", nil)
	if err != nil {
		return nil, nil, resp, err
	}
	// Method
	req.Method = pr.RequestContext.HTTP.Method
	// Body
	if len(pr.Body) > 0 {
		r := io.Reader(strings.NewReader(pr.Body))
		if pr.IsBase64Encoded {
			r = base64.NewDecoder(base64.StdEncoding, r)
		}
		req.Body = io.NopCloser(r)
		req.ContentLength = int64(len(pr.Body))
	}
	// URL values
	url := req.URL
	url.Path = pr.RawPath
	url.RawQuery = pr.RawQueryString
	// Headers
	head := req.Header
	for k, v := range pr.Headers {
		head.Add(k, v)
	}
	return w, req, resp, nil
}
