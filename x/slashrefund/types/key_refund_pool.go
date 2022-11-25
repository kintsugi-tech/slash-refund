package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// RefundPoolKeyPrefix is the prefix to retrieve all RefundPool
	RefundPoolKeyPrefix = "RefundPool/value/"
)

// RefundPoolKey returns the store key to retrieve a RefundPool from the index fields
func RefundPoolKey(
	operatorAddress string,
) []byte {
	var key []byte

	operatorAddressBytes := []byte(operatorAddress)
	key = append(key, operatorAddressBytes...)
	key = append(key, []byte("/")...)

	return key
}
