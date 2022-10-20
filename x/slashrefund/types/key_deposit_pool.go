package types

import "encoding/binary"
import sdk "github.com/cosmos/cosmos-sdk/types"

var _ binary.ByteOrder

const (
	// DepositPoolKeyPrefix is the prefix to retrieve all DepositPool
	DepositPoolKeyPrefix = "DepositPool/value/"
)

// DepositPoolKey returns the store key to retrieve a DepositPool from the index fields
func DepositPoolKey(
	operatorAddress sdk.ValAddress,
) []byte {
	var key []byte

	operatorAddressBytes := []byte(operatorAddress)
	key = append(key, operatorAddressBytes...)
	key = append(key, []byte("/")...)

	return key
}
