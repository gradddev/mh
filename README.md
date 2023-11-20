# mh

`mh` is a library written in Go (Golang) for controlling Magic Home LED Strip Controller.

## Installation

To install `mh` package, you need to install Go and set your Go workspace first.

The first need [Go](https://golang.org/) installed (**version 1.15+ is required**), then you can use the below Go command to install `mh`.

```sh
$ go get -u github.com/gradddev/mh
```
## Quick Start

```go
package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/gradddev/mh"
)

func main() {
	ip := net.ParseIP(os.Getenv("DEVICE_IP"))
	timeout := 3 * time.Second
	controller := mh.NewController(mh.Config{
		IP:      ip,
		Timeout: timeout,
	})
	rgbw, err := controller.GetRGBW()
	if err != nil {
		log.Panicln(err)
	} else {
		log.Println(rgbw)
	}
}
```

See `mh_test.go` for various usage examples.
