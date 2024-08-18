package ext4

import (
	"errors"
	"log"
	"os"

	"github.com/mdlayher/genetlink"
	"github.com/mdlayher/netlink"
)

type time struct {
	sec  uint64
	nsec uint32
}

type attrs struct {
	uid   uint32
	gid   uint32
	atime time
	mtime time
	ctime time
	mode  uint32
	ino   uint64
}

func NewConn() (*genetlink.Conn, error) {
	c, err := genetlink.Dial(nil)
	if err != nil {
		log.Fatalf("failed to dial generic netlink: %v", err)
	}
	defer c.Close()

	family, err := c.GetFamily(familyName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("%q family not available", familyName)
		}
		log.Fatalf("failed to query for family: %v", err)
	}

	log.Printf("%s: %+v", familyName, family)

	if err := c.JoinGroup(uint32(groupID)); err != nil {
		log.Fatalf("failed to join multicast group: %v", err)
	}

	return c, nil
}

func (n *time) decodeTime(ad *netlink.AttributeDecoder) error {
	for ad.Next() {
		switch ad.Type() {
		case EXT4_CHAIN_ATTR_SEC:
			n.sec = ad.Uint64()
		case EXT4_CHAIN_ATTR_NSEC:
			n.nsec = ad.Uint32()
		}
	}
	return nil
}

func Listen(c *genetlink.Conn) error {
	for {
		msgs, _, err := c.Receive()
		if err != nil {
			log.Fatalf("failed to receive message: %v", err)
		}

		for _, msg := range msgs {
			log.Printf("messsage command: %v", msg.Header.Command)
			if msg.Header.Command == EXT4_CHAIN_CMD_SET_ATTR {
				log.Printf("EXT4_CHAIN_CMD_SET_ATTR")
			} else if msg.Header.Command == EXT4_CHAIN_CMD_GET_ATTR {
				log.Printf("EXT4_CHAIN_CMD_GET_ATTR")
			} else {
				log.Printf("UNKNOWN")
			}
			ad, err := netlink.NewAttributeDecoder(msg.Data)
			if err != nil {
				log.Fatalf("failed to create attribute decoder: %v", err)
			}

			var attributes attrs
			for ad.Next() {
				switch ad.Type() {
				case EXT4_CHAIN_ATTR_UID:
					attributes.uid = ad.Uint32()
				case EXT4_CHAIN_ATTR_GID:
					attributes.gid = ad.Uint32()
				case EXT4_CHAIN_ATTR_ATIME:
					ad.Nested(attributes.atime.decodeTime)
				case EXT4_CHAIN_ATTR_MTIME:
					ad.Nested(attributes.mtime.decodeTime)
				case EXT4_CHAIN_ATTR_CTIME:
					ad.Nested(attributes.ctime.decodeTime)
				case EXT4_CHAIN_ATTR_MODE:
					attributes.mode = ad.Uint32()
				case EXT4_CHAIN_ATTR_INO:
					attributes.ino = ad.Uint64()
				}
			}
		}
	}
}
