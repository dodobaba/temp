package main

import (
	"bytes"
	"fmt"
	"log"
	"shoppingzone/mylib/mynet"
)

func init() {
	fmt.Println("Init platform")
}

func main() {
	var server mynet.Server
	l := server.Start(nil)
	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			log.Printf("[Error] APP Server Failed to accept: %s", err.Error())
		}
		log.Printf("[Info] APP Server Accept from  %s", c.RemoteAddr().String())
		defer c.Close()
		//reader := bufio.NewReader(c)
		buff := new(bytes.Buffer)
		n, _ := buff.ReadFrom(c)
		log.Println(string(buff.Bytes()), "n :", n)
	}
}
