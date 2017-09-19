# fridge [![Go Report Card](https://goreportcard.com/badge/github.com/shomali11/fridge)](https://goreportcard.com/report/github.com/shomali11/fridge) [![GoDoc](https://godoc.org/github.com/shomali11/fridge?status.svg)](https://godoc.org/github.com/shomali11/fridge) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`fridge` is a layer applied on top of a cache that makes interacting with it similar to interacting with a fridge.
Items are tagged with a "Best By" and "Use By" timestamps, stored, restocked and retrieved.

Typically when using a cache, one would store some value along with a timeout.
The value could be retrieved from the cache as long as it has not expired.
If the value had expired, then an external call _(Such as a database query)_ is usually made to retrieve the value, put it back in the cache and return it.

With `fridge`, we are taking a slightly different approach.
Before storing a value in the `fridge`, one tags its key with a **Best By** and a **Use By** durations.
When retrieving the value, a **Restock** function can be provided to refresh the value.

When attempting to retrieve a value from the `fridge`, there are multiple scenarios that could happen:
 * If the item has not passed its **Best By** duration _(it is **fresh**)_
   * Then the item is returned immediately.
 * If the item has passed its **Best By** duration but not its **Use By** duration _(Not **fresh** but not **expired** either)_
   * Then the item is returned immediately 
   * But the **Restock** function is called **asynchronously** to **refresh** the item.
 * If the item has passed its **Use By** duration _(it has **expired**)_
   * The **Restock** function is called **synchronously** to retrieve a fresh item and return it.
 * If the item was not found
   * It is treated similarly to an expired item
   * The **Restock** function is called **synchronously** to retrieve a fresh item and return it.

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

* `eventbus` [github.com/shomali11/eventbus](https://github.com/shomali11/eventbus)
* `xredis` [github.com/shomali11/xredis](https://github.com/shomali11/xredis)
* `util` [github.com/shomali11/util](https://github.com/shomali11/util)


# Examples

## Example 1

Using `NewClient` to create a new fridge client.
_Note: `NewClient` accepts an object that implements the `Cache` interface which allows the user to use `fridge ` with any underlying implementation._

```go
// Cache is a Fridge cache interface
type Cache interface {
	// Get a value by key
	Get(key string) (string, bool, error)

	// Set a key value pair
	Set(key string, value string, timeout time.Duration) error

	// Remove a key
	Remove(key string) error

	// Ping to test connectivity
	Ping() error

	// Close to close resources
	Close() error
}
```

```go
package main

import (
	"fmt"
	"github.com/shomali11/fridge"
	"time"
)

type SimpleCache struct {
	memory map[string]string
}

func (c *SimpleCache) Get(key string) (string, bool, error) {
	value, ok := c.memory[key]
	return value, ok, nil
}

func (c *SimpleCache) Set(key string, value string, timeout time.Duration) error {
	// We are not implementing the expiration to keep the example simple
	c.memory[key] = value
	return nil
}

func (c *SimpleCache) Remove(key string) error {
	delete(c.memory, key)
	return nil
}

func (c *SimpleCache) Ping() error {
	return nil
}

func (c *SimpleCache) Close() error {
	return nil
}

func main() {
	simpleCache := &SimpleCache{memory: make(map[string]string)}
	client := fridge.NewClient(simpleCache)

	fmt.Println(client.Ping())
}
```

## Example 2

Using `NewRedisCache` to use `fridge` with redis. _Note: `NewRedisCache` creates an redis client that implements the `Cache` interface_

```go
package main

import (
	"fmt"
	"github.com/shomali11/fridge"
)

func main() {
	redisCache := fridge.NewRedisCache()
	client := fridge.NewClient(redisCache)
	defer client.Close()

	fmt.Println(client.Ping())
}
```

Output

```
<nil>
```

## Example 3

Using `NewRedisCache` with modified settings for a redis client

```go
package main

import (
	"fmt"
	"github.com/shomali11/fridge"
)

func main() {
	redisCache := fridge.NewRedisCache(
		fridge.WithHost("localhost"), 
		fridge.WithPort(6379))
		
	client := fridge.NewClient(redisCache)
	defer client.Close()

	fmt.Println(client.Ping())
}
```

Output

```
<nil>
```

## Example 4

Using `NewSentinelCache` with modified settings for a redis sentinel client. _Note: `NewSentinelCache` creates an redis sentinel client that implements the `Cache` interface_


```go
package main

import (
	"fmt"
	"github.com/shomali11/fridge"
)

func main() {
	redisCache := fridge.NewSentinelCache(
		fridge.WithSentinelAddresses([]string{"localhost:26379"}),
		fridge.WithSentinelMasterName("master"))

	client := fridge.NewClient(redisCache)
	defer client.Close()

	fmt.Println(client.Ping())
}
```

Output

```
<nil>
```

## Example 5

Using `Put`, `Get` & `Remove` to show how to put, get and remove an item.
_Note: That we are using a default client that has a default Best By of 1 hour and Use By of 1 Day for all keys_

```go
package main

import (
	"fmt"
	"github.com/shomali11/fridge"
)

func main() {
	redisCache := fridge.NewRedisCache()
	client := fridge.NewClient(redisCache)
	defer client.Close()

	fmt.Println(client.Put("food", "Pizza"))
	fmt.Println(client.Get("food"))
	fmt.Println(client.Remove("food"))
	fmt.Println(client.Get("food"))
}
```

Output

```
<nil>
Pizza true <nil>
<nil>
 false <nil>
```

## Example 6

Using `WithDefaultDurations` to override the default Best By and Use By durations for all keys

```go
package main

import (
	"fmt"
	"github.com/shomali11/fridge"
	"time"
)

func main() {
	redisCache := fridge.NewRedisCache()
	client := fridge.NewClient(redisCache, fridge.WithDefaultDurations(time.Second, 2*time.Second))
	defer client.Close()

	fmt.Println(client.Put("food", "Pizza"))
	fmt.Println(client.Get("food"))

	time.Sleep(time.Second)

	fmt.Println(client.Get("food"))

	time.Sleep(2 * time.Second)

	fmt.Println(client.Get("food"))
	fmt.Println(client.Remove("food"))
}
```

Output

```
<nil>
Pizza true <nil>
Pizza true <nil>
 false <nil>
<nil>
```

## Example 7

Using `Put` to show how to put an item and override that item's durations.

```go
package main

import (
	"fmt"
	"github.com/shomali11/fridge"
	"time"
)

func main() {
	redisCache := fridge.NewRedisCache()
	client := fridge.NewClient(redisCache)
	defer client.Close()

	fmt.Println(client.Put("food", "Pizza", fridge.WithDurations(time.Second, 2*time.Second)))
	fmt.Println(client.Get("food"))

	time.Sleep(time.Second)

	fmt.Println(client.Get("food"))

	time.Sleep(2 * time.Second)

	fmt.Println(client.Get("food"))
	fmt.Println(client.Remove("food"))
}
```

Output

```
<nil>
Pizza true <nil>
Pizza true <nil>
 false <nil>
<nil>
```

## Example 8

Using `Get` to show how to retrieve an item while providing a restocking mechanism.

```go
package main

import (
	"fmt"
	"github.com/shomali11/fridge"
	"time"
)

func main() {
	redisCache := fridge.NewRedisCache()
	client := fridge.NewClient(redisCache)
	defer client.Close()

	restock := func() (string, error) {
		return "Hot Pizza", nil
	}

	fmt.Println(client.Put("food", "Pizza", fridge.WithDurations(time.Second, 2*time.Second)))
	fmt.Println(client.Get("food", fridge.WithRestock(restock)))

	time.Sleep(time.Second)

	fmt.Println(client.Get("food", fridge.WithRestock(restock)))

	time.Sleep(2 * time.Second)

	fmt.Println(client.Get("food", fridge.WithRestock(restock)))
	fmt.Println(client.Remove("food"))
}
```

Output

```
<nil>
Pizza true <nil>
Pizza true <nil>
Hot Pizza true <nil>
<nil>
```

## Example 9

Using `HandleEvent` to pass a callback to access the stream of events generated

```go
package main

import (
	"fmt"
	"github.com/shomali11/fridge"
	"time"
)

func main() {
	redisCache := fridge.NewRedisCache()
	client := fridge.NewClient(redisCache, fridge.WithDefaultDurations(time.Second, 2*time.Second))
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

	client.Put("food1", "Pizza")
	client.Put("food2", "Milk")

	client.Get("food1", fridge.WithRestock(restock))
	client.Get("food2")
	client.Get("food3")

	time.Sleep(time.Second)

	client.Get("food1", fridge.WithRestock(restock))
	client.Get("food2")
	client.Get("food3")

	time.Sleep(2 * time.Second)

	client.Get("food1", fridge.WithRestock(restock))
	client.Get("food2")
	client.Get("food3")

	client.Remove("food1")
	client.Remove("food2")
	client.Remove("food3")
}
```

Output:

```
Key: food1 - Woohoo! it is fresh!
Key: food2 - Woohoo! it is fresh!
Key: food3 - Oops! Did not find it.
Key: food1 - Not fresh! But not bad either!
Key: food1 - Yay! Getting a new one!
Key: food1 - Interesting! It has not changed.
Key: food2 - Not fresh! But not bad either!
Key: food2 - Oh no! It is out of stock.
Key: food3 - Oops! Did not find it.
Key: food1 - Sigh. It has expired!
Key: food1 - Yay! Getting a new one!
Key: food2 - Sigh. It has expired!
Key: food2 - Oh no! It is out of stock.
Key: food3 - Oops! Did not find it.
```
