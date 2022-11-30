package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ binary.ByteOrder

const (
	// RefundKeyPrefix is the prefix to retrieve all Refund
	RefundKeyPrefix = "Refund/value/"
)

// RefundKey returns the store key to retrieve a Refund from the index fields
func RefundKey(
	delegator sdk.AccAddress,
	validator sdk.ValAddress,
) []byte {
	var key []byte

	delegatorBytes := []byte(delegator)
	key = append(key, delegatorBytes...)
	key = append(key, []byte("/")...)

	validatorBytes := []byte(validator)
	key = append(key, validatorBytes...)
	key = append(key, []byte("/")...)

	return key
}
