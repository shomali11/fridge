package main

import (
	"fmt"
	"github.com/shomali11/fridge"
	"github.com/shomali11/fridge/item"
	"time"
)

func main() {
	client := fridge.NewClient()
	defer client.Close()

	client.Register("food", item.WithDurations(time.Second, 2*time.Second))

	fmt.Println(client.Put("food", "Pizza"))
	fmt.Println(client.Get("food"))

	time.Sleep(time.Second)

	fmt.Println(client.Get("food"))

	time.Sleep(2 * time.Second)

	fmt.Println(client.Get("food"))
	fmt.Println(client.Remove("food"))

	client.Deregister("food")
}
