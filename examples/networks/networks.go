package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/colevoss/go-iem-sdk"
)

func main() {
	ctx := context.Background()
	client := iem.NewClientWithOptions()

	networks, err := client.Networks().GetNetworks(ctx)

	if err != nil {
		panic(err)
	}

	str, _ := json.MarshalIndent(networks, "", "  ")
	fmt.Println(string(str))
	fmt.Printf("Network count: %d", len(networks))
}
