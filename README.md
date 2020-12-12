# mh

`mh` is a library written in Go (Golang) for controlling Magic Home LED Strip Controller.

## Installation

To install `mh` package, you need to install Go and set your Go workspace first.

The first need [Go](https://golang.org/) installed (**version 1.15+ is required**), then you can use the below Go command to install `mh`.

```sh
$ go get -u github.com/AlexeySemigradsky/mh
```
## Quick Start
```go
package main

import (
	"github.com/AlexeySemigradsky/mh"
	"log"
	"os"
)

func main() {
	address := os.Getenv("DEVICE_ADDRESS")
	controller := mh.NewController(address)
	rgbw, err := controller.GetRGBW()
	if err != nil {
		log.Panicln(err)
	}
	log.Println(rgbw)
}
```
