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
		return "Awesome", nil
	}

	fmt.Println(client.Register("name", time.Second, 2*time.Second))

	fmt.Println(client.Put("name", "Raed Shomali"))
	fmt.Println(client.Get("name", restock))

	time.Sleep(time.Second)

	fmt.Println(client.Get("name", restock))

	time.Sleep(2 * time.Second)

	fmt.Println(client.Get("name", restock))
	fmt.Println(client.Remove("name"))

	fmt.Println(client.Deregister("name"))
}
