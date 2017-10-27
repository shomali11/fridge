# maps [![Go Report Card](https://goreportcard.com/badge/github.com/shomali11/maps)](https://goreportcard.com/report/github.com/shomali11/maps) [![GoDoc](https://godoc.org/github.com/shomali11/maps?status.svg)](https://godoc.org/github.com/shomali11/maps) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A collection of various maps with concurrent versions that supports concurrent reads and writes.

## Features

* Thread safe
* Available types:
    * Maps:
        * Multi Map
    * Concurrent Maps:
        * Concurrent Map
        * Concurrent Multi Map
    * Sharded Concurrent Maps:
        * Sharded Concurrent Map
        * Sharded Concurrent Multi Map

_Note: Sharded versions provides improved performance by reducing the number of write locks_

## Usage

Using `govendor` [github.com/kardianos/govendor](https://github.com/kardianos/govendor):

```
govendor fetch github.com/shomali11/maps
```

## Dependencies

* `util` [github.com/shomali11/util](https://github.com/shomali11/util)

# Examples

## Example 1

Using `NewConcurrentMap` to create a concurrent map

```go
package main

import (
	"fmt"
	"github.com/shomali11/maps"
)

func main() {
	concurrentMap := maps.NewConcurrentMap()
	concurrentMap.Set("name", "Raed Shomali")

	fmt.Println(concurrentMap.ContainsKey("name")) // true
	fmt.Println(concurrentMap.Get("name"))      // "Raed Shomali" true
	fmt.Println(concurrentMap.Size())           // 1
	fmt.Println(concurrentMap.IsEmpty())        // false
	fmt.Println(concurrentMap.Keys())           // ["name"]

	concurrentMap.Remove("name")

	fmt.Println(concurrentMap.ContainsKey("name")) // false
	fmt.Println(concurrentMap.Get("name"))      // <nil> false
	fmt.Println(concurrentMap.Size())           // 0
	fmt.Println(concurrentMap.IsEmpty())        // true
	fmt.Println(concurrentMap.Keys())           // []

	concurrentMap.Set("name", "Raed Shomali")
	concurrentMap.Clear()

	fmt.Println(concurrentMap.ContainsKey("name")) // false
	fmt.Println(concurrentMap.Get("name"))      // <nil> false
	fmt.Println(concurrentMap.Size())           // 0
	fmt.Println(concurrentMap.IsEmpty())        // true
	fmt.Println(concurrentMap.Keys())           // []
}
```

## Example 2

Using `NewShardedConcurrentMap` to create a sharded concurrent map. _Default shards are 16_

```go
package main

import (
	"fmt"
	"github.com/shomali11/maps"
)

func main() {
	concurrentMap := maps.NewShardedConcurrentMap()
	concurrentMap.Set("name", "Raed Shomali")

	fmt.Println(concurrentMap.ContainsKey("name")) // true
	fmt.Println(concurrentMap.Get("name"))      // "Raed Shomali" true
	fmt.Println(concurrentMap.Size())           // 1
	fmt.Println(concurrentMap.IsEmpty())        // false
	fmt.Println(concurrentMap.Keys())           // ["name"]

	concurrentMap.Remove("name")

	fmt.Println(concurrentMap.ContainsKey("name")) // false
	fmt.Println(concurrentMap.Get("name"))      // <nil> false
	fmt.Println(concurrentMap.Size())           // 0
	fmt.Println(concurrentMap.IsEmpty())        // true
	fmt.Println(concurrentMap.Keys())           // []

	concurrentMap.Set("name", "Raed Shomali")
	concurrentMap.Clear()

	fmt.Println(concurrentMap.ContainsKey("name")) // false
	fmt.Println(concurrentMap.Get("name"))      // <nil> false
	fmt.Println(concurrentMap.Size())           // 0
	fmt.Println(concurrentMap.IsEmpty())        // true
	fmt.Println(concurrentMap.Keys())           // []
}
```

## Example 3

Using `WithNumberOfShards` to override default number of shards

```go
package main

import (
	"fmt"
	"github.com/shomali11/maps"
)

func main() {
	concurrentMap := maps.NewShardedConcurrentMap(maps.WithNumberOfShards(100))
	concurrentMap.Set("name", "Raed Shomali")

	fmt.Println(concurrentMap.ContainsKey("name")) // true
	fmt.Println(concurrentMap.Get("name"))      // "Raed Shomali" true
	fmt.Println(concurrentMap.Size())           // 1
	fmt.Println(concurrentMap.IsEmpty())        // false
	fmt.Println(concurrentMap.Keys())           // ["name"]

	concurrentMap.Remove("name")

	fmt.Println(concurrentMap.ContainsKey("name")) // false
	fmt.Println(concurrentMap.Get("name"))      // <nil> false
	fmt.Println(concurrentMap.Size())           // 0
	fmt.Println(concurrentMap.IsEmpty())        // true
	fmt.Println(concurrentMap.Keys())           // []

	concurrentMap.Set("name", "Raed Shomali")
	concurrentMap.Clear()

	fmt.Println(concurrentMap.ContainsKey("name")) // false
	fmt.Println(concurrentMap.Get("name"))      // <nil> false
	fmt.Println(concurrentMap.Size())           // 0
	fmt.Println(concurrentMap.IsEmpty())        // true
	fmt.Println(concurrentMap.Keys())           // []
}
```

## Example 4

Using `NewConcurrentMultiMap` to create a concurrent multi map

```go
package main

import (
	"fmt"
	"github.com/shomali11/maps"
)

func main() {
	concurrentMap := maps.NewConcurrentMultiMap()
	concurrentMap.Set("names", []interface{}{"Raed Shomali"})
	concurrentMap.Append("names", "Dwayne Johnson")

	fmt.Println(concurrentMap.ContainsKey("names")) // true
	fmt.Println(concurrentMap.Get("names"))      // ["Raed Shomali" "Dwayne Johnson"] true
	fmt.Println(concurrentMap.Size())            // 1
	fmt.Println(concurrentMap.IsEmpty())         // false
	fmt.Println(concurrentMap.Keys())            // ["name"]

	concurrentMap.Remove("names")

	fmt.Println(concurrentMap.ContainsKey("names")) // false
	fmt.Println(concurrentMap.Get("names"))      // [] false
	fmt.Println(concurrentMap.Size())            // 0
	fmt.Println(concurrentMap.IsEmpty())         // true
	fmt.Println(concurrentMap.Keys())            // []

	concurrentMap.Append("names", "Raed Shomali")
	concurrentMap.Clear()

	fmt.Println(concurrentMap.ContainsKey("names")) // false
	fmt.Println(concurrentMap.Get("names"))      // [] false
	fmt.Println(concurrentMap.Size())            // 0
	fmt.Println(concurrentMap.IsEmpty())         // true
	fmt.Println(concurrentMap.Keys())            // []
}
```

## Example 5

Using `NewShardedConcurrentMultiMap` to create a sharded concurrent multi map. _Default shards are 16_

```go
package main

import (
	"fmt"
	"github.com/shomali11/maps"
)

func main() {
	concurrentMap := maps.NewShardedConcurrentMultiMap()
	concurrentMap.Set("names", []interface{}{"Raed Shomali"})
	concurrentMap.Append("names", "Dwayne Johnson")

	fmt.Println(concurrentMap.ContainsKey("names")) // true
	fmt.Println(concurrentMap.Get("names"))      // ["Raed Shomali" "Dwayne Johnson"] true
	fmt.Println(concurrentMap.Size())            // 1
	fmt.Println(concurrentMap.IsEmpty())         // false
	fmt.Println(concurrentMap.Keys())            // ["name"]

	concurrentMap.Remove("names")

	fmt.Println(concurrentMap.ContainsKey("names")) // false
	fmt.Println(concurrentMap.Get("names"))      // [] false
	fmt.Println(concurrentMap.Size())            // 0
	fmt.Println(concurrentMap.IsEmpty())         // true
	fmt.Println(concurrentMap.Keys())            // []

	concurrentMap.Append("names", "Raed Shomali")
	concurrentMap.Clear()

	fmt.Println(concurrentMap.ContainsKey("names")) // false
	fmt.Println(concurrentMap.Get("names"))      // [] false
	fmt.Println(concurrentMap.Size())            // 0
	fmt.Println(concurrentMap.IsEmpty())         // true
	fmt.Println(concurrentMap.Keys())            // []
}
```

## Example 6

Using `WithNumberOfShards` to override default number of shards

```go
package main

import (
	"fmt"
	"github.com/shomali11/maps"
)

func main() {
	concurrentMap := maps.NewShardedConcurrentMultiMap(maps.WithNumberOfShards(100))
	concurrentMap.Set("names", []interface{}{"Raed Shomali"})
	concurrentMap.Append("names", "Dwayne Johnson")

	fmt.Println(concurrentMap.ContainsKey("names")) // true
	fmt.Println(concurrentMap.Get("names"))      // ["Raed Shomali" "Dwayne Johnson"] true
	fmt.Println(concurrentMap.Size())            // 1
	fmt.Println(concurrentMap.IsEmpty())         // false
	fmt.Println(concurrentMap.Keys())            // ["name"]

	concurrentMap.Remove("names")

	fmt.Println(concurrentMap.ContainsKey("names")) // false
	fmt.Println(concurrentMap.Get("names"))      // [] false
	fmt.Println(concurrentMap.Size())            // 0
	fmt.Println(concurrentMap.IsEmpty())         // true
	fmt.Println(concurrentMap.Keys())            // []

	concurrentMap.Append("names", "Raed Shomali")
	concurrentMap.Clear()

	fmt.Println(concurrentMap.ContainsKey("names")) // false
	fmt.Println(concurrentMap.Get("names"))      // [] false
	fmt.Println(concurrentMap.Size())            // 0
	fmt.Println(concurrentMap.IsEmpty())         // true
	fmt.Println(concurrentMap.Keys())            // []
}
```