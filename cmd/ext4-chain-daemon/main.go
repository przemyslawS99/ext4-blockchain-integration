package main

import (
	"log"

	"github.com/przemyslawS99/ext4-blockchain-integration/internal/ext4"
)

func main() {
	c, family, err := ext4.NewConn()
	if err != nil {
		log.Fatalf("failed to connect")
	}
	defer c.Close()

	err = ext4.Listen(c, family)
}
