package main

import (
	"fmt"
	"github.com/shomali11/fridge"
)

func main() {
	redisClient := fridge.NewRedisClient()
	client := fridge.NewClient(redisClient)
	defer client.Close()

	fmt.Println(client.Ping())
}
