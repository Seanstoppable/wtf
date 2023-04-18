package gitter

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/wtfutil/wtf/utils"
)

type GitterClient struct {
	apiToken string
	client   *http.Client
}

func NewGitterClient(apiToken string) *GitterClient {
	client := utils.DefaultHttpClient()
	gitter := GitterClient {
		apiToken: apiToken,
		client:   client,
	}
	return &gitter
}

func (gitter *GitterClient) GetMessages(roomId string, numberOfMessages int) ([]Message, error) {
	var messages []Message

	resp, err := gitter.apiRequest("rooms/"+roomId+"/chatMessages?limit="+strconv.Itoa(numberOfMessages))
	if err != nil {
		return nil, err
	}

	err = utils.ParseJSON(&messages, resp.Body)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (gitter *GitterClient) GetRoom(roomUri string) (*Room, error) {
	var rooms Rooms

	resp, err := gitter.apiRequest("rooms?q="+roomUri)
	if err != nil {
		return nil, err
	}

	err = utils.ParseJSON(&rooms, resp.Body)
	if err != nil {
		return nil, err
	}

	for _, room := range rooms.Results {
		if room.URI == roomUri {
			return &room, nil
		}
	}

	return nil, nil
}

/* -------------------- Unexported Functions -------------------- */

var (
	apiBaseURL = "https://api.gitter.im/v1/"
)

func (gitter *GitterClient) apiRequest(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", apiBaseURL+path, http.NoBody)
	if err != nil {
		return nil, err
	}

	bearer := fmt.Sprintf("Bearer %s", gitter.apiToken)
	req.Header.Add("Authorization", bearer)

	resp, err := gitter.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf(resp.Status)
	}

	return resp, nil
}
