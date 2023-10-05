package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"time"

	"github.com/colevoss/go-iem-sdk"
)

func main() {
	stationId := flag.String("s", "", "Station to fetch weather for")
	flag.Parse()

	if *stationId == "" {
		log.Print("Station (s) arg is required")
		return
	}

	ctx := context.Background()
	client := iem.NewClientWithOptions()

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	query := iem.NewWeatherDataQuery()
	query.Stations(*stationId)

	query.Data(iem.TempF, iem.TempC, iem.Feel)

	query.Start(yesterday)
	query.End(yesterday)

	data, err := client.Weather().Get(ctx, query)

	if err != nil {
		panic(err)
	}
	str, _ := json.MarshalIndent(data, "", "  ")
	log.Printf(string(str))
	log.Printf("Weather Records Count: %d", len(data))
}
