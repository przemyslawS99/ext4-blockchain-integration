package common

import (
	"fmt"
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

func DecodeAttributes(data []byte) (*Attrs, error) {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return nil, err
	}

	var attributes Attrs
	for ad.Next() {
		switch ad.Type() {
		case EXT4B_ATTR_UID:
			attributes.Uid = ad.Uint32()
		case EXT4B_ATTR_GID:
			attributes.Gid = ad.Uint32()
		case EXT4B_ATTR_ATIME:
			ad.Nested(attributes.Atime.DecodeTime)
		case EXT4B_ATTR_MTIME:
			ad.Nested(attributes.Mtime.DecodeTime)
		case EXT4B_ATTR_CTIME:
			ad.Nested(attributes.Ctime.DecodeTime)
		case EXT4B_ATTR_MODE:
			attributes.Mode = ad.Uint32()
		case EXT4B_ATTR_INO:
			attributes.Ino = ad.Uint64()
		}
	}
	return &attributes, nil
}

func DecodeIno(data []byte) (uint64, error) {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return 0, err
	}

	if ad.Next() && ad.Type() == EXT4B_ATTR_INO {
		return ad.Uint64(), nil
	}

	err = fmt.Errorf("expected EXT4B_ATTR_INO, but got something else or no attributes")
	return 0, err
}

func (n *Time) EncodeTime(ae *netlink.AttributeEncoder) {
	ae.Uint64(EXT4B_TIME_ATTR_SEC, n.Sec)
	ae.Uint32(EXT4B_TIME_ATTR_NSEC, n.Nsec)
}
