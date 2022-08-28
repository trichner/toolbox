package freshdesk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	httpClient *http.Client
	baseUrl    *url.URL
	token      string
}

type config struct {
	httpClient *http.Client
	baseUrl    string
	token      string
}

type ConfigHandler func(c *config)

func BaseUrl(url string) ConfigHandler {
	return func(c *config) {
		c.baseUrl = url
	}
}

func Token(token string) ConfigHandler {
	return func(c *config) {
		c.token = token
	}
}

func NewClient(handlers ...ConfigHandler) (*Client, error) {
	c := &config{
		httpClient: &http.Client{},
		baseUrl:    "",
	}

	for _, h := range handlers {
		h(c)
	}

	err := validateConfig(c)
	if err != nil {
		return nil, err
	}

	baseUrl, err := url.Parse(c.baseUrl)
	if err != nil {
		return nil, fmt.Errorf("cannot parse baseUrl: %w", err)
	}

	return &Client{
		httpClient: c.httpClient,
		baseUrl:    baseUrl,
		token:      c.token,
	}, nil
}

func validateConfig(c *config) error {
	if strings.TrimSpace(c.baseUrl) == "" {
		return fmt.Errorf("no baseurl configured")
	}

	if strings.TrimSpace(c.token) == "" {
		return fmt.Errorf("no token configured")
	}
	return nil
}

type request struct {
	method string
	url    string
	body   interface{}
}
type response struct {
	statusCode int
	body       []byte
}

func (c *Client) do(r *request) (*response, error) {
	requestData, err := marshalBody(r.body)
	if err != nil {
		return nil, err
	}

	path, err := url.Parse(r.url)
	apiUrl := c.baseUrl.ResolveReference(path)

	req, err := http.NewRequest(r.method, apiUrl.String(), requestData)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.token, "x")

	if requestData != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(res.Body)
	mappedResponse := &response{
		statusCode: res.StatusCode,
		body:       data,
	}

	return mappedResponse, nil
}

func marshalBody(body interface{}) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}
