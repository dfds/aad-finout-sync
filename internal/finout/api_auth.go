package finout

import (
	"context"
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
