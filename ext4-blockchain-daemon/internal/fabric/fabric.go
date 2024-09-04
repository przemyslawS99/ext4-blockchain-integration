package fabric

import (
	"log"

	"github.com/przemyslawS99/ext4-blockchain-integration/ext4-blockchain-daemon/internal/common"
)

func NewInode(attributes *common.Attrs) uint16 {
	log.Printf("fabric: NewInode %v", attributes.Ino)
	return common.EXT4BD_STATUS_SUCCESS
}

func SetAttributes(attributes *common.Attrs) uint16 {
	log.Printf("fabric: SetAttributes %v", attributes.Ino)
	return common.EXT4BD_STATUS_INODE_NOT_FOUND
}

func GetAttributes(ino uint64) (uint16, *common.Attrs) {
	log.Printf("fabric: GetAttributes %v", ino)
	response := common.Attrs{
		Ino: ino,
	}
	/*response := common.Attrs{
		Uid: 1,
		Gid: 1,
		Atime: common.Time{
			Sec:  1,
			Nsec: 1,
		},
		Mtime: common.Time{
			Sec:  1,
			Nsec: 1,
		},
		Ctime: common.Time{
			Sec:  1,
			Nsec: 1,
		},
		Mode: 1,
		Ino:  ino,
	}*/
	return common.EXT4BD_STATUS_INODE_NOT_FOUND, &response
}
