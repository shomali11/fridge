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

	fmt.Println(client.Put("food", "Pizza", item.WithDurations(time.Second, 2*time.Second)))
	fmt.Println(client.Get("food"))

	time.Sleep(time.Second)

	fmt.Println(client.Get("food"))

	time.Sleep(2 * time.Second)

	fmt.Println(client.Get("food"))
	fmt.Println(client.Remove("food"))
}
