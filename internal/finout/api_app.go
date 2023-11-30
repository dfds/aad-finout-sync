package finout

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ApiApp struct {
	client *Client
}

func (a *ApiApp) UpdateAccountDataAccessForGroups(ctx context.Context, accountId string, requestData UpdateAccountDataAccessForGroupsRequest) (*UpdateAccountDataAccessForGroupsResponse, error) {
	url := fmt.Sprintf("%s/account-service/account/%s", APP_API_ENDPOINT, accountId)

	serialisedPayload, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(serialisedPayload))
	if err != nil {
		return nil, err
	}

	err = a.client.prepareHttpRequest(req)
	if err != nil {
		return nil, err
	}

	rf := NewRequestFuncs()
	rf.PostResponse = func(req *http.Request, resp *http.Response) error {
		if resp.StatusCode != 200 {
			return fmt.Errorf("response returned unexpected status code: %d", resp.StatusCode)
		}
		return nil
	}
	payload, err := DoRequest[UpdateAccountDataAccessForGroupsResponse](a.client, req, rf)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (a *ApiApp) GetAccountData(ctx context.Context, accountId string) (*UpdateAccountDataAccessForGroupsResponse, error) {
	url := fmt.Sprintf("%s/account-service/account/%s", APP_API_ENDPOINT, accountId)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	err = a.client.prepareHttpRequest(req)
	if err != nil {
		return nil, err
	}

	rf := NewRequestFuncs()
	rf.PostResponse = func(req *http.Request, resp *http.Response) error {
		if resp.StatusCode != 200 {
			return fmt.Errorf("response returned unexpected status code: %d", resp.StatusCode)
		}
		return nil
	}
	payload, err := DoRequest[UpdateAccountDataAccessForGroupsResponse](a.client, req, rf)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (a *ApiApp) ListVirtualTags(ctx context.Context) (map[string]*ListVirtualTagResponseTag, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/virtual-tags-service/virtual-tag", APP_API_ENDPOINT), nil)
	if err != nil {
		return nil, err
	}
	err = a.client.prepareHttpRequest(req)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("fields", "id,name,createdBy,createdAt,updatedAt")
	req.URL.RawQuery = query.Encode()

	rf := NewRequestFuncs()
	rf.PostResponse = func(req *http.Request, resp *http.Response) error {
		if resp.StatusCode != 200 {
			return fmt.Errorf("response returned unexpected status code: %d", resp.StatusCode)
		}
		return nil
	}
	payload, err := DoRequest[[]*ListVirtualTagResponseTag](a.client, req, rf)
	if err != nil {
		return nil, err
	}

	tags := make(map[string]*ListVirtualTagResponseTag)
	for _, tag := range *payload {
		fmt.Println(strings.ToLower(tag.Name))
		tags[strings.ToLower(tag.Name)] = tag
	}

	return tags, nil
}

func (a *ApiApp) CreateVirtualTag(ctx context.Context, requestPayload CreateVirtualTagRequest) (*CreateVirtualTagResponse, error) {
	serialised, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/virtual-tags-service/virtual-tag", APP_API_ENDPOINT), bytes.NewBuffer(serialised))
	if err != nil {
		return nil, err
	}
	err = a.client.prepareJsonRequest(req)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("dataFormat", "UI")
	req.URL.RawQuery = query.Encode()

	rf := NewRequestFuncs()
	rf.PostResponse = func(req *http.Request, resp *http.Response) error {
		if resp.StatusCode != 200 {
			return fmt.Errorf("response returned unexpected status code: %d", resp.StatusCode)
		}
		return nil
	}
	payload, err := DoRequest[CreateVirtualTagResponse](a.client, req, rf)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (a *ApiApp) UpdateVirtualTag(ctx context.Context, requestPayload UpdateVirtualTagRequest, id string) (*UpdateVirtualTagResponse, error) {
	serialised, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/virtual-tags-service/virtual-tag/%s", APP_API_ENDPOINT, id), bytes.NewBuffer(serialised))
	if err != nil {
		return nil, err
	}
	err = a.client.prepareJsonRequest(req)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("dataFormat", "UI")
	req.URL.RawQuery = query.Encode()

	rf := NewRequestFuncs()
	rf.PostResponse = func(req *http.Request, resp *http.Response) error {
		if resp.StatusCode != 200 {
			return fmt.Errorf("response returned unexpected status code: %d", resp.StatusCode)
		}
		return nil
	}
	payload, err := DoRequest[UpdateVirtualTagResponse](a.client, req, rf)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (a *ApiApp) ListViews(ctx context.Context) (*ListViewsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/v1/view", APP_API_ENDPOINT), nil)
	if err != nil {
		return nil, err
	}
	err = a.client.prepareHttpRequest(req)
	if err != nil {
		return nil, err
	}

	rf := NewRequestFuncs()
	rf.PostResponse = func(req *http.Request, resp *http.Response) error {
		if resp.StatusCode != 200 {
			return fmt.Errorf("response returned unexpected status code: %d", resp.StatusCode)
		}
		return nil
	}
	payload, err := DoRequest[ListViewsResponse](a.client, req, rf)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (a *ApiApp) QueryByView(ctx context.Context, reqPayload QueryByViewRequest) (*QueryByViewResponse, error) {
	serialised, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v1/cost/query-by-view", APP_API_ENDPOINT), bytes.NewBuffer(serialised))
	if err != nil {
		return nil, err
	}
	err = a.client.prepareJsonRequest(req)
	if err != nil {
		return nil, err
	}

	rf := NewRequestFuncs()
	rf.PostResponse = func(req *http.Request, resp *http.Response) error {
		if resp.StatusCode != 200 {
			return fmt.Errorf("response returned unexpected status code: %d", resp.StatusCode)
		}
		return nil
	}
	payload, err := DoRequest[QueryByViewResponse](a.client, req, rf)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
