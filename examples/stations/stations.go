package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"

	"github.com/colevoss/go-iem-sdk"
)

func main() {
	stationId := flag.String("station", "", "Station to fetch")
	networkId := flag.String("network", "", "Network to fetch stations for")
	flag.Parse()

	client := iem.NewClientWithOptions()

	if *stationId != "" {
		getStation(client, *stationId)
		return
	}

	if *networkId != "" {
		getAllStations(client, *networkId)
		return
	}

	log.Printf("Pass a network or station id to run this example")
}

func getStation(client *iem.Client, stationId string) {
	ctx := context.Background()

	station, err := client.Stations().GetStation(ctx, stationId)

	if err != nil {
		panic(err)
	}

	str, _ := json.MarshalIndent(station, "", "  ")
	log.Printf(string(str))
}

func getAllStations(client *iem.Client, networkId string) {
	ctx := context.Background()

	stations, err := client.Stations().GetStations(ctx, networkId)

	if err != nil {
		panic(err)
	}

	str, _ := json.MarshalIndent(stations, "", "  ")
	log.Printf(string(str))
	log.Printf("Station count for network (%s): %d", networkId, len(stations))
}
