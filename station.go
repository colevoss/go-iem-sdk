package iem

import (
	"context"
	"fmt"
)

type Station struct {
	Index       int     `json:"index"`
	Id          string  `json:"id"`
	Synop       float64 `json:"synop"`
	Name        string  `json:"name"`
	State       string  `json:"state"`
	Country     string  `json:"country"`
	Elevation   float64 `json:"elevation"`
	Network     string  `json:"network"`
	Online      bool    `json:"online"`
	Params      string  `json:"params"`
	County      string  `json:"county"`
	PlotName    string  `json:"plot_name"`
	ClimateSite string  `json:"climate_site"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type IEMStationsJsonResponse struct {
	Data []*Station
}

type StationService interface {
	GetStation(ctx context.Context, stationId string) (*Station, error)
	GetStations(ctx context.Context, networkId string) ([]*Station, error)
}

type IEMStationService struct {
	client *Client
}

func (s *IEMStationService) GetStations(ctx context.Context, networkId string) ([]*Station, error) {
	url := fmt.Sprintf("/api/1/network/%s.json", networkId)
	var stationResponse IEMStationsJsonResponse

	err := s.client.getJson(ctx, url, &stationResponse)

	if err != nil {
		return nil, err
	}

	return stationResponse.Data, nil
}

func (s *IEMStationService) GetStation(ctx context.Context, stationId string) (*Station, error) {
	url := fmt.Sprintf("/api/1/station/%s.json", stationId)
	var stationResponse IEMStationsJsonResponse

	err := s.client.getJson(ctx, url, &stationResponse)

	if err != nil {
		return nil, err
	}

	return stationResponse.Data[0], nil
}
