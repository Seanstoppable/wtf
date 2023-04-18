package rollbar

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/wtfutil/wtf/utils"
)

type RollbarClient struct {
	accessToken string
	client      *http.Client
}

func NewRollbarClient(accessToken string) *RollbarClient {
	httpClient := utils.DefaultHttpClient()
	rollbar := RollbarClient {
		accessToken: accessToken,
		client: httpClient,
	}

	return &rollbar
}

func (rollbar *RollbarClient) CurrentActiveItems(assignedToName string, activeOnly bool) (*ActiveItems, error) {
	items := &ActiveItems{}

	resp, err := rollbar.rollbarItemRequest(assignedToName, activeOnly)
	if err != nil {
		return items, err
	}

	err = utils.ParseJSON(&items, resp.Body)
	if err != nil {
		return items, err
	}

	return items, nil
}

/* -------------------- Unexported Functions -------------------- */

var (
	rollbarAPIURL = &url.URL{
		Scheme: "https",
		Host: "api.rollbar.com",
		Path: "/api/1/items",
	}
)

func (rollbar *RollbarClient) rollbarItemRequest(assignedToName string, activeOnly bool) (*http.Response, error) {
	params := url.Values{}
	params.Add("access_token", rollbar.accessToken)
	params.Add("assigned_user", assignedToName)
	if activeOnly {
		params.Add("status", "active")
	}

	requestURL := rollbarAPIURL.ResolveReference(&url.URL{RawQuery: params.Encode()})
	req, _ := http.NewRequest("GET", requestURL.String(), http.NoBody)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := rollbar.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf(resp.Status)
	}

	return resp, nil
}
