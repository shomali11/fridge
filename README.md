# fridge [![Go Report Card](https://goreportcard.com/badge/github.com/shomali11/fridge)](https://goreportcard.com/report/github.com/shomali11/fridge) [![GoDoc](https://godoc.org/github.com/shomali11/fridge?status.svg)](https://godoc.org/github.com/shomali11/fridge) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`fridge` is a redis cache that resembles storing items in a fridge and retrieving them later on.

Typically when using a cache, one would store some value with a TTL.
The value could be retrieved from the cache as long as it has not expired.
If the value had expired, then a database call is usually made to retrieve the value, put it back in the cache and return it.

With `fridge`, we are taking a slightly different approach.
Before storing a value in the fridge (cache), one must register its key with a "Best By" and a "Use By" durations.
When retrieving the value from the fridge (cache), a "restock" function can be provided as input.

The idea is when one attempts to retrieve a value from the fridge (cache), if the item has not passed its "Best By" duration (it is "fresh"), then the item is returned immediately.
If the item has passed its "Best By" duration but not its "Use By" duration (Not "fresh" but not "expired" either), then the item is returned immediately but the "restock" function is called in a separate go-routine to "refresh" the item in the background.
If the item has passed its "Use By" duration (it has "expired"), the "restock" function is called to retrieve a fresh item and return it.
If the item was not found, it is treated similarly to an expired item and the "restock" function is called to retrieve a fresh item and return it.

## Usage

Using `govendor` [github.com/kardianos/govendor](https://github.com/kardianos/govendor):

```
govendor fetch github.com/shomali11/fridge
```

## Dependencies

* `xredis` [github.com/shomali11/xredis](https://github.com/shomali11/xredis)
* `util` [github.com/shomali11/util](https://github.com/shomali11/util)


# Examples

## Example 1

Using `DefaultClient` to create a client with default options

```go
package main

import (
	"fmt"
	"github.com/shomali11/fridge"
)

func main() {
	client := fridge.DefaultClient()
	defer client.Close()

	fmt.Println(client.Ping()) // <nil>
}
```

List of default options

```text
defaultHost                  = "localhost"
defaultPort                  = 6379
defaultPassword              = ""
defaultDatabase              = 0
defaultNetwork               = "tcp"
defaultConnectTimeout        = time.Second
defaultWriteTimeout          = time.Second
defaultReadTimeout           = time.Second
defaultConnectionIdleTimeout = 240 * time.Second
defaultConnectionMaxIdle     = 100
defaultConnectionMaxActive   = 10000
defaultConnectionWait        = false
defaultTlsConfig             = nil
defaultTlsSkipVerify         = false
defaultTestOnBorrowTimeout   = time.Minute
```

## Example 2

Using `SetupClient` to create a client using provided options

```go
package main

import (
	"fmt"
	"github.com/shomali11/fridge"
	"github.com/shomali11/xredis"
)

func main() {
	options := &xredis.Options{
		Host: "localhost",
		Port: 6379,
	}

	client := fridge.SetupClient(options)
	defer client.Close()

	fmt.Println(client.Ping()) // <nil>
}
```

Available options to set

```go
type Options struct {
	Host                  string
	Port                  int
	Password              string
	Database              int
	Network               string
	ConnectTimeout        time.Duration
	WriteTimeout          time.Duration
	ReadTimeout           time.Duration
	ConnectionIdleTimeout time.Duration
	ConnectionMaxIdle     int
	ConnectionMaxActive   int
	ConnectionWait        bool
	TlsConfig             *tls.Config
	TlsSkipVerify         bool
	TestOnBorrowPeriod    time.Duration
}
```

## Example 3

Using `NewClient` to create a client using `redigo`'s `redis.Pool`

```go
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

	fmt.Println(client.Ping()) // <nil>
}
```

## Example 4

Using the `Register`, `Put`, `Get`, `Remove` & `Deregister` to show how to register, put, get remove and deregister an item

```go
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

	fmt.Println(client.Register("name", time.Second, 2*time.Second)) // <nil>
	fmt.Println(client.Put("name", "Raed Shomali"))                  // <nil>
	fmt.Println(client.Get("name", restock))                         // "Raed Shomali" true <nil>

	time.Sleep(time.Second)

	fmt.Println(client.Get("name", restock))                         // "Raed Shomali" true <nil>

	time.Sleep(2 * time.Second)

	fmt.Println(client.Get("name", restock))                         // "Awesome" true <nil>
	fmt.Println(client.Remove("name"))                               // <nil>
	fmt.Println(client.Deregister("name"))                           // <nil>
}
```
