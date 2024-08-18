package main

import (
	"errors"
	"log"

	"./internal/ext4"
)

func main() {
	c, err := NewConn()
	if err != nil {
		log.Fatalf("failed to connect")
	}
	err = Listen(c)
}
