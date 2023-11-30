package finout

import (
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	authMethod AuthMethod
}

const APP_API_ENDPOINT = "https://app.finout.io"
const AUTH_API_ENDPOINT = "https://auth.finout.io"

type Config struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type authFunc func(*Client)

type AuthMethod interface {
	PrepareHttpRequest(h *http.Request) error
	Refresh() error
	AcceptedEndpoint(string) bool
}

func (c *Client) SetAuthMethod(method AuthMethod) {
	c.authMethod = method
}

func (c *Client) Auth() error {
	return c.authMethod.Refresh()
}

func (c *Client) ApiAuth() *ApiAuth {
	return &ApiAuth{client: c}
}

func (c *Client) ApiApp() *ApiApp {
	return &ApiApp{client: c}
}

func (c *Client) prepareHttpRequest(h *http.Request) error {
	err := c.authMethod.PrepareHttpRequest(h)
	if err != nil {
		return err
	}

	if !c.authMethod.AcceptedEndpoint(h.URL.String()) {
		return InvalidAuthMethodForAction.New(InvalidAuthMethodForActionMsg)
	}

	h.Header.Set("User-Agent", "aad-finout-sync - github.com/dfds/aad-finout-sync")
	requestId, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	h.Header.Set("x-request-id", requestId.String())

	return nil
}

func (c *Client) prepareJsonRequest(req *http.Request) error {
	err := c.prepareHttpRequest(req)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	return nil
}

func NewFinoutClient() *Client {
	payload := &Client{
		httpClient: http.DefaultClient,
		authMethod: nil,
	}
	return payload
}

func DoRequest[T any](client *Client, req *http.Request, rf *RequestFuncs) (*T, error) {
	if rf == nil {
		rf = NewRequestFuncs()
	}
	resp, _, err := DoRequestWithResp[T](client, req, rf)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func DoRequestWithResp[T any](client *Client, req *http.Request, rf *RequestFuncs) (*T, *http.Response, error) {
	err := rf.PreResponse(req)
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	err = rf.PostResponse(req, resp)
	if err != nil {
		return nil, resp, err
	}

	defer resp.Body.Close()

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	var payload *T

	err = rf.PreDeserialise(req, resp)
	if err != nil {
		return nil, resp, err
	}

	err = json.Unmarshal(rawData, &payload)
	if err != nil {
		return nil, resp, err
	}

	err = rf.PostDeserialise(req, resp, payload)
	if err != nil {
		return payload, resp, err
	}

	return payload, resp, nil
}

type RequestFuncs struct {
	PreResponse     func(req *http.Request) error
	PostResponse    func(req *http.Request, resp *http.Response) error
	PreDeserialise  func(req *http.Request, resp *http.Response) error
	PostDeserialise func(req *http.Request, resp *http.Response, data interface{}) error
}

func NewRequestFuncs() *RequestFuncs {
	rf := &RequestFuncs{
		PreResponse: func(req *http.Request) error {
			return nil
		},
		PostResponse: func(req *http.Request, resp *http.Response) error {
			return nil
		},
		PreDeserialise: func(req *http.Request, resp *http.Response) error {
			return nil
		},
		PostDeserialise: func(req *http.Request, resp *http.Response, data interface{}) error {
			return nil
		},
	}

	return rf
}
