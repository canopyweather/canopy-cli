package canopyapi

import (
	"errors"
	"io"
	"net/http"
)

type Client struct {
	APIKey string
	Url    string
}

type Response struct {
	Data       []byte
	StatusCode int
	Error      error
}

func NewClient(apiKey string, url string) *Client {

	return &Client{
		APIKey: apiKey,
		Url:    url,
	}
}

func (c *Client) createRequest(method string, endpoint string) (*http.Request, error) {
	url := c.Url + endpoint

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.APIKey)

	return req, nil
}

func (c *Client) executeRequest(req *http.Request) Response {

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return Response{
			Error: err,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return Response{
			Error: err,
		}
	}

	if resp.StatusCode != http.StatusOK {
		return Response{
			Data:       body,
			StatusCode: resp.StatusCode,
			Error:      errors.New("API Request failed"),
		}
	}

	return Response{
		Data:       body,
		StatusCode: resp.StatusCode,
		Error:      nil,
	}
}
