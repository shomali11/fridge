package main

import (
	"fmt"
	"github.com/shomali11/fridge"
	"time"
)

func main() {
	redisCache := fridge.NewRedisCache()
	client := fridge.NewClient(redisCache)
	defer client.Close()

	fmt.Println(client.Put("food", "Pizza", fridge.WithDurations(time.Second, 2*time.Second)))
	fmt.Println(client.Get("food"))

	time.Sleep(time.Second)

	fmt.Println(client.Get("food"))

	time.Sleep(2 * time.Second)

	fmt.Println(client.Get("food"))
	fmt.Println(client.Remove("food"))
}
