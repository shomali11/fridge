package main

import (
	"fmt"
	"github.com/shomali11/fridge"
	"time"
)

func main() {
	client := fridge.DefaultClient()
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
		case fridge.Refresh:
			fmt.Println("Yay! Getting a new one!")
		case fridge.OutOfStock:
			fmt.Println("Oh no! It is out of stock.")
		}
	})

	restock := func() (string, error) {
		return "Hot Pizza", nil
	}

	client.Register("food1", time.Second, 2*time.Second)
	client.Register("food2", time.Second, 2*time.Second)
	client.Register("food3", time.Second, 2*time.Second)

	client.Put("food1", "Pizza")
	client.Put("food2", "Milk")

	client.Get("food1", restock)
	client.Get("food2", nil)
	client.Get("food3", nil)

	time.Sleep(time.Second)

	client.Get("food1", restock)
	client.Get("food2", nil)
	client.Get("food3", nil)

	time.Sleep(2 * time.Second)

	client.Get("food1", restock)
	client.Get("food2", nil)
	client.Get("food3", nil)

	client.Remove("food1")
	client.Remove("food2")
	client.Remove("food3")

	client.Deregister("food1")
	client.Deregister("food2")
	client.Deregister("food3")
}
