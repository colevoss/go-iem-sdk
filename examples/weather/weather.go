package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/colevoss/go-iem-sdk"
)

func main() {
	stationId := flag.String("s", "", "Station to fetch weather for")
	flag.Parse()

	if *stationId == "" {
		fmt.Println("Station (s) arg is required")
		return
	}

	ctx := context.Background()
	client := iem.NewClientWithOptions()

	now := time.Now()
	yesterday := now.AddDate(0, 0, -2)

	query := iem.NewWeatherDataQuery()
	query.Stations(*stationId)
	query.Timezone("America/Chicago")
	query.ReportType(3)

	query.Data(iem.TempF, iem.TempC, iem.Feel)

	query.Start(yesterday)
	query.End(yesterday)

	data, err := client.Weather().Get(ctx, query)

	if err != nil {
		panic(err)
	}

	str, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(str))
	fmt.Printf("Weather Records Count: %d\n", len(data))
}
