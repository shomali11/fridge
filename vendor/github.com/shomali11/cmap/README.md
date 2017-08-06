# cmap [![Go Report Card](https://goreportcard.com/badge/github.com/shomali11/cmap)](https://goreportcard.com/report/github.com/shomali11/cmap) [![GoDoc](https://godoc.org/github.com/shomali11/cmap?status.svg)](https://godoc.org/github.com/shomali11/cmap) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Map that supports concurrent reads and writes.

## Features

* Thread safe
* Two Concurrent Map types:
    * Simple Concurrent Map
    * Sharded Concurrent Map
        * Provides improved performance by reducing the number of write locks

## Usage

Using `govendor` [github.com/kardianos/govendor](https://github.com/kardianos/govendor):

```
govendor fetch github.com/shomali11/cmap
```

## Dependencies

* `util` [github.com/shomali11/util](https://github.com/shomali11/util)

# Examples

## Example 1

Using `NewConcurrentMap` to create concurrent map

```go
package main

import (
	"fmt"
	"github.com/shomali11/cmap"
)

func main() {
	concurrentMap := cmap.NewConcurrentMap()
	concurrentMap.Set("name", "Raed Shomali")

	fmt.Println(concurrentMap.Contains("name")) // true
	fmt.Println(concurrentMap.Get("name"))      // "Raed Shomali" true
	fmt.Println(concurrentMap.Size())           // 1

	concurrentMap.Remove("name")

	fmt.Println(concurrentMap.Contains("name")) // false
	fmt.Println(concurrentMap.Get("name"))      // <nil> false
	fmt.Println(concurrentMap.Size())           // 0

	concurrentMap.Set("name", "Raed Shomali")
	concurrentMap.Clear()

	fmt.Println(concurrentMap.Contains("name")) // false
	fmt.Println(concurrentMap.Get("name"))      // <nil> false
	fmt.Println(concurrentMap.Size())           // 0
}
```

## Example 2

Using `NewShardedConcurrentMap` to create a sharded concurrent map. _Default shards are 16_

```go
package main

import (
	"fmt"
	"github.com/shomali11/cmap"
)

func main() {
	concurrentMap := cmap.NewShardedConcurrentMap()
	concurrentMap.Set("name", "Raed Shomali")

	fmt.Println(concurrentMap.Contains("name")) // true
	fmt.Println(concurrentMap.Get("name"))      // "Raed Shomali" true
	fmt.Println(concurrentMap.Size())           // 1

	concurrentMap.Remove("name")

	fmt.Println(concurrentMap.Contains("name")) // false
	fmt.Println(concurrentMap.Get("name"))      // <nil> false
	fmt.Println(concurrentMap.Size())           // 0

	concurrentMap.Set("name", "Raed Shomali")
	concurrentMap.Clear()

	fmt.Println(concurrentMap.Contains("name")) // false
	fmt.Println(concurrentMap.Get("name"))      // <nil> false
	fmt.Println(concurrentMap.Size())           // 0
}
```

## Example 3

Using `WithNumberOfShards` to override default number of shards

```go
package main

import (
	"fmt"
	"github.com/shomali11/cmap"
)

func main() {
	concurrentMap := cmap.NewShardedConcurrentMap(cmap.WithNumberOfShards(100))
	concurrentMap.Set("name", "Raed Shomali")

	fmt.Println(concurrentMap.Contains("name")) // true
	fmt.Println(concurrentMap.Get("name"))      // "Raed Shomali" true
	fmt.Println(concurrentMap.Size())           // 1

	concurrentMap.Remove("name")

	fmt.Println(concurrentMap.Contains("name")) // false
	fmt.Println(concurrentMap.Get("name"))      // <nil> false
	fmt.Println(concurrentMap.Size())           // 0

	concurrentMap.Set("name", "Raed Shomali")
	concurrentMap.Clear()

	fmt.Println(concurrentMap.Contains("name")) // false
	fmt.Println(concurrentMap.Get("name"))      // <nil> false
	fmt.Println(concurrentMap.Size())           // 0
}
```