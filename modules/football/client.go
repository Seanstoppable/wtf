package football

import (
	"fmt"
	"net/http"

	"github.com/wtfutil/wtf/utils"
)

var (
	footballAPIUrl = "https://api.football-data.org/v2"
)

type leagueInfo struct {
	id      int
	caption string
}

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	client := Client{
		apiKey:     apiKey,
		httpClient: utils.DefaultHttpClient(),
	}

	return &client
}

func (client *Client) footballRequest(path string, id int) (*http.Response, error) {

	url := fmt.Sprintf("%s/competitions/%d/%s", footballAPIUrl, id, path)
	req, err := http.NewRequest("GET", url, http.NoBody)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", client.apiKey)
	if err != nil {
		return nil, err
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	return resp, nil
}
