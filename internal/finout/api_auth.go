package finout

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiAuth struct {
	client *Client
}

func (a *ApiAuth) ListGroups(ctx context.Context) (*ListGroupsResponse, error) {
	url := fmt.Sprintf("%s/identity/resources/groups/v1", AUTH_API_ENDPOINT)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	err = a.client.prepareHttpRequest(req)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("_groupsRelations", "rolesAndUsers")
	req.URL.RawQuery = query.Encode()

	rf := NewRequestFuncs()
	rf.PostResponse = func(req *http.Request, resp *http.Response) error {
		if resp.StatusCode != 200 {
			return fmt.Errorf("response returned unexpected status code: %d", resp.StatusCode)
		}
		return nil
	}
	payload, err := DoRequest[ListGroupsResponse](a.client, req, rf)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (a *ApiAuth) CreateGroup(ctx context.Context, requestData CreateGroupRequest) (*CreateGroupResponse, error) {
	url := fmt.Sprintf("%s/identity/resources/groups/v1", AUTH_API_ENDPOINT)

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
	payload, err := DoRequest[CreateGroupResponse](a.client, req, rf)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (a *ApiAuth) AddUsersToGroup(ctx context.Context, groupId string, requestData AddUsersToGroupRequest) error {
	url := fmt.Sprintf("%s/identity/resources/groups/v1/%s/users", AUTH_API_ENDPOINT, groupId)

	serialisedPayload, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(serialisedPayload))
	if err != nil {
		return err
	}

	err = a.client.prepareHttpRequest(req)
	if err != nil {
		return err
	}

	resp, err := a.client.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response returned unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (a *ApiAuth) AddRolesToGroup(ctx context.Context, groupId string, requestData AddRolesToGroupRequest) error {
	url := fmt.Sprintf("%s/identity/resources/groups/v1/%s/roles", AUTH_API_ENDPOINT, groupId)

	serialisedPayload, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(serialisedPayload))
	if err != nil {
		return err
	}

	err = a.client.prepareHttpRequest(req)
	if err != nil {
		return err
	}

	resp, err := a.client.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response returned unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (a *ApiAuth) UpdateGroup(ctx context.Context, groupId string, requestData interface{}) error {
	url := fmt.Sprintf("%s/identity/resources/groups/v1/%s", AUTH_API_ENDPOINT, groupId)

	serialisedPayload, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(serialisedPayload))
	if err != nil {
		return err
	}

	err = a.client.prepareHttpRequest(req)
	if err != nil {
		return err
	}

	resp, err := a.client.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response returned unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (a *ApiAuth) DeleteGroup(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/identity/resources/groups/v1/%s", AUTH_API_ENDPOINT, id)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	err = a.client.prepareHttpRequest(req)
	if err != nil {
		return err
	}

	resp, err := a.client.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response returned unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
