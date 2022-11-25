package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ binary.ByteOrder

const (
	// RefundPoolKeyPrefix is the prefix to retrieve all RefundPool
	RefundPoolKeyPrefix = "RefundPool/value/"
)

// RefundPoolKey returns the store key to retrieve a RefundPool from the index fields
func RefundPoolKey(
	operatorAddress sdk.ValAddress,
) []byte {
	var key []byte

	operatorAddressBytes := []byte(operatorAddress)
	key = append(key, operatorAddressBytes...)
	key = append(key, []byte("/")...)

	return key
}
