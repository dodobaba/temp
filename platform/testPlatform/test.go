package main

import (
	"shoppingzone/mylib/mynet"
	"time"
)

func main() {
	var client mynet.Client
	t := time.Tick(2 * time.Microsecond)
	for tt := range t {
		m := tt.String()
		client.Send(m)
	}
}
