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

	return c, family, nil
}

func Listen(c *genetlink.Conn, family genetlink.Family) error {
	for {
		msgs, _, err := c.Receive()
		if err != nil {
			log.Fatalf("failed to receive message: %v", err)
		}

		for _, msg := range msgs {
			switch msg.Header.Command {
			case common.EXT4B_CMD_NEW_INODE_REQUEST:
				attributes, err := common.DecodeAttributes(msg.Data)
				if err != nil {
					log.Fatalf("failed to decode attributes: %v", err)
				}
				status := fabric.NewInode(attributes)
				err = sendStatusResponse(c, family, attributes.Ino, status)
				if err != nil {
					log.Printf("failed to send: ino=%v, status=%v", attributes.Ino, status)
				}

			case common.EXT4B_CMD_SETATTR_REQUEST:
				attributes, err := common.DecodeAttributes(msg.Data)
				if err != nil {
					log.Fatalf("failed to decode attributes: %v", err)
				}
				status := fabric.SetAttributes(attributes)
				err = sendStatusResponse(c, family, attributes.Ino, status)
				if err != nil {
					log.Printf("failed to send: ino=%v, status=%v", attributes.Ino, status)
				}

			case common.EXT4B_CMD_GETATTR_REQUEST:
				ino, err := common.DecodeIno(msg.Data)
				if err != nil {
					log.Fatalf("failed to decode ino: %v", err)
				}
				status, attributes := fabric.GetAttributes(ino)
				err = sendGetAttributesResponse(c, family, attributes, status)
				if err != nil {
					log.Printf("failed to send: ino=%v, status=%v", ino, status)
				}
			}
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

func sendStatusResponse(c *genetlink.Conn, family genetlink.Family, ino uint64, status uint16) error {
	ae := netlink.NewAttributeEncoder()
	ae.Uint64(common.EXT4B_ATTR_INO, ino)
	ae.Uint16(common.EXT4B_ATTR_STATUS, status)

	b, err := ae.Encode()
	if err != nil {
		log.Printf("failed to encode attributes: %v", err)
	}

	msg := genetlink.Message{
		Header: genetlink.Header{
			Command: common.EXT4B_CMD_STATUS_RESPONSE,
			Version: family.Version,
		},
		Data: b,
	}

	_, err = c.Send(msg, family.ID, netlink.Request)
	if err != nil {
		return err
	}
	log.Printf("sendStatusResponse: ino=%v, status=%v", ino, status)
	return nil
}

func sendGetAttributesResponse(c *genetlink.Conn, family genetlink.Family, response *common.Attrs, status uint16) error {
	ae := netlink.NewAttributeEncoder()

	ae.Uint16(common.EXT4B_ATTR_STATUS, status)
	ae.Uint64(common.EXT4B_ATTR_INO, response.Ino)

	if status == common.EXT4BD_STATUS_SUCCESS {
		ae.Uint32(common.EXT4B_ATTR_MODE, response.Mode)
		ae.Uint32(common.EXT4B_ATTR_UID, response.Uid)
		ae.Uint32(common.EXT4B_ATTR_GID, response.Gid)

		ae.Nested(common.EXT4B_ATTR_ATIME, func(nae *netlink.AttributeEncoder) error {
			nae.Uint64(common.EXT4B_TIME_ATTR_SEC, response.Atime.Sec)
			nae.Uint32(common.EXT4B_TIME_ATTR_NSEC, response.Atime.Nsec)
			return nil
		})

		ae.Nested(common.EXT4B_ATTR_MTIME, func(nae *netlink.AttributeEncoder) error {
			nae.Uint64(common.EXT4B_TIME_ATTR_SEC, response.Mtime.Sec)
			nae.Uint32(common.EXT4B_TIME_ATTR_NSEC, response.Mtime.Nsec)
			return nil
		})

		ae.Nested(common.EXT4B_ATTR_CTIME, func(nae *netlink.AttributeEncoder) error {
			nae.Uint64(common.EXT4B_TIME_ATTR_SEC, response.Ctime.Sec)
			nae.Uint32(common.EXT4B_TIME_ATTR_NSEC, response.Ctime.Nsec)
			return nil
		})
	}

	b, err := ae.Encode()
	if err != nil {
		log.Printf("failed to encode attributes: %v", err)
	}

	msg := genetlink.Message{
		Header: genetlink.Header{
			Command: common.EXT4B_CMD_GETATTR_RESPONSE,
			Version: family.Version,
		},
		Data: b,
	}

	_, err = c.Send(msg, family.ID, netlink.Request)
	if err != nil {
		return err
	}
	log.Printf("sendGetattrResponse: status=%v", status)
	return nil
}
