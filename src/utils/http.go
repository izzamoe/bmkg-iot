package utils

import (
	"os"

	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
)

// HTTPClient is a struct that wraps the fasthttp.Client
type HTTPClient struct {
	client *fasthttp.Client
}

// NewHTTPClient creates a new instance of HTTPClient
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &fasthttp.Client{},
	}
}

// Get sends a GET request to the specified URL and returns the response body as []byte
func (hc *HTTPClient) Get(url string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := hc.client.Do(req, resp); err != nil {
		return nil, err
	}

	return resp.Body(), nil
}

// Post sends a POST request to the specified URL with the given body and returns the response body as []byte
func (hc *HTTPClient) Post(url string, body []byte) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetBody(body)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := hc.client.Do(req, resp); err != nil {
		return nil, err
	}

	return resp.Body(), nil
}

// Put sends a PUT request to the specified URL with the given body and returns the response body as []byte
func (hc *HTTPClient) Put(url string, body []byte) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodPut)
	req.SetBody(body)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := hc.client.Do(req, resp); err != nil {
		return nil, err
	}

	return resp.Body(), nil
}

// Delete sends a DELETE request to the specified URL and returns the response body as []byte
func (hc *HTTPClient) Delete(url string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodDelete)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := hc.client.Do(req, resp); err != nil {
		return nil, err
	}

	return resp.Body(), nil
}

// GetWithHeaders sends a GET request with custom headers
func (hc *HTTPClient) GetWithHeaders(url string, headers map[string]string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)

	// Set custom headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := hc.client.Do(req, resp); err != nil {
		return nil, err
	}

	return resp.Body(), nil
}

// PostJSON sends a POST request with JSON body and appropriate headers
func (hc *HTTPClient) PostJSON(url string, data interface{}) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")
	req.SetBody(body)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := hc.client.Do(req, resp); err != nil {
		return nil, err
	}

	return resp.Body(), nil
}

// DownloadFile downloads a file from the given URL and saves it to the destination path
func (hc *HTTPClient) DownloadFile(url string, destPath string) error {
	body, err := hc.Get(url)
	if err != nil {
		return err
	}

	return os.WriteFile(destPath, body, 0644)
}
