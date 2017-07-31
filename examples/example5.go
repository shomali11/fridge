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
		fmt.Println("Key: " + event.Key)

		switch event.Type {
		case fridge.Fresh:
			fmt.Println("Woohoo! it is fresh!")
		case fridge.Cold:
			fmt.Println("Not fresh! But not bad either!")
		case fridge.Miss:
			fmt.Println("Oops! Did not find it!")
		case fridge.Expired:
			fmt.Println("Sigh. It has expired!")
		}
	})

	restock := func() (string, error) {
		return "Hot Pizza", nil
	}

	client.Register("food", time.Second, 2*time.Second)
	client.Put("food", "Pizza")
	client.Get("food", restock)

	time.Sleep(time.Second)

	client.Get("food", restock)

	time.Sleep(2 * time.Second)

	client.Get("food", restock)
	client.Remove("food")
	client.Deregister("food")
}
