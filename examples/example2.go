package main

import (
	"fmt"
	"github.com/shomali11/fridge"
	"github.com/shomali11/xredis"
)

func main() {
	options := &xredis.Options{
		Host: "localhost",
		Port: 6379,
	}

	xredisClient := xredis.SetupClient(options)
	client := fridge.NewClient(fridge.WithRedisClient(xredisClient))
	defer client.Close()

	fmt.Println(client.Ping())
}
