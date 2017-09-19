package main

import (
	"fmt"
	"github.com/shomali11/fridge"
)

func main() {
	redisCache := fridge.NewSentinelCache(
		fridge.WithSentinelAddresses([]string{"localhost:26379"}),
		fridge.WithSentinelMasterName("master"))

	client := fridge.NewClient(redisCache)
	defer client.Close()

	fmt.Println(client.Ping())
}
