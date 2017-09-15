package main

import (
	"fmt"
	"github.com/shomali11/fridge"
)

func main() {
	redisCache := fridge.NewRedisCache()
	client := fridge.NewClient(redisCache)
	defer client.Close()

	fmt.Println(client.Ping())
}
