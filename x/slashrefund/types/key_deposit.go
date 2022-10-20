package types

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ binary.ByteOrder

const (
	// DepositKeyPrefix is the prefix to retrieve all Deposit
	DepositKeyPrefix = "Deposit/value/"
)

// DepositKey returns the store key to retrieve a Deposit from the index fields
func DepositKey(
	depAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
) []byte {
	var key []byte

	depAddrBytes := []byte(depAddr)
	key = append(key, depAddrBytes...)
	key = append(key, []byte("/")...)

	valAddrBytes := []byte(valAddr)
	key = append(key, valAddrBytes...)
	key = append(key, []byte("/")...)

	return key
}
