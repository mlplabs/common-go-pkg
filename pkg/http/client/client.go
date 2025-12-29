package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	domainErrors "git.puls.ru/3pl/xpl/common-utils/pkg/domain/errors"
	httpErrors "git.puls.ru/3pl/xpl/common-utils/pkg/http/errors"
	"io"
	"net/http"
)

const Token string = "TOKEN"

type RequestParams struct {
	ProxyToken     bool
	RequestBody    interface{}
	ResponseBody   interface{}
	RequestHandler func(request *http.Request) *http.Request
}

type Client struct {
	client           *http.Client
	clientName       string // какой сервис выполняет запрос
	ownerServiceName string // какому сервису запрос
	baseURL          string
	body             []byte
}

func NewClient(clientName string, ownerServiceName string, baseURL string) *Client {
	return &Client{
		client:           &http.Client{},
		clientName:       clientName,
		ownerServiceName: ownerServiceName,
		baseURL:          baseURL,
	}
}

func (c *Client) GetBaseURL() string {
	return c.baseURL
}
func (c *Client) Get(ctx context.Context, url string, params *RequestParams) (*http.Response, error) {
	return c.get(ctx, http.MethodGet, url, params)
}
func (c *Client) Delete(ctx context.Context, url string, params *RequestParams) (*http.Response, error) {
	return c.get(ctx, http.MethodDelete, url, params)
}
func (c *Client) Post(ctx context.Context, url string, params *RequestParams) (*http.Response, error) {
	return c.post(ctx, http.MethodPost, url, params)
}
func (c *Client) Put(ctx context.Context, url string, params *RequestParams) (*http.Response, error) {
	return c.post(ctx, http.MethodPut, url, params)
}

func (c *Client) get(ctx context.Context, method string, url string, params *RequestParams) (*http.Response, error) {
	bodyBuffer := bytes.NewBuffer([]byte(""))
	if params != nil && params.RequestBody != nil {
		body, err := json.Marshal(params.RequestBody)
		if err != nil {
			return nil, err
		}
		bodyBuffer = bytes.NewBuffer(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", c.GetBaseURL(), url), bodyBuffer)
	if err != nil {
		return nil, err
	}

	if c.clientName != "" {
		req.Header.Add("X-Service-Name", c.clientName)
	}

	if params != nil {
		if params.ProxyToken {
			token, ok := c.GetToken(ctx)
			if !ok {
				return nil, fmt.Errorf("no token in context")
			}
			req.Header.Add("Authorization", "Bearer "+token)
		}

		if params.RequestHandler != nil {
			req = params.RequestHandler(req)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusBadRequest {
		if params != nil && params.ResponseBody != nil {
			err = c.ReadBody(resp)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(c.body, &params.ResponseBody)
			if err != nil {
				return nil, err
			}
		}
		return resp, err
	} else {
		err = c.ReadBody(resp)
		if err != nil {
			return nil, err
		}
		return resp, c.ParseError(resp.StatusCode)
	}
}

func (c *Client) post(ctx context.Context, method string, url string, params *RequestParams) (*http.Response, error) {
	bodyBuffer := bytes.NewBuffer([]byte(""))
	if params != nil && params.RequestBody != nil {
		body, err := json.Marshal(params.RequestBody)
		if err != nil {
			return nil, err
		}
		bodyBuffer = bytes.NewBuffer(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", c.GetBaseURL(), url), bodyBuffer)
	if err != nil {
		return nil, err
	}

	if c.clientName != "" {
		req.Header.Add("X-Service-Name", c.clientName)
	}

	if params != nil {
		if params.ProxyToken {
			token, ok := c.GetToken(ctx)
			if !ok {
				return nil, fmt.Errorf("no token in context")
			}
			req.Header.Add("Authorization", token)
		}

		if params.RequestHandler != nil {
			req = params.RequestHandler(req)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusBadRequest {
		if params != nil && params.ResponseBody != nil {
			err = c.ReadBody(resp)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(c.body, &params.ResponseBody)
			if err != nil {
				return nil, err
			}
		}
		return resp, err
	} else {
		err = c.ReadBody(resp)
		if err != nil {
			return nil, err
		}
		return resp, c.ParseError(resp.StatusCode)
	}

}

func (c *Client) ReadBody(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	c.body = body
	resp.Body.Close()
	return nil
}

func (c *Client) GetToken(ctx context.Context) (string, bool) {
	val := ctx.Value(Token)
	value, ok := val.(string)
	return value, ok
}

func (c *Client) ParseError(statusCode int) error {
	r := httpErrors.ResponseError{}
	err := json.Unmarshal(c.body, &r)
	if err != nil {
		return fmt.Errorf("client error: can not unmarshal body from %s: %s", c.ownerServiceName, err)
	}
	return domainErrors.NewCommonError(
		statusCode,
		r.Error.Code,
		nil,
		r.Error.Message,
		c.ownerServiceName,
	)
}
