package iem

import (
	"context"
	"time"
)

type Network struct {
	Index int    `json:"index"`
	Id    string `json:"id"`
	Name  string `json:"name"`
	Tz    string `json:"tzname"`
	// Extent string `json"extent"`
	WindroseUpdate time.Time `json:"windrose_update"`
}

type IEMNetowrkJsonResposne struct {
	Data []*Network `json:"data"`
}

type NetworkService interface {
	GetNetworks(ctx context.Context) ([]*Network, error)
}

type IEMNetworkService struct {
	client *Client
}

func (s *IEMNetworkService) GetNetworks(ctx context.Context) ([]*Network, error) {
	url := "/api/1/networks.json"
	var networkResponse IEMNetowrkJsonResposne
	err := s.client.getJson(ctx, url, &networkResponse)

	if err != nil {
		return nil, err
	}

	return networkResponse.Data, nil
}
