package main

import (
	"fmt"
	"github.com/shomali11/fridge"
)

func main() {
	client := fridge.DefaultClient()
	defer client.Close()

	fmt.Println(client.Ping())
}
