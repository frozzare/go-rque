# go-rque

Simple queue build with RethinkDB

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
	"github.com/frozzare/go-rque/worker"
	"github.com/frozzare/go-rque/shared"
)

func main() {
	e := emitter.New()

	e.On("hello", func(job shared.Job) {
		fmt.Println("hello")
	})

	worker.Run(shared.Config{
		Address:  "localhost:28015",
		Database: "db",
		Table:    "queue",
		Emitter:  e,
	})
}
```

# License

 MIT © [Fredrik Forsmo](https://github.com/frozzare)
