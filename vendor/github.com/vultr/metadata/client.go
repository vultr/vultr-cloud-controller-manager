package metadata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	timeout  = 3 * time.Second
	basePath = "http://169.254.169.254"
)

// Client ...
type Client struct {
	client  *http.Client
	baseURL *url.URL
}

// NewClient creates a client to interact with the metadata
func NewClient() *Client {

	u, err := url.Parse(basePath)
	if err != nil {
		panic(err)
	}

	c := &Client{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: u,
	}

	return c
}

// Metadata returns the entire contents of the instances metadata
func (c *Client) Metadata() (*MetaData, error) {
	metadata := &MetaData{}

	err := c.doRequest("/v1.json", metadata)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (c *Client) doRequest(uri string, meta *MetaData) error {
	resp, err := c.client.Get(fmt.Sprintf("%s%s", c.baseURL, uri))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		if err := json.Unmarshal(body, meta); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("error : %s", string(body))
	}

	return nil
}

// SetBaseURL lets you change the default metadata url
func (c *Client) SetBaseURL(baseURL string) error {
	updatedURL, err := url.Parse(baseURL)
	if err != nil {
		return err
	}

	c.baseURL = updatedURL
	return nil
}

// RegionCodeToID takes in a region code and returns back the region ID
func RegionCodeToID(code string) string {
	regions := map[string]string{
		"EWR": "1",
		"ORD": "2",
		"DFW": "3",
		"SEA": "4",
		"LAX": "5",
		"ATL": "6",
		"AMS": "7",
		"LHR": "8",
		"FRA": "9",
		"SJC": "12",
		"SYD": "19",
		"YTO": "22",
		"CDG": "24",
		"NRT": "25",
		"ICN": "34",
		"MIA": "39",
		"SGP": "40",
	}

	return regions[code]
}
