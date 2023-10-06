package iem

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	client  *http.Client
	baseUrl string

	networkService NetworkService
	weatherService WeatherService
	stationService StationService
}

type ClientOption func(*Client)

func WithNetworkService(service NetworkService) ClientOption {
	return func(client *Client) {
		client.networkService = service
	}
}

func WithWeatherService(service WeatherService) ClientOption {
	return func(client *Client) {
		client.weatherService = service
	}
}

const iemUrl = "https://mesonet.agron.iastate.edu"

func NewClient() *Client {
	client := &Client{
		baseUrl: iemUrl,
		client:  http.DefaultClient,
	}

	client.networkService = &IEMNetworkService{client}
	client.weatherService = &IEMWeatherService{client}
	client.stationService = &IEMStationService{client}

	return client
}

func NewClientWithOptions(opts ...ClientOption) *Client {
	client := NewClient()

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *Client) get(ctx context.Context, url string) (io.ReadCloser, error) {
	requestUrl := fmt.Sprintf("%s%s", c.baseUrl, url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl, nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	// TODO: Handle HTTP Errors

	return resp.Body, nil
}

type IEMNotFoundError struct {
	Detail string `json:"detail"`
	Code   int    `json:"code"`
}

func (err IEMNotFoundError) Error() string {
	return err.Detail
}

func (c *Client) getJson(ctx context.Context, url string, result interface{}) error {
	requestUrl := fmt.Sprintf("%s%s", c.baseUrl, url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl, nil)

	if err != nil {
		return err
	}

	res, err := c.client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	if res.StatusCode == 404 {
		var notFoundError IEMNotFoundError

		if err = json.Unmarshal(body, &notFoundError); err != nil {
			return err
		}

		notFoundError.Code = 404

		return notFoundError
	}

	if err = json.Unmarshal(body, result); err != nil {
		return err
	}

	return nil
}

func (c *Client) Networks() NetworkService {
	return c.networkService
}

func (c *Client) Weather() WeatherService {
	return c.weatherService
}

func (c *Client) Stations() StationService {
	return c.stationService
}
