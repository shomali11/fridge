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

	restock := func() (string, error) {
		return "Hot Pizza", nil
	}

	fmt.Println(client.Put("food", "Pizza", fridge.WithDurations(time.Second, 2*time.Second)))
	fmt.Println(client.Get("food", fridge.WithRestock(restock)))

	time.Sleep(time.Second)

	fmt.Println(client.Get("food", fridge.WithRestock(restock)))

	time.Sleep(2 * time.Second)

	fmt.Println(client.Get("food", fridge.WithRestock(restock)))
	fmt.Println(client.Remove("food"))
}
