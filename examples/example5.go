package main

import (
	"fmt"
	"github.com/shomali11/fridge"
	"github.com/shomali11/fridge/item"
	"github.com/shomali11/xredis"
	"time"
)

func main() {
	client := fridge.NewClient(xredis.DefaultClient())
	defer client.Close()

	restock := func() (string, error) {
		return "Hot Pizza", nil
	}

	fmt.Println(client.Put("food", "Pizza", item.WithDurations(time.Second, 2*time.Second)))
	fmt.Println(client.Get("food", item.WithRestock(restock)))

	time.Sleep(time.Second)

	fmt.Println(client.Get("food", item.WithRestock(restock)))

	time.Sleep(2 * time.Second)

	fmt.Println(client.Get("food", item.WithRestock(restock)))
	fmt.Println(client.Remove("food"))
}
