package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/shomali11/fridge"
)

func main() {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	client := fridge.NewClient(pool)
	defer client.Close()

	fmt.Println(client.Ping())
}
