package iem

import (
	"context"
	"fmt"
	"os"
)

type WeatherService interface {
	Get(ctx context.Context, query *WeatherDataQueryBuilder) ([]*IEMWeatherData, error)
}

type TestWeatherService struct{}

func (s *TestWeatherService) GetWeatherAtStation(ctx context.Context, query *WeatherDataQueryBuilder) ([]*IEMWeatherData, error) {
	// file, err := os.Open("./data/full_weather_data.csv")
	file, err := os.Open("./data/partial_weather_data.csv")
	defer file.Close()

	data, err := ParseWeatherData(file, query)

	return data, err
}

type IEMWeatherService struct {
	client *Client
}

func (s *IEMWeatherService) Get(ctx context.Context, query *WeatherDataQueryBuilder) ([]*IEMWeatherData, error) {
	v, err := query.BuildUrl()

	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("/cgi-bin/request/asos.py?%s", v.Encode())

	body, err := s.client.get(ctx, url)

	if err != nil {
		return nil, err
	}

	defer body.Close()

	weather, err := ParseWeatherData(body, query)

	return weather, err
}
