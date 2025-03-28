package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
)

var (
	ErrInvalidResponse = errors.New("invalid response")
	ErrRequestFailed   = errors.New("request failed")
)

// HTTPClientConfig holds the configuration for the HTTP client
type HTTPClientConfig struct {
	MaxConns               int
	MaxIdleConnDuration    time.Duration
	ReadTimeout            time.Duration
	WriteTimeout           time.Duration
	MaxResponseBodySize    int
	DisablePathNormalizing bool
}

// DefaultConfig returns the default HTTP client configuration
func DefaultConfig() HTTPClientConfig {
	return HTTPClientConfig{
		MaxConns:               1000,
		MaxIdleConnDuration:    10 * time.Second,
		ReadTimeout:            5 * time.Second,
		WriteTimeout:           5 * time.Second,
		MaxResponseBodySize:    10 * 1024 * 1024, // 10MB
		DisablePathNormalizing: true,
	}
}

// HTTPClient is a struct that wraps the fasthttp.Client
type HTTPClient struct {
	client *fasthttp.Client
}

// NewHTTPClient creates a new instance of HTTPClient with default configuration
func NewHTTPClient() *HTTPClient {
	config := DefaultConfig()
	return NewHTTPClientWithConfig(config)
}

// NewHTTPClientWithConfig creates a new HTTPClient with custom configuration
func NewHTTPClientWithConfig(config HTTPClientConfig) *HTTPClient {
	return &HTTPClient{
		client: &fasthttp.Client{
			MaxConnsPerHost:               config.MaxConns,
			MaxIdleConnDuration:           config.MaxIdleConnDuration,
			ReadTimeout:                   config.ReadTimeout,
			WriteTimeout:                  config.WriteTimeout,
			MaxResponseBodySize:           config.MaxResponseBodySize,
			DisablePathNormalizing:        config.DisablePathNormalizing,
			NoDefaultUserAgentHeader:      true,
			DisableHeaderNamesNormalizing: true,
		},
	}
}

// doRequest executes an HTTP request and handles the response
func (hc *HTTPClient) doRequest(ctx context.Context, req *fasthttp.Request, resp *fasthttp.Response) ([]byte, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	err := hc.client.DoTimeout(req, resp, 30*time.Second)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequestFailed, err)
	}

	statusCode := resp.StatusCode()
	if statusCode < 200 || statusCode >= 300 {
		return nil, fmt.Errorf("%w: status code %d", ErrInvalidResponse, statusCode)
	}

	return resp.Body(), nil
}

// executeRequest is a helper function to reduce code duplication
func (hc *HTTPClient) executeRequest(ctx context.Context, method, url string, headers map[string]string, body []byte) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(method)
	req.SetRequestURI(url)

	if body != nil {
		req.SetBody(body)
	}

	// Set custom headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	return hc.doRequest(ctx, req, resp)
}

// Get sends a GET request
func (hc *HTTPClient) Get(url string) ([]byte, error) {
	return hc.GetWithContext(context.Background(), url, nil)
}

// GetWithContext sends a GET request with context and optional headers
func (hc *HTTPClient) GetWithContext(ctx context.Context, url string, headers map[string]string) ([]byte, error) {
	return hc.executeRequest(ctx, fasthttp.MethodGet, url, headers, nil)
}

// Post sends a POST request
func (hc *HTTPClient) Post(url string, body []byte) ([]byte, error) {
	return hc.PostWithContext(context.Background(), url, nil, body)
}

// PostWithContext sends a POST request with context and optional headers
func (hc *HTTPClient) PostWithContext(ctx context.Context, url string, headers map[string]string, body []byte) ([]byte, error) {
	return hc.executeRequest(ctx, fasthttp.MethodPost, url, headers, body)
}

// Put sends a PUT request
func (hc *HTTPClient) Put(url string, body []byte) ([]byte, error) {
	return hc.PutWithContext(context.Background(), url, nil, body)
}

// PutWithContext sends a PUT request with context and optional headers
func (hc *HTTPClient) PutWithContext(ctx context.Context, url string, headers map[string]string, body []byte) ([]byte, error) {
	return hc.executeRequest(ctx, fasthttp.MethodPut, url, headers, body)
}

// Delete sends a DELETE request
func (hc *HTTPClient) Delete(url string) ([]byte, error) {
	return hc.DeleteWithContext(context.Background(), url, nil)
}

// DeleteWithContext sends a DELETE request with context and optional headers
func (hc *HTTPClient) DeleteWithContext(ctx context.Context, url string, headers map[string]string) ([]byte, error) {
	return hc.executeRequest(ctx, fasthttp.MethodDelete, url, headers, nil)
}

// JSON helper methods

// GetJSON sends a GET request and unmarshals the JSON response
func (hc *HTTPClient) GetJSON(url string, result interface{}) error {
	return hc.GetJSONWithContext(context.Background(), url, nil, result)
}

// GetJSONWithContext sends a GET request with context and unmarshals the JSON response
func (hc *HTTPClient) GetJSONWithContext(ctx context.Context, url string, headers map[string]string, result interface{}) error {
	resp, err := hc.GetWithContext(ctx, url, headers)
	if err != nil {
		return err
	}

	return json.Unmarshal(resp, result)
}

// PostJSON sends a POST request with JSON body and appropriate headers
func (hc *HTTPClient) PostJSON(url string, data interface{}) ([]byte, error) {
	headers := map[string]string{"Content-Type": "application/json"}
	return hc.PostJSONWithContext(context.Background(), url, headers, data)
}

// PostJSONWithContext sends a POST request with JSON body, context and headers
func (hc *HTTPClient) PostJSONWithContext(ctx context.Context, url string, headers map[string]string, data interface{}) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"

	return hc.PostWithContext(ctx, url, headers, body)
}

// PostJSONAndParseResponse sends POST with JSON and parses JSON response
func (hc *HTTPClient) PostJSONAndParseResponse(url string, data interface{}, result interface{}) error {
	resp, err := hc.PostJSON(url, data)
	if err != nil {
		return err
	}

	return json.Unmarshal(resp, result)
}

// File operations

// DownloadFile downloads a file from the given URL and saves it to the destination path
func (hc *HTTPClient) DownloadFile(url string, destPath string) error {
	return hc.DownloadFileWithContext(context.Background(), url, destPath, nil)
}

// DownloadFileWithContext downloads a file with context and optional headers
func (hc *HTTPClient) DownloadFileWithContext(ctx context.Context, url string, destPath string, headers map[string]string) error {
	body, err := hc.GetWithContext(ctx, url, headers)
	if err != nil {
		return err
	}

	return os.WriteFile(destPath, body, 0644)
}

// DownloadFileWithProgress downloads a file with progress reporting
func (hc *HTTPClient) DownloadFileWithProgress(url string, destPath string, progressCb func(downloaded, total int64)) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := hc.client.Do(req, resp); err != nil {
		return err
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("%w: status code %d", ErrInvalidResponse, resp.StatusCode())
	}

	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer file.Close()

	contentLength := resp.Header.ContentLength()

	var downloaded int64
	buffer := make([]byte, 32*1024) // 32KB buffer

	reader := bytes.NewReader(resp.Body())
	for {
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		if _, err := file.Write(buffer[:n]); err != nil {
			return err
		}

		downloaded += int64(n)
		if progressCb != nil {
			progressCb(downloaded, int64(contentLength))
		}
	}

	return nil
}
