package fabric

import (
	"log"

	"github.com/przemyslawS99/ext4-blockchain-integration/internal/common"
)

func SetAttributes(attributes *common.Attrs) uint16 {
	log.Printf("fabric: SetAttributes %v", attributes.Ino)
	return 0
}

func GetAttributes(ino uint64) uint16 {
	log.Printf("fabric: GetAttributes %v", ino)
	return 0
}
