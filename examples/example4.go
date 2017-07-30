package main

import (
	"fmt"
	"github.com/shomali11/fridge"
	"time"
)

func main() {
	client := fridge.DefaultClient()
	defer client.Close()

	restock := func() (string, error) {
		return "Hot Pizza", nil
	}

	fmt.Println(client.Register("food", time.Second, 2*time.Second))
	fmt.Println(client.Put("food", "Pizza"))
	fmt.Println(client.Get("food", restock))

	time.Sleep(time.Second)

	fmt.Println(client.Get("food", restock))

	time.Sleep(2 * time.Second)

	fmt.Println(client.Get("food", restock))
	fmt.Println(client.Remove("food"))
	fmt.Println(client.Deregister("food"))
}
