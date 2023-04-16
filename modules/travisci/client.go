package travisci

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/wtfutil/wtf/utils"
)

var TRAVIS_HOSTS = map[bool]string{
	false: "travis-ci.org",
	true:  "travis-ci.com",
}

var (
	travisAPIURL = &url.URL{Scheme: "https", Path: "/"}
)

type TravisClient struct {
	hostname string
	baseUrl  string
	apiKey   string
	client  *http.Client
}

func NewTravisClient(settings *Settings) *TravisClient {
	httpClient := utils.DefaultHttpClient()
	hostname := "api." + TRAVIS_HOSTS[settings.pro]
	baseUrl := "/api"
	if settings.baseURL != "" {
		baseUrl = settings.baseURL
	}

	travis := TravisClient {
		hostname: hostname,
		apiKey: settings.apiKey,
		baseUrl: baseUrl,
		client: &httpClient,
	}

	return &travis
}

func (travis *TravisClient) BuildsFor(limit string, sortBy string) (*Builds, error) {
	builds := &Builds{}

	resp, err := travis.buildRequest(limit, sortBy)
	if err != nil {
		return builds, err
	}

	err = utils.ParseJSON(&builds, resp.Body)
	if err != nil {
		return builds, err
	}

	return builds, nil
}

/* -------------------- Unexported Functions -------------------- */


func (travis *TravisClient) buildRequest(limit string, sortBy string) (*http.Response, error) {
	var path string = "builds"
	params := url.Values{}
	params.Add("limit", limit)
	params.Add("sort_by", sortBy)

	requestUrl := travisAPIURL.ResolveReference(&url.URL{Path: path, RawQuery: params.Encode()})

	req, err := http.NewRequest("GET", requestUrl.String(), http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Travis-API-Version", "3")

	bearer := fmt.Sprintf("token %s", travis.apiKey)
	req.Header.Add("Authorization", bearer)

	resp, err := travis.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf(resp.Status)
	}

	return resp, nil
}
