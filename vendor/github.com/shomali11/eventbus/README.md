# eventbus [![Go Report Card](https://goreportcard.com/badge/github.com/shomali11/eventbus)](https://goreportcard.com/report/github.com/shomali11/eventbus) [![GoDoc](https://godoc.org/github.com/shomali11/eventbus?status.svg)](https://godoc.org/github.com/shomali11/eventbus) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

An event bus to facilitate publishing and subscribing to events via topics

## Features

* Non blocking publishing of events

## Usage

Using `govendor` [github.com/kardianos/govendor](https://github.com/kardianos/govendor):

```
govendor fetch github.com/shomali11/eventbus
```

## Dependencies

* `maps` [github.com/shomali11/maps](https://github.com/shomali11/maps)

# Examples

## Example 1

Using `NewClient` to create an eventbus.

```go
package main

import "github.com/shomali11/eventbus"

func main() {
	client := eventbus.NewClient()
	defer client.Close()
}
```

## Example 2

Using `Publish` and `Subscribe` to publish and listen to events

```go
package main

import (
	"fmt"
	"github.com/shomali11/eventbus"
	"time"
)

func main() {
	client := eventbus.NewClient()
	defer client.Close()

	client.Publish("test", "test")

	client.Subscribe("name", func(value interface{}) {
		fmt.Println(value)
	})

	client.Subscribe("name", func(value interface{}) {
		fmt.Println(value, "is Awesome")
	})

	client.Publish("name", "Raed Shomali")

	time.Sleep(time.Second)
}
```

Output:

```
Raed Shomali
Raed Shomali is Awesome
```