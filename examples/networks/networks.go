package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/colevoss/go-iem"
)

func main() {
	ctx := context.Background()
	client := iem.NewClientWithOptions()

	networks, err := client.Networks().GetNetworks(ctx)

	if err != nil {
		panic(err)
	}

	str, _ := json.MarshalIndent(networks, "", "  ")
	log.Printf(string(str))
	log.Printf("Network count: %d", len(networks))
}
