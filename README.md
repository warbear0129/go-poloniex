go-poloniex
==========

go-poloniex is an implementation of the Poloniex API (public and private) in Golang.

## Import
	import "github.com/warbear0129/go-poloniex"
	
## Usage
~~~ go
package main

import (
	"fmt"
	"github.com/jyap808/go-poloniex"
)

const (
	API_KEY    = "YOUR_API_KEY"
	API_SECRET = "YOUR_API_SECRET"
)

func main() {
	// Poloniex client
	poloniex := poloniex.New(API_KEY, API_SECRET)

	// Get tickers
    tickers, err := poloniex.GetTickers()
	fmt.Println(err, tickers)
}
~~~
