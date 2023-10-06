package canopyapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ImpactPredictionReport struct {
	EventDate      string   `json:"eventDate"`
	EventPublished string   `json:"eventPublished"`
	StateUrl       []string `json:"stateUrl"`
	SizeUrl        []string `json:"sizeUrl"`
}

type impactPrediction struct {
	client *Client
}

func (c *Client) ImpactPrediction() *impactPrediction {
	return &impactPrediction{
		client: c,
	}
}

func (i *impactPrediction) GetByDate(date string) (*ImpactPredictionReport, error) {
	req, err := i.client.createRequest("GET", "/impact-prediction/"+date)

	if err != nil {
		log.Printf("Request creation failed. %s", err)
	}

	response := i.client.executeRequest(req)

	if response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	totalRequests := 1
	for response.StatusCode == http.StatusTooManyRequests {
		response = i.client.executeRequest(req)
		time.Sleep(10 * time.Millisecond)

		if totalRequests > 10 {
			break
		}
	}

	if response.StatusCode != http.StatusOK {
		apiError := ApiError{}

		json.Unmarshal(response.Data, &apiError)

		return nil, errors.New("Status " + fmt.Sprint(response.StatusCode) + ":" + apiError.Message)
	}

	report := ImpactPredictionReport{}

	err = json.Unmarshal(response.Data, &report)

	if err != nil {
		return nil, err
	}

	return &report, nil
}
