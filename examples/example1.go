package main

import (
	"fmt"
	"github.com/shomali11/fridge"
)

func main() {
	client := fridge.NewClient()
	defer client.Close()

	fmt.Println(client.Ping())
}
