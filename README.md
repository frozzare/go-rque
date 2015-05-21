# go-rque

Simple queue build with RethinkDB. Requires [go-emitter](http://github.com/frozzare/go-emitter).

View the [docs](http://godoc.org/github.com/frozzare/go-rque).

## Installation

```
$ go get github.com/frozzare/go-rque
```

## Example

```go
package main

import (
	"fmt"

	"github.com/frozzare/go-emitter"
	"github.com/frozzare/go-rque"
)

func main() {
	e := emitter.New()

	e.On("hello", func(job rque.Job) {
		fmt.Println("hello")
	})

	rque.Run(rque.Config{
		Address:  "localhost:28015",
		Database: "db",
		Table:    "queue",
		Emitter:  e,
	})
}
```

JSON row in the database:

```json
{
    "name": "hello",
    "data": {}
}
```

# License

 MIT © [Fredrik Forsmo](https://github.com/frozzare)
