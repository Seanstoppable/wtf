package airbrake

import (
	"fmt"
	"net/http"

	"github.com/wtfutil/wtf/utils"
)

type AirbrakeClient struct {
	authToken string
	client *http.Client
}

func NewAirbrakeClient(authToken string) *AirbrakeClient {
	httpClient := utils.DefaultHttpClient()
	airbrake := AirbrakeClient {
		client: &httpClient,
		authToken: authToken,
	}
	return &airbrake
}

func (airbrake *AirbrakeClient) project(projectID int) (*Project, error) {
	url := fmt.Sprintf(
		"https://api.airbrake.io/api/v4/projects/%d?key=%s",
		projectID, airbrake.authToken)
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := airbrake.client.Do(req)
	if err != nil {
		return nil, err
	}

	p := &ProjectJSON{}
	err = utils.ParseJSON(p, resp.Body)
	if err != nil {
		return nil, err
	}

	return &p.Project, nil
}

func (airbrake *AirbrakeClient) groups(projectID int) ([]Group, error) {
	url := fmt.Sprintf(
		"https://api.airbrake.io/api/v4/projects/%d/groups?key=%s&limit=10&order=last_notice&resolved=false",
		projectID, airbrake.authToken)
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := airbrake.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	j := &GroupJSON{}
	err = utils.ParseJSON(j, resp.Body)
	if err != nil {
		return nil, err
	}

	return j.Groups, nil
}

func (airbrake *AirbrakeClient) resolveGroup(projectID int64, groupID string) error {
	url := fmt.Sprintf(
		"https://airbrake.io/api/v4/projects/%d/groups/%s/resolved?key=%s",
		projectID, groupID, airbrake.authToken)
	req, err := http.NewRequest("PUT", url, http.NoBody)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := airbrake.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	return nil
}

func (airbrake *AirbrakeClient) muteGroup(projectID int64, groupID string) error {
	url := fmt.Sprintf(
		"https://airbrake.io/api/v4/projects/%d/groups/%s/muted?key=%s",
		projectID, groupID, airbrake.authToken)
	req, err := http.NewRequest("PUT", url, http.NoBody)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := airbrake.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	return nil
}

func (airbrake *AirbrakeClient) unmuteGroup(projectID int64, groupID string) error {
	url := fmt.Sprintf(
		"https://airbrake.io/api/v4/projects/%d/groups/%s/unmuted?key=%s",
		projectID, groupID, airbrake.authToken)
	req, err := http.NewRequest("PUT", url, http.NoBody)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := airbrake.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	return nil
}
