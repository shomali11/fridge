# fridge [![Go Report Card](https://goreportcard.com/badge/github.com/shomali11/fridge)](https://goreportcard.com/report/github.com/shomali11/fridge) [![GoDoc](https://godoc.org/github.com/shomali11/fridge?status.svg)](https://godoc.org/github.com/shomali11/fridge) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`fridge` is a redis cache that resembles storing items in a fridge and retrieving them later on.

Typically when using a cache, one would store some value with a TTL.
The value could be retrieved from the cache as long as it has not expired.
If the value had expired, then a database call is usually made to retrieve the value, put it back in the cache and return it.

With `fridge`, we are taking a slightly different approach.
Before storing a value in the fridge _(cache)_, one must register its key with a "Best By" and a "Use By" durations.
When retrieving the value from the fridge _(cache)_, a "restock" function can be provided to refresh the value.

When attempting to retrieve a value from the fridge _(cache)_, there are multiple scenarios that could happen:
 * If the item has not passed its "Best By" duration _(it is "fresh")_
   * Then the item is returned immediately.
 * If the item has passed its "Best By" duration but not its "Use By" duration _(Not "fresh" but not "expired" either)_
   * Then the item is returned immediately 
   * But the "restock" function is called **asynchronously** to "refresh" the item.
 * If the item has passed its "Use By" duration _(it has "expired")_
   * The "restock" function is called **synchronously** to retrieve a fresh item and return it.
 * If the item was not found
   * It is treated similarly to an expired item
   * The "restock" function is called **synchronously** to retrieve a fresh item and return it.

## Why?

The thinking behind `fridge` is to increase the chances for a value to be retrieved from the cache.
The longer the value stays in the cache, the better the chances are to retrieve it faster (As opposed to from the database)

The challenge, of course, is to keep the value in the cache "fresh".

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

Using `NewClient` with default options

```go
package main

import (
	"fmt"
	"github.com/shomali11/fridge"
)

func main() {
	client := fridge.NewClient()
	defer client.Close()

	fmt.Println(client.Ping()) // <nil>
}
```

## Example 2

Using `NewClient` to with a custom `xredis` client

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

	xredisClient := xredis.SetupClient(options)
	client := fridge.NewClient(fridge.WithRedisClient(xredisClient))
	defer client.Close()

	fmt.Println(client.Ping()) // <nil>
}
```

## Example 3

Using the `Put`, `Get` & `Remove` to show how to put, get and remove an item.
_Note: That we are using a default client that has a default Best By of 1 hour and Use By of 1 Day for all keys_

```go
package main

import (
	"fmt"
	"github.com/shomali11/fridge"
)

func main() {
	client := fridge.NewClient()
	defer client.Close()

	fmt.Println(client.Put("food", "Pizza")) // <nil>
	fmt.Println(client.Get("food"))          // "Pizza" true <nil>
	fmt.Println(client.Remove("food"))       // <nil>
	fmt.Println(client.Get("food"))          // "" false <nil>
}
```

## Example 4

Using the `WithDefaultDurations` to override the default Best By and Use By durations for all keys

```go
package main

import (
	"fmt"
	"github.com/shomali11/fridge"
	"time"
)

func main() {
	client := fridge.NewClient(fridge.WithDefaultDurations(time.Second, 2*time.Second))
	defer client.Close()

	fmt.Println(client.Put("food", "Pizza")) // <nil>
	fmt.Println(client.Get("food"))          // "Pizza" true <nil>

	time.Sleep(time.Second)

	fmt.Println(client.Get("food"))          // "Pizza" true <nil>

	time.Sleep(2 * time.Second)

	fmt.Println(client.Get("food"))          // "" false <nil>
	fmt.Println(client.Remove("food"))       // <nil>
}
```

## Example 5

Using the `Register` & `Deregister` to show how to register an item and override that item's durations.

```go
package main

import (
	"fmt"
	"github.com/shomali11/fridge"
	"github.com/shomali11/fridge/item"
	"time"
)

func main() {
	client := fridge.NewClient()
	defer client.Close()

	client.Register("food", item.WithDurations(time.Second, 2*time.Second))

	fmt.Println(client.Put("food", "Pizza")) // <nil>
	fmt.Println(client.Get("food"))          // "Pizza" true <nil>

	time.Sleep(time.Second)

	fmt.Println(client.Get("food"))          // "Pizza" true <nil>

	time.Sleep(2 * time.Second)

	fmt.Println(client.Get("food"))          // "" false <nil>
	fmt.Println(client.Remove("food"))       // <nil>

	client.Deregister("food")
}
```

## Example 6

Using the `Register` & `Deregister` to show how to register & deregister an item to set that item's restocking mechanism.

```go
package main

import (
	"fmt"
	"github.com/shomali11/fridge"
	"github.com/shomali11/fridge/item"
	"time"
)

func main() {
	client := fridge.NewClient()
	defer client.Close()

	restock := func() (string, error) {
		return "Hot Pizza", nil
	}

	client.Register("food", item.WithDurations(time.Second, 2*time.Second), item.WithRestock(restock))

	fmt.Println(client.Put("food", "Pizza")) // <nil>
	fmt.Println(client.Get("food"))          // "Pizza" true <nil>

	time.Sleep(time.Second)

	fmt.Println(client.Get("food"))          // "Pizza" true <nil>

	time.Sleep(2 * time.Second)

	fmt.Println(client.Get("food"))          // "Hot Pizza" true <nil>
	fmt.Println(client.Remove("food"))       // <nil>

	client.Deregister("food")
}
```

## Example 7

Using the `HandleEvent` to pass a callback to access the stream of events generated

```go
package main

import (
	"fmt"
	"github.com/shomali11/fridge"
	"github.com/shomali11/fridge/item"
	"time"
)

func main() {
	client := fridge.NewClient()
	defer client.Close()

	client.HandleEvent(func(event *fridge.Event) {
		fmt.Print("Key: " + event.Key + " - ")

		switch event.Type {
		case fridge.Fresh:
			fmt.Println("Woohoo! it is fresh!")
		case fridge.Cold:
			fmt.Println("Not fresh! But not bad either!")
		case fridge.Expired:
			fmt.Println("Sigh. It has expired!")
		case fridge.NotFound:
			fmt.Println("Oops! Did not find it.")
		case fridge.Refresh:
			fmt.Println("Yay! Getting a new one!")
		case fridge.OutOfStock:
			fmt.Println("Oh no! It is out of stock.")
		case fridge.Unchanged:
			fmt.Println("Interesting! It has not changed.")
		}
	})

	restock := func() (string, error) {
		return "Pizza", nil
	}

	client.Register("food1", item.WithDurations(time.Second, 2*time.Second), item.WithRestock(restock))
	client.Register("food2", item.WithDurations(time.Second, 2*time.Second))
	client.Register("food3", item.WithDurations(time.Second, 2*time.Second))

	client.Put("food1", "Pizza")
	client.Put("food2", "Milk")

	client.Get("food1")
	client.Get("food2")
	client.Get("food3")

	time.Sleep(time.Second)

	client.Get("food1")
	client.Get("food2")
	client.Get("food3")

	time.Sleep(2 * time.Second)

	client.Get("food1")
	client.Get("food2")
	client.Get("food3")

	client.Remove("food1")
	client.Remove("food2")
	client.Remove("food3")
	
	client.Deregister("food1")
	client.Deregister("food2")
	client.Deregister("food3")
}
```

Output:

```
Key: food1 - Woohoo! it is fresh!
Key: food2 - Woohoo! it is fresh!
Key: food3 - Oops! Did not find it.
Key: food3 - Oh no! It is out of stock.
Key: food1 - Not fresh! But not bad either!
Key: food1 - Yay! Getting a new one!
Key: food1 - Interesting! It has not changed.
Key: food2 - Not fresh! But not bad either!
Key: food2 - Oh no! It is out of stock.
Key: food3 - Oops! Did not find it.
Key: food3 - Oh no! It is out of stock.
Key: food1 - Sigh. It has expired!
Key: food1 - Yay! Getting a new one!
Key: food2 - Sigh. It has expired!
Key: food2 - Oh no! It is out of stock.
Key: food3 - Oops! Did not find it.
Key: food3 - Oh no! It is out of stock.
```
