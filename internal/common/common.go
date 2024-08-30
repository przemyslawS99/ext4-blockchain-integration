package common

import (
	"github.com/mdlayher/netlink"
)

type Time struct {
	Sec  uint64
	Nsec uint32
}

type Attrs struct {
	Uid   uint32
	Gid   uint32
	Atime Time
	Mtime Time
	Ctime Time
	Mode  uint32
	Ino   uint64
}

func (n *Time) DecodeTime(ad *netlink.AttributeDecoder) error {
	for ad.Next() {
		switch ad.Type() {
		case EXT4B_TIME_ATTR_SEC:
			n.Sec = ad.Uint64()
		case EXT4B_TIME_ATTR_NSEC:
			n.Nsec = ad.Uint32()
		}
	}
	return nil
}
