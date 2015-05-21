# rque

Simple queue built with RethinkDB.

View the [docs](http://godoc.org/github.com/frozzare/go-rque).

## Installation

```
$ go get github.com/frozzare/go-rque
```

## Example

```go
que, err := rque.New(rque.Config{
	Address:  "localhost:28015",
	Database: "test",
	Table:    "queue",
})

if err != nil {
	log.Fatalf("Failed to create que: %s", err)
}

for job := range que.Jobs() {
	fmt.Printf("Hello %s\n", job.Name)
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
