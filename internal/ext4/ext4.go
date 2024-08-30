package ext4

import (
	"errors"
	"log"
	"os"

	"github.com/mdlayher/genetlink"
	"github.com/mdlayher/netlink"
	"github.com/przemyslawS99/ext4-blockchain-integration/internal/common"
	"github.com/przemyslawS99/ext4-blockchain-integration/internal/fabric"
)

func NewConn() (*genetlink.Conn, genetlink.Family, error) {
	c, err := genetlink.Dial(nil)
	if err != nil {
		log.Fatalf("failed to dial generic netlink: %v", err)
	}

	family, err := c.GetFamily(common.FamilyName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("%q family not available", common.FamilyName)
		}
		log.Fatalf("failed to query for family: %v", err)
	}

	log.Printf("%s: %+v", common.FamilyName, family)

	err = sendSetPid(c, family)
	if err != nil {
		log.Fatalf("failed to send setpid: %v", err)
	}
	/*if err := c.JoinGroup(uint32(common.GroupID)); err != nil {
		log.Fatalf("failed to join multicast group: %v", err)
	}*/

	return c, family, nil
}

func Listen(c *genetlink.Conn, family genetlink.Family) error {
	for {
		msgs, _, err := c.Receive()
		if err != nil {
			log.Fatalf("failed to receive message: %v", err)
		}

		for _, msg := range msgs {
			ad, err := netlink.NewAttributeDecoder(msg.Data)
			if err != nil {
				log.Fatalf("failed to create attribute decoder: %v", err)
			}

			var attributes common.Attrs
			for ad.Next() {
				switch ad.Type() {
				case common.EXT4B_ATTR_UID:
					attributes.Uid = ad.Uint32()
				case common.EXT4B_ATTR_GID:
					attributes.Gid = ad.Uint32()
				case common.EXT4B_ATTR_ATIME:
					ad.Nested(attributes.Atime.DecodeTime)
				case common.EXT4B_ATTR_MTIME:
					ad.Nested(attributes.Mtime.DecodeTime)
				case common.EXT4B_ATTR_CTIME:
					ad.Nested(attributes.Ctime.DecodeTime)
				case common.EXT4B_ATTR_MODE:
					attributes.Mode = ad.Uint32()
				case common.EXT4B_ATTR_INO:
					attributes.Ino = ad.Uint64()
				}
			}

			if msg.Header.Command == common.EXT4B_CMD_SETATTR_REQUEST {
				status := fabric.SetAttributes(&attributes)
				err := sendSetattrResponse(c, family, attributes.Ino, status)
				if err != nil {
					log.Printf("failed to send: ino: %v, status: %v", attributes.Ino, status)
				}
			}
			/*if msg.Header.Command == common.EXT4_CHAIN_CMD_SETATTR_REQUEST {
				/*err := fabric.SetAttributes(&attributes)
				if err != nil {
					log.Printf("failed to set attributes")
				} else {
					err := sendSetattrSuccess(c, family, attributes.Ino)
					if err != nil {
						log.Printf("failed to send SETATTR_SUCCESS")
					}
				}
				err := sendSetattrSuccess(c, family, attributes.Ino)
				if err != nil {
					log.Printf("failed to send SETATTRSUCCESS")
				}
				//log.Printf("EXT4_CHAIN_CMD_SETATTR")
			} else if msg.Header.Command == common.EXT4_CHAIN_CMD_GETATTR {
				err := fabric.GetAttributes(attributes.Ino)
				if err != nil {
					log.Printf("failed to get attributes")
				}
			} else {
				log.Printf("UNKNOWN")
			}*/
		}
	}
}

func sendSetPid(c *genetlink.Conn, family genetlink.Family) error {
	msg := genetlink.Message{
		Header: genetlink.Header{
			Command: common.EXT4B_CMD_SETPID,
			Version: family.Version,
		},
		Data: nil,
	}

	_, err := c.Send(msg, family.ID, netlink.Request)
	if err != nil {
		return err
	}
	log.Printf("sendSetPid: request sent")
	return nil
}

func sendSetattrResponse(c *genetlink.Conn, family genetlink.Family, ino uint64, status uint16) error {
	ae := netlink.NewAttributeEncoder()
	ae.Uint64(common.EXT4B_ATTR_INO, ino)
	ae.Uint16(common.EXT4B_ATTR_STATUS, status)

	b, err := ae.Encode()
	if err != nil {
		log.Printf("failed to encode attributes: %v", err)
	}

	msg := genetlink.Message{
		Header: genetlink.Header{
			Command: common.EXT4B_CMD_SETATTR_RESPONSE,
			Version: family.Version,
		},
		Data: b,
	}

	_, err = c.Send(msg, family.ID, netlink.Request)
	if err != nil {
		return err
	}
	log.Printf("sendSetattrResponse: ino: %v, status_code: %v", ino, status)
	return nil
}

/*func sendSetattrSuccess(c *genetlink.Conn, family genetlink.Family, ino uint64) error {
	ae := netlink.NewAttributeEncoder()
	ae.Uint64(common.EXT4_CHAIN_ATTR_INO, ino)
	ae.Flag(common.EXT4_CHAIN_ATTR_SETATTR_SUCCESS, true)

	b, err := ae.Encode()
	if err != nil {
		log.Printf("failed to encode attributes: %v", err)
	}

	msg := genetlink.Message{
		Header: genetlink.Header{
			Command: common.EXT4_CHAIN_CMD_SETATTR_HANDLE_RESPONSE,
			Version: family.Version,
		},
		Data: b,
	}

	_, err = c.Send(msg, family.ID, netlink.Request)
	if err != nil {
		return err
	}
	log.Printf("sendSetattrSuccess: %v", ino)
	return nil
}*/
