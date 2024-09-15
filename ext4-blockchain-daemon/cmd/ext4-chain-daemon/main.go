package main

import (
	"log"
	"os"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	//"github.com/przemyslawS99/ext4-blockchain-integration/ext4-blockchain-daemon/internal/common"
	"github.com/przemyslawS99/ext4-blockchain-integration/ext4-blockchain-daemon/internal/ext4"
	"github.com/przemyslawS99/ext4-blockchain-integration/ext4-blockchain-daemon/internal/fabric"
)

func main() {
	clientConnection := fabric.NewGrpcConnection()
	defer clientConnection.Close()

	id := fabric.NewIdentity()
	sign := fabric.NewSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatalf("failed to connect")
	}
	defer gw.Close()

	chaincodeName := "ext4"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	channelName := "mychannel"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	connection, family, err := ext4.NewConn()
	if err != nil {
		log.Fatalf("failed to connect")
	}
	defer connection.Close()

	err = ext4.Listen(connection, family, contract)
}
