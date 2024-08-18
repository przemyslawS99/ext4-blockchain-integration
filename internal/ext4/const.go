package ext4

const (
	familyName = "ext4_chain"
	groupID    = 6
)

const (
	EXT4_CHAIN_CMD_SET_ATTR uint8 = iota
	EXT4_CHAIN_CMD_GET_ATTR
)

const (
	EXT4_CHAIN_ATTR_UNSPEC uint16 = iota
	EXT4_CHAIN_ATTR_UID
	EXT4_CHAIN_ATTR_GID
	EXT4_CHAIN_ATTR_ATIME
	EXT4_CHAIN_ATTR_MTIME
	EXT4_CHAIN_ATTR_CTIME
	EXT4_CHAIN_ATTR_SEC
	EXT4_CHAIN_ATTR_NSEC
	EXT4_CHAIN_ATTR_MODE
	EXT4_CHAIN_ATTR_INO
)
