package main

import (
	"fmt"
	"github.com/shomali11/fridge"
	"time"
)

func main() {
	redisCache := fridge.NewRedisCache()
	client := fridge.NewClient(redisCache, fridge.WithDefaultDurations(time.Second, 2*time.Second))
	defer client.Close()

	client.HandleEvent(func(event *fridge.Event) {
		fmt.Print("Key: " + event.Key + " - ")

		switch event.Type {
		case fridge.Fresh:
			fmt.Println("Woohoo! it is fresh!")
		case fridge.Cold:
			fmt.Println("Not fresh! But not bad either!")
		case fridge.Expired:
			fmt.Println("Sigh. It has expired!")
		case fridge.NotFound:
			fmt.Println("Oops! Did not find it.")
		case fridge.Restock:
			fmt.Println("Yay! Getting a new one!")
		case fridge.OutOfStock:
			fmt.Println("Oh no! It is out of stock.")
		case fridge.Unchanged:
			fmt.Println("Interesting! It has not changed.")
		}
	})

	restock := func() (string, error) {
		return "Pizza", nil
	}

	client.Put("food1", "Pizza")
	client.Put("food2", "Milk")

	client.Get("food1", fridge.WithRestock(restock))
	client.Get("food2")
	client.Get("food3")

	time.Sleep(time.Second)

	client.Get("food1", fridge.WithRestock(restock))
	client.Get("food2")
	client.Get("food3")

	time.Sleep(2 * time.Second)

	client.Get("food1", fridge.WithRestock(restock))
	client.Get("food2")
	client.Get("food3")

	client.Remove("food1")
	client.Remove("food2")
	client.Remove("food3")
}
