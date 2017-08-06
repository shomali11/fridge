package main

import (
	"fmt"
	"github.com/shomali11/fridge"
)

func main() {
	redisClient := fridge.NewRedisClient(fridge.WithHost("localhost"), fridge.WithPort(6379))
	client := fridge.NewClient(redisClient)
	defer client.Close()

	fmt.Println(client.Ping())
}
