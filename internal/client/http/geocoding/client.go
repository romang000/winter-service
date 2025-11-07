package geocoding

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	httpClient *http.Client
}

type Response struct {
	Name      string  `json:"name"`
	Country   string  `json:"country"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

func (c *Client) GetCoords(city string) (Response, error) {
	const op = "geocoding.client.GetCoords"
	res, err := c.httpClient.Get(
		fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=ru&format=json",
			city,
		),
	)
	
	if err != nil {
		return Response{}, fmt.Errorf("%s: %w", op, err)
	}
	
	defer res.Body.Close()
	
	if res.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("%s: status code: %d", op, res.StatusCode)
	}
	
	var geoResp struct {
		Results []Response `json:"results"`
	}
	
	if err = json.NewDecoder(res.Body).Decode(&geoResp); err != nil {
		return Response{}, fmt.Errorf("%s: %w", op, err)
	}
	
	return geoResp.Results[0], nil
}
