package main

import (
	"fmt"
	"github.com/shomali11/fridge"
	"github.com/shomali11/xredis"
)

func main() {
	client := fridge.NewClient(xredis.DefaultClient())
	defer client.Close()

	fmt.Println(client.Put("food", "Pizza"))
	fmt.Println(client.Get("food"))
	fmt.Println(client.Remove("food"))
	fmt.Println(client.Get("food"))
}
