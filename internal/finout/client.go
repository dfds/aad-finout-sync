package finout

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
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
}

func (c *Client) SetAuthMethod(method AuthMethod) {
	c.authMethod = method
}

func (c *Client) Auth() error {
	return c.authMethod.Refresh()
}

func (c *Client) prepareHttpRequest(h *http.Request) error {
	err := c.authMethod.PrepareHttpRequest(h)
	if err != nil {
		return err
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

func (c *Client) ListVirtualTags(ctx context.Context) (map[string]*ListVirtualTagResponseTag, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/virtual-tags-service/virtual-tag", APP_API_ENDPOINT), nil)
	if err != nil {
		return nil, err
	}
	err = c.prepareHttpRequest(req)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("fields", "id,name,createdBy,createdAt,updatedAt")
	req.URL.RawQuery = query.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("response returned unexpected status code: %d", resp.StatusCode)
	}

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var payload []*ListVirtualTagResponseTag

	err = json.Unmarshal(rawData, &payload)
	if err != nil {
		return nil, err
	}

	tags := make(map[string]*ListVirtualTagResponseTag)
	for _, tag := range payload {
		fmt.Println(strings.ToLower(tag.Name))
		tags[strings.ToLower(tag.Name)] = tag
	}

	return tags, nil
}

func (c *Client) CreateVirtualTag(ctx context.Context, requestPayload CreateVirtualTagRequest) (*CreateVirtualTagResponse, error) {
	serialised, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/virtual-tags-service/virtual-tag", APP_API_ENDPOINT), bytes.NewBuffer(serialised))
	if err != nil {
		return nil, err
	}
	err = c.prepareJsonRequest(req)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("dataFormat", "UI")
	req.URL.RawQuery = query.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected HTTP status code")
	}

	var payload *CreateVirtualTagResponse
	err = json.Unmarshal(rawData, &payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (c *Client) UpdateVirtualTag(ctx context.Context, requestPayload UpdateVirtualTagRequest, id string) (*UpdateVirtualTagResponse, error) {
	serialised, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/virtual-tags-service/virtual-tag/%s", APP_API_ENDPOINT, id), bytes.NewBuffer(serialised))
	if err != nil {
		return nil, err
	}
	err = c.prepareJsonRequest(req)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("dataFormat", "UI")
	req.URL.RawQuery = query.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(string(rawData))
		return nil, errors.New(fmt.Sprintf("unexpected HTTP status code %d", resp.StatusCode))
	}

	var payload *UpdateVirtualTagResponse
	err = json.Unmarshal(rawData, &payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (c *Client) ListViews(ctx context.Context) (*ListViewsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/v1/view", APP_API_ENDPOINT), nil)
	if err != nil {
		return nil, err
	}
	err = c.prepareHttpRequest(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("response returned unexpected status code: %d", resp.StatusCode)
	}

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var payload *ListViewsResponse

	err = json.Unmarshal(rawData, &payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (c *Client) QueryByView(ctx context.Context, reqPayload QueryByViewRequest) (*QueryByViewResponse, error) {
	serialised, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v1/cost/query-by-view", APP_API_ENDPOINT), bytes.NewBuffer(serialised))
	if err != nil {
		return nil, err
	}
	err = c.prepareJsonRequest(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected HTTP status code")
	}

	var payload *QueryByViewResponse
	err = json.Unmarshal(rawData, &payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func NewFinoutClient() *Client {
	payload := &Client{
		httpClient: http.DefaultClient,
		authMethod: nil,
	}
	return payload
}
