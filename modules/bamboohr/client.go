package bamboohr

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/wtfutil/wtf/utils"
)

// A BambooClient represents the data required to connect to the BambooHR API
type BambooClient struct {
	apiBase   string
	apiKey    string
	subdomain string
	client    *http.Client
}

// NewBambooClient creates and returns a new BambooHR client
func NewBambooClient(url string, apiKey string, subdomain string) *BambooClient {
	httpClient := utils.DefaultHttpClient()
	client := BambooClient{
		apiBase:   url,
		apiKey:    apiKey,
		subdomain: subdomain,
		client:    &httpClient,
	}

	return &client
}

/* -------------------- Public Functions -------------------- */

// Away returns a string representation of the people who are out of the office during the defined period
func (client *BambooClient) Away(itemType, startDate, endDate string) []Item {
	calendar, err := client.getWhoIsAway(startDate, endDate)
	if err != nil {
		return []Item{}
	}

	items := calendar.ItemsByType(itemType)

	return items
}

/* -------------------- Private Functions -------------------- */

// getWhoIsAway is the private interface for retrieving structural data about who will be out of the office
// This method does the actual communication with BambooHR and returns the raw Go
// data structures used by the public interface
func (client *BambooClient) getWhoIsAway(startDate, endDate string) (cal Calendar, err error) {
	apiURL := fmt.Sprintf(
		"%s/%s/v1/time_off/whos_out?start=%s&end=%s",
		client.apiBase,
		client.subdomain,
		startDate,
		endDate,
	)

	data, err := client.request(apiURL)
	if err != nil {
		return cal, err
	}
	err = xml.Unmarshal(data, &cal)

	return
}


func (client *BambooClient) request(apiURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", apiURL, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(client.apiKey, "x")

	resp, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := ParseBody(resp)
	if err != nil {
		return nil, err
	}

	return data, err
}

func ParseBody(resp *http.Response) ([]byte, error) {
	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
