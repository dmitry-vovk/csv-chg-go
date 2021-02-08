package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Client is a warehouse API client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Item represents a response from `/item/{uuid}` API endpoint
type Item struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

const (
	getItemPath   = "/item/"
	postAlertPath = "/low-stock-alert/"
)

// New returns an instance of Client initialized with baseURL
func New(baseURL string) *Client {
	return &Client{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		httpClient: http.DefaultClient,
	}
}

// GetItem performs a GET API call to `/item/{uuid}`
func (c *Client) GetItem(uuid string) (*Item, error) {
	resp, err := c.httpClient.Get(c.baseURL + getItemPath + uuid)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	switch resp.StatusCode {
	case http.StatusOK: // 200
		if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
			return nil, ErrInvalidContentType{contentType: ct}
		}
		var item Item
		decoder := json.NewDecoder(resp.Body)
		decoder.DisallowUnknownFields()
		if err = decoder.Decode(&item); err != nil {
			return nil, err
		}
		return &item, nil
	case http.StatusBadRequest: // 400
		return nil, ErrBadRequest
	case http.StatusInternalServerError: // 500
		return nil, ErrServerError
	}
	// Any other code is error
	return nil, ErrUnexpectedStatusCode{code: resp.StatusCode}
}

// PostAlert performs a POST API call to `/low-stock-alert/{uuid}`
func (c *Client) PostAlert(uuid string) error {
	resp, err := c.httpClient.Post(c.baseURL+postAlertPath+uuid, "", nil)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusCreated: // 201
		return nil
	case http.StatusBadRequest: // 400
		return ErrBadRequest
	case http.StatusInternalServerError: // 500
		return ErrServerError
	}
	return ErrUnexpectedStatusCode{resp.StatusCode}
}
