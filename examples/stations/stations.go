package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

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

	fmt.Println("Pass a network or station id to run this example")
}

func getStation(client *iem.Client, stationId string) {
	ctx := context.Background()

	station, err := client.Stations().GetStation(ctx, stationId)

	if err != nil {
		panic(err)
	}

	str, _ := json.MarshalIndent(station, "", "  ")
	fmt.Println(string(str))
}

func getAllStations(client *iem.Client, networkId string) {
	ctx := context.Background()

	stations, err := client.Stations().GetStations(ctx, networkId)

	if err != nil {
		panic(err)
	}

	str, _ := json.MarshalIndent(stations, "", "  ")
	fmt.Println(string(str))
	fmt.Printf("Station count for network (%s): %d", networkId, len(stations))
}
